package notify

import (
	"reflect"
	"strconv"

	"github.com/bir/iken/errs"
	"github.com/bir/iken/fastctx"
	"github.com/honeybadger-io/honeybadger-go"
	"github.com/valyala/fasthttp"
)

type hbNotifier struct {
	c *honeybadger.Client
}

// NewHoneybadger returns a Notifier given a honeybadger.Client.
//
// This wrapper will automatically marshall stack traces and attach RequestCtx info if available.
func NewHoneybadger(c *honeybadger.Client) Notifier {
	return &hbNotifier{c: c}
}

func (h hbNotifier) Send(msg interface{}, extra ...interface{}) (string, error) {
	return h.FastSend(nil, msg, extra...)
}

func (h hbNotifier) FastSend(ctx *fasthttp.RequestCtx, msg interface{}, extra ...interface{}) (string, error) {
	e, ok := msg.(error)
	if ok {
		err := honeybadger.Error{
			Message: e.Error(),
			Class:   reflect.TypeOf(errs.Cause(e)).String(),
			Stack:   mapStackTracerHB(e),
		}
		msg = err
	}

	if ctx != nil {
		extra = append(extra, honeybadger.Context(fastctx.DebugMap(ctx)))
	}

	return h.c.Notify(msg, extra...)
}

func (h hbNotifier) Flush() {
	h.c.Flush()
}

func mapStackTracerHB(err error) []*honeybadger.Frame {
	st := errs.ExtractStackFrame(err)
	out := make([]*honeybadger.Frame, 0, len(st))

	for _, frame := range st {
		out = append(out, &honeybadger.Frame{
			Number: strconv.Itoa(frame.Line),
			File:   frame.File,
			Method: frame.Func,
		})
	}

	return out
}
