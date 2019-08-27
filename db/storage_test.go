package db

import (
	"fmt"
	"testing"

	"github.com/eqlabs/sprawl/config"
	"github.com/eqlabs/sprawl/interfaces"
	"github.com/stretchr/testify/assert"
)

var storage interfaces.Storage = &Storage{}

func init() {
	// Load config
	var config interfaces.Config = &config.Config{}
	config.ReadConfig("../config/test")
	fmt.Println(config.GetString("database.path"))
	// Initialize storage
	storage.SetDbPath(config.GetString("database.path"))
	storage.Run()
}

func TestStorage(t *testing.T) {
	storage.Put([]byte("0"), []byte("testing"))

	testBytes, err := storage.Get([]byte("0"))
	assert.Equal(t, string(testBytes), "testing")
	assert.Equal(t, err, nil)
	assert.NotEmpty(t, testBytes)

	storage.Delete([]byte("0"))
	deleted, err := storage.Get([]byte("0"))
	assert.Empty(t, deleted)
}

func TestStorageClosing(t *testing.T) {
	storage.Close()
	assert.Panics(t, func() { storage.Put([]byte("1"), []byte("testing")) }, "The database connection should be closed")
}
