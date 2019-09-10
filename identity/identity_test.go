package identity

import(
	"testing"
	"github.com/stretchr/testify/assert"
	"crypto/rand"
)

func TestKeyPair(t *testing.T) {
	privateKey, publicKey := GenerateKeyPair(rand.Reader)
	assert.Equal(t, privateKey.GetPublic(), publicKey)
}