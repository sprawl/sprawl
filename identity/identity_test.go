package identity

import (
	"crypto/rand"
	"testing"

	"github.com/gogo/protobuf/proto"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/sprawl/sprawl/config"
	"github.com/sprawl/sprawl/database/leveldb"
	"github.com/sprawl/sprawl/errors"
	"github.com/sprawl/sprawl/interfaces"
	"github.com/sprawl/sprawl/pb"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

const testConfigPath = "../config/test"

var storage interfaces.Storage = &leveldb.Storage{}
var testConfig interfaces.Config
var logger *zap.Logger
var log *zap.SugaredLogger

func init() {
	logger = zap.NewNop()
	log = logger.Sugar()
	testConfig = &config.Config{}
	testConfig.ReadConfig(testConfigPath)
}

func TestKeyPairStorage(t *testing.T) {
	t.Logf("Database path: %s", testConfig.GetDatabasePath())
	storage.SetDbPath(testConfig.GetDatabasePath())
	storage.Run()
	defer storage.Close()
	storage.DeleteAll()
	privateKey1, publicKey1, err := GenerateKeyPair(storage, rand.Reader)
	assert.True(t, errors.IsEmpty(err))
	privateKey2, publicKey2, errStorage := getKeyPair(storage)
	assert.NoError(t, errStorage)
	assert.Equal(t, privateKey1, privateKey2)
	assert.Equal(t, publicKey1, publicKey2)
}

func TestGetIdentity(t *testing.T) {
	t.Logf("Database path: %s", testConfig.GetDatabasePath())
	storage.SetDbPath(testConfig.GetDatabasePath())
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

func TestSignAndVerify(t *testing.T) {
	t.Logf("Database path: %s", testConfig.GetDatabasePath())
	storage.SetDbPath(testConfig.GetDatabasePath())
	storage.Run()
	defer storage.Close()
	storage.DeleteAll()
	_, publicKey, err := GetIdentity(storage)
	assert.True(t, errors.IsEmpty(err))
	testOrder := &pb.Order{Asset: string("ETH"), CounterAsset: string("BTC"), Amount: 52152, Price: 0.2, Id: []byte("jgkahgkjal")}
	testOrderInBytes, err := proto.Marshal(testOrder)
	assert.NoError(t, err)
	sig, err := Sign(storage, testOrderInBytes)
	assert.NoError(t, err)
	publicKeyBytes, err := crypto.MarshalPublicKey(publicKey)
	assert.NoError(t, err)
	legit, err := Verify(publicKeyBytes, testOrderInBytes, sig)
	assert.NoError(t, err)
	assert.True(t, legit)

}
