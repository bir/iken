package httputil

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/bir/iken/validation"
)

// InternalErrorFormat is the default error message returned for unhandled errors.
const InternalErrorFormat = "Internal Server Error: %d"

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

	switch e := err.(type) { //nolint:errorlint,varnamelen // false positive
	case *json.SyntaxError:
		if err := JSONWrite(w, http.StatusBadRequest, e.Error()); err != nil {
			panic(err)
		}

		return
	case *validation.Errors:
		if err := JSONWrite(w, http.StatusBadRequest, e); err != nil {
			panic(err)
		}

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

	http.Error(w, fmt.Sprintf(InternalErrorFormat, GetID(r.Context())), http.StatusInternalServerError)
}
