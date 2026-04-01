package test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	dirPerm  = 0o755
	filePerm = 0o600
)

// RunFunc executes the CLI with the supplied arguments and stdio writers.
type RunFunc func(args []string, stdout, stderr io.Writer) int

// RunScenario describes one CLI test case that writes merged output to stdout.
type RunScenario struct {
	Name               string
	Setup              func(t *testing.T, dir string) []string
	Stdout             io.Writer
	Check              func(t *testing.T, dir, stdout, stderr string)
	WantStdout         string
	WantStderrContains []string
	WantStderrExcludes []string
	WantExit           int
	SkipStdoutCheck    bool
	WantStderrEmpty    bool
}

// FileOutputScenario describes one CLI test case that writes merged output to a
// file via `-o`.
type FileOutputScenario struct {
	Name               string
	Setup              func(t *testing.T, dir string) []string
	Check              func(t *testing.T, dir, stdout, stderr string)
	WantStderrContains []string
	WantExit           int
	WantStderrEmpty    bool
}

// FailingWriter is an io.Writer that always returns Err from Write.
type FailingWriter struct {
	Err error
}

// Write implements io.Writer.
func (w FailingWriter) Write(_ []byte) (int, error) {
	return 0, w.Err
}

// RunScenarioCase executes one RunScenario with a temporary working directory.
func RunScenarioCase(t *testing.T, run RunFunc, tt RunScenario) {
	t.Helper()

	dir := t.TempDir()
	args := []string(nil)
	if tt.Setup != nil {
		args = tt.Setup(t, dir)
	}

	exitCode, stdout, stderr := ExecuteRun(t, run, args, tt.Stdout)
	require.Equalf(t, tt.WantExit, exitCode, "stderr: %q", stderr)

	if !tt.SkipStdoutCheck {
		require.Equal(t, tt.WantStdout, stdout)
	}

	assertStderr(t, stderr, tt.WantStderrEmpty, tt.WantStderrContains, tt.WantStderrExcludes)

	if tt.Check != nil {
		tt.Check(t, dir, stdout, stderr)
	}
}

// RunFileOutputScenarioCase executes one FileOutputScenario with a temporary
// working directory.
func RunFileOutputScenarioCase(t *testing.T, run RunFunc, tt FileOutputScenario) {
	t.Helper()

	dir := t.TempDir()
	exitCode, stdout, stderr := ExecuteRun(t, run, tt.Setup(t, dir), nil)
	require.Equalf(t, tt.WantExit, exitCode, "stderr: %q", stderr)

	assertStderr(t, stderr, tt.WantStderrEmpty, tt.WantStderrContains, nil)

	if tt.Check != nil {
		tt.Check(t, dir, stdout, stderr)
	}
}

// ExecuteRun executes run with in-memory stdout and stderr buffers.
func ExecuteRun(t *testing.T, run RunFunc, args []string, stdoutWriter io.Writer) (int, string, string) {
	t.Helper()

	var stdout bytes.Buffer
	if stdoutWriter == nil {
		stdoutWriter = &stdout
	}

	var stderr bytes.Buffer
	exitCode := run(args, stdoutWriter, &stderr)

	return exitCode, stdout.String(), stderr.String()
}

// WriteProfileFile writes body to dir/name and returns the resulting path.
func WriteProfileFile(t *testing.T, dir, name, body string) string {
	t.Helper()

	path := filepath.Join(dir, name)
	WriteTextFile(t, path, body)

	return path
}

// WriteTextFile writes body to path, creating parent directories as needed.
func WriteTextFile(t *testing.T, path, body string) {
	t.Helper()

	err := os.MkdirAll(filepath.Dir(path), dirPerm)
	require.NoErrorf(t, err, "failed to create parent directory for %q", path)

	err = os.WriteFile(path, []byte(body), filePerm)
	require.NoErrorf(t, err, "failed to write file %q", path)
}

// TextProfile formats a coverage profile file body with the given mode and
// block lines.
func TextProfile(mode string, blocks ...string) string {
	return fmt.Sprintf("mode: %s\n%s\n", mode, strings.Join(blocks, "\n"))
}

// TextBlock formats one coverage profile block line.
func TextBlock(file string, startLine, startCol, endLine, endCol, numStmt, count int) string {
	return fmt.Sprintf("%s:%d.%d,%d.%d %d %d", file, startLine, startCol, endLine, endCol, numStmt, count)
}

func assertStderr(t *testing.T, stderr string, wantEmpty bool, contains, excludes []string) {
	t.Helper()

	if wantEmpty {
		require.Empty(t, stderr)
	}

	for _, want := range contains {
		require.Contains(t, stderr, want)
	}

	for _, unwanted := range excludes {
		require.NotContains(t, stderr, unwanted)
	}
}
