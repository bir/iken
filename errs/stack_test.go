package errs_test

import (
	"fmt"
	"testing"

	"github.com/bir/iken/errs"
	"github.com/pkg/errors"
)

// For testing stacks we only check the function names in the stack, to
//reduce brittleness of the tests with the version of go
var (
	stackTests = []StackTest{
		{
			// Error
			errs.WithStack(errors.New("pkgError"), -2), "pkgError", false, "",
			[]string{"Callers",
				"WithStack",
				"init",
				"doInit",
				"doInit",
				"main",
				"goexit"},
		},
		{
			// Default
			errs.WithStack(123, -1), "123", false, "",
			[]string{"WithStack",
				"init",
				"doInit",
				"doInit",
				"main",
				"goexit"},
		},
		{
			// String
			errs.WithStack("1", 0), "1", false, "",
			[]string{"init",
				"doInit",
				"doInit",
				"main",
				"goexit"},
		},
		{
			// Nil
			errs.WithStack(nil, 0), "", true, "", nil,
		},
	}
	stopTests = []StackTest{
		{
			// String
			errs.WithStack("1", 0), "1", false, "doInit",
			[]string{"init",
				"doInit"},
		}, {
			// String
			errs.WithStack("1", 0), "1", false, "No match for this stop",
			[]string{"init",
				"doInit",
				"doInit",
				"main",
				"goexit"},
		},
	}
	marshallTests = []StackTest{
		{
			// errs
			errs.WithStack("errs.WithStack", 0), "errs.WithStack", false, "",
			[]string{"init",
				"doInit",
				"doInit",
				"main",
				"goexit"},
		}, {
			// pkg.errors - uses a truncated file name
			errors.New("pkgErrors passthru"), "pkgErrors passthru", false, "",
			[]string{"init",
				"doInit",
				"doInit",
				"main",
				"goexit"},
		},
	}
)

type StackTest struct {
	err   error
	want  string
	isNil bool
	stop  string
	stack []string
}

func (test StackTest) testErr(t *testing.T) bool {
	if test.err == nil {
		if test.isNil {
			return false
		}
		t.Errorf("expected %#v, got nil", test.want)
		return false
	}
	got := test.err.Error()
	if got != test.want {
		t.Errorf("expected %#v, got %#v", test.want, got)
		return false
	}
	return true
}

func (test StackTest) testFrameString(t *testing.T, ff []errs.Frame) {
	if len(ff) != len(test.stack) {
		t.Errorf("len(stack) expected %#v, got %#v", len(test.stack), len(ff))
		return
	}
	for i, f := range ff {
		if f.Func != test.stack[i] {
			t.Errorf("stack[%d] expected `%#v`, got `%#v`", i, test.stack[i], f.String())
			return
		}
	}
}

func (test StackTest) testFramesMap(t *testing.T, ff []map[string]string) {
	if len(ff) != len(test.stack) {
		t.Errorf("len(stack) expected %#v, got %#v", len(test.stack), len(ff))
		return
	}
	for i, f := range ff {
		l := fmt.Sprintf("%s:%s %s", f["source"], f["line"], f["func"])
		if f["func"] != test.stack[i] {
			t.Errorf("stack expected `%#v`, got `%#v`", test.stack[i], l)
			return
		}
	}
}

func TestStack(t *testing.T) {

	for _, test := range stackTests {
		if test.testErr(t) {
			test.testFrameString(t, errs.ExtractStackFrame(test.err))
		}
	}

	got := errs.ExtractStackFrame(nil)
	if got != nil {
		t.Errorf("expected nil, got %#v", got)
	}
}

func TestExtractStackFrameStop(t *testing.T) {
	for _, test := range stopTests {
		if test.testErr(t) {
			test.testFrameString(t, errs.ExtractStackFrameStop(test.err, test.stop))
		}
	}
}

func TestMarshalStack(t *testing.T) {
	for _, test := range marshallTests {
		if test.testErr(t) {
			got := errs.MarshalStack(test.err)
			switch ff := got.(type) {
			case []errs.Frame:
				test.testFrameString(t, ff)
			case []map[string]string:
				test.testFramesMap(t, ff)
			default:
				t.Errorf("expected %#v, got %#v", test.stack, got)
			}
		}
	}
}

func TestFrame_String(t *testing.T) {
	f := errs.Frame{
		File: "file",
		Line: 123,
		Func: "func",
	}
	want := "file:123 func"
	got := f.String()
	if got != want {
		t.Errorf("want `%#v`, got `%#v`", want, got)
	}
}
