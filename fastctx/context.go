package fastctx

import (
	"context"

	"github.com/valyala/fasthttp"
)

type ctxKey uint8

const (
	opCtx         = "iken.op"
	errCtx        = "iken.err"
	stdCtx ctxKey = iota
)

// SetError store unhandled exceptions in the context to be processed by
// the request logger.
func SetError(ctx *fasthttp.RequestCtx, err error) {
	ctx.SetUserValue(errCtx, err)
}

// GetError returns the error logged to the context, otherwise nil.
func GetError(ctx *fasthttp.RequestCtx) error {
	e, ok := ctx.UserValue(errCtx).(error)
	if ok {
		return e
	}

	return nil
}

// SetOp stores a label for the current request.  Useful for developer friendly log messages.
func SetOp(ctx *fasthttp.RequestCtx, op string) {
	ctx.SetUserValue(opCtx, op)
}

// GetOp returns the operation label logged to the context, otherwise "".
func GetOp(ctx *fasthttp.RequestCtx) string {
	s, ok := ctx.UserValue(opCtx).(string)
	if ok {
		return s
	}

	return ""
}

// SetToStd stores a fasthttp RequestCtx in a stdlib context.Context.
func SetToStd(std context.Context, ctx *fasthttp.RequestCtx) context.Context {
	return context.WithValue(std, stdCtx, ctx)
}

// GetFromStd retrieves a fasthttp RequestCtx from a stdlib context.Context.
func GetFromStd(std context.Context) *fasthttp.RequestCtx {
	if std == nil {
		return nil
	}

	fast, ok := std.(*fasthttp.RequestCtx)
	if ok {
		return fast
	}

	ctx := std.Value(stdCtx)
	if ctx == nil {
		return nil
	}

	return ctx.(*fasthttp.RequestCtx)
}

// GetRequestID returns the requestID from the context, otherwise 0.
func GetRequestID(std context.Context) uint64 {
	ctx := GetFromStd(std)
	if ctx == nil {
		return 0
	}

	return ctx.ID()
}
