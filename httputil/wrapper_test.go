package httputil_test

import (
	"net/http/httptest"
	"testing"

	"github.com/bir/iken/httputil"
	"github.com/stretchr/testify/assert"
)

func TestNewWrapResponse(t *testing.T) {

	tests := []struct {
		name       string
		bytes      string
		header     int
		wantStatus int
		wantBytes  int
	}{
		{"basic success", "success", 0, 200, 7},
		{"header success", "success", 201, 201, 7},
		{"header success", "", 201, 201, 0},
		{"no usage", "", 0, 0, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rw := httptest.NewRecorder()

			w := httputil.NewWrapResponse(rw)

			if tt.header > 0 {
				w.WriteHeader(tt.header)
			}

			if len(tt.bytes) > 0 {
				w.Write([]byte(tt.bytes))
			}

			assert.Equalf(t, tt.wantStatus, w.Status(), "got %v want %v", tt.wantStatus, w.Status())
			assert.Equalf(t, tt.wantBytes, w.Bytes(), "got %v want %v", tt.wantBytes, w.Bytes())
		})
	}
}
