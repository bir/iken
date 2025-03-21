package httputil

import (
	"encoding/json"
	"errors"
	"io"
	"slices"

	"github.com/bir/iken/validation"
)

var ErrMissingBody = validation.Error{Message: "missing body"}

var nullBody = []byte("null")

func GetJSONBody(r io.Reader, body any) error {
	if r == nil {
		return ErrMissingBody
	}

	data, err := io.ReadAll(r)
	if err != nil {
		return validation.Error{Source: err}
	}

	// json.Unmarshal doesn't call custom unmarshalers if the body is null
	if len(data) == 0 || slices.Equal(data, nullBody) {
		return ErrMissingBody
	}

	err = json.Unmarshal(data, &body)
	if err == nil {
		return nil
	}

	var (
		validationError  validation.Error
		validationErrors *validation.Errors
	)

	switch {
	case errors.As(err, &validationError):
		return err //nolint:wrapcheck
	case errors.As(err, &validationErrors):
		return err //nolint:wrapcheck
	default:
		return validation.Error{Source: err}
	}
}
