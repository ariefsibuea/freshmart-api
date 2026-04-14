package logger

import (
	"context"
	"log/slog"
	"os"
)

type contextKey struct{}

const (
	FieldRequestID = "request_id"
	FieldMethod    = "method"
	FieldPath      = "path"
	FieldStatus    = "status"
	FieldLatencyMS = "latency_ms"
	FieldError     = "error"
)

var global *slog.Logger

func init() {
	global = slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

func Init(level slog.Level) {
	opts := &slog.HandlerOptions{Level: level}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	global = slog.New(handler)

	slog.SetDefault(global)
}

func IntoContext(ctx context.Context, log *slog.Logger) context.Context {
	return context.WithValue(ctx, contextKey{}, log)
}

func FromContext(ctx context.Context) *slog.Logger {
	if log, ok := ctx.Value(contextKey{}).(*slog.Logger); ok && log != nil {
		return log
	}
	return global
}

func ToSlogLevel(level int) slog.Level {
	switch level {
	case 0:
		return slog.LevelDebug
	case 2:
		return slog.LevelWarn
	case 3:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
