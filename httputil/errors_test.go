package httputil_test

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bir/iken/errs"
	"github.com/bir/iken/httputil"
	"github.com/bir/iken/logctx"
	"github.com/bir/iken/validation"
)

func TestErrorHandler(t *testing.T) {
	nop := "ignore"
	tests := []struct {
		name      string
		ctx       context.Context
		err       error
		requestID string
		status    int
		body      string
	}{
		{"nil error", context.Background(), nil, "", 200, ""},
		{"unknown error", context.Background(), errs.WithStack("unknown error", 0), "", 500, "Internal Server Error\n"},
		{"unknown error w/request ID", context.Background(), errs.WithStack("unknown error", 0), "FOO", 500, "Internal Server Error: Request \"FOO\"\n"},
		{"json error", context.Background(), json.Unmarshal([]byte("bad json"), &nop), "", 400, `"invalid character 'b' looking for beginning of value at offset 1"`},
		{"validation error", context.Background(), validation.New("name", "bad"), "", 400, `{"name":["bad"]}`},
		{"auth error unauthorized", context.Background(), httputil.ErrUnauthorized, "", 401, `Unauthorized` + "\n"},
		{"auth error forbidden", context.Background(), httputil.ErrForbidden, "", 403, `Forbidden` + "\n"},
		{"auth error basic", context.Background(), httputil.ErrBasicAuthenticate, "", 401, `Unauthorized` + "\n"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := logctx.SetID(test.ctx, "test")
			r := httptest.NewRequest("FOO", "/BAR", nil)
			r = r.WithContext(c)
			w := httptest.NewRecorder()

			if test.requestID != "" {
				r.Header.Set(httputil.RequestIDHeader, test.requestID)
			}

			httputil.ErrorHandler(w, r, test.err)
			assert.Equal(t, test.status, w.Code, "status")

			assert.Equal(t, test.body, w.Body.String(), "body")

		})
	}
}
