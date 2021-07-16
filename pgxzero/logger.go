package pgxzero

import (
	"context"

	"github.com/bir/iken/httputil"
	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog"
)

// Logger manages mapping pgx error messages to Zerolog.
type Logger struct {
	logger zerolog.Logger
}

// New converts pgx logging messages to zerolog.
func New(logger zerolog.Logger) *Logger {
	return &Logger{
		logger: logger.With().Str("module", "pgx").Logger(),
	}
}

// Log is the pgx Logger interface contract.
func (l *Logger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	lvl := zerolog.DebugLevel

	switch level {
	case pgx.LogLevelNone:
		lvl = zerolog.NoLevel
	case pgx.LogLevelError:
		lvl = zerolog.ErrorLevel
	case pgx.LogLevelWarn:
		lvl = zerolog.WarnLevel
	case pgx.LogLevelInfo:
		lvl = zerolog.InfoLevel
	case pgx.LogLevelDebug:
		lvl = zerolog.DebugLevel
	}

	if ctx != nil && (data == nil || data["request_id"] == nil) {
		requestID := httputil.GetID(ctx)
		if requestID != 0 {
			if data == nil {
				data = make(map[string]interface{})
			}

			data["request_id"] = requestID
		}
	}

	l.logger.WithLevel(lvl).Fields(data).Msg(msg)
}
