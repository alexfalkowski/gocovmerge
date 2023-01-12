.PHONY: vendor

download:
	go mod download

tidy:
	go mod tidy

vendor:
	go mod vendor

get:
	go get $(module)

## Setup go deps
dep: download tidy vendor

## Check outdated deps
outdated:
	go list -u -m -mod=mod -json all | go-mod-outdated -update -direct

## Update go dep
update-dep: get tidy vendor

## Lint all the code
lint:
	golangci-lint run --timeout 5m

## Fix the lint issues in the code (if possible)
fix-lint:
	golangci-lint run --timeout 5m --fix

## Build release binary
build:
	go build -mod vendor gocovmerge.go
