package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var storage = &Storage{}

func init() {
	// Initialize storage
	storage.SetDbPath("/var/lib/sprawl/test")
	storage.Run()
}

func TestStorage(t *testing.T) {
	storage.Put([]byte("0"), []byte("testing"))

	testBytes, err := storage.Get([]byte("0"))
	assert.Equal(t, string(testBytes), "testing")
	assert.Equal(t, err, nil)
	assert.NotEmpty(t, testBytes)

	storage.Delete([]byte("0"))
	deleted, err := storage.Get([]byte("0"))
	assert.Empty(t, deleted)
}
