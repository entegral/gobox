.PHONY: build test clean doc

# The name of your binary
BINARY_NAME=myapp

# The Go path
GOPATH=$(shell go env GOPATH)

# The Go files
GOFILES=$(wildcard *.go)

all: build

build: $(GOFILES)
	go build -o $(BINARY_NAME) $(GOFILES)

test:
	go test -v ./...

clean:
	go clean
	rm -f $(BINARY_NAME)

run: build
	./$(BINARY_NAME)

doc:
	godoc -http=:6061