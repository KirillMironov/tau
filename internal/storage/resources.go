package storage

import (
	"bytes"
	"encoding/gob"

	"github.com/dgraph-io/badger/v3"

	"github.com/KirillMironov/tau/resources"
)

type Resources struct {
	db *badger.DB
}

func NewResources(db *badger.DB) *Resources {
	return &Resources{db: db}
}

func (r Resources) Create(resource resources.Resource) error {
	id := resource.Kind().String() + resource.ID()

	return r.db.Update(func(txn *badger.Txn) error {
		var buf bytes.Buffer

		err := gob.NewEncoder(&buf).Encode(resource)
		if err != nil {
			return err
		}

		return txn.Set([]byte(id), buf.Bytes())
	})
}

func (r Resources) Get(name string, kind resources.Kind) (resource resources.Resource, _ error) {
	id := kind.String() + name

	return resource, r.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(id))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			return gob.NewDecoder(bytes.NewReader(val)).Decode(&resource)
		})
	})
}

func (r Resources) Delete(name string, kind resources.Kind) error {
	id := kind.String() + name

	return r.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(id))
	})
}
