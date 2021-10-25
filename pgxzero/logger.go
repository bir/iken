package pgxzero

import (
	"context"

	"github.com/bir/iken/httputil"
	"github.com/jackc/pgx/v4"
	"github.com/rs/zerolog"
)

// LevelMapper converts pgx Log levels to zerolog levels.  This allows custom overrides for the levels provided by pgx.
type LevelMapper = func(pgx.LogLevel, string) zerolog.Level

// Logger manages mapping pgx error messages to Zerolog.
type Logger struct {
	logger zerolog.Logger
	mapper LevelMapper
}

func defaultMapper(level pgx.LogLevel, _ string) zerolog.Level {
	switch level {
	case pgx.LogLevelNone:
		return zerolog.NoLevel
	case pgx.LogLevelError:
		return zerolog.ErrorLevel
	case pgx.LogLevelWarn:
		return zerolog.WarnLevel
	case pgx.LogLevelInfo:
		return zerolog.InfoLevel
	case pgx.LogLevelDebug:
		return zerolog.DebugLevel
	}

	return zerolog.DebugLevel
}

// New converts pgx logging messages to zerolog.
func New(logger zerolog.Logger) *Logger {
	return &Logger{
		logger: logger.With().Str("module", "pgx").Logger(),
		mapper: defaultMapper,
	}
}

func (l *Logger) WithMapper(m LevelMapper) *Logger {
	l.mapper = m

	return l
}

// Log is the pgx Logger interface contract.
func (l *Logger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	if ctx != nil && (data == nil || data["request_id"] == nil) {
		requestID := httputil.GetID(ctx)
		if requestID != 0 {
			if data == nil {
				data = make(map[string]interface{})
			}

			data["request_id"] = requestID
		}
	}

	l.logger.WithLevel(l.mapper(level, msg)).Fields(data).Msg(msg)
}
