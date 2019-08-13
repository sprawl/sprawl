package db

import (
	"testing"

	"github.com/eqlabs/sprawl/config"
	"github.com/stretchr/testify/assert"
)

var storage = &Storage{}

func init() {
	// Load config
	config := &config.Config{}
	config.ReadConfig("../config/test")

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
