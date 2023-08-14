package httplog

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/bir/iken/logctx"
)

func TestRecover(t *testing.T) {
	MaxRequestBodyLog = 10
	RecoverBasePath = "iken/httplog/"

	tests := []struct {
		name          string
		body          string
		next          http.Handler
		wantFirstLine string
	}{
		{"panic String", "123", readPanic("test"), "./recover_test.go:65 (iken/httplog.TestRecover.readPanic.func2)"},
		{"panic Error", "123", readPanic(errors.New("test")), "./recover_test.go:65 (iken/httplog.TestRecover.readPanic.func3)"},
		{"panic other", "123", readPanic(1), "./recover_test.go:65 (iken/httplog.TestRecover.readPanic.func4)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logOutput := bytes.NewBuffer(nil)

			h := RecoverLogger(zerolog.New(logOutput))

			w := httptest.NewRecorder()
			r := httptest.NewRequest("FOO", "/BAR", bytes.NewBufferString(tt.body))

			now = startNow
			h(tt.next).ServeHTTP(w, r)

			got := logOutput.String()

			result := make(map[string]any)
			err := json.Unmarshal([]byte(got), &result)
			assert.Nil(t, err, "json Unmarshal")

			assert.Equal(t, "value", result["key"], "log context")
			assert.Equal(t, "Panic", result["message"], "log context")

			stack, ok := result["error.stack"].([]any)
			assert.True(t, ok, "error.stack type")

			assert.Equal(t, tt.wantFirstLine, stack[0], "logs")
		})
	}
}

func readPanic(result any) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		now = endNow
		logctx.AddStrToContext(r.Context(), "key", "value")

		panic(result)
	}
}
