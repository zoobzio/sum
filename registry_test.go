//go:build testing || integration

package sum

import "github.com/zoobzio/slush"

// Reset clears all registered services and resets initialization state.
// Only available in test builds.
func Reset() {
	slush.Reset()
}

// Unregister removes a service by type.
// Only available in test builds.
func Unregister[T any]() {
	slush.Unregister[T]()
}
