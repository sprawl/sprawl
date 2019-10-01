package errors

import (
	"bytes"
	"fmt"
	"runtime"
)

// stack is a type that is embedded in an Error struct, and contains
// information about the stack that created that Error.
type stack struct {
	callers []uintptr
}

func (e *Error) populateStack() {
	e.callers = callers()
}

func (e *Error) writeStackToBuffer(buf *bytes.Buffer, callers []uintptr, printCallers []uintptr) {
	var diff bool
	for i := 0; i < len(callers); i++ {
		thisFrame := frame(callers, i)
		if !diff && i < len(printCallers) {
			if thisFrame.Func.Name() == frame(printCallers, i).Func.Name() {
				if thisFrame.Line == frame(printCallers, i).Line {
					// both stacks share this PC, skip it.
					continue
				}
			}
			// No match, don't consider printCallers again.
			diff = true
		}
		fmt.Fprintf(buf, "\n%v:%d", thisFrame.Func.Name(), thisFrame.Line)
	}
	e.writeOpToBuffer(buf)
	e.writeKindToBuffer(buf)
	e.writeErrorToBuffer(buf)

}

func (e *Error) createStackToBuffer(buf *bytes.Buffer) {
	printCallers := callers()
	callers := e.callers
	//Might be unnecessary variable in the end.
	e1 := e
	e1.writeStackToBuffer(buf, callers, printCallers)
	for {
		e2, ok := e1.Err.(*Error)
		if !ok {
			break
		}

		i := 0
		ok = false
		for ; i < len(callers) && i < len(e2.callers); i++ {
			if callers[len(callers)-1-i] != e2.callers[len(e2.callers)-1-i] {
				break
			}
			ok = true
		}
		if ok {
			head := e2.callers[:len(e2.callers)-i]
			e2.writeStackToBuffer(buf, head, printCallers)
			tail := callers
			callers = make([]uintptr, len(head)+len(tail))
			copy(callers, head)
			copy(callers[len(head):], tail)
		}
		e1 = e2
	}
}

// frame returns the nth frame, with the frame at top of stack being 0.
func frame(callers []uintptr, n int) *runtime.Frame {
	frames := runtime.CallersFrames(callers)

	var f runtime.Frame
	for i := len(callers) - 1; i >= n; i-- {
		f, _ = frames.Next()
	}
	return &f
}

// callers is a wrapper for runtime.Callers that allocates a slice.
func callers() []uintptr {
	var stack [64]uintptr
	// Skip 4 stack frames; ok for both E and Error funcs.
	/* Extra comment to clarify the skipping for those who don't understand it at all.
	* The stack frames were for E in my case (24. September):
	* runtime/extern.go Runtime.Callers()
	* debug.go runtime.Callers()
	* debug.go e.callers = callers()
	* errors.go e.populateStack()
	 */
	const skip = 4

	n := runtime.Callers(skip, stack[:])
	return stack[:n]
}
