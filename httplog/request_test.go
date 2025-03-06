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

	"github.com/bir/iken/httputil"
	"github.com/bir/iken/logctx"
)

func TestRequestLogger(t *testing.T) {
	MaxBodyLog = 10
	RecoverBasePath = "iken/httplog/"

	tests := []struct {
		name         string
		shouldLog    FnShouldLog
		body         io.Reader
		addRequestID bool
		next         http.Handler
		want         string
	}{
		{"default logs", nil, bytes.NewBufferString("DO NOT LOG ME"), true, http.HandlerFunc(emptyNext), `{"level":"info","http.method":"FOO","http.url_details.path":"/BAR","request.headers":{"FOO":"/BAR HTTP/1.1","Host":"example.com", "X-Request-Id":"default logs"},"op":"empty","http.status_code":0,"network.bytes_written":0,"duration":0.1,"message":"0 FOO /BAR", "http.request_id":"default logs"}
`},
		{"no op", nil, bytes.NewBufferString("DO NOT LOG ME"), false, http.HandlerFunc(emptyOp), `{"level":"info","http.method":"FOO","http.url_details.path":"/BAR","request.headers":{"FOO":"/BAR HTTP/1.1","Host":"example.com"},"http.status_code":0,"network.bytes_written":0,"duration":0.1,"message":"0 FOO /BAR"}
`},
		{"default warn", nil, bytes.NewBufferString("DO NOT LOG ME"), false, http.HandlerFunc(statusNext(404)), `{"level":"warn","http.method":"FOO","http.url_details.path":"/BAR","request.headers":{"FOO":"/BAR HTTP/1.1","Host":"example.com"},"http.status_code":404,"network.bytes_written":11,"duration":0.1,"message":"404 FOO /BAR"}
`},
		{"default err", nil, bytes.NewBufferString("DO NOT LOG ME"), false, http.HandlerFunc(statusNext(503)), `{"level":"error","http.method":"FOO","http.url_details.path":"/BAR","request.headers":{"FOO":"/BAR HTTP/1.1","Host":"example.com"},"http.status_code":503,"network.bytes_written":11,"duration":0.1,"message":"503 FOO /BAR"}
`},
		{"no logs", doLogs(false, false, false, StatusToLogLevel), bytes.NewBufferString("DO NOT LOG ME"), false, http.HandlerFunc(emptyNext), ""},
		{"all logs", LogAll, bytes.NewBufferString("LOG ME"), false, http.HandlerFunc(bodyNext), `{"level":"info","http.method":"FOO","http.url_details.path":"/BAR","request.headers":{"FOO":"/BAR HTTP/1.1","Host":"example.com"},"network.bytes_read":6,"request.body":"LOG ME","request.size":6,"response.body":"TEST","response.size":4,"http.status_code":200,"network.bytes_written":4,"duration":0.1,"response.body":"TEST","message":"200 FOO /BAR"}
`},
		{"disabled level", doLogs(true, true, true, func(_ *http.Request, _ int) zerolog.Level { return zerolog.Disabled }), bytes.NewBufferString("DO NOT LOG ME"), false, http.HandlerFunc(emptyNext), ""},
		{"all logs warn level", doLogs(true, true, true, func(_ *http.Request, _ int) zerolog.Level { return zerolog.WarnLevel }), bytes.NewBufferString("LOG ME"), false, http.HandlerFunc(bodyNext), `{"level":"warn","http.method":"FOO","http.url_details.path":"/BAR","request.headers":{"FOO":"/BAR HTTP/1.1","Host":"example.com"},"network.bytes_read":6,"request.body":"LOG ME","request.size":6,"response.body":"TEST","response.size":4,"http.status_code":200,"network.bytes_written":4,"duration":0.1,"response.body":"TEST","message":"200 FOO /BAR"}
`},
		{"all logs error level", doLogs(true, true, true, func(_ *http.Request, _ int) zerolog.Level { return zerolog.ErrorLevel }), bytes.NewBufferString("LOG ME"), false, http.HandlerFunc(bodyNext), `{"level":"error","http.method":"FOO","http.url_details.path":"/BAR","request.headers":{"FOO":"/BAR HTTP/1.1","Host":"example.com"},"network.bytes_read":6,"request.body":"LOG ME","request.size":6,"response.body":"TEST","response.size":4,"http.status_code":200,"network.bytes_written":4,"duration":0.1,"response.body":"TEST","message":"200 FOO /BAR"}
`},
		{"request Body", LogRequestBody, bytes.NewBufferString("LOG ME"), false, http.HandlerFunc(bodyNext), `{"level":"info","http.method":"FOO","http.url_details.path":"/BAR","request.headers":{"FOO":"/BAR HTTP/1.1","Host":"example.com"},"network.bytes_read":6,"request.body":"LOG ME","request.size":6,"http.status_code":200,"network.bytes_written":4,"duration":0.1,"message":"200 FOO /BAR"}
`},
		{"request Body read", LogRequestBody, bytes.NewBufferString("LOG ME"), false, http.HandlerFunc(readNext), `{"level":"info","http.method":"FOO","http.url_details.path":"/BAR","request.headers":{"FOO":"/BAR HTTP/1.1","Host":"example.com"},"network.bytes_read":6,"request.body":"LOG ME","request.size":6,"http.status_code":200,"network.bytes_written":6,"duration":0.1,"message":"200 FOO /BAR"}
`},
		{"response Body", doLogs(true, false, true, StatusToLogLevel), bytes.NewBufferString("LOG ME"), false, http.HandlerFunc(bodyNext), `{"level":"info","http.method":"FOO","http.url_details.path":"/BAR","request.headers":{"FOO":"/BAR HTTP/1.1","Host":"example.com"},"http.status_code":200,"network.bytes_written":4,"response.size":4,"duration":0.1,"response.body":"TEST","message":"200 FOO /BAR"}
`},
		{"request Body too big", LogRequestBody, bytes.NewBufferString("12345678901"), false, http.HandlerFunc(readNext), `{"level":"info","http.method":"FOO","http.url_details.path":"/BAR","request.headers":{"FOO":"/BAR HTTP/1.1","Host":"example.com"},"network.bytes_read":11,"request.body":"1234567890","request.size":11,"request.truncated":true,"request.truncatedSize":10,"http.status_code":200,"network.bytes_written":11,"duration":0.1,"message":"200 FOO /BAR"}
`},
		{"error body", LogRequestBody, BadReader{}, false, http.HandlerFunc(readNext), `{"level":"info","http.method":"FOO","http.url_details.path":"/BAR","request.headers":{"FOO":"/BAR HTTP/1.1","Host":"example.com"},"request.body_error":"buf.ReadFrom:BadReader","http.status_code":200,"network.bytes_written":0,"duration":0.1,"message":"200 FOO /BAR"}
`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logOutput := bytes.NewBuffer(nil)

			h := RequestLogger(tt.shouldLog)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("FOO", "/BAR", tt.body)

			if tt.addRequestID {
				r.Header.Set(httputil.RequestIDHeader, tt.name)
			}

			now = startNow
			h(tt.next).ServeHTTP(w, r.WithContext(zerolog.New(logOutput).WithContext(r.Context())))

			got := logOutput.String()

			if len(got) < 1 {
				assert.True(t, len(tt.want) < 1, "got empty data, expected logs")

				return
			}

			result := make(map[string]any)
			err := json.Unmarshal([]byte(got), &result)
			assert.Nil(t, err, "json Unmarshal got")

			want := make(map[string]any)
			err = json.Unmarshal([]byte(tt.want), &want)
			assert.Nil(t, err, "json Unmarshal want")

			assert.Equal(t, want, result, "logs")
		})
	}
}

type BadReader struct{}

func (_ BadReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("BadReader")
}

func emptyNext(_ http.ResponseWriter, r *http.Request) {
	now = endNow
	logctx.AddStrToContext(r.Context(), Operation, "empty")
}

func emptyOp(_ http.ResponseWriter, r *http.Request) {
	now = endNow
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

func statusNext(status int) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		now = endNow
		w.WriteHeader(status)
		_, _ = w.Write([]byte("TEST STATUS"))
	}
}

func doLogs(logRequest, logRequestBody, logResponseBody bool, toLogLevel FnToLogLevel) func(r *http.Request) (logRequest, logRequestBody, logResponseBody bool, toLogLevel FnToLogLevel) {
	return func(_ *http.Request) (bool, bool, bool, FnToLogLevel) {
		return logRequest, logRequestBody, logResponseBody, toLogLevel
	}
}

func startNow() time.Time {
	return time.Date(2023, 1, 1, 1, 1, 1, 0, time.UTC)
}

func endNow() time.Time {
	return time.Date(2023, 1, 1, 1, 1, 1, 100000, time.UTC)
}
