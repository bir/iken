package errs

/*
Similar to pkg/errors with tweaks for custom stack skip, stack extraction filtering and marshalling.
*/
import (
	"fmt"
	"runtime"
	"strings"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/pkgerrors"
)

// funcName removes the path prefix component (redundant to the file name) of a function's name reported by func.Name().
//
// example: "github.com/iken/router.(*Router).Handler" => "(*Router).Handler"
func funcName(name string) string {
	i := strings.LastIndex(name, "/")
	name = name[i+1:]
	i = strings.Index(name, ".")

	return name[i+1:]
}

type stack []uintptr

func (s *stack) StackTrace() []uintptr {
	return *s
}

type stackError struct {
	err error
	*stack
}

// Unwrap provides compatibility for Go 1.13 error chains.
func (s *stackError) Unwrap() error {
	return s.err
}

// Error directly returns the wrapped error's Error string.
func (s *stackError) Error() string {
	return s.err.Error()
}

const (
	maxDepth    = 32
	stackOffset = 2
)

// WithStack wraps the `e` and records the stack, skipping `skip` frames in the stack.
// 0 skip is considered the function that calls WithStack.
func WithStack(e interface{}, skip int) error {
	// Clone of pkg.errors.WithStack with support for skip.
	if e == nil {
		return nil
	}

	var err error
	switch eT := e.(type) {
	case error:
		err = eT
	case string:
		err = errors.New(eT)
	default:
		err = errors.Errorf("%v", eT)
	}

	var pcs [maxDepth]uintptr
	n := runtime.Callers(skip+stackOffset, pcs[:])

	var st stack = pcs[0:n]

	return &stackError{err: err, stack: &st}
}

// Frame is a simple view of the call stack frame.
type Frame struct {
	File string
	Line int
	Func string
}

func (f Frame) String() string {
	return fmt.Sprintf("%s:%d %s", f.File, f.Line, f.Func)
}

// ExtractStackFrame extracts an embedded array of stack pointers and converts to array of Frames.
func ExtractStackFrame(err error) []Frame {
	type stackTracer interface{ StackTrace() []uintptr }

	stacker, ok := err.(stackTracer) //nolint:errorlint // false positive
	if !ok {
		return nil
	}

	st := stacker.StackTrace()
	out := make([]Frame, 0, len(st))
	frames := runtime.CallersFrames(st)

	for {
		frame, more := frames.Next()
		if frame.Function != "" {
			f := Frame{
				File: frame.File,
				Line: frame.Line,
				Func: funcName(frame.Function),
			}
			out = append(out, f)
		}

		if !more {
			break
		}
	}

	return out
}

// ExtractStackFrameStop works the same as ExtractStackFrame, but allows defining a stop function for the stack.
// Useful for filtering out stack elements that are mostly redundant and simplifying logging and review.
func ExtractStackFrameStop(err error, stopFuncName string) []Frame {
	stack := ExtractStackFrame(err)
	out := make([]Frame, 0, len(stack))

	for _, f := range stack {
		out = append(out, f)

		if f.Func == stopFuncName {
			return out
		}
	}

	return out
}

// MarshalStack is a helper to extract the stack trace from an error and make it available as an easily
// marshalled object (array of file/line/func).
func MarshalStack(err error) interface{} {
	if out := ExtractStackFrame(err); len(out) > 0 {
		return out
	}

	return pkgerrors.MarshalStack(err)
}
