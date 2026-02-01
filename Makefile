.PHONY: build install test lint clean

# Build binary
build:
	go build -o bin/sopsctl ./cmd/sopsctl

# Install to GOPATH/bin
install:
	go install ./cmd/sopsctl

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
	./bin/sopsctl $(ARGS)
