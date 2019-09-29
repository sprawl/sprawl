package errors

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContent(t *testing.T) {
	e1 := E(Op("Get"), Placeholder, "network unreachable")
	e2 := E(Op("Set"), Placeholder, e1)
	e3 := errors.New("network unreachable")
	e4 := E(e3)
	e5 := E(321)

	assert.Equal(t, e2.(*Error).Op, Op("Set"))
	assert.NotEqual(t, e2.(*Error).Op, Op("Get"))
	assert.Equal(t, e2.(*Error).Kind, Placeholder)
	assert.NotEqual(t, e2.(*Error).Kind, Ignore)
	assert.Equal(t, e2.(*Error).Err, e1)
	assert.NotEqual(t, e2.(*Error).Err, e2)
	assert.Equal(t, fmt.Sprintf("%s", e1.(*Error).Err), fmt.Sprintf("%s", e3))
	assert.NotEqual(t, fmt.Sprintf("%s", e1.(*Error).Err), fmt.Sprintf("%s", e2))
	assert.Equal(t, fmt.Sprintf("%s", e1.(*Error).Err), fmt.Sprintf("%s", e4.(*Error).Err))
	assert.NotEqual(t, fmt.Sprintf("%s", e1.(*Error).Err), fmt.Sprintf("%s", e2.(*Error).Err))
	assert.Equal(t, e5, Errorf("unknown type %T, value %v in error call", 321, 321))
	assert.NotEqual(t, e5, Errorf("unknown type %T, value %v in error call", "321", "321"))
}

func TestIsZero(t *testing.T) {
	e1 := E(Ignore)
	assert.True(t, e1.(*Error).isZero())
	assert.Equal(t, e1.Error(), "no error")
}

func TestPad(t *testing.T) {
	buffer := new(bytes.Buffer)
	pad(buffer, "test")
	assert.Equal(t, buffer.String(), "")
	buffer.WriteString("test")
	pad(buffer, "test")
	assert.Equal(t, buffer.String(), "testtest")
}

func TestBufferWriting(t *testing.T) {
	e1 := E(Op("Get"), Placeholder, "network unreachable")
	buffer := new(bytes.Buffer)
	e1.(*Error).writeOpToBuffer(buffer)
	assert.Equal(t, buffer.String(), "Get")
	e1.(*Error).writeKindToBuffer(buffer)
	assert.Equal(t, buffer.String(), "Get: placeholder error")
	e1.(*Error).writeErrorToBuffer(buffer)
	assert.Equal(t, buffer.String(), "Get: placeholder error:\n\tnetwork unreachable")

}
