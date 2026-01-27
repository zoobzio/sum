package sum

import (
	"github.com/jmoiron/sqlx"
	"github.com/zoobzio/astql"
	"github.com/zoobzio/grub"
)

// Database wraps grub.Database and registers with scio on creation.
// Embed this type in your store structs to add custom query methods.
type Database[M any] struct {
	*grub.Database[M]
}

// NewDatabase creates a Database[M] and registers it with the scio catalog.
// Requires sum.New() to have been called first.
func NewDatabase[M any](db *sqlx.DB, table string, renderer astql.Renderer) (*Database[M], error) {
	gdb, err := grub.NewDatabase[M](db, table, renderer)
	if err != nil {
		return nil, err
	}
	svc().catalog.RegisterDatabase("db://"+table, gdb.Atomic())
	return &Database[M]{Database: gdb}, nil
}

// Store wraps grub.Store and registers with scio on creation.
// Embed this type in your store structs to add custom methods.
type Store[M any] struct {
	*grub.Store[M]
}

// NewStore creates a Store[M] and registers it with the scio catalog.
// Requires sum.New() to have been called first.
func NewStore[M any](provider grub.StoreProvider, name string) *Store[M] {
	store := grub.NewStore[M](provider)
	svc().catalog.RegisterStore("kv://"+name, store.Atomic())
	return &Store[M]{Store: store}
}

// Bucket wraps grub.Bucket and registers with scio on creation.
// Embed this type in your store structs to add custom methods.
type Bucket[M any] struct {
	*grub.Bucket[M]
}

// NewBucket creates a Bucket[M] and registers it with the scio catalog.
// Requires sum.New() to have been called first.
func NewBucket[M any](provider grub.BucketProvider, name string) *Bucket[M] {
	bucket := grub.NewBucket[M](provider)
	svc().catalog.RegisterBucket("bcs://"+name, bucket.Atomic())
	return &Bucket[M]{Bucket: bucket}
}
