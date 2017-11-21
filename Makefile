SOURCES  ?= $(shell find . -name "*.go" | grep -v vendor/)
PACKAGES ?= $(shell go list ./...)
GOTOOLS  ?= github.com/GeertJohan/fgt \
			golang.org/x/tools/cmd/goimports \
			github.com/golang/lint/golint \
			github.com/kisielk/errcheck \
			honnef.co/go/tools/cmd/staticcheck

all: lint test

dependencies:
	go get -t ./...

test: dependencies
	go test -race $(PACKAGES)

test-verbose: dependencies
	go test -race -v $(PACKAGES)

lint: tools
	fgt go fmt $(PACKAGES)
	fgt goimports -w $(SOURCES)
	echo $(PACKAGES) | xargs -L1 fgt golint
	fgt go vet $(PACKAGES)
	fgt errcheck -ignore Close $(PACKAGES)
	staticcheck $(PACKAGES)
.SILENT: lint

tools:
	go get $(GOTOOLS)
.SILENT: tools

build: dependencies
	go build -o bin/numserver

build-linux: dependencies
	GOOS=linux GOARCH=amd64 go build -o bin/numserver
