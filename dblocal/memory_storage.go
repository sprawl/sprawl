package dblocal

import (
	"strings"

	"github.com/sprawl/sprawl/errors"
)

// Storage is a struct containing a database and its address
type Storage struct {
	db map[string]string
}

var err error
var data []byte

// SetDbPath sets the path the database files are located
func (storage *Storage) SetDbPath(dbPath string) {
}

// Run starts the database connection for Storage
func (storage *Storage) Run() error {
	return nil
}

// Close closes the underlying LevelDB connection
func (storage *Storage) Close() {
}

// Has uses LevelDB's method Has to check does the data exists in LevelDB
func (storage *Storage) Has(key []byte) (bool, error) {
	_, ok := storage.db[string(key)]
	return ok, nil
}

// Get uses LevelDB's method Get to fetch data from LevelDB
func (storage *Storage) Get(key []byte) ([]byte, error) {
	value, ok := storage.db[string(key)]
	var err error
	if !ok {
		err = errors.E(errors.Op("Get value from memory database"))
	}
	return []byte(value), err
}

// Put uses LevelDB's Put method to put data into LevelDB
func (storage *Storage) Put(key []byte, data []byte) error {
	storage.db[string(key)] = string(data)
	return nil
}

// Delete uses LevelDB's Delete method to remove data from LevelDB
func (storage *Storage) Delete(key []byte) error {
	delete(storage.db, string(key))
	return nil
}

// GetAll returns all entries in the database regardless of key or prefix
func (storage *Storage) GetAll() (map[string]string, error) {
	return storage.db, nil
}

// GetAllWithPrefix returns all entries in the database with the specified prefix
func (storage *Storage) GetAllWithPrefix(prefix string) (map[string]string, error) {
	entries := make(map[string]string)
	for k, v := range storage.db {
		if strings.HasPrefix(k, prefix) {
			entries[k] = v
		}
	}
	return entries, nil
}

// DeleteAll deletes all entries from the database
// USE CAREFULLY
func (storage *Storage) DeleteAll() error {
	storage.db = make(map[string]string)
	return nil
}

// DeleteAllWithPrefix deletes all entries starting with a prefix
func (storage *Storage) DeleteAllWithPrefix(prefix string) error {
	for k := range storage.db {
		if strings.HasPrefix(k, prefix) {
			delete(storage.db, k)
		}
	}
	return nil
}
