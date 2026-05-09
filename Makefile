.PHONY: help build run dev tell lint lint-fix fmt fmt-check test tidy clean install uninstall install-tools

BINARY      := gopher-gotchi
BIN_DIR     := bin
BIN_PATH    := $(BIN_DIR)/$(BINARY)
INSTALL_DIR := /usr/local/bin
MODULE      := gopher-gotchi
SPECIES     ?= diana

# Default target
help:
	@echo "Available commands:"
	@echo "  make build         - Build the $(BINARY) binary into $(BIN_PATH)"
	@echo "  make run           - Build and run the companion (SPECIES=$(SPECIES))"
	@echo "  make dev           - Run the companion directly via 'go run'"
	@echo "  make tell CMD=...  - Send a command to the running companion (e.g. CMD=hello)"
	@echo "  make install       - Build and install $(BINARY) to $(INSTALL_DIR)"
	@echo "  make uninstall     - Remove $(BINARY) from $(INSTALL_DIR)"
	@echo "  make lint          - Run golangci-lint"
	@echo "  make lint-fix      - Run golangci-lint with auto-fix"
	@echo "  make fmt           - Format code with gofmt and goimports"
	@echo "  make fmt-check     - Check code formatting without making changes"
	@echo "  make test          - Run tests"
	@echo "  make tidy          - Tidy go.mod / go.sum"
	@echo "  make clean         - Remove build artifacts"
	@echo "  make install-tools - Install required development tools"

# Build the binary
build:
	@echo "🔨 Building $(BINARY)..."
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_PATH) .

# Install the binary to /usr/local/bin
install: build
	@echo "📦 Installing $(BINARY) to $(INSTALL_DIR)..."
	@sudo cp $(BIN_PATH) $(INSTALL_DIR)/$(BINARY)
	@echo "✅ Installation complete! You can now run '$(BINARY)' from anywhere."

# Uninstall the binary
uninstall:
	@echo "🗑️  Removing $(BINARY) from $(INSTALL_DIR)..."
	@sudo rm -f $(INSTALL_DIR)/$(BINARY)
	@echo "✅ $(BINARY) uninstalled."

# Run the application (built binary)
run: build
	@./$(BIN_PATH) -species $(SPECIES)

# Start dev server (no build artifact)
dev:
	@echo "🔧 Starting $(BINARY) in dev mode..."
	@go run . -species $(SPECIES)

# Send a command to the running companion via its local HTTP API
# Usage: make tell CMD=hello
tell:
	@if [ -z "$(CMD)" ]; then \
		echo "❌ Usage: make tell CMD=<command>"; \
		exit 1; \
	fi
	@./$(BIN_PATH) tell $(CMD)

# Run linter
lint:
	@echo "🔍 Running golangci-lint..."
	@golangci-lint run --config .golangci.yml

# Run linter with auto-fix
lint-fix:
	@echo "🔧 Running golangci-lint with auto-fix..."
	@golangci-lint run --config .golangci.yml --fix

# Format code
fmt:
	@echo "✨ Formatting code..."
	@gofmt -w .
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -w -local $(MODULE) .; \
	else \
		echo "⚠️  goimports not found. Install with: go install golang.org/x/tools/cmd/goimports@latest"; \
	fi
	@echo "✅ Code formatted!"

# Check formatting without making changes
fmt-check:
	@echo "🔍 Checking code formatting..."
	@UNFORMATTED=$$(gofmt -l .); \
	if [ -n "$$UNFORMATTED" ]; then \
		echo "❌ The following files need formatting:"; \
		echo "$$UNFORMATTED"; \
		exit 1; \
	fi
	@echo "✅ All files are properly formatted!"

# Run tests
test:
	@echo "🧪 Running tests..."
	@go test -v ./...

# Tidy modules
tidy:
	@echo "🧼 Tidying modules..."
	@go mod tidy

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	@rm -rf $(BIN_DIR)/
	@go clean

# Install development tools
install-tools:
	@echo "📦 Installing development tools..."
	@echo "Installing golangci-lint..."
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin; \
	else \
		echo "✅ golangci-lint already installed"; \
	fi
	@echo "Installing goimports..."
	@go install golang.org/x/tools/cmd/goimports@latest
	@echo "✅ All tools installed!"
