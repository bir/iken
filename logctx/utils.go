package logctx

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
)

// NewContextFrom clones the current context.  Used to branch execution (go routines).
func NewContextFrom(ctx context.Context) context.Context {
	return NewSubLoggerContext(ctx, *zerolog.Ctx(ctx))
}

// NewSubLoggerContext creates a new logger with an empty context.
func NewSubLoggerContext(ctx context.Context, log zerolog.Logger) context.Context {
	return log.With().Logger().WithContext(WithoutCancel(ctx))
}

// AddStrToContext adds the key/value to the log context.
func AddStrToContext(ctx context.Context, key, value string) {
	zerolog.Ctx(ctx).UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Str(key, value)
	})
}

// AddToContext adds the key/value to the log context.
func AddToContext(ctx context.Context, key string, value interface{}) {
	zerolog.Ctx(ctx).UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Interface(key, value)
	})
}

// AddMapToContext adds the map of key/values to the log context.
func AddMapToContext(ctx context.Context, fields map[string]interface{}) {
	zerolog.Ctx(ctx).UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Fields(fields)
	})
}

// WithoutCancel is a port from go1.21.  Placeholder until generally available in stdlib.
func WithoutCancel(parent context.Context) context.Context {
	if parent == nil {
		panic("cannot create context from nil parent")
	}

	return withoutCancelCtx{parent}
}

//nolint:containedctx
type withoutCancelCtx struct {
	ctx context.Context
}

//nolint:nonamedreturns
func (withoutCancelCtx) Deadline() (deadline time.Time, ok bool) {
	return
}

func (withoutCancelCtx) Done() <-chan struct{} {
	return nil
}

func (withoutCancelCtx) Err() error {
	return nil
}

func (c withoutCancelCtx) Value(key any) any {
	return c.ctx.Value(key)
}

func (c withoutCancelCtx) String() string {
	return fmt.Sprintf("%s", c.ctx) + ".WithoutCancel"
}
