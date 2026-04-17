# AGENTS.md

This repository is a small Go CLI tool (`gocovmerge`) that merges multiple Go coverage profiles (from `go test -coverprofile`) into a single profile.

## Shared skill

Use the shared `coding-standards` skill from `./bin/skills/coding-standards`
for cross-repository coding, review, testing, documentation, and PR
conventions. Treat this `AGENTS.md` as the repo-specific companion to that
skill.

## Repo map

- `main.go`: CLI entrypoint. Parses flags, resolves input files, constructs the output writer, runs the merge, and finalizes file output.
- `internal/cmd`: high-level merge pipeline (`cmd.Run(files, out)`).
- `internal/cover`: parse/merge/write logic on top of `golang.org/x/tools/cover`.
- `internal/flag`: CLI flags and input selection.
- `internal/path`: recursive file discovery with regexp filtering and optional excluded path.
- `internal/io`: stdout writer or buffered file writer committed on `Close()`.
- `internal/log`: `slog` logger constructor for CLI diagnostics.
- `bin/`: git submodule with shared build/CI tooling. Most `make` targets depend on it.

## Tooling

- Go version: `go 1.26.0`.
- The repo uses vendoring (`-mod vendor`).
- If the `bin/` submodule is missing, `make` targets will fail.

Useful commands:

- `make dep`: `go mod download`, `go mod tidy`, `go mod vendor`
- `make build`: build `gocovmerge`
- `make lint`: field alignment + `golangci-lint`
- `make sec`: `govulncheck -show verbose -test ./...`
- `make specs`: race-enabled test run with coverage via `gotestsum`
- `make coverage`: post-process `test/reports/profile.cov`
- `go test -mod=vendor ./...`: quick local test run

Tests currently live in:

- `main_test.go`
- `internal/cover/cover_test.go`
- `internal/test`: shared test helpers and scenario scaffolding used by the test suites above.

## CI

CircleCI config: `.circleci/config.yml`.

The main build job runs:

- `git submodule sync` / `git submodule update --init`
- `make source-key`
- `make clean`
- `make dep`
- `make clean`
- `make lint`
- `make sec`
- `make specs`
- `make build`
- `make coverage`

Artifacts and test results are stored from `test/reports`.

## Merge gotchas

- All profiles must use the same coverage mode: `set`, `count`, or `atomic`.
- Merging is per `Profile.FileName`.
- Same-position blocks must also agree on `NumStmt`.
- Overlapping or otherwise incompatible blocks fail the merge.
- When `-d` is used and `-o` points inside the scanned directory, the output path is excluded from discovered inputs so rerunning the same command does not re-merge the old output.
- When `-o` is used, file output is buffered and only written on successful completion.
- Diagnostics go to stderr; stdout is reserved for merged profile output.
- Only merge profiles generated from the same source revision.

## Conventions

- Prefer wrapped errors with context via `fmt.Errorf("...: %w", err)`.
- Sentinel errors in `internal/cover`: `ErrInvalidMode`, `ErrEmptyProfiles`.
- Go files and Makefiles use tabs (`.editorconfig`).
