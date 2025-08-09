.PHONY: default test build run clean deps test-unit test-integration test-coverage test-watch health test-fresh test-race coverage-html

# Default target
default: test

# Test only the test package (integration tests)
test-integration:
	go test -v ./test/...

# Test only unit tests (excluding integration tests)
test-unit:
	go test -v ./internal/... -short

# Test everything with proper setup
test:
	@echo "Running integration tests..."
	go test -v ./test/...
	@echo "Running unit tests..."
	go test -v ./internal/... -short

# Test with coverage
test-coverage:
	go test -v -covermode=atomic -coverprofile=coverage.out ./internal/... ./test/...
	@echo "Coverage summary:"
	go tool cover -func=coverage.out
	@echo "Tip: generate HTML with 'make coverage-html'"

coverage-html:
	@[ -f coverage.out ] || (echo "coverage.out not found. Run 'make test-coverage' first."; exit 1)
	go tool cover -html=coverage.out -o coverage.html
	@echo "Wrote coverage.html"

# Build the application
build:
	go build -o bin/weather-api cmd/main.go

# Run the application
run:
	go run cmd/main.go

# Install dependencies
deps:
	go mod tidy
	go mod download

# Clean build artifacts
clean:
	rm -rf bin/

# Run only integration tests in watch mode (requires entr)
test-watch:
	find ./test -name "*.go" | entr -c make test-integration

# Quick health check
health:
	curl -s http://localhost:8080/health | jq .

# Run tests with verbose output and no cache
test-fresh:
	go clean -testcache
	go test -v ./test/... ./internal/...

test-race:
	go test -race ./internal/... ./test/...
