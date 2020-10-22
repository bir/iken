package fastutil

import (
	"encoding/json"
	"errors"
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
	case AuthError:
		if errors.Is(e, ErrForbidden) {
			if err := JSONWrite(ctx, fasthttp.StatusForbidden, e); err != nil {
				panic(err)
			}

			return
		}

		if errors.Is(e, ErrBasicAuthenticate) {
			ctx.Response.Header.Set("WWW-Authenticate", "Basic realm=Restricted")
			ctx.Error(fasthttp.StatusMessage(fasthttp.StatusUnauthorized), fasthttp.StatusUnauthorized)

			return
		}

		if err := JSONWrite(ctx, fasthttp.StatusUnauthorized, e); err != nil {
			panic(err)
		}

		return
	}

	fastctx.SetError(ctx, err)
	ctx.Error(fmt.Sprintf(InternalErrorFormat, ctx.ID()), fasthttp.StatusInternalServerError)
}
