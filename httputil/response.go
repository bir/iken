package httputil

import (
	"encoding/json"
	"net/http"
)

// JSONWrite is a simple helper utility to return the json encoded obj with appropriate content-type and code.
func JSONWrite(w http.ResponseWriter, code int, obj interface{}) error {
	w.Header().Set(ContentType, ApplicationJSON)
	w.WriteHeader(code)

	return json.NewEncoder(w).Encode(obj)
}
