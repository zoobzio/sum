package sum

import (
	"github.com/jmoiron/sqlx"
	"github.com/zoobzio/astql"
	"github.com/zoobzio/grub"
)

// Database creates a grub.Database[M], passes it to the factory to obtain a C implementation,
// registers C with the service locator, and registers the atomic with scio.
func Database[C, M any](k Key, db *sqlx.DB, table, keyCol string, renderer astql.Renderer, factory func(*grub.Database[M]) C) error {
	gdb, err := grub.NewDatabase[M](db, table, keyCol, renderer)
	if err != nil {
		return err
	}

	impl := factory(gdb)
	Register[C](k, impl)

	s := svc()
	return s.catalog.RegisterDatabase("db://"+table, gdb.Atomic())
}

// Store creates a grub.Store[M], passes it to the factory to obtain a C implementation,
// registers C with the service locator, and registers the atomic with scio.
func Store[C, M any](k Key, provider grub.StoreProvider, name string, factory func(*grub.Store[M]) C) error {
	store := grub.NewStore[M](provider)

	impl := factory(store)
	Register[C](k, impl)

	s := svc()
	return s.catalog.RegisterStore("kv://"+name, store.Atomic())
}

// Bucket creates a grub.Bucket[M], passes it to the factory to obtain a C implementation,
// registers C with the service locator, and registers the atomic with scio.
func Bucket[C, M any](k Key, provider grub.BucketProvider, name string, factory func(*grub.Bucket[M]) C) error {
	bucket := grub.NewBucket[M](provider)

	impl := factory(bucket)
	Register[C](k, impl)

	s := svc()
	return s.catalog.RegisterBucket("bcs://"+name, bucket.Atomic())
}
