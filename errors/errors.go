package errors

import (
	"bytes"
	"fmt"
)

type Error struct {
	Op   Op
	Kind Kind
	Err  error
}

type Op string
type Kind uint8

var Separator = ":\n\t"

const (
	Ignore Kind = iota //Unclassified
	Placeholder
)

func (e *Error) isZero() bool {
	return e.Op == "" && e.Kind == 0 && e.Err == nil
}

func (k Kind) String() string {
	switch k {
	case Ignore:
		return "ignored kind of error"
	case Placeholder:
		return "placeholder error"
	}
	return "unknown error kind"
}

func E(args ...interface{}) error {
	if len(args) == 0 {
		panic("call to errors.E with no arguments")
	}
	
	e := &Error{}
	for _, arg := range args {
		switch arg := arg.(type) {
		case Op:
			e.Op = arg
		case Kind:
			e.Kind = arg
		case string:
			e.Err = StringToError(arg)
		case *Error:
			// Make a copy
			copy := *arg
			e.Err = &copy
		case error:
			e.Err = arg
		default:
			return Errorf("unknown type %T, value %v in error call", arg, arg)
		}
	}
	return e
}

// pad appends str to the buffer if the buffer already has some data.
func pad(b *bytes.Buffer, str string) {
	if b.Len() == 0 {
		return
	}
	b.WriteString(str)
}

func (e *Error) Error() string {
	b := new(bytes.Buffer)
	if e.Op != "" {
		pad(b, ": ")
		b.WriteString(string(e.Op))
	}
	//Ignore if Ignore kind
	if e.Kind != 0 {
		pad(b, ": ")
		b.WriteString(e.Kind.String())
	}
	if e.Err != nil {
		// Indent on new line if we are cascading non-empty errors.
		if prevErr, ok := e.Err.(*Error); ok {
			if !prevErr.isZero() {
				pad(b, Separator)
				b.WriteString(e.Err.Error())
			}
		} else {
			pad(b, ": ")
			b.WriteString(e.Err.Error())
		}
	}
	if b.Len() == 0 {
		return "no error"
	}
	return b.String()
}

// Str returns an error that formats as the given text. It is intended to
// be used as the error-typed argument to the E function.
func StringToError(text string) error {
	return &errorString{text}
}

// errorString is a trivial implementation of error.
type errorString struct {
	s string
}

func (e *errorString) Error() string {
	return e.s
}

// package for all error handling.
func Errorf(format string, args ...interface{}) error {
	return &errorString{fmt.Sprintf(format, args...)}
}
