package logctx

import (
	"context"
)

type ContextKey string

const (
	operationKey ContextKey = "operation"
	opID         ContextKey = "request_id"
	opMessage    ContextKey = "request_message"
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

// SetOperation sets the operation logged to the context.
func SetOperation(ctx context.Context, operation string) context.Context {
	return context.WithValue(ctx, operationKey, operation)
}

// GetOperation returns the operation logged to the context, otherwise an empty string.
func GetOperation(ctx context.Context) string {
	s, ok := ctx.Value(operationKey).(string)
	if ok {
		return s
	}

	return ""
}
