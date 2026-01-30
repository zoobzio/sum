//go:build testing

package sum

import (
	"context"
	"testing"
)

func resetRegistry(t *testing.T) {
	t.Helper()
	Reset()
	t.Cleanup(Reset)
}

func TestStart(t *testing.T) {
	resetRegistry(t)

	k := Start()
	if k == (Key{}) {
		t.Fatal("expected non-zero key from Start")
	}
}

func TestStartPanicsOnDoubleCall(t *testing.T) {
	resetRegistry(t)

	Start()
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic on second Start call")
		}
	}()
	Start()
}

type testGreeter struct{}

func (testGreeter) Greet() string { return "hello" }

type testGreeterIface interface{ Greet() string }

func TestRegisterAndUse(t *testing.T) {
	resetRegistry(t)

	k := Start()
	Register[testGreeterIface](k, testGreeter{})
	Freeze(k)

	got, err := Use[testGreeterIface](context.Background())
	if err != nil {
		t.Fatalf("Use returned error: %v", err)
	}
	if got == nil {
		t.Fatal("expected non-nil service")
	}
}

func TestMustUse(t *testing.T) {
	resetRegistry(t)

	k := Start()
	Register[testGreeterIface](k, testGreeter{})
	Freeze(k)

	got := MustUse[testGreeterIface](context.Background())
	if got == nil {
		t.Fatal("expected non-nil service from MustUse")
	}
}

type testMissing interface{ Missing() }

func TestMustUsePanicsWhenNotRegistered(t *testing.T) {
	resetRegistry(t)

	k := Start()
	Freeze(k)

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic from MustUse on unregistered service")
		}
	}()
	MustUse[testMissing](context.Background())
}

func TestUseReturnsErrNotFound(t *testing.T) {
	resetRegistry(t)

	k := Start()
	Freeze(k)

	_, err := Use[testMissing](context.Background())
	if err == nil {
		t.Fatal("expected error for unregistered service")
	}
	if err != ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

type testLate interface{ Late() }
type testLateImpl struct{}

func (testLateImpl) Late() {}

func TestFreezePreventsRegistration(t *testing.T) {
	resetRegistry(t)

	k := Start()
	Freeze(k)

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic when registering after Freeze")
		}
	}()
	Register[testLate](k, testLateImpl{})
}

type testSvc interface{ Do() }
type testSvcImpl struct{}

func (testSvcImpl) Do() {}

func TestServices(t *testing.T) {
	resetRegistry(t)

	k := Start()
	Register[testSvc](k, testSvcImpl{})

	infos, err := Services(k)
	if err != nil {
		t.Fatalf("Services returned error: %v", err)
	}
	if len(infos) == 0 {
		t.Error("expected at least one registered service")
	}
}

type testGuarded interface{ Secret() }
type testGuardedImpl struct{}

func (testGuardedImpl) Secret() {}

func TestGuardDeniesAccess(t *testing.T) {
	resetRegistry(t)

	k := Start()
	h := Register[testGuarded](k, testGuardedImpl{})
	h.Guard(func(_ context.Context) error {
		return ErrAccessDenied
	})
	Freeze(k)

	_, err := Use[testGuarded](context.Background())
	if err == nil {
		t.Fatal("expected error from guarded service")
	}
}

func TestSignalsExported(t *testing.T) {
	if SignalRegistered.Name() == "" {
		t.Error("SignalRegistered should have a name")
	}
	if SignalAccessed.Name() == "" {
		t.Error("SignalAccessed should have a name")
	}
	if SignalDenied.Name() == "" {
		t.Error("SignalDenied should have a name")
	}
	if SignalNotFound.Name() == "" {
		t.Error("SignalNotFound should have a name")
	}
}

func TestFieldKeysExported(t *testing.T) {
	if KeyInterface.Name() == "" {
		t.Error("KeyInterface should have a name")
	}
	if KeyImpl.Name() == "" {
		t.Error("KeyImpl should have a name")
	}
	if KeyError.Name() == "" {
		t.Error("KeyError should have a name")
	}
}
