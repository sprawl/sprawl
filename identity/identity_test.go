package identity

import (
	"crypto/rand"
	"testing"

	"github.com/eqlabs/sprawl/config"
	"github.com/eqlabs/sprawl/db"
	"github.com/eqlabs/sprawl/interfaces"
	"github.com/eqlabs/sprawl/errors"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

const dbPathVar = "database.path"
const testConfigPath = "../config/test"

var storage interfaces.Storage = &db.Storage{}
var testConfig interfaces.Config
var logger *zap.Logger
var log *zap.SugaredLogger

func init() {
	logger, _ = zap.NewProduction()
	log = logger.Sugar()
	testConfig = &config.Config{Logger: log}
	testConfig.ReadConfig(testConfigPath)
}

func TestKeyPairMatching(t *testing.T) {
	privateKey, publicKey, err := GenerateKeyPair(rand.Reader)
	assert.NoError(t, err)
	assert.Equal(t, privateKey.GetPublic(), publicKey)
}

func TestKeyPairStorage(t *testing.T) {
	t.Logf("Database path: %s", testConfig.GetString(dbPathVar))
	storage.SetDbPath(testConfig.GetString(dbPathVar))
	storage.Run()
	defer storage.Close()
	storage.DeleteAll()
	privateKey1, publicKey1, err := GenerateKeyPair(rand.Reader)
	assert.NoError(t, err)
	storeKeyPair(storage, privateKey1, publicKey1)
	privateKey2, publicKey2, errStorage := getKeyPair(storage)
	assert.NoError(t, errStorage)
	assert.Equal(t, privateKey1, privateKey2)
	assert.Equal(t, publicKey1, publicKey2)
}

func TestGetIdentity(t *testing.T) {
	t.Logf("Database path: %s", testConfig.GetString(dbPathVar))
	storage.SetDbPath(testConfig.GetString(dbPathVar))
	storage.Run()
	defer storage.Close()
	storage.DeleteAll()
	privateKey1, publicKey1, err := GetIdentity(storage)
	assert.True(t, errors.IsEmpty(err))
	assert.NotNil(t, privateKey1)
	assert.NotNil(t, publicKey1)
	privateKey2, publicKey2, err := GetIdentity(storage)
	assert.True(t, errors.IsEmpty(err))
	assert.Equal(t, privateKey1, privateKey2)
	assert.Equal(t, publicKey1, publicKey2)
}
