package db

import (
	"github.com/sprawl/sprawl/errors"
	"github.com/syndtr/goleveldb/leveldb"
	util "github.com/syndtr/goleveldb/leveldb/util"
)

// Storage is a struct containing a database and its address
type Storage struct {
	dbPath string
	db     *leveldb.DB
}

var err error
var data []byte

// SetDbPath sets the path the database files are located
func (storage *Storage) SetDbPath(dbPath string) {
	storage.dbPath = dbPath
}

// Run starts the database connection for Storage
func (storage *Storage) Run() error {
	storage.db, err = leveldb.OpenFile(storage.dbPath, nil)
	return err
}

// Close closes the underlying LevelDB connection
func (storage *Storage) Close() {
	storage.db.Close()
}

// Has uses LevelDB's method Has to check does the data exists in LevelDB
func (storage *Storage) Has(key []byte) (bool, error) {
	return storage.db.Has(key, nil)
}

// Get uses LevelDB's method Get to fetch data from LevelDB
func (storage *Storage) Get(key []byte) ([]byte, error) {
	return storage.db.Get(key, nil)
}

// Put uses LevelDB's Put method to put data into LevelDB
func (storage *Storage) Put(key []byte, data []byte) error {
	return storage.db.Put(key, data, nil)
}

// Delete uses LevelDB's Delete method to remove data from LevelDB
func (storage *Storage) Delete(key []byte) error {
	return storage.db.Delete(key, nil)
}

// GetAll returns all entries in the database regardless of key or prefix
func (storage *Storage) GetAll() (map[string]string, error) {
	entries := make(map[string]string)
	iter := storage.db.NewIterator(nil, nil)

	// Iterate over every key in the database, append to entries
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		entries[string(key)] = string(value)
	}

	iter.Release()
	err = errors.E(errors.Op("Get all using iterator"), iter.Error())

	return entries, err
}

// GetAllWithPrefix returns all entries in the database with the specified prefix
func (storage *Storage) GetAllWithPrefix(prefix string) (map[string]string, error) {
	entries := make(map[string]string)
	iter := storage.db.NewIterator(util.BytesPrefix([]byte(prefix)), nil)

	// Iterate over every key in the database, append to entries
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		entries[string(key)] = string(value)
	}

	iter.Release()
	err = errors.E(errors.Op("Get all with prefix using iterator"), iter.Error())

	return entries, err
}

// DeleteAll deletes all entries from the database
// USE CAREFULLY
func (storage *Storage) DeleteAll() error {
	iter := storage.db.NewIterator(nil, nil)

	// Iterate over every key in the database, append to entries
	for iter.Next() {
		key := iter.Key()
		err = errors.E(errors.Op("Delete from storage"), storage.Delete(key))
	}

	iter.Release()
	err = errors.E(errors.Op("Delete all from storage"), iter.Error())

	return err
}

// DeleteAllWithPrefix deletes all entries starting with a prefix
func (storage *Storage) DeleteAllWithPrefix(prefix string) error {
	iter := storage.db.NewIterator(util.BytesPrefix([]byte(prefix)), nil)

	// Iterate over every key in the database, append to entries
	for iter.Next() {
		key := iter.Key()
		err = errors.E(errors.Op("Delete with prefix from storage"), storage.Delete(key))
	}

	iter.Release()
	err = errors.E(errors.Op("Delete all with prefix from storage"), iter.Error())

	return err
}
