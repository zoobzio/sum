package sum

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

// Token is an unforgeable capability for service access.
type Token struct {
	id   string // random UUID, unexported
	name string // for debugging/logging
}

// NewToken creates a new access token with the given name.
func NewToken(name string) Token {
	return Token{
		id:   uuid.New().String(),
		name: name,
	}
}

// String returns the token name for debugging.
func (t Token) String() string {
	return t.name
}

type tokenKey struct{}

// WithToken injects a token into the context.
func WithToken(ctx context.Context, t Token) context.Context {
	return context.WithValue(ctx, tokenKey{}, t)
}

func tokenFromContext(ctx context.Context) (Token, bool) {
	t, ok := ctx.Value(tokenKey{}).(Token)
	return t, ok
}

// ErrTokenRequired indicates a service requires a token but none was provided.
var ErrTokenRequired = fmt.Errorf("token required")

// Require returns a guard that checks for any of the provided tokens.
// If the service has Require, the context must contain a matching token.
func Require(tokens ...Token) Guard {
	return func(ctx context.Context) error {
		ctxToken, ok := tokenFromContext(ctx)
		if !ok {
			return ErrTokenRequired
		}
		for _, t := range tokens {
			if t.id == ctxToken.id {
				return nil
			}
		}
		return fmt.Errorf("%w: token %q does not grant access", ErrAccessDenied, ctxToken.name)
	}
}
