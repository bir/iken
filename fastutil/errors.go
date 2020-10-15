package fastutil

import (
	"encoding/json"
	"fmt"

	"github.com/bir/iken/fastctx"
	"github.com/bir/iken/validation"
	"github.com/valyala/fasthttp"
)

// InternalErrorFormat is the default error message returned for unhandled errors.
const InternalErrorFormat = "Internal Server Error: %d"

// ErrorHandlerFunc is useful to standardize the exception management of
// requests.
type ErrorHandlerFunc func(*fasthttp.RequestCtx, error)

// ErrorHandler provides some standard handling for errors in an http request
// flow.
//
// Maps json.SyntaxError to "BadRequest", body is the JSON string of the error
// message.
// Maps validation.Errors to "BadRequest", body is the JSON of the error
// object (map of field name to list of errors).
// Unhandled errors are added to the ctx and return "Internal Server Error" with
// the request ID to aid with troubleshooting.
func ErrorHandler(ctx *fasthttp.RequestCtx, err error) {
	if err == nil {
		return
	}

	switch e := err.(type) {
	case *json.SyntaxError:
		if err := JSONWrite(ctx, fasthttp.StatusBadRequest, e.Error()); err != nil {
			panic(err)
		}

		return

	case *validation.Errors:
		if err := JSONWrite(ctx, fasthttp.StatusBadRequest, e); err != nil {
			panic(err)
		}

		return
	}

	fastctx.SetError(ctx, err)
	ctx.Error(fmt.Sprintf(InternalErrorFormat, ctx.ID()), fasthttp.StatusInternalServerError)
}
