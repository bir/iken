package logctx_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/bir/iken/logctx"
)

func TestLogContext(t *testing.T) {
	logBuffer := bytes.NewBuffer(nil)

	id := "42"
	ctx := logctx.SetID(context.Background(), id)
	ctx, cancel := context.WithCancel(logctx.NewSubLoggerContext(ctx, zerolog.New(logBuffer)))

	ctx2 := logctx.NewContextFrom(ctx)
	logctx.AddToContext(ctx, "key", 1)
	logctx.AddMapToContext(ctx, map[string]any{"test": "value", "test2": "value2"})

	logctx.AddStrToContext(ctx2, "key", "value")

	logctx.AddBytesToContext(ctx2, "small", []byte("12345"), 5)
	logctx.AddBytesToContext(ctx2, "large", []byte("123456789"), 5)

	zerolog.Ctx(ctx).Log().Msg("ctx")

	zerolog.Ctx(ctx2).Log().Msg("ctx2")

	assert.Equal(t, `{"key":1,"test":"value","test2":"value2","message":"ctx"}
{"key":"value","small.size":5,"small.body":"12345","large.size":9,"large.body":"12345","large.truncated":true,"large.truncatedSize":5,"message":"ctx2"}
`, logBuffer.String())
	// ensure NewContext does not propagate cancel
	cancel()
	assert.Error(t, ctx.Err())
	assert.NoError(t, ctx2.Err())

	// ensure Context variables are shared
	assert.Equal(t, id, logctx.GetID(ctx))
	assert.Equal(t, id, logctx.GetID(ctx2))
}
