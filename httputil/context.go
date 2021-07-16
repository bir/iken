package httputil

import (
	"context"
)

const (
	opCtx = "iken.op"
	opID  = "request_id"
)

// SetOp stores a label for the current request.  Useful for developer friendly log messages.
func SetOp(ctx context.Context, op string) context.Context {
	return context.WithValue(ctx, opCtx, op)
}

// GetOp returns the operation label logged to the context, otherwise "".
func GetOp(ctx context.Context) string {
	s, ok := ctx.Value(opCtx).(string)
	if ok {
		return s
	}

	return ""
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
