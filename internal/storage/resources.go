package storage

import (
	"bytes"
	"encoding/gob"
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

		buf := new(bytes.Buffer)
		if err = gob.NewEncoder(buf).Encode(resource); err != nil {
			return err
		}

		return bucket.Put([]byte(name), buf.Bytes())
	})
}

func (r Resources) Get(descriptor resources.Descriptor) (resource resources.Resource, _ error) {
	kind := descriptor.Kind.String()
	name := descriptor.Name

	return resource, r.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(kind))
		if bucket == nil {
			return fmt.Errorf("bucket %q not found", kind)
		}

		data := bucket.Get([]byte(name))
		if data == nil {
			return fmt.Errorf("resource with name %q not found", name)
		}

		return gob.NewDecoder(bytes.NewReader(data)).Decode(&resource)
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
