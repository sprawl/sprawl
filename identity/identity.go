package identity

import (
	"crypto/rand"
	"io"

	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/sprawl/sprawl/errors"
	"github.com/sprawl/sprawl/interfaces"
)

const privateKeyDbKey = "private_key"
const publicKeyDbKey = "public_key"

func GenerateKeyPair(reader io.Reader) (crypto.PrivKey, crypto.PubKey, error) {
	privateKey, publicKey, err := crypto.GenerateEd25519Key(reader)
	return privateKey, publicKey, err
}

func storeKeyPair(storage interfaces.Storage, privateKey crypto.PrivKey, publicKey crypto.PubKey) error {
	privateKeyBytes, err := crypto.MarshalPrivateKey(privateKey)
	if !errors.IsEmpty(err) {
		return errors.E(errors.Op("Marshal Private Key"), err)
	}
	publicKeyBytes, err := crypto.MarshalPublicKey(publicKey)
	if !errors.IsEmpty(err) {
		return errors.E(errors.Op("Marshal Public Key"), err)
	}

	err = storage.Put([]byte(privateKeyDbKey), privateKeyBytes)
	if !errors.IsEmpty(err) {
		return errors.E(errors.Op("Store Private Key"), err)
	}

	err = storage.Put([]byte(publicKeyDbKey), publicKeyBytes)
	if !errors.IsEmpty(err) {
		return errors.E(errors.Op("Store Public Key"), err)
	}

	return nil
}

func getKeyPair(storage interfaces.Storage) (crypto.PrivKey, crypto.PubKey, error) {
	var err error
	hasPrivateKey, err := storage.Has([]byte(privateKeyDbKey))
	if !errors.IsEmpty(err) {
		return nil, nil, errors.E(errors.Op("Check private key from storage"), err)
	}

	if !hasPrivateKey {
		return nil, nil, nil
	}

	hasPublicKey, err := storage.Has([]byte(publicKeyDbKey))
	if !errors.IsEmpty(err) {
		return nil, nil, errors.E(errors.Op("Check public key from storage"), err)
	}

	if !hasPublicKey {
		return nil, nil, nil
	}

	privateKeyBytes, err := storage.Get([]byte(privateKeyDbKey))
	if !errors.IsEmpty(err) {
		return nil, nil, errors.E(errors.Op("Get private key from storage"), err)
	}
	publicKeyBytes, err := storage.Get([]byte(publicKeyDbKey))
	if !errors.IsEmpty(err) {
		return nil, nil, errors.E(errors.Op("Get public key from storage"), err)
	}

	privateKey, err := crypto.UnmarshalPrivateKey(privateKeyBytes)
	if !errors.IsEmpty(err) {
		return nil, nil, errors.E(errors.Op("Unmarshal private key"), err)
	}

	publicKey, err := crypto.UnmarshalPublicKey(publicKeyBytes)
	if !errors.IsEmpty(err) {
		return nil, nil, errors.E(errors.Op("Unmarshal public key"), err)
	}

	return privateKey, publicKey, nil
}

func GetIdentity(storage interfaces.Storage) (crypto.PrivKey, crypto.PubKey, error) {
	privateKey, publicKey, err := getKeyPair(storage)
	if !errors.IsEmpty(err) {
		return privateKey, publicKey, errors.E(errors.Op("Get key pair"), err)
	}

	if privateKey == nil || publicKey == nil {
		privateKey, publicKey, err := GenerateKeyPair(rand.Reader)
		if !errors.IsEmpty(err) {
			return privateKey, publicKey, errors.E(errors.Op("Generate key pair"), err)
		}
		err = storeKeyPair(storage, privateKey, publicKey)
		return privateKey, publicKey, errors.E(errors.Op("Store key pair"), err)
	}
	return privateKey, publicKey, err
}
