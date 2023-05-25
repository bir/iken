package httputil

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// JSONWrite is a simple helper utility to return the json encoded obj with appropriate content-type and code.
func JSONWrite(w http.ResponseWriter, r *http.Request, code int, obj interface{}) {
	b, err := json.Marshal(obj)
	if err != nil {
		ErrorHandler(w, r, fmt.Errorf("JSONWrite:%w", err))
	}

	Write(w, r, ApplicationJSON, code, b)
}

func HTMLWrite(w http.ResponseWriter, r *http.Request, code int, data string) {
	Write(w, r, TextHTML, code, []byte(data))
}

func Write(w http.ResponseWriter, r *http.Request, contentType string, code int, data []byte) {
	w.WriteHeader(code)
	w.Header().Set(ContentType, contentType)

	if _, err := w.Write(data); err != nil {
		ErrorHandler(w, r, err)
	}
}

func AddHeaders(w http.ResponseWriter, headers http.Header) {
	for k, v := range headers {
		for _, h := range v {
			w.Header().Add(k, h)
		}
	}
}
