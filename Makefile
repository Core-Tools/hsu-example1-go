# HSU Example1 Go - Cross-Platform Makefile
# Replaces build-cli.bat, build-srv.bat, run-srv.bat, run-cli-port.bat, run-cli-srvpath.bat

# Detect OS and Shell Environment - Robust cross-platform detection
ifeq ($(OS),Windows_NT)
    # Check if we're in MSYS/Git Bash environment
    ifneq ($(findstring msys,$(shell uname -s 2>/dev/null | tr A-Z a-z)),)
        # MSYS environment - use Unix-style commands
        DETECTED_OS := Windows-MSYS
        RM_RF := rm -rf
        RM := rm -f
        MKDIR := mkdir -p
        COPY := cp
        PATH_SEP := /
        EXE_EXT := .exe
        NULL_DEV := /dev/null
        SHELL_TYPE := MSYS
    else
        # Native Windows
        DETECTED_OS := Windows
        RM_RF := rmdir /s /q
        RM := del /Q /F
        MKDIR := mkdir
        COPY := copy
        PATH_SEP := \\
        EXE_EXT := .exe
        NULL_DEV := NUL
        SHELL_TYPE := CMD
    endif
else
    DETECTED_OS := $(shell uname -s)
    RM_RF := rm -rf
    RM := rm -f
    MKDIR := mkdir -p
    COPY := cp
    PATH_SEP := /
    EXE_EXT := 
    NULL_DEV := /dev/null
    SHELL_TYPE := UNIX
endif

# Compatibility aliases
EXECUTABLE_EXT := $(EXE_EXT)
RMDIR := $(RM_RF)

# Project configuration
PROJECT_NAME := hsu-example1-go
MODULE_NAME := github.com/core-tools/$(PROJECT_NAME)

# Build directories and targets
CLI_DIR := cmd$(PATH_SEP)cli$(PATH_SEP)echogrpccli
SRV_DIR := cmd$(PATH_SEP)srv$(PATH_SEP)echogrpcsrv
CLI_TARGET := $(CLI_DIR)$(PATH_SEP)echogrpccli$(EXECUTABLE_EXT)
SRV_TARGET := $(SRV_DIR)$(PATH_SEP)echogrpcsrv$(EXECUTABLE_EXT)

# Default port for testing
DEFAULT_PORT := 50055

# Go build flags
GO_BUILD_FLAGS := -v
GO_MOD_FLAGS := -mod=readonly

# Default target
.DEFAULT_GOAL := help

.PHONY: help clean build build-all build-cli build-srv run-srv run-cli run-cli-port run-cli-srvpath test tidy check lint-diag lint-fix deps proto setup info

## Help - Show available targets
help:
	@echo "HSU Example1 Go - Available Make Targets:"
	@echo ""
	@echo "Building:"
	@echo "  build-all    - Build both CLI and server"
	@echo "  build-cli    - Build CLI client"
	@echo "  build-srv    - Build server"
	@echo "  build        - Alias for build-all"
	@echo ""
	@echo "Running:"
	@echo "  run-srv      - Run server on port $(DEFAULT_PORT)"
	@echo "  run-cli      - Run CLI client (connects to port $(DEFAULT_PORT))"
	@echo "  run-cli-port - Run CLI client with custom port (PORT=xxxx)"
	@echo "  run-cli-srvpath - Run CLI client with server path"
	@echo ""
	@echo "Development:"
	@echo "  test         - Run all tests"
	@echo "  tidy         - Clean up go.mod and go.sum"
	@echo "  check        - Run go vet and other checks"
	@echo "  lint-diag    - Diagnose linter issues with domain imports"
	@echo "  lint-fix     - Attempt common fixes for linter issues"
	@echo "  deps         - Download dependencies"
	@echo "  clean        - Clean built binaries"
	@echo ""
	@echo "Examples:"
	@echo "  make build && make run-srv"
	@echo "  make run-cli-port PORT=50056"

## Setup - Initialize development environment
setup: deps tidy
	@echo "Development environment setup complete"

## Dependencies - Download Go modules
deps:
	@echo "Downloading dependencies..."
	go mod download
	go mod verify

## Tidy - Clean up go.mod and go.sum
tidy:
	@echo "Tidying Go modules..."
	go mod tidy

## Build - Build both CLI and server
build: build-all

## Build All - Build both CLI and server
build-all: build-cli build-srv

## Build CLI - Build the CLI client
build-cli:
	@echo "Building CLI client..."
	@$(MKDIR) $(CLI_DIR) 2>$(NULL_DEV) || true
	go build $(GO_BUILD_FLAGS) -o $(CLI_TARGET) ./cmd/cli/echogrpccli/
	@echo "✓ CLI built: $(CLI_TARGET)"

## Build Server - Build the server
build-srv:
	@echo "Building server..."
	@$(MKDIR) $(SRV_DIR) 2>$(NULL_DEV) || true
	go build $(GO_BUILD_FLAGS) -o $(SRV_TARGET) ./cmd/srv/echogrpcsrv/
	@echo "✓ Server built: $(SRV_TARGET)"

## Run Server - Start the server on default port
run-srv: build-srv
	@echo "Starting server on port $(DEFAULT_PORT)..."
	./$(SRV_TARGET) --port $(DEFAULT_PORT)

