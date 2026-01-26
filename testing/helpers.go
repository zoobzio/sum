//go:build testing

// Package testing provides test helpers for sum.
package testing

import (
	"context"
	"testing"
	"time"
)

// TestContext returns a context with a test-appropriate timeout.
func TestContext(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	t.Cleanup(cancel)
	return ctx
}

// ShortContext returns a context with a short timeout for quick tests.
func ShortContext(t *testing.T) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	t.Cleanup(cancel)
	return ctx
}
