package httputil

import (
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

	switch e := err.(type) { //nolint:errorlint,varnamelen // false positive
	case *json.SyntaxError:
		JSONWrite(w, r, http.StatusBadRequest, fmt.Sprintf("%s at offset %d", e.Error(), e.Offset))

		return
	case *validation.Errors:
		JSONWrite(w, r, http.StatusBadRequest, e)

		return
	case AuthError:
		if errors.Is(e, ErrForbidden) {
			http.Error(w, e.Error(), http.StatusForbidden)

			return
		}

		msg := e.Error()

		if errors.Is(e, ErrBasicAuthenticate) {
			w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")

			msg = "Unauthorized"
		}

		http.Error(w, msg, http.StatusUnauthorized)

		return
	}

	HTTPInternalServerError(w, r)
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
