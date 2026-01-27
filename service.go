package sum

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/zoobzio/rocco"
	"github.com/zoobzio/scio"
)

// service is the singleton instance.
var (
	instance *Service
	once     sync.Once
)

// Service wraps a rocco engine and scio catalog, providing application lifecycle.
type Service struct {
	engine  *rocco.Engine
	catalog *scio.Scio
}

// New creates or returns the singleton Service.
// Subsequent calls return the existing instance.
func New() *Service {
	once.Do(func() {
		instance = &Service{
			engine:  rocco.NewEngine(),
			catalog: scio.New(),
		}
	})
	return instance
}

// svc returns the singleton, panicking if not initialized.
func svc() *Service {
	if instance == nil {
		panic("sum: service not initialized, call New() first")
	}
	return instance
}

// Handle registers endpoints with the underlying engine.
func (s *Service) Handle(endpoints ...rocco.Endpoint) {
	s.engine.WithHandlers(endpoints...)
}

// Tag registers an OpenAPI tag with a description.
func (s *Service) Tag(name, description string) {
	s.engine.WithTag(name, description)
}

// Engine returns the underlying rocco engine for advanced usage.
func (s *Service) Engine() *rocco.Engine {
	return s.engine
}

// Catalog returns the scio data catalog for advanced usage.
func (s *Service) Catalog() *scio.Scio {
	return s.catalog
}

// Start begins serving. This method blocks until shutdown.
func (s *Service) Start(host string, port int) error {
	return s.engine.Start(host, port)
}

// Shutdown gracefully stops the service.
func (s *Service) Shutdown(ctx context.Context) error {
	if s.engine == nil {
		return fmt.Errorf("service not started")
	}
	return s.engine.Shutdown(ctx)
}

// Run starts the service and blocks until a shutdown signal is received.
// Handles SIGINT and SIGTERM, then performs graceful shutdown with a 30 second timeout.
func (s *Service) Run(host string, port int) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	errCh := make(chan error, 1)
	go func() {
		errCh <- s.engine.Start(host, port)
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer shutdownCancel()
		return s.Shutdown(shutdownCtx)
	}
}
