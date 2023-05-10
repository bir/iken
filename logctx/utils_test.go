package logctx_test

import (
	"bytes"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/bir/iken/logctx"
)

func TestLogContext(t *testing.T) {
	original := bytes.NewBuffer(nil)
	ctx := logctx.NewSubLoggerContext(zerolog.New(original))
	ctx2 := logctx.NewContextFrom(ctx)
	logctx.AddToContext(ctx, "key", 1)
	logctx.AddMapToContext(ctx, map[string]interface{}{"test": "value", "test2": "value2"})

	logctx.AddStrToContext(ctx2, "key", "value")

	zerolog.Ctx(ctx).Log().Msg("ctx")

	zerolog.Ctx(ctx2).Log().Msg("ctx2")

	assert.Equal(t, original.String(), `{"key":1,"test":"value","test2":"value2","message":"ctx"}
{"key":"value","message":"ctx2"}
`)
}
