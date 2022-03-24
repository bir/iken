package validation

import (
	"fmt"
	"sort"
	"strings"

	"github.com/bir/iken/strutil"
)

// Messages are the validation failures for a given field.
type Messages []string //nolint: errname

func (mm Messages) Error() string {
	return strings.Join(mm, ", ")
}

// Errors maps fields to the list of validation failures.
type Errors map[string]Messages //nolint: errname

// Error returns the error string of Errors.
func (ee Errors) Error() string {
	if len(ee) == 0 {
		return ""
	}

	return strutil.Join(ee.Keys(), "", "; ", ".", func(key string) string {
		return fmt.Sprintf("%v: %v", key, ee[key].Error())
	})
}

// Add appends the field and msg to the current list of errors.  Add will initialize the Errors
// object if it is not initialized.
func (ee *Errors) Add(field, msg string) *Errors {
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

func (ee Errors) Keys() []string {
	keys := make([]string, len(ee))
	i := 0

	for key := range ee {
		keys[i] = key
		i++
	}

	sort.Strings(keys)

	return keys
}

// New returns a single validation error for the field with msg.
func New(field, msg string) error {
	return (&Errors{}).Add(field, msg)
}
