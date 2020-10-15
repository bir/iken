package notify_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/bir/iken/errs"
	"github.com/bir/iken/notify"
	"github.com/rs/zerolog"
)

func crashAndRecover(msg interface{}, n notify.Notifier) {
	// Simulates outer final panic handler
	defer func() {
		if r := recover(); r != nil {
			if fmt.Sprintf("%v", r) != fmt.Sprintf("%v", msg) {
				panic(r)
			}
		}
	}()

	defer notify.Monitor(n)

	panic(msg)
}

type LogMessage struct {
	Level   string
	Message string
	Error   string
	Stack   []interface{}
	Extra   interface{}
	Ctx     struct {
		ConnId     int
		Ip         string
		RequestID  uint64
		RequestNum int
	}
}

func (l LogMessage) WithExtra(extra interface{}) LogMessage {
	l.Extra = extra
	return l
}

func BasicLog(level, message, error string) LogMessage {
	return LogMessage{
		Level:   level,
		Message: message,
		Error:   error,
		Stack:   nil,
		Extra:   nil,
		Ctx: struct {
			ConnId     int
			Ip         string
			RequestID  uint64
			RequestNum int
		}{},
	}
}

func TestMonitor(t *testing.T) {
	var logBuf bytes.Buffer
	l := zerolog.New(&logBuf)
	errWithStack := errs.WithStack("foo", 0)

	n := notify.NewZerolog(l)
	tests := []struct {
		notifier  notify.Notifier
		msg       interface{}
		want      LogMessage
		wantStack bool
	}{
		{n, "simple", BasicLog("error", "notify", "simple"), false},
		{n, errWithStack, BasicLog("error", "notify", errWithStack.Error()), true},
		{nil, "empty", LogMessage{}, false},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test.msg), func(t *testing.T) {
			logBuf.Reset()
			crashAndRecover(test.msg, test.notifier)

			if test.notifier == nil {
				if logBuf.Len() > 0 {
					t.Errorf("expected no output, got `%v`", logBuf.String())
				}
				return
			}
			got := LogMessage{}
			err := json.Unmarshal(logBuf.Bytes(), &got)
			if err != nil {
				t.Errorf("json err `%v`", err)
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

			if test.wantStack && len(got.Stack) < 1 {
				t.Error("wantStack")
			}

		})
	}
}

func ExampleMonitor() {
	var logBuf bytes.Buffer
	l := zerolog.New(&logBuf)
	n := notify.NewZerolog(l)

	_, _ = n.Send("test")

	func() {
		defer func() {
			if r := recover(); r != nil {
				// Squashed in example
			}
		}()
		defer notify.Monitor(n)
		panic("alert")
	}()

	// Dumping log to demonstrate the notify capture
	fmt.Println(logBuf.String())

	// Output: {"level":"warn","msg":"test","message":"notify"}
	// {"level":"error","error":"alert","message":"notify"}
}
