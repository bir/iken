package httputil

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddHeaders(t *testing.T) {
	rw := httptest.NewRecorder()

	AddHeaders(rw, http.Header{"foo": []string{"a", "b"}})

	result := rw.Result()
	h := result.Header

	assert.Equal(t, "a", h.Get("foo"))
	assert.Equal(t, "a", h.Get("fOO"))
	assert.Equal(t, "", h.Get("a"))
	assert.Equal(t, []string{"a", "b"}, h.Values("foo"))
}

func TestHTMLWrite(t *testing.T) {
	rw := httptest.NewRecorder()
	r := httptest.NewRequest("FOO", "/BAR", nil)

	HTMLWrite(rw, r, http.StatusTeapot, http.StatusText(http.StatusTeapot))

	result := rw.Result()
	b, _ := io.ReadAll(result.Body)

	assert.Equal(t, TextHTML, result.Header.Get(ContentType))
	assert.Equal(t, http.StatusTeapot, result.StatusCode)
	assert.Equal(t, http.StatusText(http.StatusTeapot), string(b))
}

func TestJSONWrite(t *testing.T) {
	rw := httptest.NewRecorder()
	r := httptest.NewRequest("FOO", "/BAR", nil)

	JSONWrite(rw, r, http.StatusTeapot, http.StatusText(http.StatusTeapot))

	result := rw.Result()
	b, _ := io.ReadAll(result.Body)

	assert.Equal(t, ApplicationJSON, result.Header.Get(ContentType))
	assert.Equal(t, http.StatusTeapot, result.StatusCode)
	assert.Equal(t, fmt.Sprintf("%q", http.StatusText(http.StatusTeapot)), string(b))

	// Fail
	rw = httptest.NewRecorder()

	JSONWrite(rw, r, 412, badJson{})

	result = rw.Result()
	b, _ = io.ReadAll(result.Body)

	assert.Equal(t, TextPlain, result.Header.Get(ContentType))
	assert.Equal(t, http.StatusInternalServerError, result.StatusCode)
	assert.Equal(t, fmt.Sprintf("%s\n", http.StatusText(http.StatusInternalServerError)), string(b))

}

type badJson struct {
}

func (badJson) MarshalJSON() ([]byte, error) {
	return nil, errors.New("bad")
}
