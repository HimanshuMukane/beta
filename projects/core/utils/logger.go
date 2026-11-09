package utils

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"
	"sync"
)

type loggerContextKey struct{}

var (
	defaultLoggerOnce sync.Once
	rootLogger        *slog.Logger
)

// InitGlobalLogger initializes a production-grade slog JSON or Text logger based on LOG_FORMAT and LOG_LEVEL.
func InitGlobalLogger(env string, w io.Writer) *slog.Logger {
	if w == nil {
		w = os.Stdout
	}

	var minLevel slog.Level
	switch strings.ToLower(os.Getenv("LOG_LEVEL")) {
	case "debug":
		minLevel = slog.LevelDebug
	case "warn", "warning":
		minLevel = slog.LevelWarn
	case "error":
		minLevel = slog.LevelError
	default:
		minLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level:     minLevel,
		AddSource: strings.ToLower(os.Getenv("LOG_ADD_SOURCE")) == "true" || env == "development",
	}

	var handler slog.Handler
	if strings.ToLower(os.Getenv("LOG_FORMAT")) == "text" || (env == "development" && os.Getenv("LOG_FORMAT") == "") {
		handler = slog.NewTextHandler(w, opts)
	} else {
		handler = slog.NewJSONHandler(w, opts)
	}

	rootLogger = slog.New(handler)
	slog.SetDefault(rootLogger)
	return rootLogger
}

// GetGlobalLogger returns the configured root logger, initializing default if not yet called.
func GetGlobalLogger() *slog.Logger {
	defaultLoggerOnce.Do(func() {
		if rootLogger == nil {
			rootLogger = InitGlobalLogger(os.Getenv("APP_ENV"), os.Stdout)
		}
	})
	return rootLogger
}

// WithLogger stores a contextual request-scoped logger inside context.Context.
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey{}, logger)
}

// LoggerFromContext retrieves the contextual logger or falls back to root logger.
func LoggerFromContext(ctx context.Context) *slog.Logger {
	if ctx != nil {
		if l, ok := ctx.Value(loggerContextKey{}).(*slog.Logger); ok && l != nil {
			return l
		}
	}
	return GetGlobalLogger()
}

