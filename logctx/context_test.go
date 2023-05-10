package logctx_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bir/iken/logctx"
)

func TestOp(t *testing.T) {
	tests := []struct {
		name string
		ctx  context.Context
		op   string
	}{
		{"empty", context.Background(), ""},
		{"opName", context.Background(), "opName"},
		{"override", logctx.SetOp(context.Background(), "foo"), "opName2"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := logctx.SetOp(test.ctx, test.op)
			if op := logctx.GetOp(c); op != test.op {
				t.Errorf("GetOp() got = `%v`, want `%v`", op, test.op)
			}
		})
	}

	if op := logctx.GetOp(context.Background()); op != "" {
		t.Errorf("GetOp() got = `%v`, want `%v`", op, "")
	}
}

func TestID(t *testing.T) {
	ctx := context.Background()

	assert.Zero(t, logctx.GetID(ctx))

	ctx = logctx.SetID(ctx, 111)
	assert.Equal(t, int64(111), logctx.GetID(ctx))
}
