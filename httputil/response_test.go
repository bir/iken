package httputil

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddHeaders(t *testing.T) {

	w := httptest.NewRecorder()
	AddHeaders(w, http.Header{"foo": []string{"a", "b"}})
	assert.Equal(t, w.Header().Get("foo"), "a")
	assert.Equal(t, w.Header().Get("fOO"), "a")
	assert.Equal(t, w.Header().Get("a"), "")
	assert.Equal(t, w.Header().Values("foo"), []string{"a", "b"})
}

func TestHTMLWrite(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("FOO", "/BAR", nil)

	HTMLWrite(w, r, http.StatusTeapot, http.StatusText(http.StatusTeapot))

	assert.Equal(t, w.Header().Get(ContentType), TextHTML)
	assert.Equal(t, w.Code, http.StatusTeapot)
	assert.Equal(t, w.Body.String(), http.StatusText(http.StatusTeapot))
}

func TestJSONWrite(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest("FOO", "/BAR", nil)

	JSONWrite(w, r, http.StatusTeapot, http.StatusText(http.StatusTeapot))

	assert.Equal(t, w.Header().Get(ContentType), ApplicationJSON)
	assert.Equal(t, w.Code, http.StatusTeapot)
	assert.Equal(t, w.Body.String(), fmt.Sprintf("%q", http.StatusText(http.StatusTeapot)))

	// Fail
	w = httptest.NewRecorder()

	JSONWrite(w, r, 412, badJson{})

	assert.Equal(t, w.Header().Get(ContentType), ApplicationJSON)
	assert.Equal(t, w.Code, 500)
	assert.Equal(t, w.Body.String(), `Internal Server Error
`)

}

type badJson struct {
}

func (badJson) MarshalJSON() ([]byte, error) {
	return nil, errors.New("bad")
}
