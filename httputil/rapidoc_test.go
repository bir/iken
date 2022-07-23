package httputil_test

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/bir/iken/httputil"
	"github.com/stretchr/testify/assert"
)

func TestRapiDoc(t *testing.T) {
	rw := httptest.NewRecorder()

	h := httputil.RapiDoc(httputil.RapiDocOpts{})(nil)
	h.ServeHTTP(rw, httptest.NewRequest("GET", "http://test/docs", nil))
	result := rw.Result()
	b, _ := io.ReadAll(result.Body)

	assert.Equal(t, 200, result.StatusCode, "status")
	assert.Equal(t, 354, len(b), "body size")

	h.ServeHTTP(rw, httptest.NewRequest("GET", "http://test/notDocs", nil))
	result = rw.Result()
	b, _ = io.ReadAll(result.Body)

	assert.Equal(t, 200, result.StatusCode, "status")
	assert.Equal(t, 0, len(b), "body size")
}
