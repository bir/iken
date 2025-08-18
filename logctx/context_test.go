package logctx_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bir/iken/logctx"
)

func TestID(t *testing.T) {
	ctx := context.Background()

	assert.Empty(t, logctx.GetID(ctx))

	ctx = logctx.SetID(ctx, "123")
	assert.Equal(t, "123", logctx.GetID(ctx))
}

func TestOperation(t *testing.T) {
	ctx := context.Background()

	assert.Empty(t, logctx.GetOperation(ctx))

	ctx = logctx.SetOperation(ctx, "getTest")
	assert.Equal(t, "getTest", logctx.GetOperation(ctx))
}
