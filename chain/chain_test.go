package chain_test

import (
	"testing"

	"github.com/bir/iken/chain"
	"github.com/valyala/fasthttp"
)

func prefixLetter(letter string) chain.Constructor {
	return func(h fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			_, _ = ctx.WriteString(letter)
			h(ctx)
		}
	}
}

func nop(*fasthttp.RequestCtx) {
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
	ctx := &fasthttp.RequestCtx{}
	c2.Handler(nop)(ctx)
	return string(ctx.Response.Body())
}
