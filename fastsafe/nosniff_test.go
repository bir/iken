package fastsafe_test

import (
	"testing"

	"github.com/bir/iken/fastsafe"
	"github.com/valyala/fasthttp"
)

func TestNoSniff(t *testing.T) {
	ctx := &fasthttp.RequestCtx{}

	h := ctx.Response.Header.Peek(fasthttp.HeaderXContentTypeOptions)
	if len(h) > 0 {
		t.Errorf("expected no header by default, got %s", h)
	}

	fastsafe.NoSniff(ctx)

	h = ctx.Response.Header.Peek(fasthttp.HeaderXContentTypeOptions)
	if string(h) != "nosniff" {
		t.Errorf("expected nosniff, got %s", h)
	}
}
