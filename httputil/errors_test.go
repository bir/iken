package httputil_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/bir/iken/errs"
	"github.com/bir/iken/httputil"
	"github.com/bir/iken/logctx"
	"github.com/bir/iken/validation"
)

type errorLog struct {
	Msg string `json:"error.message"`
}

func TestErrorHandler(t *testing.T) {
	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	nop := "ignore"
	tests := []struct {
		name       string
		ctx        context.Context
		err        error
		requestID  string
		status     int
		body       string
		logMessage string
	}{
		{"nil error", context.Background(), nil, "", 200, "", ""},
		{"not found error", context.Background(), errs.WithStack(httputil.ErrNotFound, 0), "", 404, "Not Found\n", "not found"},
		{"unknown error", context.Background(), errs.WithStack("unknown error", 0), "", 500, "Internal Server Error\n", "unknown error"},
		{"unknown error w/request ID", context.Background(), errs.WithStack("unknown error", 0), "FOO", 500, "Internal Server Error: Request \"FOO\"\n", "unknown error"},
		{"validation errors", context.Background(), validation.New("name", "bad"), "", 400, `{"code":400,"message":"validation errors","fields":{"name":["bad"]}}`, "name: bad."},
		{"validation errors public", context.Background(), validation.New("name", validation.Error{Message: "public message", Source: errors.New("private error")}), "", 400, `{"code":400,"message":"validation errors","fields":{"name":["public message"]}}`, "name: public message: private error."},
		{"validation errors json", context.Background(), validation.New("name", validation.Error{Message: "json error", Source: json.Unmarshal([]byte("bad json"), &nop)}), "", 400, `{"code":400,"message":"validation errors","fields":{"name":["json error"]}}`, "name: json error: invalid character 'b' looking for beginning of value."},
		{"validation error", context.Background(), validation.Error{Message: "public message", Source: errors.New("private error")}, "", 400, `{"code":400,"message":"public message"}`, "public message: private error"},
		{"validation message only", context.Background(), validation.Error{Message: "bad"}, "", 400, `{"code":400,"message":"bad"}`, "bad"},
		{"validation error only", context.Background(), validation.Error{Source: errors.New("test")}, "", 400, `{"code":400,"message":"test"}`, "test"},
		{"auth error unauthorized", context.Background(), httputil.ErrUnauthorized, "", 401, `Unauthorized` + "\n", "Unauthorized"},
		{"auth error forbidden", context.Background(), httputil.ErrForbidden, "", 403, `Forbidden` + "\n", "Forbidden"},
		{"auth error basic", context.Background(), httputil.ErrBasicAuthenticate, "", 401, `Unauthorized` + "\n", "ErrWWWAuthenticate"},
		{"nested", context.Background(), fmt.Errorf("wrap:%w", fmt.Errorf("wrap2:%w", httputil.ErrBasicAuthenticate)), "", 401, `Unauthorized` + "\n", "wrap:wrap2:ErrWWWAuthenticate"},
		{"canceled", canceledCtx, canceledCtx.Err(), "FOO", 499, "canceled\n", "context canceled"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			logOutput := bytes.NewBuffer(nil)

			c := zerolog.New(logOutput).WithContext(logctx.SetID(test.ctx, "test"))
			r := httptest.NewRequest("FOO", "/BAR", nil)
			r = r.WithContext(c)
			w := httptest.NewRecorder()

			if test.requestID != "" {
				r.Header.Set(httputil.RequestIDHeader, test.requestID)
			}

			httputil.ErrorHandler(w, r, test.err)

			zerolog.Ctx(c).Log().Msg("test")

			result := w.Result()
			b, _ := io.ReadAll(result.Body)

			assert.Equal(t, test.status, result.StatusCode, "status")
			assert.Equal(t, test.body, string(b), "body")

			var log errorLog
			err := json.Unmarshal(logOutput.Bytes(), &log)
			assert.Nil(t, err)
			assert.Equal(t, test.logMessage, log.Msg, logOutput.String())
		})
	}
}
