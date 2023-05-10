package logctx

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// NewContextFrom clones the current context.  Used to branch execution (go routines).
func NewContextFrom(ctx context.Context) context.Context {
	return NewSubLoggerContext(*log.Ctx(ctx))
}

// NewSubLoggerContext creates a new logger with an empty context.
func NewSubLoggerContext(log zerolog.Logger) context.Context {
	l := log.With().Logger()

	return l.WithContext(context.Background())
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
