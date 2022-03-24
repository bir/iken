package httputil_test

import (
	"context"
	"testing"

	"github.com/bir/iken/httputil"
	"github.com/stretchr/testify/assert"
)

func TestOp(t *testing.T) {
	tests := []struct {
		name string
		ctx  context.Context
		op   string
	}{
		{"empty", context.Background(), ""},
		{"opName", context.Background(), "opName"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := httputil.SetOp(test.ctx, test.op)
			if op := httputil.GetOp(c); op != test.op {
				t.Errorf("GetOp() got = `%v`, want `%v`", op, test.op)
			}
		})
	}

	if op := httputil.GetOp(context.Background()); op != "" {
		t.Errorf("GetOp() got = `%v`, want `%v`", op, "")
	}
}

func TestID(t *testing.T) {
	ctx := context.Background()

	assert.Zero(t, httputil.GetID(ctx))

	ctx = httputil.SetID(ctx, 111)
	assert.Equal(t, int64(111), httputil.GetID(ctx))
}
