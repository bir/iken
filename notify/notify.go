package notify

import (
	"fmt"
	"io"

	"github.com/bir/iken/errs"
	"github.com/bir/iken/fastctx"
	"github.com/rs/zerolog"
	"github.com/valyala/fasthttp"
)

// Notifier is used mostly for exception tracking services, it creates a simplified interface for any of the providers.
type Notifier interface {
	// Send packages params for notification.  If msg is an error, stack will be attached to the notice.
	Send(msg interface{}, extra ...interface{}) (string, error)
	// FastSend packages params for notification.  If msg is an error, stack will be attached to the notice.
	FastSend(ctx *fasthttp.RequestCtx, msg interface{}, extra ...interface{}) (string, error)
	// Flush is used to clear any send buffer (useful on shutdown or panic).
	Flush()
}

type zerologNotifier struct {
	l zerolog.Logger
}

// NewZerolog creates a notifier that logs a "notify" message to zerolog.
func NewZerolog(l zerolog.Logger) Notifier {
	return &zerologNotifier{l: l}
}

func (z zerologNotifier) Send(msg interface{}, extra ...interface{}) (string, error) {
	return z.FastSend(nil, msg, extra...)
}

func (z zerologNotifier) FastSend(ctx *fasthttp.RequestCtx, msg interface{}, extra ...interface{}) (string, error) {
	if msg == nil {
		return "", nil
	}

	var log *zerolog.Event
	switch e := msg.(type) {
	case error:
		log = z.l.Err(e)

		ss := errs.ExtractStackFrameStop(e, "(*Router).Handler")
		if len(ss) > 0 {
			log = log.Interface("stack", ss)
		}
	default:
		log = z.l.Warn().Interface("msg", msg)
	}

	if extra != nil {
		log = log.Interface("extra", extra)
	}

	if ctx != nil {
		log = log.Interface("ctx", fastctx.DebugMap(ctx))
	}

	log.Msg("notify")

	return "", nil
}

func (zerologNotifier) Flush() {
	// Nop for zerologNotifier
}

type debug struct {
	w io.Writer
}

// NewDebug creates a notifier designed for local testing, all Notifications are written to the supplied io.Writer.
func NewDebug(w io.Writer) Notifier {
	return &debug{w: w}
}

func (d debug) Send(msg interface{}, extra ...interface{}) (string, error) {
	return d.FastSend(nil, msg, extra...)
}

func (d debug) FastSend(ctx *fasthttp.RequestCtx, msg interface{}, extra ...interface{}) (string, error) {
	if msg == nil {
		return "", nil
	}

	_, err := fmt.Fprintf(d.w, "NOTIFY:\n%+v\n", msg)
	if err != nil {
		return "", err //nolint:wrapcheck // false positive
	}

	e, ok := msg.(error)
	if ok {
		ss := errs.ExtractStackFrameStop(e, "(*Router).Handler")
		if len(ss) > 0 {
			_, err := fmt.Fprintf(d.w, "STACK:\n%+v\n", ss)
			if err != nil {
				return "", err //nolint:wrapcheck // false positive
			}
		}
	}

	if extra != nil {
		_, err := fmt.Fprintf(d.w, "Extra:\n%+v\n", extra)
		if err != nil {
			return "", err //nolint:wrapcheck // false positive
		}
	}

	if ctx != nil {
		_, err := fmt.Fprintf(d.w, "Context:\n%+v\n", ctx)
		if err != nil {
			return "", err //nolint:wrapcheck // false positive
		}
	}

	return "", nil
}

func (debug) Flush() {
	// Nop for zerologNotifier
}
