# AGENTS.md

This repository is a small Go CLI tool (`gocovmerge`) that merges multiple Go coverage profiles (from `go test -coverprofile`) into a single profile.

## Quick orientation

- Entry point: `main.go` (parses flags, resolves input files, opens output writer, runs merge).
- Main logic: `internal/cover` (parsing/merging/writing profiles) and `internal/cmd` (orchestrates parse/merge/write).
- Build/lint tooling is largely provided via the `bin/` git submodule (see `.gitmodules`).

## Repository layout

- `main.go`: wires `internal/flag`, `internal/path`, `internal/io`, `internal/cmd`, `internal/log`.
- `internal/cmd/cmd.go`: `cmd.Run(files []string, out io.Writer) error` (high-level merge pipeline).
- `internal/cover/cover.go`: coverage merge implementation on top of `golang.org/x/tools/cover`.
- `internal/flag/flag.go`: CLI flags (`-o`, `-d`, `-p`) and positional args.
- `internal/path/path.go`: recursive file discovery with optional regexp filter.
- `internal/io/io.go`: output writer (file or stdout).
- `internal/log/log.go`: `slog`-based logger with `Fatal` exiting process.
- `vendor/`: vendored dependencies (the build uses `-mod vendor`).
- `bin/`: git submodule with shared build/CI tooling.

## Tooling and dependencies

### Go version

- `go.mod` specifies `go 1.26.0`.

### Git submodule (required for most `make` targets)

The root `Makefile` includes `bin/build/make/*.mak`, which live in the `bin/` submodule.

- Submodule definition: `.gitmodules` (`bin` uses an SSH URL).
- CI runs: `git submodule sync` and `git submodule update --init`.

If the submodule is missing/uninitialized, `make` will fail due to missing included `.mak` files.

## Common commands

All commands below are taken from the checked-in Makefiles/CI config.

### Setup deps / vendoring

- `make dep`
  - runs: `go mod download`, `go mod tidy`, `go mod vendor` (from `bin/build/make/go.mak`).

### Build

- `make build`
  - runs: `go build -race -mod vendor -o gocovmerge` (root `Makefile`).

### Lint

- `make lint`
  - runs `field-alignment` and `golangci-lint` (from `bin/build/make/go.mak`).
  - `golangci-lint` is invoked via `$(PWD)/bin/build/go/lint run --timeout 5m`.

- `make fix-lint`
  - runs the same checks with `--fix` where supported.

Lint configuration: `.golangci.yml`.

### Formatting

- `make format`
  - runs: `go fmt ./...` (from `bin/build/make/go.mak`).

### Security

- `make sec`
  - runs: `govulncheck -show verbose -test ./...` (from `bin/build/make/go.mak`).

### Tests

- There are currently no `*_test.go` files in this repository.
- CI does not run `go test` directly; it runs `make lint`, `make sec`, and `make build`.

If you add tests, a reasonable local command consistent with this repo’s vendoring is:

- `go test -mod=vendor ./...`

(Observed that build uses `-mod vendor` and `vendor/` exists.)

## CI (CircleCI)

CircleCI configuration: `.circleci/config.yml`.

The build job runs (in order):

- `git submodule sync` / `git submodule update --init`
- `make source-key`
- `make clean`
- `make dep`
- `make lint`
- `make sec`
- `make build`

## Code conventions and patterns

### Error handling

- Prefer wrapping errors with context using `fmt.Errorf("...: %w", err)` (see `internal/cmd/cmd.go`).
- Sentinel errors exist in `internal/cover/cover.go`:
  - `cover.ErrInvalidMode`
  - `cover.ErrEmptyProfiles`

### Coverage merge behavior (important gotchas)

- Profiles are merged per `Profile.FileName` and must use the same `mode` (`set`, `count`, `atomic`), otherwise `ErrInvalidMode`.
- Overlapping or incompatible blocks return errors (e.g. `overlap merge`, `overlap before`, `overlap after`).
- The repository README notes: only merge profiles generated from the same source code, and the tool exits non-zero when it cannot merge.

### Logging

- Logger is `log/slog` with a text handler writing to stdout (`internal/log/log.go`).
- `Logger.Fatal` logs at error level and then `os.Exit(1)`.

### Formatting/style

- `.editorconfig`:
  - Go files use tabs (`indent_style = tab`).
  - Makefiles use tabs.
- `.golangci.yml` enables formatters (`gofmt`, `gofumpt`, `goimports`, `gci`) and sets some linter thresholds (e.g. `lll` line length 150).

## Release

GoReleaser configuration: `.goreleaser.yml`.

- Runs `go mod tidy` as a pre-hook.
- Builds with `CGO_ENABLED=0`.
