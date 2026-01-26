//go:build testing

package testing

import (
	"testing"
)

func TestTestContext(t *testing.T) {
	ctx := TestContext(t)
	if ctx == nil {
		t.Fatal("expected non-nil context")
	}
	select {
	case <-ctx.Done():
		t.Fatal("context should not be done immediately")
	default:
		// Expected
	}
}

func TestShortContext(t *testing.T) {
	ctx := ShortContext(t)
	if ctx == nil {
		t.Fatal("expected non-nil context")
	}
	select {
	case <-ctx.Done():
		t.Fatal("context should not be done immediately")
	default:
		// Expected
	}
}
