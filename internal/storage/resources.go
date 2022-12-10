package storage

import (
	"fmt"

	"github.com/boltdb/bolt"

	"github.com/KirillMironov/tau/resources"
)

type Resources struct {
	db *bolt.DB
}

func NewResources(db *bolt.DB) *Resources {
	return &Resources{db: db}
}

func (r Resources) Put(resource resources.Resource) error {
	descriptor := resource.Descriptor()
	kind := descriptor.Kind.String()
	name := descriptor.Name

	return r.db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte(kind))
		if err != nil {
			return err
		}

		data, err := resource.MarshalBinary()
		if err != nil {
			return err
		}

		return bucket.Put([]byte(name), data)
	})
}

func (r Resources) Get(descriptor resources.Descriptor) (resource resources.Resource, _ error) {
	kind := descriptor.Kind.String()
	name := descriptor.Name

	switch descriptor.Kind {
	case resources.KindContainer:
		resource = &resources.Container{}
	case resources.KindPod:
		resource = &resources.Pod{}
	default:
		return nil, fmt.Errorf("unexpected resource kind %q", descriptor.Kind)
	}

	return resource, r.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(kind))
		if bucket == nil {
			return fmt.Errorf("bucket %q not found", kind)
		}

		data := bucket.Get([]byte(name))
		if data == nil {
			return fmt.Errorf("resource with name %q not found", name)
		}

		return resource.UnmarshalBinary(data)
	})
}

func (r Resources) Delete(descriptor resources.Descriptor) error {
	kind := descriptor.Kind.String()
	name := descriptor.Name

	return r.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(kind))
		if bucket == nil {
			return fmt.Errorf("bucket %q not found", kind)
		}

		return bucket.Delete([]byte(name))
	})
}
