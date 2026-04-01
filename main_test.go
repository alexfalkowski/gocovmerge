package main

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/alexfalkowski/gocovmerge/v2/internal/test"
)

var errWriteFailed = errors.New("write failed")

func TestRunScenarios(t *testing.T) {
	for _, tt := range runScenarioCases {
		t.Run(tt.Name, func(t *testing.T) {
			test.RunScenarioCase(t, run, tt)
		})
	}
}

func TestMainEntrypoint(t *testing.T) {
	t.Helper()

	if os.Getenv("GO_WANT_MAIN") == "1" {
		os.Args = []string{"gocovmerge", os.Getenv("GO_WANT_MAIN_PROFILE")}
		main()
		return
	}

	dir := t.TempDir()
	profilePath := test.WriteProfileFile(t, dir, "a.out", test.TextProfile("set",
		test.TextBlock("foo.go", 1, 1, 1, 2, 1, 1),
	))

	executable, err := os.Executable()
	if err != nil {
		t.Fatalf("failed to resolve test binary path: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, executable, "-test.run=^TestMainEntrypoint$")
	cmd.Env = append(os.Environ(),
		"GO_WANT_MAIN=1",
		"GO_WANT_MAIN_PROFILE="+profilePath,
	)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		t.Fatalf("expected main to succeed, got %v with stderr %q", err, stderr.String())
	}

	want := test.TextProfile("set", test.TextBlock("foo.go", 1, 1, 1, 2, 1, 1))
	if stdout.String() != want {
		t.Fatalf("expected stdout %q, got %q", want, stdout.String())
	}

	if stderr.Len() != 0 {
		t.Fatalf("expected stderr to be empty, got %q", stderr.String())
	}
}

var runScenarioCases = []test.RunScenario{
	{
		Name:            "help exits zero",
		Setup:           func(_ *testing.T, _ string) []string { return []string{"-help"} },
		WantExit:        0,
		WantStdout:      "",
		WantStderrEmpty: false,
		WantStderrContains: []string{
			"Usage of gocovmerge:",
		},
		WantStderrExcludes: []string{
			"flag: help requested",
		},
	},
	{
		Name:            "invalid flag writes a single diagnostic",
		Setup:           func(_ *testing.T, _ string) []string { return []string{"-nope"} },
		WantExit:        1,
		WantStdout:      "",
		WantStderrEmpty: false,
		WantStderrContains: []string{
			"Usage of gocovmerge:",
			"flag provided but not defined: -nope",
		},
		WantStderrExcludes: []string{
			"level=ERROR",
		},
		Check: func(t *testing.T, _, _, stderr string) {
			t.Helper()

			if count := strings.Count(stderr, "flag provided but not defined: -nope"); count != 1 {
				t.Fatalf("expected one invalid flag diagnostic, got %d in %q", count, stderr)
			}
		},
	},
	{
		Name:            "empty input fails without writing stdout",
		WantExit:        1,
		WantStdout:      "",
		WantStderrEmpty: false,
		WantStderrContains: []string{
			"failed to write profile: empty profiles",
		},
	},
	{
		Name: "missing input file reports parse error",
		Setup: func(_ *testing.T, dir string) []string {
			return []string{filepath.Join(dir, "missing.out")}
		},
		WantExit:        1,
		WantStdout:      "",
		WantStderrEmpty: false,
		WantStderrContains: []string{
			"failed to parse profiles",
			"missing.out",
		},
	},
	{
		Name: "invalid directory pattern reports discovery error",
		Setup: func(_ *testing.T, dir string) []string {
			return []string{"-d", dir, "-p", "["}
		},
		WantExit:        1,
		WantStdout:      "",
		WantStderrEmpty: false,
		WantStderrContains: []string{
			"error parsing regexp",
		},
	},
	{
		Name: "missing directory reports discovery error",
		Setup: func(_ *testing.T, dir string) []string {
			return []string{"-d", filepath.Join(dir, "missing")}
		},
		WantExit:        1,
		WantStdout:      "",
		WantStderrEmpty: false,
		WantStderrContains: []string{
			"missing",
			"no such file or directory",
		},
	},
	{
		Name: "single file merges to stdout",
		Setup: func(t *testing.T, dir string) []string {
			t.Helper()

			file := test.WriteProfileFile(t, dir, "a.out", test.TextProfile("set",
				test.TextBlock("foo.go", 1, 1, 1, 2, 1, 1),
			))

			return []string{file}
		},
		WantExit:        0,
		WantStdout:      test.TextProfile("set", test.TextBlock("foo.go", 1, 1, 1, 2, 1, 1)),
		WantStderrEmpty: true,
	},
	{
		Name: "set mode merges counts with bitwise or",
		Setup: func(t *testing.T, dir string) []string {
			t.Helper()

			first := test.WriteProfileFile(t, dir, "a.out", test.TextProfile("set",
				test.TextBlock("foo.go", 1, 1, 1, 2, 1, 0),
			))
			second := test.WriteProfileFile(t, dir, "b.out", test.TextProfile("set",
				test.TextBlock("foo.go", 1, 1, 1, 2, 1, 1),
			))

			return []string{first, second}
		},
		WantExit:        0,
		WantStdout:      test.TextProfile("set", test.TextBlock("foo.go", 1, 1, 1, 2, 1, 1)),
		WantStderrEmpty: true,
	},
	{
		Name: "count mode merges counts by addition",
		Setup: func(t *testing.T, dir string) []string {
			t.Helper()

			first := test.WriteProfileFile(t, dir, "a.out", test.TextProfile("count",
				test.TextBlock("foo.go", 1, 1, 1, 2, 1, 1),
			))
			second := test.WriteProfileFile(t, dir, "b.out", test.TextProfile("count",
				test.TextBlock("foo.go", 1, 1, 1, 2, 1, 2),
			))

			return []string{first, second}
		},
		WantExit:        0,
		WantStdout:      test.TextProfile("count", test.TextBlock("foo.go", 1, 1, 1, 2, 1, 3)),
		WantStderrEmpty: true,
	},
	{
		Name: "non overlapping blocks are inserted in order",
		Setup: func(t *testing.T, dir string) []string {
			t.Helper()

			first := test.WriteProfileFile(t, dir, "a.out", test.TextProfile("set",
				test.TextBlock("foo.go", 1, 1, 1, 2, 1, 1),
				test.TextBlock("foo.go", 3, 1, 3, 2, 1, 1),
			))
			second := test.WriteProfileFile(t, dir, "b.out", test.TextProfile("set",
				test.TextBlock("foo.go", 2, 1, 2, 2, 1, 1),
			))

			return []string{first, second}
		},
		WantExit: 0,
		WantStdout: test.TextProfile("set",
			test.TextBlock("foo.go", 1, 1, 1, 2, 1, 1),
			test.TextBlock("foo.go", 2, 1, 2, 2, 1, 1),
			test.TextBlock("foo.go", 3, 1, 3, 2, 1, 1),
		),
		WantStderrEmpty: true,
	},
	{
		Name: "directory input filters files by pattern",
		Setup: func(t *testing.T, dir string) []string {
			t.Helper()

			test.WriteProfileFile(t, dir, "a.out", test.TextProfile("set",
				test.TextBlock("a.go", 1, 1, 1, 2, 1, 1),
			))
			test.WriteProfileFile(t, dir, "b.out", test.TextProfile("set",
				test.TextBlock("b.go", 2, 1, 2, 2, 1, 1),
			))
			test.WriteTextFile(t, filepath.Join(dir, "ignored.txt"), "not a coverage profile\n")

			return []string{"-d", dir, "-p", `\.out$`}
		},
		WantExit: 0,
		WantStdout: test.TextProfile("set",
			test.TextBlock("a.go", 1, 1, 1, 2, 1, 1),
			test.TextBlock("b.go", 2, 1, 2, 2, 1, 1),
		),
		WantStderrEmpty: true,
	},
	{
		Name: "mixed modes fail during merge",
		Setup: func(t *testing.T, dir string) []string {
			t.Helper()

			first := test.WriteProfileFile(t, dir, "a.out", test.TextProfile("set",
				test.TextBlock("foo.go", 1, 1, 1, 2, 1, 1),
			))
			second := test.WriteProfileFile(t, dir, "b.out", test.TextProfile("atomic",
				test.TextBlock("bar.go", 1, 1, 1, 2, 1, 1),
			))

			return []string{first, second}
		},
		WantExit:        1,
		WantStdout:      "",
		WantStderrEmpty: false,
		WantStderrContains: []string{
			"failed to add profile: invalid profiles merge with different modes",
		},
	},
	{
		Name: "overlapping blocks fail during merge",
		Setup: func(t *testing.T, dir string) []string {
			t.Helper()

			first := test.WriteProfileFile(t, dir, "a.out", test.TextProfile("set",
				test.TextBlock("foo.go", 1, 1, 1, 5, 1, 1),
			))
			second := test.WriteProfileFile(t, dir, "b.out", test.TextProfile("set",
				test.TextBlock("foo.go", 1, 4, 1, 10, 1, 1),
			))

			return []string{first, second}
		},
		WantExit:        1,
		WantStdout:      "",
		WantStderrEmpty: false,
		WantStderrContains: []string{
			"failed to add profile",
			"overlap before",
		},
	},
	{
		Name: "same block with different numstmt fails during merge",
		Setup: func(t *testing.T, dir string) []string {
			t.Helper()

			first := test.WriteProfileFile(t, dir, "a.out", test.TextProfile("count",
				test.TextBlock("foo.go", 1, 1, 1, 2, 1, 1),
			))
			second := test.WriteProfileFile(t, dir, "b.out", test.TextProfile("count",
				test.TextBlock("foo.go", 1, 1, 1, 2, 2, 1),
			))

			return []string{first, second}
		},
		WantExit:        1,
		WantStdout:      "",
		WantStderrEmpty: false,
		WantStderrContains: []string{
			"failed to add profile",
			"inconsistent NumStmt",
		},
	},
	{
		Name: "stdout write failure is reported",
		Setup: func(t *testing.T, dir string) []string {
			t.Helper()

			file := test.WriteProfileFile(t, dir, "a.out", test.TextProfile("set",
				test.TextBlock("foo.go", 1, 1, 1, 2, 1, 1),
			))

			return []string{file}
		},
		Stdout:          test.FailingWriter{Err: errWriteFailed},
		WantExit:        1,
		SkipStdoutCheck: true,
		WantStderrEmpty: false,
		WantStderrContains: []string{
			"failed to write profile",
			errWriteFailed.Error(),
		},
	},
}

func TestRunFileOutputScenarios(t *testing.T) {
	for _, tt := range fileOutputScenarioCases {
		t.Run(tt.Name, func(t *testing.T) {
			test.RunFileOutputScenarioCase(t, run, tt)
		})
	}
}

var fileOutputScenarioCases = []test.FileOutputScenario{
	{
		Name: "parse errors do not truncate the output file",
		Setup: func(t *testing.T, dir string) []string {
			t.Helper()

			out := filepath.Join(dir, "merged.out")
			test.WriteTextFile(t, out, "previous data\n")

			return []string{"-o", out, filepath.Join(dir, "missing.out")}
		},
		WantExit:        1,
		WantStderrEmpty: false,
		WantStderrContains: []string{
			"failed to parse profiles",
		},
		Check: func(t *testing.T, dir, stdout, _ string) {
			t.Helper()

			if stdout != "" {
				t.Fatalf("expected file output to keep stdout empty, got %q", stdout)
			}

			got, err := os.ReadFile(filepath.Join(dir, "merged.out"))
			if err != nil {
				t.Fatalf("failed to read output file: %v", err)
			}

			if string(got) != "previous data\n" {
				t.Fatalf("expected output file to be preserved, got %q", string(got))
			}
		},
	},
	{
		Name: "the same path can be used as both input and output",
		Setup: func(t *testing.T, dir string) []string {
			t.Helper()

			cover := test.WriteProfileFile(t, dir, "cover.out", test.TextProfile("set",
				test.TextBlock("foo.go", 1, 1, 1, 2, 1, 1),
			))

			return []string{"-o", cover, cover}
		},
		WantExit:        0,
		WantStderrEmpty: true,
		Check: func(t *testing.T, dir, stdout, _ string) {
			t.Helper()

			if stdout != "" {
				t.Fatalf("expected file output to keep stdout empty, got %q", stdout)
			}

			got, err := os.ReadFile(filepath.Join(dir, "cover.out"))
			if err != nil {
				t.Fatalf("failed to read merged profile: %v", err)
			}

			want := test.TextProfile("set", test.TextBlock("foo.go", 1, 1, 1, 2, 1, 1))
			if string(got) != want {
				t.Fatalf("expected merged profile %q, got %q", want, string(got))
			}
		},
	},
	{
		Name: "successful file output is committed on close",
		Setup: func(t *testing.T, dir string) []string {
			t.Helper()

			first := test.WriteProfileFile(t, dir, "a.out", test.TextProfile("count",
				test.TextBlock("foo.go", 1, 1, 1, 2, 1, 1),
			))
			second := test.WriteProfileFile(t, dir, "b.out", test.TextProfile("count",
				test.TextBlock("foo.go", 1, 1, 1, 2, 1, 2),
			))

			return []string{"-o", filepath.Join(dir, "merged.out"), first, second}
		},
		WantExit:        0,
		WantStderrEmpty: true,
		Check: func(t *testing.T, dir, stdout, _ string) {
			t.Helper()

			if stdout != "" {
				t.Fatalf("expected file output to keep stdout empty, got %q", stdout)
			}

			got, err := os.ReadFile(filepath.Join(dir, "merged.out"))
			if err != nil {
				t.Fatalf("failed to read output file: %v", err)
			}

			want := test.TextProfile("count", test.TextBlock("foo.go", 1, 1, 1, 2, 1, 3))
			if string(got) != want {
				t.Fatalf("expected merged output %q, got %q", want, string(got))
			}
		},
	},
	{
		Name: "output close failures are reported after a successful merge",
		Setup: func(t *testing.T, dir string) []string {
			t.Helper()

			input := test.WriteProfileFile(t, dir, "a.out", test.TextProfile("set",
				test.TextBlock("foo.go", 1, 1, 1, 2, 1, 1),
			))

			return []string{"-o", filepath.Join(dir, "missing", "merged.out"), input}
		},
		WantExit:        1,
		WantStderrEmpty: false,
		WantStderrContains: []string{
			"no such file or directory",
		},
		Check: func(t *testing.T, dir, stdout, _ string) {
			t.Helper()

			if stdout != "" {
				t.Fatalf("expected file output to keep stdout empty, got %q", stdout)
			}

			if _, err := os.Stat(filepath.Join(dir, "missing", "merged.out")); !errors.Is(err, os.ErrNotExist) {
				t.Fatalf("expected output file to be absent after close failure, got %v", err)
			}
		},
	},
}

func TestRunDirectoryInputExcludesExistingOutputFile(t *testing.T) {
	dir := t.TempDir()
	input := test.WriteProfileFile(t, dir, "a.out", test.TextProfile("count",
		test.TextBlock("foo.go", 1, 1, 1, 2, 1, 1),
	))
	output := filepath.Join(dir, "merged.out")

	for i := range 2 {
		exitCode, stdout, stderr := test.ExecuteRun(t, run, []string{"-d", dir, "-o", output}, nil)
		if exitCode != 0 {
			t.Fatalf("run %d expected success, got %d with stderr %q", i+1, exitCode, stderr)
		}

		if stdout != "" {
			t.Fatalf("run %d expected file output to keep stdout empty, got %q", i+1, stdout)
		}
	}

	got, err := os.ReadFile(output)
	if err != nil {
		t.Fatalf("failed to read merged output: %v", err)
	}

	want, err := os.ReadFile(input)
	if err != nil {
		t.Fatalf("failed to read input profile: %v", err)
	}

	if string(got) != string(want) {
		t.Fatalf("expected output file to be excluded from later runs, got %q", string(got))
	}
}
