package storage

import (
	"fmt"

	"github.com/boltdb/bolt"

	"github.com/KirillMironov/tau"
	"github.com/KirillMironov/tau/resources"
)

type Resources struct {
	db *bolt.DB
}

func NewResources(db *bolt.DB) *Resources {
	return &Resources{db: db}
}

func (r Resources) Put(resource tau.Resource) error {
	descriptor := resource.Descriptor()

	return r.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(descriptor.Kind))
		if err != nil {
			return err
		}

		data, err := resource.MarshalBinary()
		if err != nil {
			return err
		}

		return bucket.Put([]byte(descriptor.Name), data)
	})
}

func (r Resources) Get(descriptor tau.Descriptor) (tau.Resource, error) {
	resource, err := resources.ResourceByKind(descriptor.Kind)
	if err != nil {
		return nil, err
	}

	return resource, r.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(descriptor.Kind))
		if bucket == nil {
			return fmt.Errorf("bucket %q not found", descriptor.Kind)
		}

		data := bucket.Get([]byte(descriptor.Name))
		if data == nil {
			return fmt.Errorf("resource with name %q not found", descriptor.Name)
		}

		return resource.UnmarshalBinary(data)
	})
}

func (r Resources) Delete(descriptor tau.Descriptor) error {
	return r.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(descriptor.Kind))
		if bucket == nil {
			return fmt.Errorf("bucket %q not found", descriptor.Kind)
		}

		return bucket.Delete([]byte(descriptor.Name))
	})
}
