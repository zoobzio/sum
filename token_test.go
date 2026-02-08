//go:build testing

package sum

import (
	"context"
	"errors"
	"testing"
)

func TestNewTokenUnique(t *testing.T) {
	t1 := NewToken("test")
	t2 := NewToken("test")

	if t1.id == t2.id {
		t.Error("expected different tokens to have different IDs")
	}
}

func TestTokenString(t *testing.T) {
	tok := NewToken("handlers")
	if got := tok.String(); got != "handlers" {
		t.Errorf("String() = %q, want %q", got, "handlers")
	}
}

func TestWithTokenRoundTrip(t *testing.T) {
	tok := NewToken("test")
	ctx := WithToken(context.Background(), tok)

	got, ok := tokenFromContext(ctx)
	if !ok {
		t.Fatal("expected token in context")
	}
	if got.id != tok.id {
		t.Error("token IDs do not match")
	}
}

func TestTokenFromContextMissing(t *testing.T) {
	_, ok := tokenFromContext(context.Background())
	if ok {
		t.Error("expected no token in empty context")
	}
}

func TestRequireAllowsMatchingToken(t *testing.T) {
	tok := NewToken("handlers")
	guard := Require(tok)

	ctx := WithToken(context.Background(), tok)
	if err := guard(ctx); err != nil {
		t.Errorf("expected access granted, got %v", err)
	}
}

func TestRequireDeniesWrongToken(t *testing.T) {
	allowed := NewToken("handlers")
	wrong := NewToken("ingest")
	guard := Require(allowed)

	ctx := WithToken(context.Background(), wrong)
	err := guard(ctx)
	if err == nil {
		t.Fatal("expected access denied")
	}
	if !errors.Is(err, ErrAccessDenied) {
		t.Errorf("expected ErrAccessDenied, got %v", err)
	}
}

func TestRequireDeniesNoToken(t *testing.T) {
	tok := NewToken("handlers")
	guard := Require(tok)

	err := guard(context.Background())
	if err == nil {
		t.Fatal("expected error when no token provided")
	}
	if !errors.Is(err, ErrTokenRequired) {
		t.Errorf("expected ErrTokenRequired, got %v", err)
	}
}

func TestRequireMultipleTokensAnyGrantsAccess(t *testing.T) {
	handlers := NewToken("handlers")
	ingest := NewToken("ingest")
	guard := Require(handlers, ingest)

	// handlers token grants access
	ctx := WithToken(context.Background(), handlers)
	if err := guard(ctx); err != nil {
		t.Errorf("handlers token should grant access: %v", err)
	}

	// ingest token also grants access
	ctx = WithToken(context.Background(), ingest)
	if err := guard(ctx); err != nil {
		t.Errorf("ingest token should grant access: %v", err)
	}

	// unrelated token denied
	other := NewToken("other")
	ctx = WithToken(context.Background(), other)
	if err := guard(ctx); err == nil {
		t.Error("other token should be denied")
	}
}

type testTokenSvc interface{ Op() }
type testTokenSvcImpl struct{}

func (testTokenSvcImpl) Op() {}

func TestRequireIntegrationWithRegistry(t *testing.T) {
	resetRegistry(t)

	k := Start()
	tok := NewToken("test")
	Register[testTokenSvc](k, testTokenSvcImpl{}).Guard(Require(tok))
	Freeze(k)

	// Without token: denied
	_, err := Use[testTokenSvc](context.Background())
	if err == nil {
		t.Error("expected error without token")
	}

	// With token: allowed
	ctx := WithToken(context.Background(), tok)
	svc, err := Use[testTokenSvc](ctx)
	if err != nil {
		t.Fatalf("expected access with token, got %v", err)
	}
	if svc == nil {
		t.Error("expected non-nil service")
	}
}

type testForSvc interface{ DoWork() }
type testForSvcImpl struct{}

func (testForSvcImpl) DoWork() {}

func TestForMethodIntegration(t *testing.T) {
	resetRegistry(t)

	k := Start()
	tok := NewToken("handlers")
	Register[testForSvc](k, testForSvcImpl{}).For(tok)
	Freeze(k)

	// Without token: denied
	_, err := Use[testForSvc](context.Background())
	if err == nil {
		t.Error("expected error without token")
	}

	// With token: allowed
	ctx := WithToken(context.Background(), tok)
	svc, err := Use[testForSvc](ctx)
	if err != nil {
		t.Fatalf("expected access with token, got %v", err)
	}
	if svc == nil {
		t.Error("expected non-nil service")
	}
}

func TestForMethodMultipleTokens(t *testing.T) {
	resetRegistry(t)

	k := Start()
	handlers := NewToken("handlers")
	ingest := NewToken("ingest")
	Register[testForSvc](k, testForSvcImpl{}).For(handlers, ingest)
	Freeze(k)

	// handlers token grants access
	ctx := WithToken(context.Background(), handlers)
	if _, err := Use[testForSvc](ctx); err != nil {
		t.Errorf("handlers token should grant access: %v", err)
	}

	// ingest token also grants access
	ctx = WithToken(context.Background(), ingest)
	if _, err := Use[testForSvc](ctx); err != nil {
		t.Errorf("ingest token should grant access: %v", err)
	}

	// unrelated token denied
	other := NewToken("other")
	ctx = WithToken(context.Background(), other)
	if _, err := Use[testForSvc](ctx); err == nil {
		t.Error("other token should be denied")
	}
}
