package fastsafe

import "github.com/valyala/fasthttp"

// NoSniff adds the X-Content-Type-Options header.
//
// X-Content-Type-Options: nosniff - tells browsers to not to sniff the content type of the response.
// See https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/X-Content-Type-Options.
// You MUST correctly set your content-type in responses, otherwise your web pages will fail to process.
func NoSniff(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set(fasthttp.HeaderXContentTypeOptions, "nosniff")
}
