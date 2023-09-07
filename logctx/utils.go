package logctx

import (
	"context"

	"github.com/rs/zerolog"
)

// NewContextFrom returns a child context without cancel and a sub-logger attached.
// If no logger is associated with the given ctx as the parent logger DefaultContextLogger is used if not nil,
// otherwise a disabled logger is used.
func NewContextFrom(ctx context.Context) context.Context {
	return NewSubLoggerContext(ctx, *zerolog.Ctx(ctx))
}

// NewSubLoggerContext creates a new logger with an empty context.
func NewSubLoggerContext(ctx context.Context, log zerolog.Logger) context.Context {
	return log.With().Logger().WithContext(context.WithoutCancel(ctx))
}

// AddStrToContext adds the key/value to the log context.
func AddStrToContext(ctx context.Context, key, value string) {
	zerolog.Ctx(ctx).UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Str(key, value)
	})
}

// AddToContext adds the key/value to the log context.
func AddToContext(ctx context.Context, key string, value any) {
	zerolog.Ctx(ctx).UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Interface(key, value)
	})
}

// AddMapToContext adds the map of key/values to the log context.
func AddMapToContext(ctx context.Context, fields map[string]any) {
	zerolog.Ctx(ctx).UpdateContext(func(c zerolog.Context) zerolog.Context {
		return c.Fields(fields)
	})
}

// AddBytesToContext adds the key/value to the log context.
func AddBytesToContext(ctx context.Context, key string, value []byte, maxSize uint32) {
	zerolog.Ctx(ctx).UpdateContext(func(c zerolog.Context) zerolog.Context {
		return AddBytes(c, key, value, maxSize)
	})
}

// AddBytes adds the key/value (truncated by maxSize) to the log context.
func AddBytes(ctx zerolog.Context, key string, value []byte, maxSize uint32) zerolog.Context {
	size := len(value)

	ctx = ctx.Int(key+".size", size)

	if size > int(maxSize) {
		ctx = ctx.Bytes(key+".body", value[:maxSize])
		ctx = ctx.Bool(key+".truncated", true)
		ctx = ctx.Uint32(key+".truncatedSize", maxSize)
	} else {
		ctx = ctx.Bytes(key+".body", value)
	}

	return ctx
}
