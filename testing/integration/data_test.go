//go:build testing

package integration

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/zoobzio/astql/postgres"
	"github.com/zoobzio/grub"
	"github.com/zoobzio/sum"
	sumtest "github.com/zoobzio/sum/testing"
)

// testModel represents a simple entity for testing database operations.
type testModel struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

// testClient wraps a grub.Database for testing.
type testClient struct {
	db *grub.Database[testModel]
}

func TestDatabaseIntegration(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set - skipping database integration test")
	}

	ctx := sumtest.TestContext(t)
	_ = ctx

	sum.Reset()
	t.Cleanup(sum.Reset)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("failed to open database: %v", err)
	}
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "postgres")

	k := sum.Start()

	// Create test table
	_, err = sqlxDB.Exec(`
		CREATE TABLE IF NOT EXISTS test_models (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL
		)
	`)
	if err != nil {
		t.Fatalf("failed to create test table: %v", err)
	}
	t.Cleanup(func() {
		sqlxDB.Exec(`DROP TABLE IF EXISTS test_models`)
	})

	err = sum.Database[testClient, testModel](
		k,
		sqlxDB,
		"test_models",
		"id",
		postgres.New(),
		func(db *grub.Database[testModel]) testClient {
			return testClient{db: db}
		},
	)
	if err != nil {
		t.Fatalf("Database registration failed: %v", err)
	}

	client, err := sum.Use[testClient](ctx)
	if err != nil {
		t.Fatalf("Use[testClient] failed: %v", err)
	}

	if client.db == nil {
		t.Error("expected non-nil database in client")
	}
}

// mockStoreProvider is a test implementation of grub.StoreProvider.
type mockStoreProvider struct{}

func (m mockStoreProvider) Get(ctx context.Context, key string) ([]byte, error) {
	return nil, nil
}

func (m mockStoreProvider) Set(ctx context.Context, key string, value []byte, ttl time.Duration) error {
	return nil
}

func (m mockStoreProvider) Delete(ctx context.Context, key string) error {
	return nil
}

func (m mockStoreProvider) Exists(ctx context.Context, key string) (bool, error) {
	return false, nil
}

func (m mockStoreProvider) List(ctx context.Context, prefix string, limit int) ([]string, error) {
	return nil, nil
}

func (m mockStoreProvider) GetBatch(ctx context.Context, keys []string) (map[string][]byte, error) {
	return nil, nil
}

func (m mockStoreProvider) SetBatch(ctx context.Context, items map[string][]byte, ttl time.Duration) error {
	return nil
}

type testStoreClient struct {
	store *grub.Store[testModel]
}

func TestStoreIntegration(t *testing.T) {
	// This test uses a mock provider. For actual KV store testing,
	// set TEST_STORE_PROVIDER and implement real provider initialization.
	if os.Getenv("TEST_STORE_PROVIDER") == "" {
		t.Skip("TEST_STORE_PROVIDER not set - skipping store integration test")
	}

	ctx := sumtest.TestContext(t)
	_ = ctx

	sum.Reset()
	t.Cleanup(sum.Reset)

	k := sum.Start()

	err := sum.Store[testStoreClient, testModel](
		k,
		mockStoreProvider{},
		"test-store",
		func(s *grub.Store[testModel]) testStoreClient {
			return testStoreClient{store: s}
		},
	)
	if err != nil {
		t.Fatalf("Store registration failed: %v", err)
	}

	client, err := sum.Use[testStoreClient](ctx)
	if err != nil {
		t.Fatalf("Use[testStoreClient] failed: %v", err)
	}

	if client.store == nil {
		t.Error("expected non-nil store in client")
	}
}

// mockBucketProvider is a test implementation of grub.BucketProvider.
type mockBucketProvider struct{}

func (m mockBucketProvider) Get(ctx context.Context, key string) ([]byte, *grub.ObjectInfo, error) {
	return nil, nil, nil
}

func (m mockBucketProvider) Put(ctx context.Context, key string, data []byte, info *grub.ObjectInfo) error {
	return nil
}

func (m mockBucketProvider) Delete(ctx context.Context, key string) error {
	return nil
}

func (m mockBucketProvider) Exists(ctx context.Context, key string) (bool, error) {
	return false, nil
}

func (m mockBucketProvider) List(ctx context.Context, prefix string, limit int) ([]grub.ObjectInfo, error) {
	return nil, nil
}

type testBucketClient struct {
	bucket *grub.Bucket[testModel]
}

func TestBucketIntegration(t *testing.T) {
	// This test uses a mock provider. For actual blob storage testing,
	// set TEST_BUCKET_PROVIDER and implement real provider initialization.
	if os.Getenv("TEST_BUCKET_PROVIDER") == "" {
		t.Skip("TEST_BUCKET_PROVIDER not set - skipping bucket integration test")
	}

	ctx := sumtest.TestContext(t)
	_ = ctx

	sum.Reset()
	t.Cleanup(sum.Reset)

	k := sum.Start()

	err := sum.Bucket[testBucketClient, testModel](
		k,
		mockBucketProvider{},
		"test-bucket",
		func(b *grub.Bucket[testModel]) testBucketClient {
			return testBucketClient{bucket: b}
		},
	)
	if err != nil {
		t.Fatalf("Bucket registration failed: %v", err)
	}

	client, err := sum.Use[testBucketClient](ctx)
	if err != nil {
		t.Fatalf("Use[testBucketClient] failed: %v", err)
	}

	if client.bucket == nil {
		t.Error("expected non-nil bucket in client")
	}
}
