package logctx

import (
	"context"
)

type ContextKey string

const (
	opID      ContextKey = "request_id"
	opMessage ContextKey = "request_message"
)

// SetID sets the request ID logged to the context.
func SetID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, opID, id)
}

// GetID returns the request ID logged to the context, otherwise 0.
func GetID(ctx context.Context) string {
	s, ok := ctx.Value(opID).(string)
	if ok {
		return s
	}

	return ""
}
