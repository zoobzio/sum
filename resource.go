package sum

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/zoobzio/atom"
	"github.com/zoobzio/capitan"
	"github.com/zoobzio/flume"
	"github.com/zoobzio/pipz"
	"github.com/zoobzio/rocco"
	"github.com/zoobzio/sentinel"
)

// Resource represents a type-safe resource with auto-generated CRUD endpoints.
// T must implement pipz.Cloner[T] for pipeline support.
type Resource[T pipz.Cloner[T]] struct {
	atomizer       *atom.Atomizer[T]
	basePath       string
	handlers       []rocco.Endpoint
	tag            string
	tagDescription string
	factory        *flume.Factory[T]
	hooks          *Hooks[T]
}

// New creates a new Resource[T] and registers CRUD endpoints with the singleton server.
// Panics if the type T cannot be atomized.
func New[T pipz.Cloner[T]]() *Resource[T] {
	ensureInitialized()

	atomizer, err := atom.Use[T]()
	if err != nil {
		panic(fmt.Sprintf("sum: failed to create atomizer for %T: %v", *new(T), err))
	}

	typeName := atomizer.Spec().TypeName
	basePath := "/" + strings.ToLower(typeName) + "s"
	tag := strings.ToLower(typeName) + "s"

	factory := flume.New[T]()

	r := &Resource[T]{
		atomizer:       atomizer,
		basePath:       basePath,
		handlers:       make([]rocco.Endpoint, 0, 5),
		tag:            tag,
		tagDescription: fmt.Sprintf("Operations for %s resources", typeName),
		factory:        factory,
		hooks:          NewHooks[T](factory, typeName),
	}

	r.generateCRUD()
	registerResource(r)

	return r
}

// WithTag sets the OpenAPI tag name and description for this resource's endpoints.
// Returns the resource for chaining.
func (r *Resource[T]) WithTag(name, description string) *Resource[T] {
	r.tag = name
	r.tagDescription = description
	r.generateCRUD()
	return r
}

// Atomizer returns the atomizer for this resource's type.
func (r *Resource[T]) Atomizer() *atom.Atomizer[T] {
	return r.atomizer
}

// Meta returns the sentinel metadata for this resource's type.
func (r *Resource[T]) Meta() sentinel.Metadata {
	return r.atomizer.Spec()
}

// BasePath returns the base path for this resource's endpoints.
func (r *Resource[T]) BasePath() string {
	return r.basePath
}

// Tag returns the OpenAPI tag name and description for this resource.
func (r *Resource[T]) Tag() (name, description string) {
	return r.tag, r.tagDescription
}

// endpoints implements the resource interface.
func (r *Resource[T]) endpoints() []rocco.Endpoint {
	return r.handlers
}

// tagInfo implements the resource interface.
func (r *Resource[T]) tagInfo() (name, description string) {
	return r.tag, r.tagDescription
}

// Factory returns the flume.Factory[T] for full pipeline control.
func (r *Resource[T]) Factory() *flume.Factory[T] {
	return r.factory
}

// HooksManager returns the underlying Hooks[T] for hook configuration.
func (r *Resource[T]) HooksManager() *Hooks[T] {
	return r.hooks
}

// HookID returns the pipz.Identity for a specific hook key.
func (r *Resource[T]) HookID(key HookKey) pipz.Identity {
	return r.hooks.Identity(key)
}

// BeforeCreate registers a hook to run before create operations.
func (r *Resource[T]) BeforeCreate(fn HookFunc[T]) *Resource[T] {
	r.hooks.Register(BeforeCreate, fn)
	return r
}

// AfterCreate registers a hook to run after create operations.
func (r *Resource[T]) AfterCreate(fn HookFunc[T]) *Resource[T] {
	r.hooks.Register(AfterCreate, fn)
	return r
}

// BeforeUpdate registers a hook to run before update operations.
func (r *Resource[T]) BeforeUpdate(fn HookFunc[T]) *Resource[T] {
	r.hooks.Register(BeforeUpdate, fn)
	return r
}

// AfterUpdate registers a hook to run after update operations.
func (r *Resource[T]) AfterUpdate(fn HookFunc[T]) *Resource[T] {
	r.hooks.Register(AfterUpdate, fn)
	return r
}

// BeforeDelete registers a hook to run before delete operations.
func (r *Resource[T]) BeforeDelete(fn HookFunc[T]) *Resource[T] {
	r.hooks.Register(BeforeDelete, fn)
	return r
}

// AfterDelete registers a hook to run after delete operations.
func (r *Resource[T]) AfterDelete(fn HookFunc[T]) *Resource[T] {
	r.hooks.Register(AfterDelete, fn)
	return r
}

