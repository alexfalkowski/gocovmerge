.PHONY: vendor

include bin/build/make/go.mak

## Build release binary.
build:
	go build -mod vendor gocovmerge.go
