include bin/build/make/go.mak
include bin/build/make/git.mak

# Build release binary.
build:
	go build -mod vendor gocovmerge.go
