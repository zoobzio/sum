package sum

import (
	"context"
	"fmt"
	"sync"

	"github.com/zoobzio/rocco"
)

// server is the internal singleton managing the rocco engine and resources.
type server struct {
	mu        sync.RWMutex
	config    Config
	engine    *rocco.Engine
	resources []resource
}

// resource is the internal interface for registered resources.
type resource interface {
	endpoints() []rocco.Endpoint
	tagInfo() (name, description string)
}

var (
	instance *server
	once     sync.Once
)

// ensureInitialized lazily initializes the singleton server.
func ensureInitialized() *server {
	once.Do(func() {
		instance = &server{
			config:    DefaultConfig(),
			resources: make([]resource, 0),
		}
	})
	return instance
}

// Configure sets the server configuration.
// Must be called before New[T]() or Start().
func Configure(cfg Config) {
	s := ensureInitialized()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.config = cfg
}

// GetConfig returns a copy of the current configuration.
func GetConfig() Config {
	s := ensureInitialized()
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.config
}

// registerResource adds a resource to the server.
func registerResource(r resource) {
	s := ensureInitialized()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.resources = append(s.resources, r)
}

// Start initializes the rocco engine with all registered resources and begins serving.
// This method blocks until the server is shutdown.
func Start() error {
	s := ensureInitialized()
	s.mu.Lock()

	// Create the rocco engine
	s.engine = rocco.NewEngine(s.config.Host, s.config.Port, nil)

	// Register tags and collect endpoints from resources
	endpoints := make([]rocco.Endpoint, 0)
	for _, r := range s.resources {
		// Apply tag description to engine
		name, description := r.tagInfo()
		s.engine.WithTag(name, description)

		endpoints = append(endpoints, r.endpoints()...)
	}

	// Register all endpoints with the engine
	s.engine.WithHandlers(endpoints...)

	s.mu.Unlock()

	return s.engine.Start()
}

// Shutdown performs a graceful shutdown of the server.
func Shutdown(ctx context.Context) error {
	s := ensureInitialized()
	s.mu.RLock()
	engine := s.engine
	s.mu.RUnlock()

	if engine == nil {
		return fmt.Errorf("server not started")
	}
	return engine.Shutdown(ctx)
}
