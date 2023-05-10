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

	next := func(rw http.ResponseWriter, r *http.Request) {}

	h := httputil.RapiDoc(httputil.RapiDocOpts{})(http.HandlerFunc(next))
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
