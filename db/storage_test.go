package db

import (
	"testing"

	"github.com/sprawl/sprawl/config"
	"github.com/sprawl/sprawl/errors"
	"github.com/sprawl/sprawl/interfaces"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

const testConfigPath = "../config/test"
const dbPathVar = "database.path"
const testID = "0"
const testMessage = "testing"
const orderPrefix = "order-"
const channelPrefix = "channel-"

var testMessages = make(map[string]string)

var storage interfaces.Storage = &Storage{}
var logger *zap.Logger
var log *zap.SugaredLogger

func init() {
	initTestMessages()
	logger = zap.NewNop()
	log = logger.Sugar()
	// Load config
	var config interfaces.Config = &config.Config{Logger: log}
	config.ReadConfig(testConfigPath)
	log.Info(config.GetString(dbPathVar))
	// Initialize storage
	storage.SetDbPath(config.GetString(dbPathVar))
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
	storage.Run()
	defer storage.Close()
	deleteAllFromDatabase()

	storage.Put([]byte(testID), []byte(testMessage))

	testBytes, err := storage.Get([]byte(testID))
	testBool, err := storage.Has([]byte(testID))
	assert.True(t, testBool)
	assert.Equal(t, string(testBytes), testMessage)
	assert.NoError(t, err)
	assert.NotEmpty(t, testBytes)

	storage.Delete([]byte(testID))
	deleted, err := storage.Get([]byte(testID))
	testBool, err = storage.Has([]byte(testID))
	assert.False(t, testBool)
	assert.Empty(t, deleted)
}

func TestStorageGetAll(t *testing.T) {
	storage.Run()
	defer storage.Close()
	deleteAllFromDatabase()

	for key, value := range testMessages {
		storage.Put([]byte(key), []byte(value))
	}

	var allItems map[string]string
	allItems, err = storage.GetAll()

	assert.True(t, errors.IsEmpty(err))
	assert.Equal(t, len(allItems), len(testMessages))
}

func TestStorageGetAllWithPrefix(t *testing.T) {
	storage.Run()
	defer storage.Close()
	deleteAllFromDatabase()

	for key, value := range testMessages {
		key = orderPrefix + key
		storage.Put([]byte(key), []byte(value))
	}

	for key, value := range testMessages {
		key = channelPrefix + key
		storage.Put([]byte(key), []byte(value))
	}

	var prefixedItems map[string]string
	prefixedItems, err = storage.GetAllWithPrefix(orderPrefix)
	var allItems map[string]string
	allItems, err = storage.GetAll()

	assert.True(t, errors.IsEmpty(err))
	assert.Equal(t, len(prefixedItems), len(testMessages))
	assert.Equal(t, len(allItems), len(testMessages)*2)
}

func TestStorageDeleteAllWithPrefix(t *testing.T) {
	storage.Run()
	defer storage.Close()
	deleteAllFromDatabase()

	for key, value := range testMessages {
		key = orderPrefix + key
		storage.Put([]byte(key), []byte(value))
	}

	storage.DeleteAllWithPrefix(orderPrefix)

	var prefixedItems map[string]string
	prefixedItems, err = storage.GetAllWithPrefix(orderPrefix)

	assert.Zero(t, len(prefixedItems))
}

func BenchmarkAdd(b *testing.B) {
	storage.Run()
	defer storage.Close()
	deleteAllFromDatabase()

	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		storage.Put([]byte(string(i)), []byte(testMessage+string(i)))
	}
}

func BenchmarkRead(b *testing.B) {
	storage.Run()
	defer storage.Close()

	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		storage.Get([]byte(string(i)))
	}
}
