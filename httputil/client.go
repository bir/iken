package httputil

import (
	"fmt"
	"net/http"
	"strconv"
)

type UnexpectedResponseError struct {
	Resp *http.Response
	URL  string
	Body []byte
}

func (e UnexpectedResponseError) Error() string {
	status := "unknown"
	if e.Resp != nil {
		status = strconv.Itoa(e.Resp.StatusCode)
	}

	return fmt.Sprintf("url: %q, status: %s, body: %q: unexpected response status code", e.URL, status, e.Body)
}
