//go:build testing

package sum

import (
	"context"
	"sync"
	"testing"
	"time"
)

func TestNewService(t *testing.T) {
	// Reset singleton state
	instance = nil
	once = sync.Once{}
	t.Cleanup(func() {
		instance = nil
		once = sync.Once{}
	})

	cfg := ServiceConfig{
		Host: "localhost",
		Port: 8080,
	}

	svc := New(cfg)
	if svc == nil {
		t.Fatal("expected non-nil service")
	}

	// Verify singleton behavior
	svc2 := New(ServiceConfig{Host: "other", Port: 9090})
	if svc != svc2 {
		t.Error("expected same instance from second New call")
	}

	// Verify original config is preserved
	if svc2.Config().Host != "localhost" {
		t.Errorf("expected host 'localhost', got '%s'", svc2.Config().Host)
	}
}

func TestServiceConfig(t *testing.T) {
	instance = nil
	once = sync.Once{}
	t.Cleanup(func() {
		instance = nil
		once = sync.Once{}
	})

	cfg := ServiceConfig{
		Host: "0.0.0.0",
		Port: 3000,
	}

	svc := New(cfg)
	got := svc.Config()

	if got.Host != cfg.Host {
		t.Errorf("expected host '%s', got '%s'", cfg.Host, got.Host)
	}
	if got.Port != cfg.Port {
		t.Errorf("expected port %d, got %d", cfg.Port, got.Port)
	}
}

func TestServiceEngineAndCatalog(t *testing.T) {
	instance = nil
	once = sync.Once{}
	t.Cleanup(func() {
		instance = nil
		once = sync.Once{}
	})

	svc := New(ServiceConfig{Host: "localhost", Port: 8080})

	if svc.Engine() == nil {
		t.Error("expected non-nil engine")
	}
	if svc.Catalog() == nil {
		t.Error("expected non-nil catalog")
	}
}

func TestSvcPanicsWithoutInit(t *testing.T) {
	instance = nil
	once = sync.Once{}
	t.Cleanup(func() {
		instance = nil
		once = sync.Once{}
	})

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic when calling svc() without initialization")
		}
	}()

	svc()
}

func TestServiceShutdownWithoutStart(t *testing.T) {
	instance = nil
	once = sync.Once{}
	t.Cleanup(func() {
		instance = nil
		once = sync.Once{}
	})

	svc := New(ServiceConfig{Host: "localhost", Port: 8080})

	// Temporarily nil the engine to simulate not started state
	originalEngine := svc.engine
	svc.engine = nil
	defer func() { svc.engine = originalEngine }()

	ctx := context.Background()
	err := svc.Shutdown(ctx)
	if err == nil {
		t.Error("expected error when shutting down non-started service")
	}
}

func TestServiceTagAndHandle(t *testing.T) {
	instance = nil
	once = sync.Once{}
	t.Cleanup(func() {
		instance = nil
		once = sync.Once{}
	})

	svc := New(ServiceConfig{Host: "localhost", Port: 8080})

	// These should not panic
	svc.Tag("users", "User management endpoints")
	svc.Handle() // empty handlers
}

func TestServiceRunWithCancel(t *testing.T) {
	instance = nil
	once = sync.Once{}
	t.Cleanup(func() {
		instance = nil
		once = sync.Once{}
	})

	t.Skip("requires running server - see testing/integration/service_test.go")
}

func TestServiceStartStop(t *testing.T) {
	instance = nil
	once = sync.Once{}
	t.Cleanup(func() {
		instance = nil
		once = sync.Once{}
	})

	svc := New(ServiceConfig{Host: "localhost", Port: 0}) // Port 0 for random available port

	errCh := make(chan error, 1)
	go func() {
		errCh <- svc.Start()
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := svc.Shutdown(ctx); err != nil {
		t.Errorf("shutdown failed: %v", err)
	}

	select {
	case err := <-errCh:
		// Server should return nil or http.ErrServerClosed on clean shutdown
		if err != nil && err.Error() != "http: Server closed" {
			t.Errorf("unexpected start error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Error("timeout waiting for server to stop")
	}
}
