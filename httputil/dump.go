package httputil

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"strings"
)

// Ported from go stdlib httputil.dump.  Tweaked to split the header and body into separate functions for more
// flexible logging. Header is returned as a map[string]string for ease of handling.  Strictly a logging utility.

func DumpHeader(req *http.Request) map[string]string {
	out := map[string]string{}

	reqURI := req.RequestURI
	if reqURI == "" {
		reqURI = req.URL.RequestURI()
	}

	out[valueOrDefault(req.Method, "GET")] = fmt.Sprintf("%s HTTP/%d.%d", reqURI, req.ProtoMajor, req.ProtoMinor)

	absRequestURI := strings.HasPrefix(req.RequestURI, "http://") || strings.HasPrefix(req.RequestURI, "https://")
	if !absRequestURI {
		host := req.Host
		if host == "" && req.URL != nil {
			host = req.URL.Host
		}

		if host != "" {
			out["Host"] = host
		}
	}

	if len(req.TransferEncoding) > 0 {
		out["Transfer-Encoding"] = strings.Join(req.TransferEncoding, ",")
	}

	if req.Close {
		out["Connection"] = "close"
	}

	for name, values := range req.Header {
		if reqWriteExcludeHeaderDump[name] {
			continue
		}

		out[name] = strings.Join(values, ",")
	}

	return out
}

func DumpBody(req *http.Request) ([]byte, error) {
	if req.Body == nil {
		return nil, nil
	}

	var (
		err  error
		save io.ReadCloser
	)

	save, req.Body, err = drainBody(req.Body)
	if err != nil {
		return nil, err
	}

	chunked := len(req.TransferEncoding) > 0 && req.TransferEncoding[0] == "chunked"

	var b bytes.Buffer

	err = dumpBody(req, &b, chunked)
	if err != nil {
		return nil, err
	}

	req.Body = save

	return b.Bytes(), nil
}

func dumpBody(req *http.Request, body io.Writer, chunked bool) error {
	w := body

	if chunked {
		w = httputil.NewChunkedWriter(body)
	}

	_, err := io.Copy(w, req.Body)
	if err != nil {
		return fmt.Errorf("copy:%w", err)
	}

	if chunked {
		if closer, ok := w.(io.Closer); ok {
			err = closer.Close()
		}

		if err != nil {
			return fmt.Errorf("close:%w", err)
		}

		_, err = io.WriteString(body, "\r\n")
		if err != nil {
			return fmt.Errorf("io.WriteString:%w", err)
		}
	}

	return nil
}

func valueOrDefault(value, def string) string {
	if value != "" {
		return value
	}

	return def
}

var reqWriteExcludeHeaderDump = map[string]bool{
	"Host":              true, // not in Header map anyway
	"Transfer-Encoding": true,
	"Trailer":           true,
}

func drainBody(body io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	if body == nil || body == http.NoBody {
		// No copying needed. Preserve the magic sentinel meaning of NoBody.
		return http.NoBody, http.NoBody, nil
	}

	var buf bytes.Buffer

	if _, err = buf.ReadFrom(body); err != nil {
		return nil, body, fmt.Errorf("buf.ReadFrom:%w", err)
	}

	if err = body.Close(); err != nil {
		return nil, body, fmt.Errorf("body.Close:%w", err)
	}

	return io.NopCloser(&buf), io.NopCloser(bytes.NewReader(buf.Bytes())), nil
}
