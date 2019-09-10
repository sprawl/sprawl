package identity

import (
	"io"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/prometheus/common/log"
)

func GenerateKeyPair(reader io.Reader) (crypto.PrivKey, crypto.PubKey) {
	var err error
	privateKey, publicKey, err := crypto.GenerateEd25519Key(reader)
	log.Error(err)
	return privateKey, publicKey
}