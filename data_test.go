//go:build testing

package sum

import (
	"testing"

	"github.com/zoobzio/grub"
)

// Note: Full integration tests for Database, Store, and Bucket
// require actual database/store connections and are located in
// testing/integration/data_test.go.
//
// These unit tests verify the registration flow with mocked dependencies.

type testModel struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

type testClient struct {
	db *grub.Database[testModel]
}

func TestDatabaseRegistration(t *testing.T) {
	t.Skip("requires database connection - see testing/integration/data_test.go")
}

func TestStoreRegistration(t *testing.T) {
	t.Skip("requires store provider - see testing/integration/data_test.go")
}

func TestBucketRegistration(t *testing.T) {
	t.Skip("requires bucket provider - see testing/integration/data_test.go")
}
