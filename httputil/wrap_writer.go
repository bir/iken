package httputil

import (
	"bufio"
	"errors"
	"io"
	"net"
	"net/http"
)

// Cloned from https://github.com/rs/zerolog/blob/master/hlog/internal/mutil/writer_proxy.go
// Used to allow tracking status, bytes, and optionally copying of response payloads.

// WriterProxy is a proxy around an http.ResponseWriter that allows you to hook
// into various parts of the response process.
type WriterProxy interface {
	http.ResponseWriter
	// Status returns the HTTP status of the request, or 0 if one has not
	// yet been sent.
	Status() int
	// BytesWritten returns the total number of bytes sent to the client.
	BytesWritten() int
	// Tee causes the response body to be written to the given io.Writer in
	// addition to proxying the writes through. Only one io.Writer can be
	// tee'd to at once: setting a second one will overwrite the first.
	// Writes will be sent to the proxy before being written to this
	// io.Writer. It is illegal for the tee'd writer to be modified
	// concurrently with writes.
	Tee(w io.Writer)
	// Unwrap returns the original proxied target.
	Unwrap() http.ResponseWriter
}

// WrapWriter wraps an http.ResponseWriter, returning a proxy that allows you to
// hook into various parts of the response process.
func WrapWriter(w http.ResponseWriter) WriterProxy {
	_, fl := w.(http.Flusher)
	_, hj := w.(http.Hijacker)
	_, rf := w.(io.ReaderFrom)

	bw := basicWriter{ResponseWriter: w}

	if fl && hj && rf {
		return &fancyWriter{bw}
	}

	return &bw
}

// basicWriter wraps a http.ResponseWriter that implements the minimal
// http.ResponseWriter interface.
type basicWriter struct {
	http.ResponseWriter
	wroteHeader bool
	code        int
	bytes       int
	tee         io.Writer
}

func (b *basicWriter) WriteHeader(code int) {
	if !b.wroteHeader {
		b.code = code
		b.wroteHeader = true
		b.ResponseWriter.WriteHeader(code)
	}
}

func (b *basicWriter) Write(buf []byte) (int, error) {
	b.WriteHeader(http.StatusOK)
	n, err := b.ResponseWriter.Write(buf)

	if b.tee != nil {
		_, err2 := b.tee.Write(buf[:n])
		// Prefer errors generated by the proxied writer.
		if err == nil {
			err = err2
		}
	}

	b.bytes += n

	return n, err
}

func (b *basicWriter) maybeWriteHeader() {
	if !b.wroteHeader {
		b.WriteHeader(http.StatusOK)
	}
}

func (b *basicWriter) Status() int {
	return b.code
}

func (b *basicWriter) BytesWritten() int {
	return b.bytes
}

func (b *basicWriter) Tee(w io.Writer) {
	b.tee = w
}

func (b *basicWriter) Unwrap() http.ResponseWriter {
	return b.ResponseWriter
}

func (b *basicWriter) Flush() {
	if fl, ok := b.ResponseWriter.(http.Flusher); ok {
		fl.Flush()
	}
}

// fancyWriter is a writer that additionally satisfies http.CloseNotifier,
// http.Flusher, http.Hijacker, and io.ReaderFrom. It exists for the common case
// of wrapping the http.ResponseWriter that package http gives you, in order to
// make the proxied object support the full method set of the proxied object.
type fancyWriter struct {
	basicWriter
}

func (f *fancyWriter) Flush() {
	f.basicWriter.Flush()
}

var ErrUnsupported = errors.New("not implemented")

func (f *fancyWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj, ok := f.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, ErrUnsupported // Ignore coverage - unlikely to error
	}

	return hj.Hijack() //nolint:wrapcheck // just a proxy
}

func (f *fancyWriter) ReadFrom(r io.Reader) (int64, error) {
	if f.tee != nil {
		return io.Copy(&f.basicWriter, r) //nolint:wrapcheck // just a proxy
	}

	rf, ok := f.ResponseWriter.(io.ReaderFrom)
	if !ok {
		return 0, ErrUnsupported // Ignore coverage - unlikely to error
	}

	f.maybeWriteHeader()

	n, err := rf.ReadFrom(r)
	f.bytes += int(n)

	return n, err //nolint:wrapcheck // just a proxy
}

var (
	_ http.Flusher  = &fancyWriter{}
	_ http.Hijacker = &fancyWriter{}
	_ io.ReaderFrom = &fancyWriter{}
	_ http.Flusher  = &basicWriter{}
)
