package errors

import (
	"os"
	"regexp"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testRegexStringStart = "\nsprawl/errors/debug_test.go:"
const testRegexStringSecond = ": Set\nsprawl/errors/debug_test.go:"
const testRegexStringThird = ": Get\nsprawl/errors/debug_test.go:"
const testRegexStringEnd = ": Set: placeholder error:\n\tnetwork unreachable"

var testRegexPattern = regexp.QuoteMeta(testRegexStringStart) + "\\d+" + regexp.QuoteMeta(testRegexStringSecond) + "\\d+" + regexp.QuoteMeta(testRegexStringThird) + "\\d+" + regexp.QuoteMeta(testRegexStringEnd)

func TestFileClean(t *testing.T) {
	dir, _ := os.Getwd()
	dir = strings.ReplaceAll(dir, "\\", "/")
	assert.Equal(t, tryCleanDirPath(dir), mainDirName+"/errors")
}

func TestStackPopulation(t *testing.T) {
	e := E(testOpGet).(*Error)
	e.populateStack()
	assert.Equal(t, e.callers[1], callers()[0])
}

func TestFrames(t *testing.T) {
	callers := callers()
	frames := runtime.CallersFrames(callers)
	f, _ := frames.Next()
	assert.Equal(t, f, *frame(callers, 0))
}

func TestDebug(t *testing.T) {
	e1 := E(testOpSet, Placeholder, testStringNetworkUnreachable)
	assert.NotNil(t, e1)
	e2 := E(testOpGet, e1)
	assert.NotNil(t, e2)
	e3 := E(testOpSet, e2)
	assert.NotNil(t, e3)
	e4 := E(testOpSet, e2)
	assert.NotNil(t, e4)

	assert.Equal(t, e3.Error(), e4.Error())
	debug = true

	d1 := E(testOpSet, Placeholder, testStringNetworkUnreachable)
	assert.NotNil(t, d1)
	d2 := E(testOpGet, d1)
	assert.NotNil(t, d2)
	d3 := E(testOpSet, d2)
	assert.NotNil(t, d3)
	d4 := E(testOpSet, d2)
	assert.NotNil(t, d4)
	d3Matches, err := regexp.MatchString(testRegexPattern, d3.Error())
	assert.NoError(t, err)
	assert.True(t, d3Matches)
	d4Matches, err := regexp.MatchString(testRegexPattern, d4.Error())
	assert.NoError(t, err)
	assert.True(t, d4Matches)
	assert.NotEqual(t, d3.Error(), d4.Error())
	assert.NotEqual(t, e1.Error(), d1.Error())
	assert.NotEqual(t, e2.Error(), d2.Error())
	assert.NotEqual(t, e3.Error(), d3.Error())
	debug = false
}
