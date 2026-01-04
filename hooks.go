package sum

import (
	"context"
	"time"

	"github.com/zoobzio/capitan"
	"github.com/zoobzio/flume"
	"github.com/zoobzio/pipz"
)

// HookKey identifies a lifecycle hook point.
type HookKey string

// Lifecycle hook keys.
const (
	BeforeCreate HookKey = "beforeCreate"
	AfterCreate  HookKey = "afterCreate"
	BeforeUpdate HookKey = "beforeUpdate"
	AfterUpdate  HookKey = "afterUpdate"
	BeforeDelete HookKey = "beforeDelete"
	AfterDelete  HookKey = "afterDelete"
)

// allHookKeys lists all defined hook keys for initialization.
var allHookKeys = []HookKey{
	BeforeCreate, AfterCreate,
	BeforeUpdate, AfterUpdate,
	BeforeDelete, AfterDelete,
}

// HookFunc is the signature for simple hook functions.
// The function receives the entity and can modify it.
// Return the (possibly modified) entity or an error to abort.
type HookFunc[T pipz.Cloner[T]] func(ctx context.Context, entity T) (T, error)

// Hooks manages lifecycle hook pipelines for a Resource.
type Hooks[T pipz.Cloner[T]] struct {
	factory    *flume.Factory[T]
	identities map[HookKey]pipz.Identity
	typeName   string
	counter    int // For generating unique processor names
}

// NewHooks creates a new Hooks manager using the provided factory.
func NewHooks[T pipz.Cloner[T]](factory *flume.Factory[T], typeName string) *Hooks[T] {
	h := &Hooks[T]{
		factory:    factory,
		identities: make(map[HookKey]pipz.Identity),
		typeName:   typeName,
	}

	// Pre-register identities for all hook keys
	for _, key := range allHookKeys {
		h.identities[key] = h.factory.Identity(
			string(key),
			descriptionFor(key, typeName),
		)
	}

	return h
}

// descriptionFor generates a human-readable description for a hook key.
func descriptionFor(key HookKey, typeName string) string {
	switch key {
	case BeforeCreate:
		return "Executes before creating a " + typeName
	case AfterCreate:
		return "Executes after creating a " + typeName
	case BeforeUpdate:
		return "Executes before updating a " + typeName
	case AfterUpdate:
		return "Executes after updating a " + typeName
	case BeforeDelete:
		return "Executes before deleting a " + typeName
	case AfterDelete:
		return "Executes after deleting a " + typeName
	default:
		return "Hook for " + typeName
	}
}

// Identity returns the pipz.Identity for a specific hook key.
func (h *Hooks[T]) Identity(key HookKey) pipz.Identity {
	return h.identities[key]
}

// Register adds a HookFunc to the specified hook point.
// Multiple hooks for the same key are executed in sequence.
// Bindings are created lazily on first Execute() or explicitly via Build().
func (h *Hooks[T]) Register(key HookKey, fn HookFunc[T]) {
	h.counter++
	name := string(key) + "-" + string(rune('a'+h.counter-1))

	// Create identity for this processor
	id := h.factory.Identity(name, "User hook for "+string(key))

	// Wrap HookFunc in a pipz processor
	processor := pipz.Apply(id, func(ctx context.Context, data T) (T, error) {
		return fn(ctx, data)
	})

	// Add to factory
	h.factory.Add(processor)

	// Emit registration signal
	capitan.Debug(context.Background(), HookRegistered,
		KeyHook.Field(string(key)),
		KeyResource.Field(h.typeName))
}

// Build creates schemas and bindings for all hook points with registered processors.
// Call this after registration to finalize pipelines, or let Execute() build lazily.
// After Build(), users can modify schemas via Factory().SetSchema() for full flume control.
func (h *Hooks[T]) Build() {
	for _, key := range allHookKeys {
		h.buildBinding(key)
	}
}

// buildBinding creates the binding for a hook key if processors exist.
// It creates a sequence of all registered processors for that key.
func (h *Hooks[T]) buildBinding(key HookKey) {
	// Get all processors that start with this key prefix
	processors := h.factory.ListProcessors()
	var refs []flume.Node
	prefix := string(key) + "-"

	for _, name := range processors {
		if len(name) > len(prefix) && name[:len(prefix)] == prefix {
			refs = append(refs, flume.Node{Ref: name})
		}
	}

	if len(refs) == 0 {
		return
	}

	// Build schema with sequence of processors
	schemaID := string(key)
	schema := flume.Schema{
		Version: "1",
		Node: flume.Node{
			Type:     "sequence",
			Name:     string(key),
			Children: refs,
		},
	}

	// Register/update the schema (auto-sync bindings will rebuild)
	_ = h.factory.SetSchema(schemaID, schema)

	// Create binding if it doesn't exist (with auto-sync for future updates)
	id := h.identities[key]
	if h.factory.Get(id) == nil {
		_, _ = h.factory.Bind(id, schemaID, flume.WithAutoSync[T]())
	}
}

// Execute runs the hook pipeline for the given key.
// Returns the processed data or an error.
// If no hooks are registered for this key, returns the data unchanged.
// Lazily builds the binding on first call if not already built.
func (h *Hooks[T]) Execute(ctx context.Context, key HookKey, data T) (T, error) {
	id := h.identities[key]
	binding := h.factory.Get(id)
	if binding == nil {
		// Try lazy build
		h.buildBinding(key)
		binding = h.factory.Get(id)
		if binding == nil {
			// No hooks registered for this key
			return data, nil
		}
	}

	start := time.Now()

	capitan.Debug(ctx, HookExecuting,
		KeyHook.Field(string(key)),
		KeyResource.Field(h.typeName))

	result, err := binding.Process(ctx, data)
	duration := time.Since(start).Milliseconds()

	if err != nil {
		capitan.Warn(ctx, HookFailed,
			KeyHook.Field(string(key)),
			KeyResource.Field(h.typeName),
			KeyDurationMs.Field(duration),
			KeyError.Field(err.Error()))
		return data, err
	}

	capitan.Debug(ctx, HookCompleted,
		KeyHook.Field(string(key)),
		KeyResource.Field(h.typeName),
		KeyDurationMs.Field(duration))

	return result, nil
}

// ExecuteAfter runs an afterX hook with error-tolerant semantics.
// Errors are logged via capitan but do not propagate.
// Returns the original data on error, or the processed data on success.
func (h *Hooks[T]) ExecuteAfter(ctx context.Context, key HookKey, data T) T {
	result, err := h.Execute(ctx, key, data)
	if err != nil {
		capitan.Warn(ctx, HookAfterError,
			KeyHook.Field(string(key)),
			KeyResource.Field(h.typeName),
			KeyError.Field(err.Error()))
		return data
	}
	return result
}

// HasHooks returns true if any hooks are registered for the given key.
func (h *Hooks[T]) HasHooks(key HookKey) bool {
	id := h.identities[key]
	return h.factory.Get(id) != nil
}