## Run CLI - Run CLI client (connects to localhost:DEFAULT_PORT)
run-cli: build-cli
	@echo "Running CLI client (connecting to localhost:$(DEFAULT_PORT))..."
	./$(CLI_TARGET) --port $(DEFAULT_PORT)

## Run CLI Port - Run CLI client with custom port (use PORT=xxxx)
run-cli-port: build-cli
	@echo "Running CLI client (connecting to localhost:$(or $(PORT),$(DEFAULT_PORT)))..."
	./$(CLI_TARGET) --port $(or $(PORT),$(DEFAULT_PORT))

## Run CLI Server Path - Run CLI client with server executable path
run-cli-srvpath: build-all
	@echo "Running CLI client with server path..."
	./$(CLI_TARGET) --server "./$(SRV_TARGET)"

## Test - Run all tests
test:
	@echo "Running tests..."
	go test $(GO_MOD_FLAGS) -v ./...

## Check - Run static analysis
check:
	@echo "Running static analysis..."
	go vet ./...
	go fmt ./...
	@echo "✓ Static analysis complete"

## Lint Diagnostics - Diagnose linter issues with domain-based imports
lint-diag:
	@echo "=== Linter Diagnostics ==="
	@echo "Checking Go module configuration..."
	@echo ""
	@echo "1. Current go.mod content:"
	@cat go.mod
	@echo ""
	@echo "2. Go module list:"
	@go list -m all
	@echo ""
	@echo "3. Checking replace directive effectiveness:"
	@echo "   Trying to resolve hsu-echo module..."
	@go list -m github.com/core-tools/hsu-echo || echo "   ❌ hsu-echo module not found"
	@echo ""
	@echo "4. Checking local package structure:"
	@echo "   Looking for go.mod in current directory..."
	@if [ -f "go.mod" ]; then echo "   ✓ go.mod exists"; else echo "   ❌ go.mod missing"; fi
	@echo "   Looking for domain packages..."
	@find . -name "*.go" -path "*/domain/*" | head -5 | sed 's/^/   /'
	@echo ""
	@echo "5. Import analysis:"
	@echo "   Checking imports of hsu-echo in source files..."
	@grep -r "github.com/core-tools/hsu-echo" --include="*.go" . | head -5 | sed 's/^/   /' || echo "   No hsu-echo imports found"
	@echo ""
	@echo "6. Suggested fixes:"
	@echo "   - Ensure your current repo structure matches what hsu-echo should provide"
	@echo "   - Check if gopls workspace is properly configured"
	@echo "   - Try: go mod tidy && go clean -modcache"
	@echo "   - For VSCode: reload window or restart Go language server"

## Lint Fix Attempt - Try common fixes for linter issues
lint-fix:
	@echo "=== Attempting Linter Fixes ==="
	@echo "1. Cleaning module cache..."
	go clean -modcache
	@echo "2. Tidying modules..."
	go mod tidy
	@echo "3. Re-downloading dependencies..."
	go mod download
	@echo "4. Verifying modules..."
	go mod verify
	@echo "5. Building to verify imports..."
	go build ./... || echo "Build still has issues - check lint-diag output"
	@echo "✓ Linter fix attempt complete"
	@echo "If issues persist, try restarting your IDE/editor"

## Clean - Remove built binaries
clean:
	@echo "Cleaning built binaries..."
	-$(RM) $(CLI_TARGET) 2>$(NULL_DEV) || true
	-$(RM) $(SRV_TARGET) 2>$(NULL_DEV) || true
	@echo "✓ Clean complete"

## Protocol Buffers - Regenerate protobuf code (if needed)
proto:
	@echo "Generating protobuf code..."
	@cd api/proto && ./generate-go.sh
	@echo "✓ Protobuf generation complete"

# Development workflow targets
## Dev Server - Run server with auto-restart on changes (requires entr or similar)
dev-srv: build-srv
	@echo "Starting development server (manual restart)..."
	./$(SRV_TARGET) --port $(DEFAULT_PORT)

## Quick Test - Fast build and test cycle
quick: build-all test
	@echo "✓ Quick build and test complete"

# Docker targets (if Docker support is added later)
## Docker Build - Build Docker image (future)
docker-build:
	@echo "Docker build not yet implemented"

## Docker Run - Run in Docker container (future)  
docker-run:
	@echo "Docker run not yet implemented"

# Debug information
## Info - Show build environment info
info:
	@echo "Build Environment Information:"
	@echo "  Detected OS: $(DETECTED_OS)"
	@echo "  Shell Type: $(SHELL_TYPE)"
	@echo "  Native OS: $(OS)"
	@echo "  Executable Extension: $(EXECUTABLE_EXT)"
	@echo "  Path Separator: $(PATH_SEP)"
	@echo "  Remove Command: $(RM_RF)"
	@echo "  Null Device: $(NULL_DEV)"
	@echo ""
	@echo "Project Configuration:"
	@echo "  Project: $(PROJECT_NAME)"
	@echo "  Module: $(MODULE_NAME)"
	@echo "  CLI Target: $(CLI_TARGET)"
	@echo "  Server Target: $(SRV_TARGET)"
	@echo "  Default Port: $(DEFAULT_PORT)"
	@echo ""
	@echo "Go Environment:"
	@echo "  Go Version: $$(go version)"
	@echo "  Go Modules: $$(go list -m all | wc -l) modules" 