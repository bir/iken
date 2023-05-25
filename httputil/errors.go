package httputil

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/bir/iken/logctx"
	"github.com/bir/iken/validation"
)

// ErrorHandlerFunc is useful to standardize the exception management of
// requests.
type ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)

// StatusContextCancelled - reported when the context is cancelled.  Most likely caused by lost connections.
const StatusContextCancelled = 499

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

	logctx.AddToContext(r.Context(), LogErrorMessage, err)

	var (
		jsonErr       *json.SyntaxError
		validationErr *validation.Errors
	)

	switch {
	case errors.Is(err, context.Canceled):
		http.Error(w, err.Error(), StatusContextCancelled)
	case errors.Is(err, ErrForbidden):
		http.Error(w, err.Error(), http.StatusForbidden)
	case errors.Is(err, ErrBasicAuthenticate):
		w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	case errors.Is(err, ErrUnauthorized):
		http.Error(w, err.Error(), http.StatusUnauthorized)
	case errors.As(err, &jsonErr):
		JSONWrite(w, r, http.StatusBadRequest, fmt.Sprintf("%s at offset %d", jsonErr.Error(), jsonErr.Offset))
	case errors.As(err, &validationErr):
		JSONWrite(w, r, http.StatusBadRequest, validationErr)
	default:
		HTTPInternalServerError(w, r)
	}
}

const (
	// LogErrorMessage is used to report internal errors to the logging service.
	LogErrorMessage = "error.message"

	RequestIDHeader = "X-Request-Id"

	// InternalErrorFormat is the default error message returned for unhandled errors if the request ID is available.
	InternalErrorFormat = "Internal Server Error: Request %q"
)

func HTTPInternalServerError(w http.ResponseWriter, r *http.Request) {
	reqID := r.Header.Get(RequestIDHeader)

	if reqID != "" {
		http.Error(w, fmt.Sprintf(InternalErrorFormat, reqID), http.StatusInternalServerError)
	} else {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
