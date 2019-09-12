package identity

import (
	"crypto/rand"
	"io"

	"github.com/eqlabs/sprawl/interfaces"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/prometheus/common/log"
)

const privateKeyDbKey = "private_key"
const publicKeyDbKey = "public_key"

func GenerateKeyPair(reader io.Reader) (crypto.PrivKey, crypto.PubKey) {
	var err error
	privateKey, publicKey, err := crypto.GenerateEd25519Key(reader)
	if err != nil {
		log.Error(err)
	}
	return privateKey, publicKey
}

func storeKeyPair(storage interfaces.Storage, privateKey crypto.PrivKey, publicKey crypto.PubKey) {
	privateKeyBytes, err := crypto.MarshalPrivateKey(privateKey)
	if err != nil {
		log.Error(err)
	}

	publicKeyBytes, err := crypto.MarshalPublicKey(publicKey)
	if err != nil {
		log.Error(err)
	}

	err = storage.Put([]byte(privateKeyDbKey), privateKeyBytes)
	if err != nil {
		log.Error(err)
	}

	err = storage.Put([]byte(publicKeyDbKey), publicKeyBytes)
	if err != nil {
		log.Error(err)
	}
}

func getKeyPair(storage interfaces.Storage) (crypto.PrivKey, crypto.PubKey) {
	var err error
	privateKeyBytes, err := storage.Get([]byte(privateKeyDbKey))
	if err != nil {
		log.Error(err)
		return nil, nil
	}
	publicKeyBytes, err := storage.Get([]byte(publicKeyDbKey))
	if err != nil {
		log.Error(err)
		return nil, nil
	}

	privateKey, err := crypto.UnmarshalPrivateKey(privateKeyBytes)
	if err != nil {
		log.Error(err)
		return nil, nil
	}

	publicKey, err := crypto.UnmarshalPublicKey(publicKeyBytes)
	if err != nil {
		log.Error(err)
		return nil, nil
	}

	return privateKey, publicKey
}

func GetIdentity(storage interfaces.Storage) (crypto.PrivKey, crypto.PubKey) {
	privateKey, publicKey := getKeyPair(storage)
	if privateKey == nil || publicKey == nil {
		privateKey, publicKey = GenerateKeyPair(rand.Reader)
		storeKeyPair(storage, privateKey, publicKey)
		return privateKey, publicKey
	}
	return privateKey, publicKey
}
