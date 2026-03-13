.PHONY: build build-all clean install test lint

BINARY_NAME=wzjk-cli
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS=-ldflags "-X wzjkctl/internal/version.Version=$(VERSION)"

build:
	go build $(LDFLAGS) -o bin/$(BINARY_NAME) main.go

build-all:
	mkdir -p bin
	# Linux AMD64
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-linux-amd64 main.go
	# Linux ARM64
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-linux-arm64 main.go
	# macOS AMD64
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-amd64 main.go
	# macOS ARM64 (Apple Silicon)
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-arm64 main.go
	# Windows AMD64
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/$(BINARY_NAME)-windows-amd64.exe main.go

install: build
	cp bin/$(BINARY_NAME) $(GOPATH)/bin/ 2>/dev/null || cp bin/$(BINARY_NAME) ~/go/bin/ 2>/dev/null || echo "Please copy bin/$(BINARY_NAME) to your PATH"

clean:
	rm -rf bin/

test:
	go test ./...

lint:
	golangci-lint run

deps:
	go mod download
	go mod tidy

run: build
	./bin/$(BINARY_NAME)
