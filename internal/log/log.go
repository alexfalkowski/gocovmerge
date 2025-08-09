package log

import (
	"log/slog"
	"os"
)

// NewLogger creates a new text logger.
func NewLogger() Logger {
	return Logger{slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))}
}

// Logger is an slog.Logger.
type Logger struct {
	*slog.Logger
}

// Fatal logs an Error and exists.
func (l *Logger) Fatal(msg string, args ...any) {
	l.Error(msg, args...)
	os.Exit(1)
}
