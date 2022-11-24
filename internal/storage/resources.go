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
	return r.db.Update(func(txn *badger.Txn) error {
		var buf bytes.Buffer

		err := gob.NewEncoder(&buf).Encode(resource)
		if err != nil {
			return err
		}

		return txn.Set([]byte(resource.ID()), buf.Bytes())
	})
}

func (r Resources) GetByID(id string) (resource resources.Resource, _ error) {
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

func (r Resources) Delete(id string) error {
	return r.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(id))
	})
}
