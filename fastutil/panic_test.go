package fastutil_test

import (
	"testing"

	"github.com/bir/iken/errs"
	"github.com/bir/iken/fastutil"
	"github.com/valyala/fasthttp"
)

func TestPanicHandler(t *testing.T) {
	tests := []struct {
		name   string
		ctx    *fasthttp.RequestCtx
		err    error
		status int
		body   string
	}{
		{"nil error", &fasthttp.RequestCtx{}, nil, 500, "Internal Server Error: 0"},
		{"panic error", &fasthttp.RequestCtx{}, errs.WithStack("panic error", 0), 500, "Internal Server Error: 0"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fastutil.PanicHandler(test.ctx, test.err)
			if test.ctx.Response.StatusCode() != test.status {
				t.Errorf("ErrorHandler status got = `%v`, wantLog `%v`", test.ctx.Response.StatusCode(), test.status)
			}
			if string(test.ctx.Response.Body()) != test.body {
				t.Errorf("ErrorHandler body got = `%v`, wantLog `%v`", string(test.ctx.Response.Body()), test.body)
			}
		})
	}
}
