BINARY_NAME=git-esque
PKG=github.com/rayhanadev/git-esque

all: build

build:
	go build -o $(BINARY_NAME)
	chmod +x $(BINARY_NAME)

run: build
	./$(BINARY_NAME)

test:
	go test ./...

clean:
	go clean
	rm -f $(BINARY_NAME)

fmt:
	go fmt ./...

lint:
	golint ./...

vet:
	go vet ./...

deps:
	go mod tidy
	go mod vendor

.PHONY: all build run test clean fmt lint vet deps
