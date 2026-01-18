# Ralph Codex - Makefile

.PHONY: build install test test-integration lint clean help all

# Variables
BINARY=ralph
CMD_PATH=./cmd/ralph
GO=go
GOFLAGS=-v

# Default target
all: build test lint

## build: Build the ralph binary
build:
	$(GO) build $(GOFLAGS) -o $(BINARY) $(CMD_PATH)

## install: Install ralph to $GOPATH/bin or $HOME/go/bin
install: install-bin install-templates install-sdk

## install-bin: Install Go binary
install-bin:
	$(GO) install $(GOFLAGS) $(CMD_PATH)

## install-templates: Install project templates to ~/.ralph/templates
install-templates:
	@echo "Installing templates to ~/.ralph/templates..."
	@mkdir -p ~/.ralph/templates
	@cp -r templates/* ~/.ralph/templates/
	@echo "Templates installed successfully"

## install-sdk: Install Codex SDK runner and npm dependencies
install-sdk:
	@echo "Installing Codex SDK runner..."
	@mkdir -p ~/.ralph/bin
	@cp src/codex_runner.js ~/.ralph/bin/
	@chmod +x ~/.ralph/bin/codex_runner.js
	@if [ -f package.json ]; then \
		npm install; \
		echo "Codex SDK dependencies installed"; \
	else \
		echo "Warning: package.json not found, skipping npm install"; \
	fi
	@echo "Codex SDK runner installed to ~/.ralph/bin/codex_runner.js"

## test: Run all tests
test:
	$(GO) test ./...

## test-integration: Run integration tests
test-integration:
	$(GO) test -tags=integration ./...

## test-verbose: Run tests with verbose output
test-verbose:
	$(GO) test -v ./...

## test-coverage: Run tests with coverage
test-coverage:
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

## lint: Run golangci-lint
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not installed. Install with: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.50.1"; \
		exit 1; \
	fi

## fmt: Format code with gofmt
fmt:
	$(GO) fmt ./...

## vet: Run go vet
vet:
	$(GO) vet ./...

## clean: Clean build artifacts
clean:
	rm -f $(BINARY)
	rm -f coverage.out coverage.html
	rm -f /tmp/ralph-test-*

## run: Build and run ralph
run: build
	./$(BINARY) --help

## setup-test: Create test project
setup-test:
	./$(BINARY) --command setup --name test-project

## import-test: Test import command
import-test: build
	./$(BINARY) --command import --source test.md

## deps: Download dependencies
deps:
	$(GO) mod download
	$(GO) mod tidy

## deps-update: Update dependencies
deps-update:
	$(GO) get -u ./...
	$(GO) mod tidy

## help: Show this help message
help:
	@echo "Ralph Codex - Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  all              Build, test, and lint"
	@echo "  build            Build the ralph binary"
	@echo "  install          Install ralph, templates, and SDK runner"
	@echo "  install-bin      Install Go binary only"
	@echo "  install-templates Install project templates only"
	@echo "  install-sdk      Install Codex SDK runner only"
	@echo "  test             Run all tests"
	@echo "  test-integration Run integration tests"
	@echo "  test-verbose     Run tests with verbose output"
	@echo "  test-coverage    Run tests with coverage report"
	@echo "  lint             Run golangci-lint"
	@echo "  fmt              Format code with gofmt"
	@echo "  vet              Run go vet"
	@echo "  clean            Clean build artifacts"
	@echo "  run              Build and run ralph"
	@echo "  setup-test       Create test project"
	@echo "  import-test      Test import command"
	@echo "  deps             Download dependencies"
	@echo "  deps-update      Update dependencies"
	@echo "  help             Show this help message"
