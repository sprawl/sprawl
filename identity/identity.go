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

// GenerateKeyPair generates a private and a public key to use with libp2p peer and stores it
func GenerateKeyPair(storage interfaces.Storage, reader io.Reader) (crypto.PrivKey, crypto.PubKey, error) {
	privateKey, publicKey, err := crypto.GenerateEd25519Key(reader)
	if !errors.IsEmpty(err) {
		return privateKey, publicKey, errors.E(errors.Op("Generate key pair"), err)
	}
	return privateKey, publicKey, storeKeyPair(storage, privateKey, publicKey)
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

func hasKeyPair(storage interfaces.Storage) (bool, error) {
	hasPrivateKey, err := storage.Has([]byte(privateKeyDbKey))
	if !errors.IsEmpty(err) {
		return false, errors.E(errors.Op("Check private key from storage"), err)
	}
	if !hasPrivateKey {
		return false, nil
	}
	hasPublicKey, err := storage.Has([]byte(publicKeyDbKey))
	if !errors.IsEmpty(err) {
		return false, errors.E(errors.Op("Check public key from storage"), err)
	}
	return hasPublicKey, nil
}

func getKeyPair(storage interfaces.Storage) (crypto.PrivKey, crypto.PubKey, error) {
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

// GetIdentity returns the created private and public key from storage
func GetIdentity(storage interfaces.Storage) (crypto.PrivKey, crypto.PubKey, error) {
	hasKeyPair, err := hasKeyPair(storage)
	if !errors.IsEmpty(err) {
		return nil, nil, errors.E(errors.Op("Has key pair"), err)
	}
	if hasKeyPair {
		privateKey, publicKey, err := getKeyPair(storage)
		if !errors.IsEmpty(err) {
			return privateKey, publicKey, errors.E(errors.Op("Get key pair"), err)
		}
		return privateKey, publicKey, nil
	} else {
		privateKey, publicKey, err := GenerateKeyPair(storage, rand.Reader)
		return privateKey, publicKey, errors.E(errors.Op("Generate key pair"), err)
	}
}

// Sign returns a signature for given data with this node's identity
func Sign(storage interfaces.Storage, data []byte) (signature []byte, err error) {
	privateKey, _, err := GetIdentity(storage)
	if !errors.IsEmpty(err) {
		return nil, errors.E(errors.Op("Sign"), err)
	}
	return privateKey.Sign(data)
}

// Verify verifies data and its signature with a public key
func Verify(publicKey []byte, data []byte, signature []byte) (success bool, err error) {
	pubKey, err := crypto.UnmarshalEd25519PublicKey(publicKey)
	if !errors.IsEmpty(err) {
		return false, errors.E(errors.Op("Unmarshal Ed25519 public key in Verify"), err)
	}
	return pubKey.Verify(data, signature)
}
