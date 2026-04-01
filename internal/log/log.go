package log

import (
	"io"
	"log/slog"
)

// NewLogger creates a text logger that writes structured logs to out.
//
// The gocovmerge CLI passes stderr so diagnostics never contaminate the merged
// coverage profile stream.
func NewLogger(out io.Writer) *slog.Logger {
	return slog.New(slog.NewTextHandler(out, &slog.HandlerOptions{Level: slog.LevelInfo}))
}
