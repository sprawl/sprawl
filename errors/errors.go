package errors

import (
	"bytes"
	"fmt"
)

type Error struct {
	Op   Op
	Kind Kind
	Err  error
	stack
}

type Op string
type Kind uint8

const separator = ":\n\t"
const colon = ": "

var debug = false

const (
	Ignore Kind = iota //Unclassified
	Placeholder
)

func (e *Error) isZero() bool {
	return e.Op == "" && e.Kind == 0 && e.Err == nil
}

func SetDebug(debugOn bool) {
	debug = debugOn
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

func E(argument interface{}, arguments ...interface{}) error {
	e := &Error{}
	for _, arg := range append([]interface{}{argument}, arguments...) {
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
		case nil:
			e.Err = nil
		default:
			return Errorf("unknown type %T, value %v in error call", arg, arg)
		}
	}
	// Populate stack information (only in debug mode).
	if debug {
		e.populateStack()
	}
	return e
}

func IsEmpty(err error) bool {
	if err == nil {
		return true
	}
	e, ok := err.(*Error)
	if ok == false {
		return false
	}
	for {
		if e.Err == nil {
			return true
		}
		e, ok = e.Err.(*Error)
		if ok == false {
			return false
		}
	}
}

// pad appends str to the buffer if the buffer already has some data.
func pad(b *bytes.Buffer, str string) {
	if b.Len() == 0 {
		return
	}
	b.WriteString(str)
}

func (e *Error) writeOpToBuffer(buf *bytes.Buffer) {
	if e.Op != "" {
		pad(buf, colon)
		buf.WriteString(string(e.Op))
	}
}

func (e *Error) writeKindToBuffer(buf *bytes.Buffer) {
	if e.Kind != 0 {
		pad(buf, colon)
		buf.WriteString(e.Kind.String())
	}
}

func (e *Error) writeErrorToBuffer(buf *bytes.Buffer) {
	if e.Err == nil {
		return
	}
	if prevErr, ok := e.Err.(*Error); ok {
		if debug {
			return
		}
		if prevErr.isZero() {
			return
		}
	}
	pad(buf, separator)
	buf.WriteString(e.Err.Error())
}

func (e *Error) Error() string {
	b := new(bytes.Buffer)
	if debug {
		e.createStackToBuffer(b)
	} else {
		e.writeOpToBuffer(b)
		e.writeKindToBuffer(b)
		e.writeErrorToBuffer(b)
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
