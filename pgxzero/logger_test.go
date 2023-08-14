package pgxzero_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/jackc/pgx/v5/tracelog"
	"github.com/rs/zerolog"

	"github.com/bir/iken/logctx"
	"github.com/bir/iken/pgxzero"
)

func TestLogger_Log(t *testing.T) {

	dataWithRequest := map[string]any{"request_id": 123}
	dataWithoutRequest := map[string]any{"other": 123}

	ctx := logctx.SetID(context.Background(), "121")
	tests := []struct {
		name  string
		ctx   context.Context
		level tracelog.LogLevel
		msg   string
		data  map[string]any
		want  string
	}{
		{"no level", nil, 0, "no level", nil, "{\"level\":\"debug\",\"module\":\"tracelog\",\"message\":\"no level\"}\n"},
		{"none", nil, tracelog.LogLevelNone, "default", nil, "{\"module\":\"tracelog\",\"message\":\"default\"}\n"},
		{"error", nil, tracelog.LogLevelError, "error", nil, "{\"level\":\"error\",\"module\":\"tracelog\",\"message\":\"error\"}\n"},
		{"warn", nil, tracelog.LogLevelWarn, "warn", nil, "{\"level\":\"warn\",\"module\":\"tracelog\",\"message\":\"warn\"}\n"},
		{"info", nil, tracelog.LogLevelInfo, "info", nil, "{\"level\":\"info\",\"module\":\"tracelog\",\"message\":\"info\"}\n"},
		{"debug", nil, tracelog.LogLevelDebug, "debug", nil, "{\"level\":\"debug\",\"module\":\"tracelog\",\"message\":\"debug\"}\n"},
		{"trace", nil, tracelog.LogLevelTrace, "trace", nil, "{\"level\":\"trace\",\"module\":\"tracelog\",\"message\":\"trace\"}\n"},
		{"withID in Data", ctx, tracelog.LogLevelWarn, "ctx", dataWithRequest, "{\"level\":\"warn\",\"module\":\"tracelog\",\"request_id\":123,\"message\":\"ctx\"}\n"},
		{"withID in Ctx", ctx, tracelog.LogLevelWarn, "ctx", dataWithoutRequest, "{\"level\":\"warn\",\"module\":\"tracelog\",\"http.request_id\":\"121\",\"other\":123,\"message\":\"ctx\"}\n"},
		{"withID in Ctx no data", ctx, tracelog.LogLevelWarn, "ctx", nil, "{\"level\":\"warn\",\"module\":\"tracelog\",\"http.request_id\":\"121\",\"message\":\"ctx\"}\n"},
	}

	var logBuf bytes.Buffer
	l := zerolog.New(&logBuf)
	pgxLogger := pgxzero.New(l)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			logBuf.Reset()
			pgxLogger.Log(test.ctx, test.level, test.msg, test.data)

			got := logBuf.String()
			if got != test.want {
				t.Errorf("got `%v`, want `%v`", got, test.want)
			}
		})
	}

	logBuf.Reset()

	pgxLogger.WithMapper(func(level tracelog.LogLevel, s string) zerolog.Level {
		return zerolog.FatalLevel
	})

	pgxLogger.Log(context.Background(), tracelog.LogLevelDebug, "test fatal mapping", nil)
	got := logBuf.String()
	want := `{"level":"fatal","module":"tracelog","message":"test fatal mapping"}
`
	if got != want {
		t.Errorf("got `%v`, want `%v`", got, want)
	}
}
