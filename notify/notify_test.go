package notify_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"unsafe"

	"github.com/bir/iken/errs"
	"github.com/bir/iken/notify"
	"github.com/rs/zerolog"
	"github.com/valyala/fasthttp"
)

// We use this to set internal values in the RequestCtx for testing logging
func mockCtx(connID, requestNum uint64) *fasthttp.RequestCtx {
	ctx := &fasthttp.RequestCtx{}
	pointerVal := reflect.ValueOf(ctx)
	val := reflect.Indirect(pointerVal)

	member := val.FieldByName("connRequestNum")
	ptrToY := unsafe.Pointer(member.UnsafeAddr())
	realPtrToY := (*uint64)(ptrToY)
	*realPtrToY = requestNum

	member = val.FieldByName("connID")
	ptrToY = unsafe.Pointer(member.UnsafeAddr())
	realPtrToY = (*uint64)(ptrToY)
	*realPtrToY = connID

	return ctx
}

func TestNewZerolog(t *testing.T) {
	var logBuf bytes.Buffer
	l := zerolog.New(&logBuf)

	errWithStack := errs.WithStack("foo", 0)

	ctx := mockCtx(1, 2)

	testExtra := []interface{}{"test"}
	n := notify.NewZerolog(l)
	tests := []struct {
		name      string
		ctx       *fasthttp.RequestCtx
		msg       interface{}
		want      LogMessage
		wantStack bool
		wantCtx   bool
	}{
		{"simple", nil, "simple", BasicLog("warn", "notify", ""), false, false},
		{"errWithStack", nil, errWithStack, BasicLog("error", "notify", errWithStack.Error()), true, false},
		{"withExtra", nil, "withExtra", BasicLog("warn", "notify", "").WithExtra(testExtra), false, false},
		{"withCtx", ctx, errWithStack, BasicLog("error", "notify", errWithStack.Error()), true, true},
		{"empty", nil, nil, BasicLog("warn", "notify", ""), false, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			logBuf.Reset()
			var err error
			if test.ctx != nil {
				if test.want.Extra != nil {
					_, err = n.FastSend(test.ctx, test.msg, test.want.Extra)
				} else {
					_, err = n.FastSend(test.ctx, test.msg)
				}
			} else {
				if test.want.Extra != nil {
					_, err = n.Send(test.msg, test.want.Extra)
				} else {
					_, err = n.Send(test.msg)
				}
			}
			if err != nil {
				t.Errorf("Send err `%v`", err)
				return
			}

			got := LogMessage{}

			fmt.Println(logBuf.String())
			err = json.Unmarshal(logBuf.Bytes(), &got)
			if err != nil {
				if test.msg != nil {
					t.Errorf("json err `%v`", err)
				}
				return
			}

			if got.Level != test.want.Level {
				t.Errorf("gotLevel `%v`, want `%v`", got.Level, test.want.Level)
			}

			if got.Message != test.want.Message {
				t.Errorf("gotMessage `%v`, want `%v`", got.Message, test.want.Message)
			}

			if got.Error != test.want.Error {
				t.Errorf("gotError `%v`, want `%v`", got.Error, test.want.Error)
			}

			if test.want.Extra != nil {
				g := fmt.Sprintf("%v", got.Extra)
				w := fmt.Sprintf("%v", test.want.Extra)
				if g != "["+w+"]" {
					t.Errorf("Extra get: `%v`\nwant:`%v`", g, w)
				}
			} else if got.Extra != nil {
				t.Errorf("got unexpected Extra: %v", got.Extra)
			}

			if test.wantStack && len(got.Stack) < 1 {
				t.Error("wantStack")
			}

			if test.wantCtx && got.Ctx.RequestNum == 0 {
				t.Error("wantCtx")
			}
		})
	}
}

func TestNewDebug(t *testing.T) {
	var logBuf bytes.Buffer

	errWithStack := errs.WithStack("foo", 0)

	ctx := mockCtx(1, 2)

	testExtra := []interface{}{"test"}
	n := notify.NewDebug(&logBuf)
	tests := []struct {
		name     string
		ctx      *fasthttp.RequestCtx
		msg      interface{}
		extra    interface{}
		want     string
		hasStack bool
	}{
		{"simple", nil, "simple", nil, "NOTIFY:\nsimple\n", false},
		{"errWithStack", nil, errWithStack, nil, "", true},
		{"withExtra", nil, "withExtra", testExtra, "NOTIFY:\nwithExtra\nExtra:\n[[test]]\n", false},
		{"withCtx", ctx, "withCtx", nil, "NOTIFY:\nwithCtx\nContext:\n#0000000100000002 - 0.0.0.0:0<->0.0.0.0:0 - GET http:///\n", false},
		{"empty", nil, nil, nil, "", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			logBuf.Reset()
			var err error
			if test.ctx != nil {
				if test.extra != nil {
					_, err = n.FastSend(test.ctx, test.msg, test.extra)
				} else {
					_, err = n.FastSend(test.ctx, test.msg)
				}
			} else {
				if test.extra != nil {
					_, err = n.Send(test.msg, test.extra)
				} else {
					_, err = n.Send(test.msg)
				}
			}

			if err != nil {
				t.Errorf("Send err `%v`", err)
				return
			}

			got := logBuf.String()

			if got != test.want {
				if test.hasStack {
					if !strings.Contains(got, "STACK:") {
						t.Errorf("got `%v`, want `%v`", got, test.want)
					}
					return
				}

				t.Errorf("got `%v`, want `%v`", got, test.want)
			}
		})
	}
}
