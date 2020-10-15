package fastutil

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/bir/iken/fastctx"
	"github.com/bir/iken/notify"
	"github.com/google/go-cmp/cmp"
	"github.com/rs/zerolog"
	"github.com/valyala/fasthttp"
)

func TestRequestLogger(t *testing.T) {

	nop := func(ctx *fasthttp.RequestCtx) {
	}
	failure := func(ctx *fasthttp.RequestCtx) {
		ctx.Error("foo", fasthttp.StatusInternalServerError)
	}
	success := func(ctx *fasthttp.RequestCtx) {
		ctx.SuccessString("TestContent", "TestBody")
	}
	var logBuf bytes.Buffer
	l := zerolog.New(&logBuf)

	var notifyBuf bytes.Buffer
	n := notify.NewDebug(&notifyBuf)

	msgOk := LogMessage{
		Level:     "info",
		RequestID: 0,
		Op:        "GET:/",
		Code:      200,
		IP:        "0.0.0.0:0",
		Message:   "request",
	}
	msgWarn := msgOk
	msgWarn.Level = "warn"
	msgWarn.Code = 500

	msgErr := msgWarn
	msgErr.Level = "error"
	testErr := errors.New("bad")
	errNotify := `NOTIFY:
bad
Context:
#0000000000000000 - 0.0.0.0:0<->0.0.0.0:0 - GET http:///
`

	emptyContent := LogMessageWithContent{}
	logContent := LogMessageWithContent{
		Header:         "GET / HTTP/1.1\r\n\r\n",
		Body:           "",
		ResponseHeader: "HTTP/1.1 200 OK\r\nDate: Fri, 02 Oct 2020 19:08:46 GMT\r\n\r\n",
		ResponseBody:   "",
	}
	logContentErr := LogMessageWithContent{
		Header:         "GET / HTTP/1.1\r\n\r\n",
		Body:           "",
		ResponseHeader: "HTTP/1.1 500 Internal Server Error\r\nDate: Fri, 02 Oct 2020 19:08:46 GMT\r\n\r\n",
		ResponseBody:   "",
	}
	logBody := logContent
	logBody.ResponseBody = "TestBody"
	logBodyErr := logContentErr
	logBodyErr.ResponseBody = "foo"

	tests := []struct {
		name        string
		logRequest  bool
		logResponse bool
		includeBody bool
		requestErr  error
		h           fasthttp.RequestHandler
		wantLog     LogMessage
		wantNotify  string
		wantContent LogMessageWithContent
	}{
		// No Request/Response
		{"NOP", false, false, false, nil, nop, msgOk, "", emptyContent},
		{"Failure", false, false, false, nil, failure, msgWarn, "", emptyContent},
		{"FailureWithErr", false, false, false, testErr, failure, msgErr, errNotify, emptyContent},
		{"Success", false, false, false, nil, success, msgOk, "", emptyContent},
		// Log Request/Response Headers
		{"NOP w/Headers", true, true, false, nil, nop, msgOk, "", logContent},
		{"Failure w/Headers", true, true, false, nil, failure, msgWarn, "", logContentErr},
		{"FailureWithErr w/Headers", true, true, false, testErr, failure, msgErr, errNotify, logContentErr},
		{"Success w/Headers", true, true, false, nil, success, msgOk, "", logContent},
		// Log Request/Response Headers & Body
		{"NOP w/Body", true, true, true, nil, nop, msgOk, "", logContent},
		{"Failure w/Body", true, true, true, nil, failure, msgWarn, "", logBodyErr},
		{"FailureWithErr w/Body", true, true, true, testErr, failure, msgErr, errNotify, logBodyErr},
		{"Success w/Body", true, true, true, nil, success, msgOk, "", logBody},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h := RequestLogger(l, n, test.logRequest, test.logResponse, test.includeBody)(test.h)
			r := &fasthttp.RequestCtx{}
			if test.requestErr != nil {
				fastctx.SetError(r, test.requestErr)
			}
			r.Response.Header.DisableNormalizing()
			logBuf.Reset()
			notifyBuf.Reset()
			h(r)
			if !test.wantLog.EqualsJSON(logBuf.Bytes()) {
				t.Errorf("Log = %v, wantLog %v", logBuf.String(), test.wantLog)
				return
			}

			if notifyBuf.String() != test.wantNotify {
				t.Errorf("Notify = %v, wantNotify %v", notifyBuf.String(), test.wantNotify)
				return
			}

			var gotContent LogMessageWithContent
			err := json.Unmarshal(logBuf.Bytes(), &gotContent)
			if err != nil {
				t.Error(err)
				return
			}

			if !gotContent.HeaderEqual(test.wantContent) {
				t.Errorf("Header = `%v`, want `%v`", gotContent.Header, test.wantContent.Header)
			}
			if gotContent.Body != test.wantContent.Body {
				t.Errorf("Body = `%v`, want `%v`", gotContent.Body, test.wantContent.Body)
			}
			if !gotContent.ResponseHeaderEqual(test.wantContent) {
				t.Errorf("ResponseHeader = `%v`, want `%v`", gotContent.ResponseHeader, test.wantContent.ResponseHeader)
			}
			if gotContent.ResponseBody != test.wantContent.ResponseBody {
				t.Errorf("ResponseBody = `%v`, want `%v`", gotContent.ResponseBody, test.wantContent.ResponseBody)
			}

		})
	}
	return
}

type LogMessage struct {
	Level     string `json:"level"`
	RequestID int    `json:"requestID"`
	Op        string `json:"op"`
	Code      int    `json:"code"`
	IP        string `json:"ip"`
	Message   string `json:"message"`
	//Duration  time.Duration `json:"duration"` // Ignored for comparisons
}

func (l LogMessage) EqualsJSON(in []byte) bool {
	//fmt.Println(string(in))
	var r LogMessage
	err := json.Unmarshal(in, &r)
	if err != nil {
		return false
	}
	return cmp.Equal(r, l)
}

func (l LogMessage) String() string {
	b, err := json.Marshal(l)
	if err != nil {
		return fmt.Sprintf("%v", err)
	}
	return string(b)
}

type LogMessageWithContent struct {
	Header         string `json:"header"`
	Body           string `json:"body"`
	ResponseHeader string `json:"responseHeader"`
	ResponseBody   string `json:"responseBody"`
}

func (l LogMessageWithContent) HeaderEqual(r LogMessageWithContent) bool {
	// We only compare the first line of the header
	lH := strings.Split(l.Header, "\n")[0]
	rH := strings.Split(r.Header, "\n")[0]
	return lH == rH
}

func (l LogMessageWithContent) ResponseHeaderEqual(r LogMessageWithContent) bool {
	// We only compare the first line of the header
	lRH := strings.Split(l.ResponseHeader, "\n")[0]
	rRH := strings.Split(r.ResponseHeader, "\n")[0]
	return lRH == rRH
}
