package util

import "reflect"

// IsInstanceOf checks if objectPtr matches typePtr. Thanks to https://stackoverflow.com/a/48145372
func IsInstanceOf(objectPtr, typePtr interface{}) bool {
	return reflect.TypeOf(objectPtr) == reflect.TypeOf(typePtr)
}
