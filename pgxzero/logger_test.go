package pgxzero_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/bir/iken/httputil"
	"github.com/bir/iken/pgxzero"
	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog"
)

func TestLogger_Log(t *testing.T) {

	dataWithRequest := map[string]interface{}{"request_id": 123}
	dataWithoutRequest := map[string]interface{}{"other": 123}

	ctx := httputil.SetID(context.Background(), 121)
	tests := []struct {
		name  string
		ctx   context.Context
		level pgx.LogLevel
		msg   string
		data  map[string]interface{}
		want  string
	}{
		{"none", nil, pgx.LogLevelNone, "default", nil, "{\"module\":\"pgx\",\"message\":\"default\"}\n"},
		{"error", nil, pgx.LogLevelError, "error", nil, "{\"level\":\"error\",\"module\":\"pgx\",\"message\":\"error\"}\n"},
		{"warn", nil, pgx.LogLevelWarn, "warn", nil, "{\"level\":\"warn\",\"module\":\"pgx\",\"message\":\"warn\"}\n"},
		{"info", nil, pgx.LogLevelInfo, "info", nil, "{\"level\":\"info\",\"module\":\"pgx\",\"message\":\"info\"}\n"},
		{"debug", nil, pgx.LogLevelDebug, "debug", nil, "{\"level\":\"debug\",\"module\":\"pgx\",\"message\":\"debug\"}\n"},
		{"trace", nil, pgx.LogLevelTrace, "trace", nil, "{\"level\":\"debug\",\"module\":\"pgx\",\"message\":\"trace\"}\n"},
		{"withID in Data", ctx, pgx.LogLevelWarn, "ctx", dataWithRequest, "{\"level\":\"warn\",\"module\":\"pgx\",\"request_id\":123,\"message\":\"ctx\"}\n"},
		{"withID in Ctx", ctx, pgx.LogLevelWarn, "ctx", dataWithoutRequest, "{\"level\":\"warn\",\"module\":\"pgx\",\"other\":123,\"request_id\":121,\"message\":\"ctx\"}\n"},
		{"withID in Ctx no data", ctx, pgx.LogLevelWarn, "ctx", nil, "{\"level\":\"warn\",\"module\":\"pgx\",\"request_id\":121,\"message\":\"ctx\"}\n"},
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

	pgxLogger.WithMapper(func(level pgx.LogLevel, s string) zerolog.Level {
		return zerolog.FatalLevel
	})

	pgxLogger.Log(context.Background(), pgx.LogLevelDebug, "test fatal mapping", nil)
	got := logBuf.String()
	want := `{"level":"fatal","module":"pgx","message":"test fatal mapping"}
`
	if got != want {
		t.Errorf("got `%v`, want `%v`", got, want)
	}
}
