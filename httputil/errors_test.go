package httputil_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bir/iken/errs"
	"github.com/bir/iken/httputil"
	"github.com/bir/iken/validation"
)

func TestErrorHandler(t *testing.T) {
	nop := "ignore"
	tests := []struct {
		name   string
		ctx    context.Context
		err    error
		status int
		body   string
	}{
		{"nil error", context.Background(), nil, 200, ""},
		{"unknown error", context.Background(), errs.WithStack("unknown error", 0), 500, "Internal Server Error: 0\n"},
		{"json error", context.Background(), json.Unmarshal([]byte("bad json"), &nop), 400, `"invalid character 'b' looking for beginning of value"` + "\n"},
		{"validation error", context.Background(), (&validation.Errors{}).Add("name", "bad"), 400, `{"name":["bad"]}` + "\n"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := httputil.SetID(test.ctx, 0)
			r := &http.Request{}
			r = r.WithContext(c)
			w := httptest.NewRecorder()
			httputil.ErrorHandler(w, r, test.err)
			if w.Code != test.status {
				t.Errorf("ErrorHandler status got = `%v`, wantLog `%v`", w.Code, test.status)
			}

			if w.Body.String() != test.body {
				t.Errorf("ErrorHandler body got = `%v`, wantLog `%v`", w.Body.String(), test.body)
			}
		})
	}
}
