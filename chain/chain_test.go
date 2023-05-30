package chain_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bir/iken/chain"
)

func prefixLetter(letter string) chain.Constructor {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(letter))
			h.ServeHTTP(w, r)
		})
	}
}

func nop(w http.ResponseWriter, r *http.Request) {
}

func TestNew(t *testing.T) {
	c := chain.New(prefixLetter("a"),
		prefixLetter("b"),
		prefixLetter("c"),
	)

	got := testHandler(c)
	want := "abc"
	if got != want {
		t.Errorf("got = `%v`, want `%v`", got, want)
	}

	c2 := c.Append(prefixLetter("d"))

	// Validate c is unmodified
	got = testHandler(c)
	want = "abc"
	if got != want {
		t.Errorf("got = `%v`, want `%v`", got, want)
	}

	got = testHandler(c2)
	want = "abcd"
	if got != want {
		t.Errorf("got = `%v`, want `%v`", got, want)
	}

	c3 := c2.Prepend(prefixLetter("X"))
	got = testHandler(c3)
	want = "Xabcd"
	if got != want {
		t.Errorf("got = `%v`, want `%v`", got, want)
	}
}

func testHandler(c2 chain.Chain) string {
	r := httptest.NewRecorder()
	c2.Handler(http.HandlerFunc(nop)).ServeHTTP(r, nil)

	result := r.Result()
	b, _ := io.ReadAll(result.Body)

	return string(b)
}
