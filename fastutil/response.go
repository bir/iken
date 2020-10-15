package fastutil

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
)

// JSONWrite is a simple helper utility to return the json encoded obj with appropriate content-type and code.
func JSONWrite(ctx *fasthttp.RequestCtx, code int, obj interface{}) error {
	ctx.Response.Header.SetContentTypeBytes(ApplicationJSON)
	ctx.Response.SetStatusCode(code)

	return json.NewEncoder(ctx).Encode(obj)
}
