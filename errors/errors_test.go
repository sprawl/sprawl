package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatching(t *testing.T) {
	e1 := E(Op("Get"), Placeholder, "network unreachable")
	e2 := E(Op("Get"), Placeholder, "network unreachable")
	e3 := E(Op("Set"), Placeholder, "network unreachable")
	e4 := E(Op("Get"), Ignore, "network unreachable")
	e5 := E(Op("Get"), Placeholder, "network reachable")
	e6 := E(Op("Get"), "network unreachable")
	e7 := E(Op("Get"), "network unreachable")

	assert.True(t, Match(e1, e2))
	assert.False(t, Match(e1, e3))
	assert.False(t, Match(e1, e4))
	assert.False(t, Match(e1, e5))
	assert.False(t, Match(e1, e6))
	assert.True(t, Match(e6, e7))
}

// tests whether err is an Error with Kind=Permission and User=joe@schmoe.com.
func Match(err1, err2 error) bool {
	e1, ok := err1.(*Error)
	if !ok {
		return false
	}
	e2, ok := err2.(*Error)
	if !ok {
		return false
	}
	if e1.Op != "" && e2.Op != e1.Op {
		return false
	}
	if e1.Kind != Ignore && e2.Kind != e1.Kind {
		return false
	}
	if e1.Err != nil {
		if _, ok := e1.Err.(*Error); ok {
			return Match(e1.Err, e2.Err)
		}
		if e2.Err == nil || e2.Err.Error() != e1.Err.Error() {
			return false
		}
	}
	return true
}
