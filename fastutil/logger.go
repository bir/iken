package fastutil

import (
	"fmt"
	"time"

	"github.com/bir/iken/chain"
	"github.com/bir/iken/fastctx"
	"github.com/rs/zerolog"
	"github.com/valyala/fasthttp"
)

var (
	// ApplicationJSON content-type.
	ApplicationJSON = []byte("application/json")
	// TextHTML content-type.
	TextHTML = []byte("text/html")
)

// Notifier is designed to provide standardized hooks for exception managers.
type Notifier interface {
	FastSend(ctx *fasthttp.RequestCtx, msg interface{}, extra ...interface{}) (string, error)
}

func isSuccess(code int) bool {
	return code < 400
}

// MaxBodySize defines the maximum body that will be logged if logRequest/logResponse is true.
var MaxBodySize = 10 * 1024

// RequestLogger returns a handler Constructor that logs all requests to the provider logger.
// If a Notifier is provided, then any request with an error logged to the context will be sent
// to the Notifier.
// logRequest will log the request header.
// logResponse will log the response header.
// includeBody will log the request and/or response body to the log as determined by the appropriate params.
// includeBody has no effect if logRequest and logResponse are both false.
func RequestLogger(l zerolog.Logger, n Notifier, logRequest, logResponse, includeBody bool) chain.Constructor {
	return func(h fasthttp.RequestHandler) fasthttp.RequestHandler {
		return requestLogger(l, n, logRequest, logResponse, includeBody, h)
	}
}

func requestLogger(l zerolog.Logger, n Notifier, logRequest, logResponse, includeBody bool,
	h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		h(ctx)

		var log *zerolog.Event

		code := ctx.Response.StatusCode()
		err := fastctx.GetError(ctx)

		switch {
		case err != nil:
			log = l.Error().Stack().Err(err)

			if n != nil {
				_, _ = n.FastSend(ctx, err)
			}
		case isSuccess(code):
			log = l.Info()
		default:
			log = l.Warn()
		}

		op := fastctx.GetOp(ctx)
		if op == "" {
			op = fmt.Sprintf("%s:%s", ctx.Method(), ctx.URI().RequestURI())
		}

		log = log.Uint64("requestID", ctx.ID()).
			Str("op", op).
			Int("code", code).
			Str("ip", ctx.RemoteAddr().String()).
			Dur("duration", time.Since(ctx.Time()))

		if logRequest {
			log = log.Bytes("header", ctx.Request.Header.Header())

			if includeBody {
				b := ctx.Request.Body()
				if len(b) < MaxBodySize {
					log = log.Bytes("body", b)
				}
			}
		}

		if logResponse {
			log = log.Bytes("responseHeader", ctx.Response.Header.Header())

			if includeBody {
				b := ctx.Response.Body()
				if len(b) < MaxBodySize {
					log = log.Bytes("responseBody", b)
				}
			}
		}

		log.Msg("request")
	}
}
