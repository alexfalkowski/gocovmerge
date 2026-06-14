include bin/build/make/go.mak
include bin/build/make/git.mak

# Build race-enabled binary.
build:
	@go build -race -mod vendor -o gocovmerge
