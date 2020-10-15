package fastutil_test

import (
	"encoding/json"
	"testing"

	"github.com/bir/iken/errs"
	"github.com/bir/iken/fastctx"
	"github.com/bir/iken/fastutil"
	"github.com/bir/iken/validation"
	"github.com/valyala/fasthttp"
)

func TestErrorHandler(t *testing.T) {
	nop := "ignore"
	tests := []struct {
		name   string
		ctx    *fasthttp.RequestCtx
		err    error
		status int
		body   string
		logged bool
	}{
		{"nil error", &fasthttp.RequestCtx{}, nil, 200, "", false},
		{"unknown error", &fasthttp.RequestCtx{}, errs.WithStack("unknown error", 0), 500, "Internal Server Error: 0", true},
		{"json error", &fasthttp.RequestCtx{}, json.Unmarshal([]byte("bad json"), &nop), 400, `"invalid character 'b' looking for beginning of value"` + "\n", false},
		{"validation error", &fasthttp.RequestCtx{}, (&validation.Errors{}).Add("name", "bad"), 400, `{"name":["bad"]}` + "\n", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fastutil.ErrorHandler(test.ctx, test.err)
			if test.ctx.Response.StatusCode() != test.status {
				t.Errorf("ErrorHandler status got = `%v`, wantLog `%v`", test.ctx.Response.StatusCode(), test.status)
			}
			if string(test.ctx.Response.Body()) != test.body {
				t.Errorf("ErrorHandler body got = `%v`, wantLog `%v`", string(test.ctx.Response.Body()), test.body)
			}
			e := fastctx.GetError(test.ctx)
			if e != nil && !test.logged {
				t.Errorf("ErrorHandler logged err got = `%v`, wantLog `%v`", e, test.err)
			} else if test.logged && e != test.err {
				t.Errorf("ErrorHandler logged err got = `%v`, wantLog `%v`", e, test.err)
			}
		})
	}
}
