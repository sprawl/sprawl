package db

import (
	"github.com/syndtr/goleveldb/leveldb"
)

// Storage is a struct containing a database and its address
type Storage struct {
	dbPath string
	db     *leveldb.DB
}

type Entry struct {
	key   []byte
	value []byte
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

// GetAll spews out all that's saved in the database regardless of key or prefix
func (storage *Storage) GetAll() ([]Entry, error) {
	entries := make([]Entry, 0)
	iter := storage.db.NewIterator(nil, nil)

	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		entries = append(entries, Entry{key: key, value: value})
	}

	iter.Release()
	err = iter.Error()

	return entries, err
}
