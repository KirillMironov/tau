package storage

import (
	"fmt"

	"github.com/boltdb/bolt"

	"github.com/KirillMironov/tau"
	"github.com/KirillMironov/tau/internal/domain"
	"github.com/KirillMironov/tau/resources"
)

const resourcesBucket = "resources"

type Resources struct {
	db *bolt.DB
}

func NewResources(db *bolt.DB) (*Resources, error) {
	err := db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(resourcesBucket))
		return err
	})

	return &Resources{db: db}, err
}

func (r Resources) Put(resource tau.Resource) error {
	descriptor := resource.Descriptor()

	return r.db.Update(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte(resourcesBucket))
		if root == nil {
			return fmt.Errorf("bucket %q not found", resourcesBucket)
		}

		bucket, err := root.CreateBucketIfNotExists([]byte(descriptor.Kind))
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
		root := tx.Bucket([]byte(resourcesBucket))
		if root == nil {
			return fmt.Errorf("bucket %q not found", resourcesBucket)
		}

		bucket := root.Bucket([]byte(descriptor.Kind))
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
		root := tx.Bucket([]byte(resourcesBucket))
		if root == nil {
			return fmt.Errorf("bucket %q not found", resourcesBucket)
		}

		bucket := root.Bucket([]byte(descriptor.Kind))
		if bucket == nil {
			return fmt.Errorf("bucket %q not found", descriptor.Kind)
		}

		return bucket.Delete([]byte(descriptor.Name))
	})
}

func (r Resources) List() ([]tau.Resource, error) {
	var result []tau.Resource

	err := r.db.View(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte(resourcesBucket))
		if root == nil {
			return fmt.Errorf("bucket %q not found", resourcesBucket)
		}

		return root.ForEach(func(kind, _ []byte) error {
			bucket := root.Bucket(kind)
			if bucket == nil {
				return fmt.Errorf("bucket %q not found", kind)
			}

			return bucket.ForEach(func(_, data []byte) error {
				resource, err := resources.ResourceByKind(tau.Kind(kind))
				if err != nil {
					return err
				}

				if err = resource.UnmarshalBinary(data); err != nil {
					return err
				}

				result = append(result, resource)

				return nil
			})
		})
	})
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, domain.ErrNoResources
	}

	return result, nil
}

func (r Resources) ListByKind(kind tau.Kind) ([]tau.Resource, error) {
	var result []tau.Resource

	err := r.db.View(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte(resourcesBucket))
		if root == nil {
			return fmt.Errorf("bucket %q not found", resourcesBucket)
		}

		bucket := root.Bucket([]byte(kind))
		if bucket == nil {
			return fmt.Errorf("bucket %q not found", kind)
		}

		return bucket.ForEach(func(_, data []byte) error {
			resource, err := resources.ResourceByKind(kind)
			if err != nil {
				return err
			}

			if err = resource.UnmarshalBinary(data); err != nil {
				return err
			}

			result = append(result, resource)

			return nil
		})
	})
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return nil, domain.ErrNoResources
	}

	return result, nil
}
