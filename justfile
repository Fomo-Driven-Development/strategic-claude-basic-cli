# Strategic Claude Basic CLI - justfile
# Equivalent to Makefile with just's cleaner syntax

# Variables
binary_name := "strategic-claude"
cmd_pkg := "strategic-claude"
build_dir := "bin"

# Version information (can be overridden with environment variables)
version := env_var_or_default("VERSION", "0.1.0")
commit := `git rev-parse --short HEAD 2>/dev/null || echo "dev"`
date := `date -u +"%Y-%m-%dT%H:%M:%SZ"`

# Build flags with version injection
ldflags := "-ldflags \"-X main.version=" + version + " -X main.commit=" + commit + " -X main.date=" + date + "\""

# Default recipe - show help
default: help

# Show available recipes
help:
    @echo "Available recipes:"
    @echo ""
    @echo "Building & Running:"
    @echo "  build         - Build the application"
    @echo "  run           - Build and run the application"
    @echo "  install       - Install the binary to GOPATH/bin"
    @echo "  clean         - Clean build artifacts"
    @echo "  deps          - Download and tidy dependencies"
    @echo ""
    @echo "Code Quality:"
    @echo "  fmt           - Format code (goimports + mod tidy + whitespace)"
    @echo "  fmt-check     - Check formatting without making changes"
    @echo "  lint          - Run golangci-lint (basic)"
    @echo "  lint-strict   - Run comprehensive linting (matches pre-commit)"
    @echo "  pre-commit-check - Run all pre-commit validations locally"
    @echo ""
    @echo "Testing:"
    @echo "  test          - Run tests"
    @echo "  test-coverage - Run tests with coverage report"
    @echo ""
    @echo "  help          - Show this help message"
    @echo ""
    @echo "Tip: Run 'just --list' to see all available recipes"

# Create build directory
_create-build-dir:
    @mkdir -p {{build_dir}}

# Install dependencies
deps:
    go mod download
    go mod tidy

# Build the application
build: _create-build-dir deps
    go build {{ldflags}} -o {{build_dir}}/{{binary_name}} -v ./cmd/{{cmd_pkg}}

# Run the application
run: build
    ./{{build_dir}}/{{binary_name}}

# Install the binary to GOPATH/bin or $HOME/go/bin
install: build
    #!/usr/bin/env bash
    set -euo pipefail
    if [ -z "${GOPATH:-}" ]; then
        echo "Installing to $HOME/go/bin/{{binary_name}}"
        mkdir -p "$HOME/go/bin"
        cp {{build_dir}}/{{binary_name}} "$HOME/go/bin/{{binary_name}}"
    else
        echo "Installing to $GOPATH/bin/{{binary_name}}"
        mkdir -p "$GOPATH/bin"
        cp {{build_dir}}/{{binary_name}} "$GOPATH/bin/{{binary_name}}"
    fi

# Clean build artifacts
clean:
    go clean
    rm -f {{build_dir}}/{{binary_name}}
    rm -f coverage.out

# Run tests
test:
    go test -v ./...

# Run tests with coverage report
test-coverage:
    go test -v -coverprofile=coverage.out ./...
    go tool cover -html=coverage.out

# Format code (matches pre-commit formatting)
fmt:
    @echo "Formatting Go code with goimports..."
    goimports -w .
    @echo "Tidying Go modules..."
    go mod tidy
    @echo "Fixing trailing whitespace..."
    @find . -name "*.go" -exec sed -i 's/[[:space:]]*$$//' {} \;

# Check formatting without making changes
fmt-check:
    #!/usr/bin/env bash
    set -euo pipefail
    echo "Checking Go formatting..."
    if [ -n "$(goimports -l .)" ]; then
        echo "The following files need formatting:"
        goimports -l .
        echo "Run 'just fmt' to fix formatting issues."
        exit 1
    fi
    echo "Checking Go modules..."
    if ! go mod tidy -diff; then
        echo "Go modules need tidying. Run 'just fmt' to fix."
        exit 1
    fi

# Run linter (basic)
lint:
    golangci-lint run

# Run comprehensive linting (matches pre-commit strictness)
lint-strict: fmt-check
    #!/usr/bin/env bash
    set -euo pipefail
    echo "Running comprehensive linting checks..."
    echo "1. Formatting checks passed (from fmt-check dependency)"
    echo "2. Running golangci-lint..."
    golangci-lint run
    echo "3. Checking for merge conflicts..."
    if find . -name "*.go" -exec grep -l "<<<<<<< HEAD\|=======" {} \; | head -1 | grep -q .; then
        echo "Merge conflict markers found in Go files!"
        exit 1
    fi
    echo "All linting checks passed!"

# Run all pre-commit checks locally
pre-commit-check: fmt-check build test
    #!/usr/bin/env bash
    set -euo pipefail
    echo "Running all pre-commit validation checks..."
    echo "1. Formatting checks passed (from fmt-check dependency)"
    echo "2. Build passed (from build dependency)"
    echo "3. Tests passed (from test dependency)"
    echo "4. Running golangci-lint..."
    golangci-lint run --timeout=5m
    echo "5. Additional validation checks..."
    if find . -name "*.go" -exec grep -l "<<<<<<< HEAD\|=======" {} \; | head -1 | grep -q .; then
        echo "Merge conflict markers found in Go files!"
        exit 1
    fi
    echo "All pre-commit checks passed! ✅"
