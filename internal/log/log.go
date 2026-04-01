package log

import (
	"io"
	"log/slog"
	"os"
)

// NewLogger creates a new text logger that writes to out.
func NewLogger(out io.Writer) Logger {
	return Logger{slog.New(slog.NewTextHandler(out, &slog.HandlerOptions{Level: slog.LevelInfo}))}
}

// Logger wraps an slog.Logger.
type Logger struct {
	*slog.Logger
}

// Fatal logs an error message and exits the process with status code 1.
func (l *Logger) Fatal(msg string, args ...any) {
	l.Error(msg, args...)
	os.Exit(1)
}
