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

var storage interfaces.Storage = &Storage{}

func init() {
	// Load config
	var config interfaces.Config = &config.Config{}
	config.ReadConfig(testConfigPath)
	fmt.Println(config.GetString(dbPathVar))
	// Initialize storage
	storage.SetDbPath(config.GetString(dbPathVar))
	storage.Run()
}

func TestStorage(t *testing.T) {
	storage.Put([]byte(testID), []byte(testMessage))

	testBytes, err := storage.Get([]byte(testID))
	assert.Equal(t, string(testBytes), testMessage)
	assert.Equal(t, err, nil)
	assert.NotEmpty(t, testBytes)

	storage.Delete([]byte(testID))
	deleted, err := storage.Get([]byte(testID))
	assert.Empty(t, deleted)
}
