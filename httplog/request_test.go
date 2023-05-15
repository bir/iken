package httplog

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/bir/iken/logctx"
)

func TestRequestLogger(t *testing.T) {
	MaxRequestBodyLog = 10
	RecoverBasePath = "iken/httplog/"

	tests := []struct {
		name      string
		shouldLog FnShouldLog
		body      string
		next      http.Handler
		want      string
	}{
		{"default logs", nil, "DO NOT LOG ME", http.HandlerFunc(emptyNext), `{"level":"info","http.method":"FOO","http.url_details.path":"/BAR","http.headers":{"FOO":"/BAR HTTP/1.1","Host":"example.com"},"op":"empty","http.status_code":0,"network.bytes_written":0,"duration":0.1,"message":"empty"}
`},
		{"no op", nil, "DO NOT LOG ME", http.HandlerFunc(emptyOp), `{"level":"info","http.method":"FOO","http.url_details.path":"/BAR","http.headers":{"FOO":"/BAR HTTP/1.1","Host":"example.com"},"op":"","http.status_code":0,"network.bytes_written":0,"duration":0.1,"message":"[FOO] /BAR"}
`},
		{"default warn", nil, "DO NOT LOG ME", http.HandlerFunc(statusNext(404)), `{"level":"warn","http.method":"FOO","http.url_details.path":"/BAR","http.headers":{"FOO":"/BAR HTTP/1.1","Host":"example.com"},"op":"[FOO] /BAR","http.status_code":404,"network.bytes_written":11,"duration":0.1,"message":"[FOO] /BAR"}
`},
		{"default err", nil, "DO NOT LOG ME", http.HandlerFunc(statusNext(503)), `{"level":"error","http.method":"FOO","http.url_details.path":"/BAR","http.headers":{"FOO":"/BAR HTTP/1.1","Host":"example.com"},"op":"[FOO] /BAR","http.status_code":503,"network.bytes_written":11,"duration":0.1,"message":"[FOO] /BAR"}
`},
		{"no logs", doLogs(false, false, false), "DO NOT LOG ME", http.HandlerFunc(emptyNext), ""},
		{"all logs", doLogs(true, true, true), "LOG ME", http.HandlerFunc(bodyNext), `{"level":"info","http.method":"FOO","http.url_details.path":"/BAR","http.headers":{"FOO":"/BAR HTTP/1.1","Host":"example.com"},"network.bytes_read":6,"request":"LOG ME","op":"[FOO] /BAR","http.status_code":200,"network.bytes_written":4,"duration":0.1,"response":"TEST","message":"[FOO] /BAR"}
`},
		{"request Body", doLogs(true, true, false), "LOG ME", http.HandlerFunc(bodyNext), `{"level":"info","http.method":"FOO","http.url_details.path":"/BAR","http.headers":{"FOO":"/BAR HTTP/1.1","Host":"example.com"},"network.bytes_read":6,"request":"LOG ME","op":"[FOO] /BAR","http.status_code":200,"network.bytes_written":4,"duration":0.1,"message":"[FOO] /BAR"}
`},
		{"request Body read", doLogs(true, true, false), "LOG ME", http.HandlerFunc(readNext), `{"level":"info","http.method":"FOO","http.url_details.path":"/BAR","http.headers":{"FOO":"/BAR HTTP/1.1","Host":"example.com"},"network.bytes_read":6,"request":"LOG ME","op":"[FOO] /BAR","http.status_code":200,"network.bytes_written":6,"duration":0.1,"message":"[FOO] /BAR"}
`},
		{"response Body", doLogs(true, false, true), "LOG ME", http.HandlerFunc(bodyNext), `{"level":"info","http.method":"FOO","http.url_details.path":"/BAR","http.headers":{"FOO":"/BAR HTTP/1.1","Host":"example.com"},"op":"[FOO] /BAR","http.status_code":200,"network.bytes_written":4,"duration":0.1,"response":"TEST","message":"[FOO] /BAR"}
`},
		{"request Body too big", doLogs(true, true, false), "12345678901", http.HandlerFunc(readNext), `{"level":"info","http.method":"FOO","http.url_details.path":"/BAR","http.headers":{"FOO":"/BAR HTTP/1.1","Host":"example.com"},"network.bytes_read":11,"request":"1234567890","op":"[FOO] /BAR","http.status_code":200,"network.bytes_written":11,"duration":0.1,"message":"[FOO] /BAR"}
`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logOutput := bytes.NewBuffer(nil)

			h := RequestLogger(zerolog.New(logOutput), tt.shouldLog)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("FOO", "/BAR", bytes.NewBufferString(tt.body))

			now = startNow
			h(tt.next).ServeHTTP(w, r)

			got := logOutput.String()

			assert.Equal(t, tt.want, got, "logs")
		})
	}
}

func TestRequestLoggerPanic(t *testing.T) {
	MaxRequestBodyLog = 10
	RecoverBasePath = "iken/httplog/"

	tests := []struct {
		name          string
		shouldLog     FnShouldLog
		body          string
		next          http.Handler
		wantMessage   string
		wantFirstLine string
	}{
		{"panic String", doLogs(true, true, true), "123", readPanic("test"), "test: internal error", "./request_test.go:137 (iken/httplog.TestRequestLoggerPanic.func3)"},
		{"panic Error", doLogs(true, true, true), "123", readPanic(errors.New("test")), "test", "./request_test.go:137 (iken/httplog.TestRequestLoggerPanic.func5)"},
		{"panic other", doLogs(true, true, true), "123", readPanic(1), "internal error", "./request_test.go:137 (iken/httplog.TestRequestLoggerPanic.func7)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logOutput := bytes.NewBuffer(nil)

			h := RequestLogger(zerolog.New(logOutput), tt.shouldLog)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("FOO", "/BAR", bytes.NewBufferString(tt.body))

			now = startNow
			h(tt.next).ServeHTTP(w, r)

			got := logOutput.String()

			result := make(map[string]any)
			err := json.Unmarshal([]byte(got), &result)
			assert.Nil(t, err, "json Unmarshal")

			stack, ok := result["error.stack"].([]any)
			assert.True(t, ok, "error.stack type")

			assert.Equal(t, tt.wantFirstLine, stack[0], "logs")
		})
	}
}

func emptyNext(_ http.ResponseWriter, r *http.Request) {
	now = endNow
	logctx.SetOp(r.Context(), "empty")
}

func emptyOp(_ http.ResponseWriter, r *http.Request) {
	now = endNow
	logctx.SetOp(r.Context(), "")
}

func bodyNext(w http.ResponseWriter, r *http.Request) {
	now = endNow
	_, _ = w.Write([]byte("TEST"))
}

func readNext(w http.ResponseWriter, r *http.Request) {
	now = endNow
	buf := bytes.NewBuffer(nil)
	_, _ = io.Copy(buf, r.Body)
	_, _ = w.Write(buf.Bytes())
}

func readPanic(result any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		now = endNow

		panic(result)
	}
}

func statusNext(status int) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		now = endNow
		w.WriteHeader(status)
		_, _ = w.Write([]byte("TEST STATUS"))
	}
}

func doLogs(logRequest, logRequestBody, logResponseBody bool) func(r *http.Request) (logRequest, logRequestBody, logResponseBody bool) {
	return func(_ *http.Request) (bool, bool, bool) {
		return logRequest, logRequestBody, logResponseBody
	}
}

func startNow() time.Time {
	return time.Date(2023, 1, 1, 1, 1, 1, 0, time.UTC)
}

func endNow() time.Time {
	return time.Date(2023, 1, 1, 1, 1, 1, 100000, time.UTC)
}
