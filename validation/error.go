package validation

import (
	"fmt"
)

type Error struct {
	Message string
	Source  error
}

func (e Error) Error() string {
	if e.Source == nil {
		return e.Message
	}

	if e.Message == "" {
		return e.Source.Error()
	}

	return fmt.Sprintf("%s: %s", e.Message, e.Source)
}

func (e Error) Unwrap() error {
	return e.Source
}

func (e Error) UserError() string {
	if e.Message != "" {
		return e.Message
	}

	if e.Source != nil {
		return e.Source.Error()
	}

	return ""
}
