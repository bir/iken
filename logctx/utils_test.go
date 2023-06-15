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
	original := bytes.NewBuffer(nil)

	//used to verify sub-logger creation doesn't throw away the whole context
	id := "42"
	ctx := logctx.SetID(context.Background(), id)
	assert.Equal(t, id, logctx.GetID(ctx))

	ctx = logctx.NewSubLoggerContext(ctx, zerolog.New(original))

	assert.Equal(t, id, logctx.GetID(ctx), "")

	ctx2 := logctx.NewContextFrom(ctx)

	assert.Equal(t, id, logctx.GetID(ctx2))

	logctx.AddToContext(ctx, "key", 1)
	logctx.AddMapToContext(ctx, map[string]interface{}{"test": "value", "test2": "value2"})

	logctx.AddStrToContext(ctx2, "key", "value")

	zerolog.Ctx(ctx).Log().Msg("ctx")

	zerolog.Ctx(ctx2).Log().Msg("ctx2")

	assert.Equal(t, original.String(), `{"key":1,"test":"value","test2":"value2","message":"ctx"}
{"key":"value","message":"ctx2"}
`)
}
