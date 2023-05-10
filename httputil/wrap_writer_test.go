package httputil_test

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bir/iken/httputil"
)

func TestNewWrapResponse(t *testing.T) {

	tests := []struct {
		name       string
		bytes      string
		writer     http.ResponseWriter
		tee        bool
		header     int
		wantStatus int
		wantBytes  int
	}{
		{"basic success", "success", httptest.NewRecorder(), false, 0, 200, 7},
		{"header success", "success", httptest.NewRecorder(), false, 201, 201, 7},
		{"header success", "", httptest.NewRecorder(), false, 201, 201, 0},
		{"no usage", "", httptest.NewRecorder(), false, 0, 0, 0},
		{"Fancy", "1234", NewFancy(), false, 400, 400, 4},
		{"Basic Tee", "123456", httptest.NewRecorder(), true, 400, 400, 6},
		{"Fancy Tee", "12345", NewFancy(), true, 400, 400, 5},
		{"Fancy No Header", "12345", NewFancy(), false, 0, 200, 5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httputil.WrapWriter(tt.writer)

			b := bytes.Buffer{}
			if tt.tee {
				w.Tee(&b)
			}

			if tt.header > 0 {
				w.WriteHeader(tt.header)
			}

			if len(tt.bytes) > 0 {
				if rf, ok := w.(io.ReaderFrom); ok {
					r := bytes.NewBufferString(tt.bytes)
					_, err := rf.ReadFrom(r)
					assert.Nil(t, err, "ReadFrom Error")
				} else {
					_, _ = w.Write([]byte(tt.bytes))
				}

				if f, ok := w.(http.Flusher); ok {
					f.Flush()
				}
			}

			assert.Equal(t, tt.wantStatus, w.Status(), "Status")
			assert.Equal(t, tt.wantBytes, w.BytesWritten(), "Bytes")
			assert.Equal(t, w.Unwrap(), tt.writer, "unwrap returned unknown")
			if tt.tee {
				assert.Equal(t, tt.wantBytes, b.Len(), "Tee Bytes")
			}
		})
	}

	w := httputil.WrapWriter(NewFancy())

	hj, ok := w.(http.Hijacker)

	assert.True(t, ok, "Hijacker")
	_, _, err := hj.Hijack()
	assert.Error(t, err, "Hijack Not Implemented")
}

func NewFancy() fancyWriter {
	return fancyWriter{ResponseRecorder: httptest.NewRecorder()}
}

type fancyWriter struct {
	*httptest.ResponseRecorder
}

func (_ fancyWriter) Flush() {}

func (_ fancyWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return nil, nil, fmt.Errorf("not implemented")
}

func (w fancyWriter) ReadFrom(r io.Reader) (n int64, err error) {
	return io.Copy(w.Body, r)
}
