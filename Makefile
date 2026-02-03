.PHONY: build install test lint clean

# Build binary
build:
	go build -o bin/sopsy ./cmd/sopsy

# Install to GOPATH/bin
install:
	go install ./cmd/sopsy

# Run tests
test:
	go test -v ./...

# Run linter
lint:
	golangci-lint run ./...

# Clean build artifacts
clean:
	rm -rf bin/

# Tidy dependencies
tidy:
	go mod tidy

# Build and run
run: build
	./bin/sopsy $(ARGS)
