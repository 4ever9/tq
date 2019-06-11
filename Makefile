
SHELL := /bin/bash
CURRENT_PATH = $(shell pwd)
APP_NAME = tq
APP_VERSION = 0.1.0

# build with verison infos
VERSION_DIR = github.com/4ever9/${APP_NAME}
BUILD_DATE = $(shell date +%FT%T)
GIT_COMMIT = $(shell git log --pretty=format:'%h' -n 1)
GIT_BRANCH = $(shell git rev-parse --abbrev-ref HEAD)

LDFLAGS += -X "${VERSION_DIR}.BuildDate=${BUILD_DATE}"
LDFLAGS += -X "${VERSION_DIR}.CurrentCommit=${GIT_COMMIT}"
LDFLAGS += -X "${VERSION_DIR}.CurrentBranch=${GIT_BRANCH}"
LDFLAGS += -X "${VERSION_DIR}.CurrentVersion=${APP_VERSION}"

help: Makefile
	@echo "Choose a command run:"
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/    /'

## make test: Run go unittest
test:
	export GO111MODULE=on && go test ./... -count=1

## make test-cover: Test project with cover
test-cover:
	export GO111MODULE=on && go test -coverprofile cover.out ./...
	go tool cover -html=cover.out -o cover.html

## make install: Go install the project
install:
	export GO111MODULE=on && go install -ldflags '${LDFLAGS}' ./cmd/${APP_NAME}
	@echo "Build tq successfully!"

## make build-all: Go build all os executable file
build-all: build-darwin build-linux build-windows
	@echo "build all binary for darwin, linux and windows"

## make build-linux: Go build linux executable file
build-linux:
	export GO111MODULE=on && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build '-ldflags=${LDFLAGS}' -o bin/${APP_NAME}-linux ./cmd/${APP_NAME}

## make build-windows: Go build windows executable file
build-windows:
	export GO111MODULE=on && CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build '-ldflags=${LDFLAGS}' -o bin/${APP_NAME}-windows ./cmd/${APP_NAME}

## make build-darwin: Go build darwin executable file
build-darwin:
	export GO111MODULE=on && CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build '-ldflags=${LDFLAGS}' -o bin/${APP_NAME}-darwin ./cmd/${APP_NAME}