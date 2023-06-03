package httputil_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/bir/iken/httputil"
)

func TestRapiDoc(t *testing.T) {
	rw := httptest.NewRecorder()

	h := http.HandlerFunc(httputil.RapiDoc(httputil.RapiDocOpts{}))
	h.ServeHTTP(rw, nil)
	result := rw.Result()
	b, _ := io.ReadAll(result.Body)

	assert.Equal(t, 200, result.StatusCode, "status")
	assert.Equal(t, 354, len(b), "body size")
}
