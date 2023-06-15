package logctx

import (
	"context"

	"github.com/rs/zerolog"
)

// NewContextFrom returns a child context with a sub-logger attached.
// If no logger is associated with the given ctx as the parent logger,
// DefaultContextLogger is used if not nil, otherwise a disabled logger is used.
func NewContextFrom(ctx context.Context) context.Context {
	return NewSubLoggerContext(ctx, *zerolog.Ctx(ctx))
}

// NewSubLoggerContext returns a child context with a sub-logger attached.
func NewSubLoggerContext(ctx context.Context, log zerolog.Logger) context.Context {
	return log.With().Logger().WithContext(ctx)
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
