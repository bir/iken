package httputil

var (
	// ContentType header value.
	ContentType = "Content-Type"
	// ApplicationJSON content-type.
	ApplicationJSON = "application/json"
	// TextHTML content-type.
	TextHTML = "text/html"
	// TextPlain content-type.
	TextPlain = "text/plain; charset=utf-8"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	// ErrNotFound represents failure when authenticating a request.
	ErrNotFound = Error("not found")
)
