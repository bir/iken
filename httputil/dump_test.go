package httputil

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"
)

type eofReader struct{}

func (n eofReader) Close() error { return nil }

func (n eofReader) Read([]byte) (int, error) { return 0, io.EOF }

type dumpTest struct {
	Name string
	// Either Req or GetReq can be set/nil but not both.
	Req    *http.Request
	GetReq func() *http.Request

	Body any // optional []byte or func() io.ReadCloser to populate Req.Body

	WantBody   string
	WantHeader map[string]string
	MustError  bool // if true, the test is expected to throw an error
}

type errReader struct{}

func (n errReader) Close() error { return nil }

func (n errReader) Read([]byte) (int, error) { return 0, fmt.Errorf("errReader") }

type errCloser struct{}

func (n errCloser) Close() error { return fmt.Errorf("errCloser") }

func (n errCloser) Read([]byte) (int, error) { return 0, io.EOF }

var dumpTests = []dumpTest{
	{
		Name: "HTTP/1.1 => chunked coding; body; empty trailer",
		Req: &http.Request{
			Method: "GET",
			URL: &url.URL{
				Scheme: "http",
				Host:   "www.google.com",
				Path:   "/search",
			},
			ProtoMajor:       1,
			ProtoMinor:       1,
			TransferEncoding: []string{"chunked"},
		},

		Body: []byte("abcdef"),

		WantHeader: map[string]string{
			"GET":               "/search HTTP/1.1",
			"Host":              "www.google.com",
			"Transfer-Encoding": "chunked",
		},
		WantBody: chunk("abcdef") + chunk(""),
	},
	{
		Name: "Verify that DumpRequest preserves the HTTP version number, doesn't add a Host",
		Req: &http.Request{
			Method:     "GET",
			URL:        mustParseURL("/foo"),
			ProtoMajor: 1,
			ProtoMinor: 0,
			Header: http.Header{
				"X-Foo": []string{"X-Bar"},
			},
		},

		WantHeader: map[string]string{
			"GET":   "/foo HTTP/1.0",
			"X-Foo": "X-Bar",
		},
	},
	{
		Name: "ErrReader Body",
		Req: &http.Request{
			URL: &url.URL{
				Scheme: "http",
				Host:   "www.google.com",
				Path:   "/search",
			},
			ProtoMajor: 1,
			ProtoMinor: 1,
			Body:       errReader{},
		},

		WantHeader: map[string]string{
			"GET":  "/search HTTP/1.1",
			"Host": "www.google.com",
		},
		MustError: true,
	},
	{
		Name: "ErrCloser Body",
		Req: &http.Request{
			URL: &url.URL{
				Scheme: "http",
				Host:   "www.google.com",
				Path:   "/search",
			},
			ProtoMajor: 1,
			ProtoMinor: 1,
			Body:       errCloser{},
		},

		WantHeader: map[string]string{
			"GET":  "/search HTTP/1.1",
			"Host": "www.google.com",
		},
		MustError: true,
	},
	{
		Name: "Request with Body > 8196 (default buffer size)",
		Req: &http.Request{
			Method: "POST",
			URL: &url.URL{
				Scheme: "http",
				Host:   "post.tld",
				Path:   "/",
			},
			Header: http.Header{
				"Content-Length": []string{"8193"},
			},

			ContentLength: 8193,
			ProtoMajor:    1,
			ProtoMinor:    1,
		},

		Body: bytes.Repeat([]byte("a"), 8193),
		WantHeader: map[string]string{
			"POST":           "/ HTTP/1.1",
			"Host":           "post.tld",
			"Content-Length": "8193",
		},

		WantBody: strings.Repeat("a", 8193),
	},

	{
		Name: "User-Agent dumped",
		GetReq: func() *http.Request {
			return mustReadRequest("GET http://foo.com/ HTTP/1.1\r\n" +
				"User-Agent: blah\r\n\r\n")
		},
		WantHeader: map[string]string{
			"GET":        "http://foo.com/ HTTP/1.1",
			"User-Agent": "blah",
		},
	},

	{
		Name: "DumpRequest should return the \"Content-Length\" when set",
		GetReq: func() *http.Request {
			return mustReadRequest("POST /v2/api/?login HTTP/1.1\r\n" +
				"Host: passport.myhost.com\r\n" +
				"Content-Length: 3\r\n" +
				"\r\nkey1=name1&key2=name2")
		},
		WantHeader: map[string]string{
			"POST":           "/v2/api/?login HTTP/1.1",
			"Host":           "passport.myhost.com",
			"Content-Length": "3",
		},
		WantBody: "key",
	},
	{
		Name: "Issue #7215. DumpRequest should return the \"Content-Length\" in ReadRequest",
		GetReq: func() *http.Request {
			return mustReadRequest("POST /v2/api/?login HTTP/1.1\r\n" +
				"Host: passport.myhost.com\r\n" +
				"Content-Length: 0\r\n" +
				"\r\nkey1=name1&key2=name2")
		},
		WantHeader: map[string]string{
			"POST":           "/v2/api/?login HTTP/1.1",
			"Host":           "passport.myhost.com",
			"Content-Length": "0",
		},
	},

	{
		Name: "Issue #7215. DumpRequest should not return the \"Content-Length\" if unset",
		GetReq: func() *http.Request {
			return mustReadRequest("POST /v2/api/?login HTTP/1.1\r\n" +
				"Host: passport.myhost.com\r\n" +
				"Trailer: Expires\r\n" +
				"\r\nkey1=name1&key2=name2")
		},
		WantHeader: map[string]string{
			"POST": "/v2/api/?login HTTP/1.1",
			"Host": "passport.myhost.com",
		},
	},
	{
		Name: "Issue #7215. DumpRequest should not return the \"Content-Length\" if unset",
		GetReq: func() *http.Request {
			return mustReadRequest("POST /v2/api/?login HTTP/1.1\r\n" +
				"Host: passport.myhost.com\r\n" +
				"\r\nkey1=name1&key2=name2")
		},
		WantHeader: map[string]string{
			"POST": "/v2/api/?login HTTP/1.1",
			"Host": "passport.myhost.com",
		},
	},
}

