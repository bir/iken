package fastctx

import "github.com/valyala/fasthttp"

// Field names used for logging ctx params.
var (
	ConnID     = "connID"
	RequestNum = "requestNum"
	RequestID  = "requestID"
	IP         = "ip"
)

// DebugMap small helper util to extract standard info from a RequestCtx into a hashmap for logging.
//
// Currently ConnID, ConnRequestNum, ID, RemoteAddr are logged.
func DebugMap(ctx *fasthttp.RequestCtx) map[string]interface{} {
	if ctx == nil {
		return nil
	}

	r := make(map[string]interface{})
	r[ConnID] = ctx.ConnID()
	r[RequestNum] = ctx.ConnRequestNum()
	r[RequestID] = ctx.ID()
	r[IP] = ctx.RemoteAddr().String()

	return r
}
