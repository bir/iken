package fastctx_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"unsafe"

	"github.com/bir/iken/errs"
	"github.com/bir/iken/fastctx"
	"github.com/valyala/fasthttp"
)

func TestError(t *testing.T) {
	tests := []struct {
		name string
		ctx  *fasthttp.RequestCtx
		err  error
	}{
		{"nil error", &fasthttp.RequestCtx{}, nil},
		{"wrapped error", &fasthttp.RequestCtx{}, errs.WithStack("wrapped error", 0)},
		{"basic error", &fasthttp.RequestCtx{}, errors.New("basic error")},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fastctx.SetError(test.ctx, test.err)
			if err := fastctx.GetError(test.ctx); err != test.err {
				t.Errorf("GetError() got = `%v`, want `%v`", err, test.err)
			}
		})
	}
}
func TestOp(t *testing.T) {
	tests := []struct {
		name string
		ctx  *fasthttp.RequestCtx
		op   string
	}{
		{"empty", &fasthttp.RequestCtx{}, ""},
		{"opName", &fasthttp.RequestCtx{}, "opName"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fastctx.SetOp(test.ctx, test.op)
			if op := fastctx.GetOp(test.ctx); op != test.op {
				t.Errorf("GetOp() got = `%v`, want `%v`", op, test.op)
			}
		})
	}

	ctx := &fasthttp.RequestCtx{}
	if op := fastctx.GetOp(ctx); op != "" {
		t.Errorf("GetOp() got = `%v`, want `%v`", op, "")
	}
}

func TestStd(t *testing.T) {
	tests := []struct {
		name string
		ctx  *fasthttp.RequestCtx
	}{
		{"empty", nil},
		{"default", &fasthttp.RequestCtx{}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			std := context.TODO()
			if got := fastctx.GetFromStd(fastctx.SetToStd(std, test.ctx)); got != test.ctx {
				t.Errorf("TestStd got = `%v`, want `%v`", got, test.ctx)
			}
		})
	}

	if got := fastctx.GetFromStd(nil); got != nil {
		t.Errorf("TestStd got = `%v`, want `%v`", got, nil)
	}

	if got := fastctx.GetFromStd(context.TODO()); got != nil {
		t.Errorf("TestStd got = `%v`, want `%v`", got, nil)
	}
}

func setRequestID(ctx *fasthttp.RequestCtx) {
	pointerVal := reflect.ValueOf(ctx)
	val := reflect.Indirect(pointerVal)

	member := val.FieldByName("connRequestNum")
	ptrToY := unsafe.Pointer(member.UnsafeAddr())
	realPtrToY := (*uint64)(ptrToY)
	*realPtrToY = 121
}

func TestGetRequestID(t *testing.T) {
	fastCtx := &fasthttp.RequestCtx{}
	setRequestID(fastCtx)
	ctx := fastctx.SetToStd(context.TODO(), fastCtx)

	tests := []struct {
		name string
		ctx  context.Context
		want uint64
	}{
		{"empty", nil, 0},
		{"uninitialized", &fasthttp.RequestCtx{}, 0},
		{"fastCtx", fastCtx, 121},
		{"default", ctx, 121},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := fastctx.GetRequestID(test.ctx); got != test.want {
				t.Errorf("got = `%v`, want `%v`", got, test.want)
			}
		})
	}
}
