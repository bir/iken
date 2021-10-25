package httputil

import (
	"net/http"
)

type ResponseWrapper struct {
	http.ResponseWriter
	status int
	bytes  int
}

// NewWrapResponse wraps http.ResponseWriter and tracks status and bytes written.
func NewWrapResponse(w http.ResponseWriter) *ResponseWrapper {
	return &ResponseWrapper{w, 0, 0}
}

func (rw *ResponseWrapper) WriteHeader(code int) {
	if rw.status == 0 {
		rw.status = code
		rw.ResponseWriter.WriteHeader(code)
	}
}

func (rw *ResponseWrapper) Write(buf []byte) (int, error) {
	rw.WriteHeader(http.StatusOK)
	n, err := rw.ResponseWriter.Write(buf)
	rw.bytes += n

	return n, err // nolint
}

func (rw *ResponseWrapper) Status() int {
	return rw.status
}

func (rw *ResponseWrapper) Bytes() int {
	return rw.bytes
}
