package db

import (
	"fmt"
	"testing"

	"github.com/eqlabs/sprawl/config"
	"github.com/eqlabs/sprawl/interfaces"
	"github.com/stretchr/testify/assert"
)

const testConfigPath = "../config/test"
const dbPathVar = "database.path"
const testID = "0"
const testMessage = "testing"
const prefix1 = "order-"
const prefix2 = "channel-"

var testMessages = make(map[string]string)

var storage interfaces.Storage = &Storage{}

func init() {
	initTestMessages()
	// Load config
	var config interfaces.Config = &config.Config{}
	config.ReadConfig(testConfigPath)
	fmt.Println(config.GetString(dbPathVar))
	// Initialize storage
	storage.SetDbPath(config.GetString(dbPathVar))
	storage.Run()
}

func initTestMessages() {
	testMessages["test1"] = "test1"
	testMessages["test2"] = "test2"
	testMessages["test3"] = "test3"
	testMessages["test4"] = "test4"
}

func deleteAllFromDatabase() {
	storage.DeleteAll()
}

func TestStorageCRUD(t *testing.T) {
	deleteAllFromDatabase()
	storage.Put([]byte(testID), []byte(testMessage))

	testBytes, err := storage.Get([]byte(testID))
	assert.Equal(t, string(testBytes), testMessage)
	assert.Equal(t, err, nil)
	assert.NotEmpty(t, testBytes)

	storage.Delete([]byte(testID))
	deleted, err := storage.Get([]byte(testID))
	assert.Empty(t, deleted)
}

func TestStorageGetAll(t *testing.T) {
	deleteAllFromDatabase()
	for key, value := range testMessages {
		storage.Put([]byte(key), []byte(value))
	}

	var allItems map[string]string
	allItems, err = storage.GetAll()

	assert.Equal(t, err, nil)
	assert.Equal(t, len(allItems), len(testMessages))
}

func TestStorageGetAllByPrefix(t *testing.T) {
	deleteAllFromDatabase()
	for key, value := range testMessages {
		key = prefix1 + key
		storage.Put([]byte(key), []byte(value))
	}

	for key, value := range testMessages {
		key = prefix2 + key
		storage.Put([]byte(key), []byte(value))
	}

	var prefixedItems map[string]string
	prefixedItems, err = storage.GetAllWithPrefix(prefix1)
	var allItems map[string]string
	allItems, err = storage.GetAll()

	assert.Equal(t, err, nil)
	assert.Equal(t, len(prefixedItems), len(testMessages))
	assert.Equal(t, len(allItems), len(testMessages)*2)
}
