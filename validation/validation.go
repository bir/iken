package validation

import (
	"fmt"
	"sort"
	"strings"
)

type Messages []string

func (mm Messages) Error() string {
	return strings.Join(mm, ", ")
}

type Errors map[string]Messages

func Validate(object interface{}) error {
	if v, ok := object.(interface{ Validate() error }); ok {
		return v.Validate()
	}

	return nil
}

// Error returns the error string of Errors.
func (ee Errors) Error() string {
	if len(ee) == 0 {
		return ""
	}

	keys := make([]string, len(ee))

	i := 0

	for key := range ee {
		keys[i] = key
		i++
	}

	sort.Strings(keys)

	var s strings.Builder

	for i, key := range keys {
		if i > 0 {
			s.WriteString("; ")
		}

		_, _ = fmt.Fprintf(&s, "%v: %v", key, (ee)[key].Error())
	}

	s.WriteString(".")

	return s.String()
}

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

func New(field, msg string) error {
	return (&Errors{}).Add(field, msg)
}

const (
	MsgInvalidMaximum = "maximum"
	MsgInvalidMinimum = "minimum"
	MsgMultipleOf     = "multiple"
	MsgMaxLength      = "max length"
	MsgMinLength      = "min length"
	MsgPatten         = "pattern"
	MsgMaxItems       = "max items"
	MsgMinItems       = "min items"
	MsgRequired       = "required"
	MsgEnum           = "enum"
)
