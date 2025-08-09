package log

import (
	"log/slog"
	"os"
)

// NewLogger creates a new logger instance with the specified level.
func NewLogger() Logger {
	return Logger{slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))}
}

// Logger is a wrapper around slog.Logger.
type Logger struct {
	*slog.Logger
}

// Fatal logs at [LevelError] and exists.
func (l *Logger) Fatal(msg string, args ...any) {
	l.Error(msg, args...)
	os.Exit(1)
}
