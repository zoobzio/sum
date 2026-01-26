//go:build testing

package sum

import (
	"context"
	"testing"
)

type testConfig struct {
	Name    string `env:"TEST_CONFIG_NAME" default:"default-name"`
	Enabled bool   `env:"TEST_CONFIG_ENABLED" default:"false"`
	Port    int    `env:"TEST_CONFIG_PORT" default:"8080"`
}

func TestConfig(t *testing.T) {
	Reset()
	t.Cleanup(Reset)

	k := Start()
	ctx := context.Background()

	// Config loads from env with defaults
	err := Config[testConfig](ctx, k, nil)
	if err != nil {
		t.Fatalf("Config failed: %v", err)
	}

	cfg, err := Use[testConfig](ctx)
	if err != nil {
		t.Fatalf("Use failed: %v", err)
	}

	// Should have default values
	if cfg.Name != "default-name" {
		t.Errorf("expected name 'default-name', got '%s'", cfg.Name)
	}
	if cfg.Enabled != false {
		t.Error("expected enabled to be false")
	}
	if cfg.Port != 8080 {
		t.Errorf("expected port 8080, got %d", cfg.Port)
	}
}

func TestConfigWithEnv(t *testing.T) {
	Reset()
	t.Cleanup(Reset)

	t.Setenv("TEST_CONFIG_NAME", "custom-service")
	t.Setenv("TEST_CONFIG_ENABLED", "true")
	t.Setenv("TEST_CONFIG_PORT", "9000")

	k := Start()
	ctx := context.Background()

	err := Config[testConfig](ctx, k, nil)
	if err != nil {
		t.Fatalf("Config failed: %v", err)
	}

	cfg, err := Use[testConfig](ctx)
	if err != nil {
		t.Fatalf("Use failed: %v", err)
	}

	if cfg.Name != "custom-service" {
		t.Errorf("expected name 'custom-service', got '%s'", cfg.Name)
	}
	if cfg.Enabled != true {
		t.Error("expected enabled to be true")
	}
	if cfg.Port != 9000 {
		t.Errorf("expected port 9000, got %d", cfg.Port)
	}
}

type requiredConfig struct {
	Required string `env:"TEST_REQUIRED_VALUE" required:"true"`
}

func TestConfigRequiredMissing(t *testing.T) {
	Reset()
	t.Cleanup(Reset)

	k := Start()
	ctx := context.Background()

	err := Config[requiredConfig](ctx, k, nil)
	if err == nil {
		t.Error("expected error for missing required field")
	}
}
