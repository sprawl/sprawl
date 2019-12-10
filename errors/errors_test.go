package errors

import (
	"bytes"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testInt123 = 123
const testString123 = "123"
const testStringTest = "test"
const testStringNoError = "no error"
const testStringNetworkReachable = "network reachable"
const testStringNetworkUnreachable = "network unreachable"
const testStringGet = "Get"
const testStringSet = "Set"
const testOpGet = Op(testStringGet)
const testOpSet = Op(testStringSet)

func TestSetDebug(t *testing.T) {
	assert.False(t, debug)
	SetDebug(true)
	assert.True(t, debug)
	SetDebug(false)
	assert.False(t, debug)
}

func TestIsEmpty(t *testing.T) {
	var e2 error
	e1 := errors.New(testStringNetworkUnreachable)
	e2 = nil
	e3 := E(testOpGet, Placeholder, testStringNetworkUnreachable)
	e4 := E(testOpGet, Placeholder)
	e5 := E(testOpGet, Placeholder, e3)
	e6 := E(testOpGet, Placeholder, e4)
	e7 := E(testOpGet, Placeholder, e5)
	e8 := E(testOpGet, Placeholder, e6)

	assert.False(t, IsEmpty(e1))
	assert.True(t, IsEmpty(e2))
	assert.False(t, IsEmpty(e3))
	assert.True(t, IsEmpty(e4))
	assert.False(t, IsEmpty(e5))
	assert.True(t, IsEmpty(e6))
	assert.False(t, IsEmpty(e7))
	assert.True(t, IsEmpty(e8))
}

func TestContent(t *testing.T) {
	e1 := E(testOpGet, Placeholder, testStringNetworkUnreachable)
	e2 := E(testOpSet, Placeholder, e1)
	e3 := errors.New(testStringNetworkUnreachable)
	e4 := E(e3)
	e5 := E(testInt123)

	assert.Equal(t, e2.(*Error).Op, testOpSet)
	assert.NotEqual(t, e2.(*Error).Op, testOpGet)
	assert.Equal(t, e2.(*Error).Kind, Placeholder)
	assert.NotEqual(t, e2.(*Error).Kind, Ignore)
	assert.Equal(t, e2.(*Error).Err, e1)
	assert.NotEqual(t, e2.(*Error).Err, e2)
	assert.Equal(t, fmt.Sprintf("%s", e1.(*Error).Err), fmt.Sprintf("%s", e3))
	assert.NotEqual(t, fmt.Sprintf("%s", e1.(*Error).Err), fmt.Sprintf("%s", e2))
	assert.Equal(t, fmt.Sprintf("%s", e1.(*Error).Err), fmt.Sprintf("%s", e4.(*Error).Err))
	assert.NotEqual(t, fmt.Sprintf("%s", e1.(*Error).Err), fmt.Sprintf("%s", e2.(*Error).Err))
	assert.Equal(t, e5, Errorf("unknown type %T, value %v in error call", testInt123, testInt123))
	assert.NotEqual(t, e5, Errorf("unknown type %T, value %v in error call", testString123, testString123))
}

func TestIsZero(t *testing.T) {
	e1 := E(Ignore)
	assert.True(t, e1.(*Error).isZero())
	assert.Equal(t, e1.Error(), testStringNoError)
}

func TestPad(t *testing.T) {
	buffer := new(bytes.Buffer)
	pad(buffer, testStringTest)
	assert.Equal(t, buffer.String(), "")
	buffer.WriteString(testStringTest)
	pad(buffer, testStringTest)
	assert.Equal(t, buffer.String(), testStringTest+testStringTest)
}

func TestBufferWriting(t *testing.T) {
	e1 := E(testOpGet, Placeholder, testStringNetworkUnreachable)
	buffer := new(bytes.Buffer)
	e1.(*Error).writeOpToBuffer(buffer)
	assert.Equal(t, buffer.String(), "Get")
	e1.(*Error).writeKindToBuffer(buffer)
	assert.Equal(t, buffer.String(), "Get: placeholder error")
	e1.(*Error).writeErrorToBuffer(buffer)
	assert.Equal(t, buffer.String(), "Get: placeholder error:\n\tnetwork unreachable")

}
