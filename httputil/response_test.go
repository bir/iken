package httputil

import (
	"bytes"
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

func TestReaderWrite(t *testing.T) {
	rw := httptest.NewRecorder()
	r := httptest.NewRequest("FOO", "/BAR", nil)

	ReaderWrite(rw, r, TextPlain, http.StatusTeapot, bytes.NewReader([]byte(http.StatusText(http.StatusTeapot))))

	result := rw.Result()
	b, _ := io.ReadAll(result.Body)

	assert.Equal(t, TextPlain, result.Header.Get(ContentType))
	assert.Equal(t, http.StatusTeapot, result.StatusCode)
	assert.Equal(t, http.StatusText(http.StatusTeapot), string(b))
}

func TestReaderWriteError(t *testing.T) {
	rw := httptest.NewRecorder()
	r := httptest.NewRequest("FOO", "/BAR", nil)

	ReaderWrite(rw, r, TextPlain, http.StatusTeapot, errReader{})

	result := rw.Result()
	b, _ := io.ReadAll(result.Body)

	assert.Equal(t, TextPlain, result.Header.Get(ContentType))
	assert.Equal(t, http.StatusTeapot, result.StatusCode)
	assert.Equal(t, http.StatusText(http.StatusInternalServerError)+"\n", string(b))
}

type errorResponseWriter struct {
	*httptest.ResponseRecorder
}

var errorBody = []byte("error me")

func (e *errorResponseWriter) Write(b []byte) (int, error) {
	if bytes.Equal(b, errorBody) {
		// Return an error when Write is called
		return 0, errors.New("mock error")
	}

	return e.ResponseRecorder.Write(b)
}

func TestWriteError(t *testing.T) {
	rw := httptest.NewRecorder()
	r := httptest.NewRequest("FOO", "/BAR", nil)

	Write(&errorResponseWriter{ResponseRecorder: rw}, r, TextPlain, http.StatusTeapot, errorBody)

	result := rw.Result()
	b, _ := io.ReadAll(result.Body)

	assert.Equal(t, TextPlain, result.Header.Get(ContentType))
	assert.Equal(t, http.StatusTeapot, result.StatusCode)
	assert.Equal(t, http.StatusText(http.StatusInternalServerError)+"\n", string(b))
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

type badJson struct{}

func (badJson) MarshalJSON() ([]byte, error) {
	return nil, errors.New("bad")
}
