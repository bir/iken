package httputil

import (
	"encoding/json"
	"errors"
	"io"

	"github.com/bir/iken/validation"
)

var ErrMissingBody = validation.Error{Message: "missing body"}

func GetJSONBody(r io.Reader, body any) error {
	if r == nil {
		return ErrMissingBody
	}

	err := json.NewDecoder(r).Decode(body)
	if err == nil {
		return nil
	}

	var (
		validationError  validation.Error
		validationErrors *validation.Errors
	)

	switch {
	case err == io.EOF:
		return ErrMissingBody
	case errors.As(err, &validationError):
		return err //nolint:wrapcheck
	case errors.As(err, &validationErrors):
		return err //nolint:wrapcheck
	default:
		return validation.Error{Source: err}
	}
}
