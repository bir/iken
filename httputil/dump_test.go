package httputil

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
)

type dumpTest struct {
	Name   string
	Req    *http.Request
	GetReq func() *http.Request

	Body interface{} // optional []byte or func() io.ReadCloser to populate Req.Body

	WantHeader map[string]string
	WantBody   string
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
		Name: "Chunked coding",
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
		Name: "Default Method",
		Req: &http.Request{
			URL: &url.URL{
				Scheme: "http",
				Host:   "www.google.com",
				Path:   "/search",
			},
			ProtoMajor: 1,
			ProtoMinor: 1,
			Close:      true,
		},

		WantHeader: map[string]string{
			"GET":        "/search HTTP/1.1",
			"Host":       "www.google.com",
			"Connection": "close",
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
		Name: "Verify that DumpRequest preserves the HTTP version number, doesn't add a Host, and doesn't add a User-Agent.",
		Req: &http.Request{
			Method:     "GET",
			URL:        mustParseURL("/foo"),
			ProtoMajor: 1,
			ProtoMinor: 0,
			Header: http.Header{
				"X-Foo":   []string{"X-Bar"},
				"Trailer": []string{"Foo"},
			},
		},

		WantHeader: map[string]string{"GET": "/foo HTTP/1.0",
			"X-Foo": "X-Bar",
		},
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
		Name: "UserAgent",
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
		Name: "DumpRequest should return the Content-Length when set",
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
		Name: "DumpRequest should return the Content-Length when set to 0",
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
		Name: "DumpRequest should not return the Content-Length when not set",
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
	for i, tt := range dumpTests {
		t.Run(tt.Name, func(t *testing.T) {
			if tt.Req != nil && tt.GetReq != nil || tt.Req == nil && tt.GetReq == nil {
				t.Errorf("#%d: either .Req(%p) or .GetReq(%p) can be set/nil but not both", i, tt.Req, tt.GetReq)
				return
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
					t.Fatalf("Test %d: unsupported Body of %T", i, ti.Body)
				}
				return req
			}

			req := freshReq(tt)
			dump := DumpHeader(req)
			if !reflect.DeepEqual(dump, tt.WantHeader) {
				t.Errorf("DumpRequest %d, expecting:\n%s\nGot:\n%s\n", i, tt.WantHeader, dump)
				return
			}

			b, err := DumpBody(req)
			if tt.MustError && err == nil {
				t.Errorf("DumpRequest #%d: expected error", i)
				return
			}

			if tt.WantBody != string(b) {
				t.Errorf("DumpRequest %d, expecting body:\n`%s`\nGot:\n`%s`\n", i, tt.WantBody, b)
				return
			}
		})
	}
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
