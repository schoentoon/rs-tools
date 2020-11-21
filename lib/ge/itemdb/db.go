package itemdb

import (
	"encoding/json"
	"io"
	"sync"

	"gitlab.com/schoentoon/rs-tools/lib/ge"
)

type DB struct {
	mutex       sync.RWMutex
	idToItems   map[int64]*ge.Item
	nameToItems map[string]*ge.Item

	Writer io.Writer
}

func (db *DB) add(item *ge.Item) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()
	db.idToItems[item.ID] = item
	db.nameToItems[item.Name] = item

	if db.Writer != nil {
		return json.NewEncoder(db.Writer).Encode(*item)
	}

	return nil
}

func (db *DB) Serialize(w io.Writer) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	encoder := json.NewEncoder(w)

	for _, item := range db.idToItems {
		err := encoder.Encode(*item)
		if err != nil {
			return err
		}
	}
	return nil
}

func New() *DB {
	return &DB{
		idToItems:   make(map[int64]*ge.Item),
		nameToItems: make(map[string]*ge.Item),
		Writer:      nil,
	}
}

func NewFromReader(r io.Reader) (*DB, error) {
	db := New()

	decoder := json.NewDecoder(r)

	for decoder.More() {
		item := &ge.Item{}
		err := decoder.Decode(item)
		if err != nil {
			return nil, err
		}
		err = db.add(item)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}
