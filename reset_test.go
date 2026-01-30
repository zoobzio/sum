//go:build testing

package sum

import (
	"context"
	"testing"
)

func TestReset(t *testing.T) {
	Reset()
	t.Cleanup(Reset)

	k := Start()
	Register[testSvc](k, testSvcImpl{})
	Freeze(k)

	_, err := Use[testSvc](context.Background())
	if err != nil {
		t.Fatalf("expected service before reset: %v", err)
	}

	Reset()

	_, err = Use[testSvc](context.Background())
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound after reset, got %v", err)
	}
}

func TestResetClearsSingleton(t *testing.T) {
	Reset()
	t.Cleanup(Reset)

	svc1 := New()
	Reset()
	svc2 := New()

	if svc1 == svc2 {
		t.Error("expected different instances after Reset")
	}
}

type testRemovable interface{ Remove() }
type testRemovableImpl struct{}

func (testRemovableImpl) Remove() {}

func TestUnregister(t *testing.T) {
	Reset()
	t.Cleanup(Reset)

	k := Start()
	Register[testRemovable](k, testRemovableImpl{})
	Freeze(k)

	_, err := Use[testRemovable](context.Background())
	if err != nil {
		t.Fatalf("expected service before unregister: %v", err)
	}

	Unregister[testRemovable]()

	_, err = Use[testRemovable](context.Background())
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound after Unregister, got %v", err)
	}
}
