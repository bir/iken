package httputil

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/bir/iken/errs"
	"github.com/bir/iken/logctx"
	"github.com/bir/iken/validation"
)

// ErrorHandlerFunc is useful to standardize the exception management of
// requests.
type ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)

// StatusContextCancelled - reported when the context is cancelled.  Most likely caused by lost connections.
const StatusContextCancelled = 499

type ClientValidationError struct {
	Code    int                 `json:"code,omitempty"`
	Message string              `json:"message"`
	Fields  map[string][]string `json:"fields,omitempty"`
}

// ErrorHandler provides some standard handling for errors in an http request
// flow.
//
// Maps:
// json.SyntaxError to "BadRequest", body is the JSON string of the error
// message.
// validation.Errors to "BadRequest", body is the JSON of the error
// object (map of field name to list of errors).
// AuthError to "Forbidden" or "Unauthorized" as defined by the err instance.  In addition
// ErrBasicAuthenticate issues a basic auth challenge using default realm of "Restricted".
// To override handle in your custom error handlers instead.
//
// Unhandled errors are added to the ctx and return "Internal Server Error" with
// the request ID to aid with troubleshooting.
func ErrorHandler(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}

	logctx.AddStrToContext(r.Context(), LogErrorMessage, err.Error())

	if stack := errs.MarshalStack(err); stack != nil {
		logctx.AddToContext(r.Context(), LogStack, stack)
	}

	var (
		jsonErr        *json.SyntaxError
		validationErrs *validation.Errors
		validationErr  validation.Error
	)

	switch {
	case errors.Is(err, context.Canceled):
		http.Error(w, "canceled", StatusContextCancelled)

	case errors.Is(err, ErrForbidden):
		HTTPError(w, http.StatusForbidden)

	case errors.Is(err, ErrBasicAuthenticate):
		w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
		HTTPError(w, http.StatusUnauthorized)

	case errors.Is(err, ErrUnauthorized):
		HTTPError(w, http.StatusUnauthorized)

	case errors.As(err, &jsonErr):
		JSONWrite(w, r, http.StatusBadRequest, fmt.Sprintf("%s at offset %d", jsonErr.Error(), jsonErr.Offset))

	case errors.As(err, &validationErrs):
		JSONWrite(w, r, http.StatusBadRequest,
			ClientValidationError{http.StatusBadRequest, "validation errors", validationErrs.Fields()})

	case errors.As(err, &validationErr):
		JSONWrite(w, r, http.StatusBadRequest,
			ClientValidationError{http.StatusBadRequest, validationErr.UserError(), nil})

	default:
		HTTPInternalServerError(w, r)
	}
}

const (
	// LogErrorMessage is used to report internal errors to the logging service.
	LogErrorMessage = "error.message"
	// LogStack is used to report available error stacks to logging.
	LogStack = "error.stack"

	RequestIDHeader = "X-Request-Id"

	// InternalErrorFormat is the default error message returned for unhandled errors if the request ID is available.
	InternalErrorFormat = "Internal Server Error: Request %q"
)

func HTTPInternalServerError(w http.ResponseWriter, r *http.Request) {
	reqID := r.Header.Get(RequestIDHeader)

	if reqID != "" {
		http.Error(w, fmt.Sprintf(InternalErrorFormat, reqID), http.StatusInternalServerError)
	} else {
		HTTPError(w, http.StatusInternalServerError)
	}
}

func HTTPError(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}