func TestDumpRequest(t *testing.T) {
	numGoroutine := runtime.NumGoroutine()
	for _, tt := range dumpTests {
		if tt.Req != nil && tt.GetReq != nil || tt.Req == nil && tt.GetReq == nil {
			t.Errorf("%q: either .Req(%p) or .GetReq(%p) can be set/nil but not both", tt.Name, tt.Req, tt.GetReq)
			continue
		}

		freshReq := func(ti dumpTest) *http.Request {
			req := ti.Req
			if req == nil {
				req = ti.GetReq()
			}

			if req.Header == nil {
				req.Header = make(http.Header)
			}

			if ti.Body == nil {
				return req
			}
			switch b := ti.Body.(type) {
			case []byte:
				req.Body = io.NopCloser(bytes.NewReader(b))
			case func() io.ReadCloser:
				req.Body = b()
			default:
				t.Fatalf("Test %q: unsupported Body of %T", tt.Name, ti.Body)
			}
			return req
		}

		req := freshReq(tt)
		got := DumpHeader(req)
		if !reflect.DeepEqual(got, tt.WantHeader) {
			t.Errorf("DumpHeader %q, expecting:\n%s\nGot:\n%s\n", tt.Name, tt.WantHeader, got)
			continue
		}

		req = freshReq(tt)
		dump, err := DumpBody(req)
		if err != nil && !tt.MustError {
			t.Errorf("DumpBody %q: %s\nWantDump:\n%s", tt.Name, err, tt.WantBody)
			continue
		}
		if string(dump) != tt.WantBody {
			t.Errorf("DumpBody %q, expecting:\n%s\nGot:\n%s\n", tt.Name, tt.WantBody, string(dump))
			continue
		}

		if tt.MustError {
			req := freshReq(tt)
			_, err := DumpBody(req)
			if err == nil {
				t.Errorf("DumpBody %q: expected an error, got nil", tt.Name)
			}
			continue
		}

	}

	// Validate we haven't leaked any goroutines.
	var dg int
	dl := deadline(t, 5*time.Second, time.Second)
	for time.Now().Before(dl) {
		if dg = runtime.NumGoroutine() - numGoroutine; dg <= 4 {
			// No unexpected goroutines.
			return
		}

		// Allow goroutines to schedule and die off.
		runtime.Gosched()
	}

	buf := make([]byte, 4096)
	buf = buf[:runtime.Stack(buf, true)]
	t.Errorf("Unexpectedly large number of new goroutines: %d new: %s", dg, buf)
}

// deadline returns the time which is needed before t.Deadline()
// if one is configured, and it is s greater than needed in the future,
// otherwise defaultDelay from the current time.
func deadline(t *testing.T, defaultDelay, needed time.Duration) time.Time {
	if dl, ok := t.Deadline(); ok {
		if dl = dl.Add(-needed); dl.After(time.Now()) {
			// Allow an arbitrarily long delay.
			return dl
		}
	}

	// No deadline configured or its closer than needed from now
	// so just use the default.
	return time.Now().Add(defaultDelay)
}

func chunk(s string) string {
	return fmt.Sprintf("%x\r\n%s\r\n", len(s), s)
}

func mustParseURL(s string) *url.URL {
	u, err := url.Parse(s)
	if err != nil {
		panic(fmt.Sprintf("Error parsing URL %q: %v", s, err))
	}
	return u
}

func mustReadRequest(s string) *http.Request {
	req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(s)))
	if err != nil {
		panic(err)
	}
	return req
}
