# Accepted design for the CLI

Part of the [accepted-design index](../accepted-design.md).

## Input selection and output exclusion

Directory mode (`-d`) takes precedence over positional paths. It discovers
files recursively and uses `-p`, when supplied, as a Go regular expression
against each walked path. Callers use that filter to keep unrelated files out
of the merge.

With `-d` and `-o`, discovery compares canonical paths to exclude the output file
and keep reruns from consuming the previous result. This handles relative paths
and symlinked path spellings. For a missing path, discovery uses the absolute
path and resolves symlinks in its parent when possible. Other canonicalization
errors abort discovery because skipping a file could silently produce
incomplete coverage.

[Flag handling](../../../internal/flag/flag.go) ·
[Path discovery and exclusion](../../../internal/path/path.go) ·
[Path tests](../../../internal/path/path_test.go) ·
[CLI tests](../../../main_test.go)

## Output and diagnostics

Stdout carries merged coverage data. Stderr carries diagnostics and flag usage.
Help exits successfully. Flag parse and runtime errors exit with status `1`.

File output (`-o`) buffers the merged profile until the merge succeeds. The
writer then stages the data in a temporary file in the existing destination
directory and renames it over the destination. Staging protects an existing
result from truncation or partial writes during a failed merge or write. In
positional mode, callers can use an input path as the output path because the
merge finishes reading the inputs before file replacement.

[CLI entry point](../../../main.go) ·
[File writer](../../../internal/io/io.go) ·
[File writer tests](../../../internal/io/io_test.go) ·
[CLI tests](../../../main_test.go)
