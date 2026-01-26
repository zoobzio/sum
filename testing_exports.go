//go:build testing || integration

package sum

import (
	"sync"

	"github.com/zoobzio/slush"
)

// Reset clears all registered services and resets initialization state.
// Also resets the service singleton so New() can be called again.
// Only available in test builds.
func Reset() {
	slush.Reset()
	instance = nil
	once = sync.Once{}
}

// Unregister removes a service by type.
// Only available in test builds.
func Unregister[T any]() {
	slush.Unregister[T]()
}
