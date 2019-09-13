package identity

import (
	"crypto/rand"
	"io"

	"github.com/eqlabs/sprawl/interfaces"
	"github.com/libp2p/go-libp2p-core/crypto"
)

const privateKeyDbKey = "private_key"
const publicKeyDbKey = "public_key"

func GenerateKeyPair(reader io.Reader) (crypto.PrivKey, crypto.PubKey, error) {
	privateKey, publicKey, err := crypto.GenerateEd25519Key(reader)
	return privateKey, publicKey, err
}

func storeKeyPair(storage interfaces.Storage, privateKey crypto.PrivKey, publicKey crypto.PubKey) error {
	privateKeyBytes, err := crypto.MarshalPrivateKey(privateKey)
	if err != nil {
		return err
	}
	publicKeyBytes, err := crypto.MarshalPublicKey(publicKey)
	if err != nil {
		return err
	}

	err = storage.Put([]byte(privateKeyDbKey), privateKeyBytes)
	if err != nil {
		return err
	}

	err = storage.Put([]byte(publicKeyDbKey), publicKeyBytes)
	if err != nil {
		return err
	}

	return nil
}

func getKeyPair(storage interfaces.Storage) (crypto.PrivKey, crypto.PubKey, error) {
	var err error
	hasPrivateKey, err := storage.Has([]byte(privateKeyDbKey))
	if err != nil {
		return nil, nil, err
	}

	if !hasPrivateKey {
		return nil, nil, nil
	}

	hasPublicKey, err := storage.Has([]byte(publicKeyDbKey))
	if err != nil {
		return nil, nil, err
	}

	if !hasPublicKey {
		return nil, nil, nil
	}

	privateKeyBytes, err := storage.Get([]byte(privateKeyDbKey))
	if err != nil {
		return nil, nil, err
	}
	publicKeyBytes, err := storage.Get([]byte(publicKeyDbKey))
	if err != nil {
		return nil, nil, err
	}

	privateKey, err := crypto.UnmarshalPrivateKey(privateKeyBytes)
	if err != nil {
		return nil, nil, err
	}

	publicKey, err := crypto.UnmarshalPublicKey(publicKeyBytes)
	if err != nil {
		return nil, nil, err
	}

	return privateKey, publicKey, nil
}

func GetIdentity(storage interfaces.Storage) (crypto.PrivKey, crypto.PubKey, error) {
	privateKey, publicKey, err := getKeyPair(storage)
	if err != nil {
		return privateKey, publicKey, err
	}

	if privateKey == nil || publicKey == nil {
		privateKey, publicKey, err := GenerateKeyPair(rand.Reader)
		if err != nil {
			return privateKey, publicKey, err
		}
		err = storeKeyPair(storage, privateKey, publicKey)
		return privateKey, publicKey, err
	}
	return privateKey, publicKey, err
}
