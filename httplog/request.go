package httplog

import (
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"

	"github.com/bir/iken/httputil"
	"github.com/bir/iken/logctx"
)

const (
	Duration            = "duration" // In nanoseconds
	HTTPStatusCode      = "http.status_code"
	HTTPMethod          = "http.method"
	HTTPURLDetailsPath  = "http.url_details.path"
	NetworkBytesWritten = "network.bytes_written"
	Operation           = "op"
	ReqID               = "req_id"
	Request             = "request_payload"
	RequestHeaders      = "http.headers"
	TraceID             = "trace_id"
	UserID              = "usr.id"
	Stack               = "error.stack"
)

func requestLogger(r *http.Request, status, size int, duration time.Duration) {
	op := logctx.GetOp(r.Context())

	l := hlog.FromRequest(r).With().
		// Request
		Str(Operation, op).
		Str(HTTPMethod, r.Method).
		Str(HTTPURLDetailsPath, r.URL.Path).
		Interface(RequestHeaders, httputil.DumpHeader(r)).
		// Response
		Int(HTTPStatusCode, status).
		Int(NetworkBytesWritten, size).
		Dur(Duration, duration).Logger()

	var event *zerolog.Event

	switch {
	case status >= http.StatusInternalServerError:
		event = l.Error()
	case status >= http.StatusBadRequest:
		event = l.Warn()
	default:
		event = l.Info()
	}

	if op != "" {
		event.Msg(op)
	} else {
		event.Msgf("[%s] %s", r.Method, r.URL)
	}
}

// RequestLogger returns a handler that call initializes Op in the context, and logs each request.
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lw := httputil.WrapWriter(w)
		next.ServeHTTP(w, r.WithContext(logctx.SetOp(r.Context(), fmt.Sprintf("[%s] %s", r.Method, r.URL))))
		requestLogger(r, lw.Status(), lw.BytesWritten(), time.Since(start))
	})
}
