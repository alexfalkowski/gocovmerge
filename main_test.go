package main

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunHelpReturnsZeroAndWritesUsageToStderr(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := run([]string{"-help"}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("expected help to exit 0, got %d", exitCode)
	}

	if stdout.Len() != 0 {
		t.Fatalf("expected help to keep stdout empty, got %q", stdout.String())
	}

	got := stderr.String()
	if !strings.Contains(got, "Usage of gocovmerge:") {
		t.Fatalf("expected usage text on stderr, got %q", got)
	}

	if strings.Contains(got, "flag: help requested") {
		t.Fatalf("expected help not to be logged as an error, got %q", got)
	}
}

func TestRunParseErrorWritesOnlyToStderr(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing.out")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := run([]string{missing}, &stdout, &stderr)
	if exitCode != 1 {
		t.Fatalf("expected parse failure to exit 1, got %d", exitCode)
	}

	if stdout.Len() != 0 {
		t.Fatalf("expected parse failure to keep stdout empty, got %q", stdout.String())
	}

	got := stderr.String()
	if !strings.Contains(got, "failed to parse profiles") {
		t.Fatalf("expected parse failure to be written to stderr, got %q", got)
	}

	if !strings.Contains(got, missing) {
		t.Fatalf("expected missing file path in stderr, got %q", got)
	}
}
