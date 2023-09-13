package errs_test

import (
	"testing"

	"github.com/pkg/errors"

	"github.com/bir/iken/errs"
)

type CauseTest struct {
	err   error
	want  string
	isNil bool
}

type NilErr struct{}

func (t NilErr) Cause() error {
	return nil
}

func (t NilErr) Error() string {
	return ""
}

func TestCause(t *testing.T) {
	err1 := errors.New("1")
	nilErr := NilErr{}

	tests := []CauseTest{
		{
			// Shallow
			err1, "1", false,
		}, {
			// Deep
			errors.Wrap(errors.Wrap(errors.Wrap(err1, "X"), "Y"), "Z"), "1", false,
		}, {
			// Nil Cause
			nilErr, "", false,
		}, {
			// Stack
			errs.WithStack("stacked", 0), "stacked", false,
		}, {
			// Nil
			nil, "", true,
		},
	}

	for _, test := range tests {
		e := errs.RootCause(test.err)
		if e == nil {
			if test.isNil {
				continue
			}
			t.Errorf("expected %#v, got nil", test.want)
			continue
		}
		got := e.Error()
		if got != test.want {
			t.Errorf("expected %#v, got %#v", test.want, got)
		}
	}
}
