package main

import (
	"bytes"
	"os"
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

func TestRunInvalidFlagDoesNotLogTwice(t *testing.T) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := run([]string{"-nope"}, &stdout, &stderr)
	if exitCode != 1 {
		t.Fatalf("expected invalid flag to exit 1, got %d", exitCode)
	}

	if stdout.Len() != 0 {
		t.Fatalf("expected invalid flag to keep stdout empty, got %q", stdout.String())
	}

	got := stderr.String()
	if !strings.Contains(got, "Usage of gocovmerge:") {
		t.Fatalf("expected usage text on stderr, got %q", got)
	}

	if count := strings.Count(got, "flag provided but not defined: -nope"); count != 1 {
		t.Fatalf("expected single invalid flag diagnostic, got %d in %q", count, got)
	}

	if strings.Contains(got, "level=ERROR") {
		t.Fatalf("expected invalid flag not to be logged, got %q", got)
	}
}

func TestRunParseErrorDoesNotTruncateOutputFile(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "merged.out")
	const previous = "previous data\n"

	if err := os.WriteFile(out, []byte(previous), 0o600); err != nil {
		t.Fatalf("failed to seed output file: %v", err)
	}

	missing := filepath.Join(dir, "missing.out")

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := run([]string{"-o", out, missing}, &stdout, &stderr)
	if exitCode != 1 {
		t.Fatalf("expected parse failure to exit 1, got %d", exitCode)
	}

	got, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	if string(got) != previous {
		t.Fatalf("expected output file to be preserved, got %q", string(got))
	}
}

func TestRunCanUseInputFileAsOutput(t *testing.T) {
	dir := t.TempDir()
	cover := filepath.Join(dir, "cover.out")
	const profile = "mode: set\nfoo.go:1.1,1.2 1 1\n"

	if err := os.WriteFile(cover, []byte(profile), 0o600); err != nil {
		t.Fatalf("failed to seed profile: %v", err)
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	exitCode := run([]string{"-o", cover, cover}, &stdout, &stderr)
	if exitCode != 0 {
		t.Fatalf("expected run to succeed, got %d with stderr %q", exitCode, stderr.String())
	}

	got, err := os.ReadFile(cover)
	if err != nil {
		t.Fatalf("failed to read merged profile: %v", err)
	}

	if string(got) != profile {
		t.Fatalf("expected merged profile to be preserved, got %q", string(got))
	}

	if stdout.Len() != 0 {
		t.Fatalf("expected file output to keep stdout empty, got %q", stdout.String())
	}

	if stderr.Len() != 0 {
		t.Fatalf("expected successful run to keep stderr empty, got %q", stderr.String())
	}
}

func TestRunDirectoryInputExcludesExistingOutputFile(t *testing.T) {
	dir := t.TempDir()
	input := filepath.Join(dir, "a.out")
	output := filepath.Join(dir, "merged.out")
	const profile = "mode: count\nfoo.go:1.1,1.2 1 1\n"

	if err := os.WriteFile(input, []byte(profile), 0o600); err != nil {
		t.Fatalf("failed to seed input profile: %v", err)
	}

	for i := range 2 {
		var stdout bytes.Buffer
		var stderr bytes.Buffer

		exitCode := run([]string{"-d", dir, "-o", output}, &stdout, &stderr)
		if exitCode != 0 {
			t.Fatalf("run %d expected success, got %d with stderr %q", i+1, exitCode, stderr.String())
		}

		if stdout.Len() != 0 {
			t.Fatalf("run %d expected file output to keep stdout empty, got %q", i+1, stdout.String())
		}
	}

	got, err := os.ReadFile(output)
	if err != nil {
		t.Fatalf("failed to read merged output: %v", err)
	}

	if string(got) != profile {
		t.Fatalf("expected output file to be excluded from later runs, got %q", string(got))
	}
}
