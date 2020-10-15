package fastutil

import (
	"fmt"

	"github.com/bir/iken/errs"
	"github.com/bir/iken/fastctx"
	"github.com/valyala/fasthttp"
)

// PanicHandler saves the PanicErr to the context for logging/alerting.
func PanicHandler(ctx *fasthttp.RequestCtx, panicErr interface{}) {
	// Skips: runtime.Callers/WithStack/PanicHandler/fasthttp.recv/gopanic
	fastctx.SetError(ctx, errs.WithStack(panicErr, 5))
	// Never disclose cause of unknown errors, provide request ID for support
	ctx.Error(fmt.Sprintf(InternalErrorFormat, ctx.ID()), fasthttp.StatusInternalServerError)
}
