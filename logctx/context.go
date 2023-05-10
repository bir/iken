package logctx

import (
	"context"
)

type ContextKey string

const (
	opCtx ContextKey = "iken.op"
	opID  ContextKey = "request_id"
)

// SetOp stores a label for the current request.  Useful for developer friendly log messages.
// This is mutable within the context.
func SetOp(ctx context.Context, op string) context.Context {
	o := getOp(ctx)
	if o != nil {
		*o = op

		return ctx
	}

	v := op

	return context.WithValue(ctx, opCtx, &v)
}

// GetOp returns the operation label logged to the context, otherwise "".
func GetOp(ctx context.Context) string {
	s := getOp(ctx)
	if s != nil {
		return *s
	}

	return ""
}

// getOp returns the operation label logged to the context, otherwise "".
func getOp(ctx context.Context) *string {
	s, ok := ctx.Value(opCtx).(*string)
	if ok {
		return s
	}

	return nil
}

// SetID sets the request ID logged to the context.
func SetID(ctx context.Context, id int64) context.Context {
	return context.WithValue(ctx, opID, id)
}

// GetID returns the request ID logged to the context, otherwise 0.
func GetID(ctx context.Context) int64 {
	s, ok := ctx.Value(opID).(int64)
	if ok {
		return s
	}

	return 0
}