// generateCRUD creates the standard CRUD endpoints for this resource.
func (r *Resource[T]) generateCRUD() {
	typeName := r.atomizer.Spec().TypeName
	lowerName := strings.ToLower(typeName)

	// POST /resources - Create
	create := rocco.NewHandler[T, T](
		fmt.Sprintf("create-%s", lowerName),
		http.MethodPost,
		r.basePath,
		func(req *rocco.Request[T]) (T, error) {
			// Execute beforeCreate hooks (fail-fast)
			entity, err := r.hooks.Execute(req.Context, BeforeCreate, req.Body)
			if err != nil {
				var zero T
				return zero, err
			}

			// Stub: storage operation would go here
			result := entity

			// Execute afterCreate hooks (error-tolerant)
			result = r.hooks.ExecuteAfter(req.Context, AfterCreate, result)

			// Emit resource created signal
			capitan.Info(req.Context, ResourceCreated,
				KeyResource.Field(typeName),
				KeyOperation.Field("create"))

			return result, nil
		},
	).WithSummary(fmt.Sprintf("Create a new %s", typeName)).
		WithDescription(fmt.Sprintf("Creates a new %s resource with the provided data.", typeName)).
		WithTags(r.tag).
		WithSuccessStatus(http.StatusCreated)

	// GET /resources - List
	list := rocco.NewHandler[rocco.NoBody, []T](
		fmt.Sprintf("list-%ss", lowerName),
		http.MethodGet,
		r.basePath,
		func(req *rocco.Request[rocco.NoBody]) ([]T, error) {
			// Stub: return empty list
			result := []T{}

			// Emit resource listed signal
			capitan.Debug(req.Context, ResourceListed,
				KeyResource.Field(typeName),
				KeyOperation.Field("list"),
				KeyCount.Field(int64(len(result))))

			return result, nil
		},
	).WithSummary(fmt.Sprintf("List all %ss", typeName)).
		WithDescription(fmt.Sprintf("Retrieves a list of all %s resources.", typeName)).
		WithTags(r.tag)

	// GET /resources/{id} - Get
	get := rocco.NewHandler[rocco.NoBody, T](
		fmt.Sprintf("get-%s", lowerName),
		http.MethodGet,
		r.basePath+"/{id}",
		func(req *rocco.Request[rocco.NoBody]) (T, error) {
			// Stub: return zero value
			var zero T

			// Emit resource retrieved signal
			capitan.Debug(req.Context, ResourceRetrieved,
				KeyResource.Field(typeName),
				KeyOperation.Field("get"),
				KeyEntityID.Field(req.Params.Path["id"]))

			return zero, nil
		},
	).WithSummary(fmt.Sprintf("Get a %s by ID", typeName)).
		WithDescription(fmt.Sprintf("Retrieves a single %s resource by its unique identifier.", typeName)).
		WithTags(r.tag).
		WithPathParams("id")

	// PUT /resources/{id} - Update
	update := rocco.NewHandler[T, T](
		fmt.Sprintf("update-%s", lowerName),
		http.MethodPut,
		r.basePath+"/{id}",
		func(req *rocco.Request[T]) (T, error) {
			// Execute beforeUpdate hooks (fail-fast)
			entity, err := r.hooks.Execute(req.Context, BeforeUpdate, req.Body)
			if err != nil {
				var zero T
				return zero, err
			}

			// Stub: storage operation would go here
			result := entity

			// Execute afterUpdate hooks (error-tolerant)
			result = r.hooks.ExecuteAfter(req.Context, AfterUpdate, result)

			// Emit resource updated signal
			capitan.Info(req.Context, ResourceUpdated,
				KeyResource.Field(typeName),
				KeyOperation.Field("update"),
				KeyEntityID.Field(req.Params.Path["id"]))

			return result, nil
		},
	).WithSummary(fmt.Sprintf("Update a %s", typeName)).
		WithDescription(fmt.Sprintf("Replaces an existing %s resource with the provided data.", typeName)).
		WithTags(r.tag).
		WithPathParams("id")

	// DELETE /resources/{id} - Delete
	del := rocco.NewHandler[rocco.NoBody, struct{}](
		fmt.Sprintf("delete-%s", lowerName),
		http.MethodDelete,
		r.basePath+"/{id}",
		func(req *rocco.Request[rocco.NoBody]) (struct{}, error) {
			// For delete hooks, pass zero value T (storage would load entity first)
			var entity T

			// Execute beforeDelete hooks (fail-fast)
			_, err := r.hooks.Execute(req.Context, BeforeDelete, entity)
			if err != nil {
				return struct{}{}, err
			}

			// Stub: delete operation would go here

			// Execute afterDelete hooks (error-tolerant)
			r.hooks.ExecuteAfter(req.Context, AfterDelete, entity)

			// Emit resource deleted signal
			capitan.Info(req.Context, ResourceDeleted,
				KeyResource.Field(typeName),
				KeyOperation.Field("delete"),
				KeyEntityID.Field(req.Params.Path["id"]))

			return struct{}{}, nil
		},
	).WithSummary(fmt.Sprintf("Delete a %s", typeName)).
		WithDescription(fmt.Sprintf("Permanently removes a %s resource by its unique identifier.", typeName)).
		WithTags(r.tag).
		WithPathParams("id").
		WithSuccessStatus(http.StatusNoContent)

	r.handlers = []rocco.Endpoint{create, list, get, update, del}
}
