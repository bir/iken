package validation

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/bir/iken/strutil"
)

// Messages are the validation failures for a given field.
type Messages []error //nolint: errname

func (mm Messages) Error() string {
	return Join(mm, ", ")
}

// Errors maps fields to the list of validation failures.
type Errors map[string]Messages //nolint: errname

// Error returns the error string of Errors.
func (ee *Errors) Error() string {
	if len(*ee) == 0 {
		return ""
	}

	return strutil.Join(ee.Keys(), "", "; ", ".", func(key string) string {
		return fmt.Sprintf("%v: %v", key, (*ee)[key].Error())
	})
}

// Add appends the field and msg to the current list of errors.  Add will initialize the Errors
// object if it is not initialized.
func (ee *Errors) Add(field string, msg error) *Errors {
	if *ee == nil {
		*ee = Errors{}
	}

	fE, ok := (*ee)[field]
	if !ok || len(fE) == 0 {
		fE = Messages{msg}
	} else {
		fE = append(fE, msg)
	}

	(*ee)[field] = fE

	return ee
}

// GetErr allows you to use a nil Errors object and return directly.  If there are no validation errors it returns nil.
func (ee *Errors) GetErr() error {
	if *ee == nil {
		return nil
	}

	return ee
}

func (ee *Errors) Keys() []string {
	keys := make([]string, len(*ee))
	i := 0

	for key := range *ee {
		keys[i] = key
		i++
	}

	sort.Strings(keys)

	return keys
}

func (ee *Errors) Fields() map[string][]string {
	out := make(map[string][]string)

	for k, errs := range *ee {
		for _, e := range errs {
			out[k] = append(out[k], getError(e))
		}
	}

	return out
}

// New returns a single validation error for the field with msg.
func New(field string, msg string) error {
	return (&Errors{}).Add(field, errors.New(msg)) //nolint: goerr113
}

// NewError returns a single validation error for the field with the embedded error.
func NewError(field string, msg error) error {
	return (&Errors{}).Add(field, msg)
}

type UserError interface {
	error
	UserError() string
}

func getError(err error) string {
	var u UserError
	if errors.As(err, &u) {
		return u.UserError()
	}

	return err.Error()
}

func Join(elems []error, sep string) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return elems[0].Error()
	}

	n := len(sep) * (len(elems) - 1)

	for i := 0; i < len(elems); i++ {
		n += len(elems[i].Error())
	}

	var b strings.Builder

	b.Grow(n)
	b.WriteString(elems[0].Error())

	for _, s := range elems[1:] {
		b.WriteString(sep)
		b.WriteString(s.Error())
	}

	return b.String()
}
