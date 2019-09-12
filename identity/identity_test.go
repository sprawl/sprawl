package identity

import (
	"crypto/rand"
	"testing"

	"github.com/eqlabs/sprawl/config"
	"github.com/eqlabs/sprawl/db"
	"github.com/eqlabs/sprawl/interfaces"
	"github.com/stretchr/testify/assert"
)

var storage interfaces.Storage = &db.Storage{}

const dbPathVar = "database.path"
const testConfigPath = "../config/test"

func TestKeyPairMatching(t *testing.T) {
	privateKey, publicKey, err := GenerateKeyPair(rand.Reader)
	assert.NoError(t, err)
	assert.Equal(t, privateKey.GetPublic(), publicKey)
}

func TestKeyPairStorage(t *testing.T) {
	// Load config
	var config interfaces.Config = &config.Config{}
	config.ReadConfig(testConfigPath)
	t.Log(config.GetString(dbPathVar))
	// Initialize storage
	storage.SetDbPath(config.GetString(dbPathVar))
	storage.Run()
	defer storage.Close()
	storage.DeleteAll()
	privateKey1, publicKey1, err := GenerateKeyPair(rand.Reader)
	assert.NoError(t, err)
	storeKeyPair(storage, privateKey1, publicKey1)
	privateKey2, publicKey2, err_storage := getKeyPair(storage)
	assert.NoError(t, err_storage)
	assert.Equal(t, privateKey1, privateKey2)
	assert.Equal(t, publicKey1, publicKey2)
}

func TestGetIdentity(t *testing.T) {
	// Load config
	var config interfaces.Config = &config.Config{}
	config.ReadConfig(testConfigPath)
	t.Log(config.GetString(dbPathVar))
	// Initialize storage
	storage.SetDbPath(config.GetString(dbPathVar))
	storage.Run()
	defer storage.Close()
	storage.DeleteAll()
	privateKey1, publicKey1, err_storage, err := GetIdentity(storage)
	assert.Error(t, err_storage)
	assert.NoError(t, err)
	assert.NotNil(t, privateKey1)
	assert.NotNil(t, publicKey1)
	privateKey2, publicKey2, err_storage, err := GetIdentity(storage)
	assert.NoError(t, err_storage)
	assert.NoError(t, err)
	assert.Equal(t, privateKey1, privateKey2)
	assert.Equal(t, publicKey1, publicKey2)
}
