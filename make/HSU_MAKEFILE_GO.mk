# HSU Go-Specific Makefile
# Handles Go operations for HSU Repository Portability Framework
# HSU Makefile System Version: 1.0.0

# Go Configuration (with defaults)
GO_MODULE_NAME ?= github.com/core-tools/$(PROJECT_NAME)
GO_BUILD_FLAGS ?= -v
GO_MOD_FLAGS ?= -mod=readonly
GO_TEST_FLAGS ?= -v -race
GO_TEST_TIMEOUT ?= 10m

# Domain-specific import handling
DOMAIN_IMPORT_PREFIX ?= github.com/core-tools/hsu-$(PROJECT_DOMAIN)
DOMAIN_REPLACE_TARGET ?= .

# Build directory configuration (relative to GO_DIR) - Language-specific
GO_CLI_BUILD_DIR := $(or $(GO_CLI_BUILD_DIR),cmd/cli)
GO_SRV_BUILD_DIR := $(or $(GO_SRV_BUILD_DIR),cmd/srv)
GO_LIB_BUILD_DIR := $(or $(GO_LIB_BUILD_DIR),pkg)

# Auto-detect Go executables (cross-platform compatible)
ifeq ($(SHELL_TYPE),MSYS)
    GO_CLI_TARGETS := $(shell bash -c "find $(GO_DIR)/$(GO_CLI_BUILD_DIR) -name 'main.go' -exec dirname {} \;" 2>$(NULL_DEV) | sed 's|$(GO_DIR)/||')
    GO_SRV_TARGETS := $(shell bash -c "find $(GO_DIR)/$(GO_SRV_BUILD_DIR) -name 'main.go' -exec dirname {} \;" 2>$(NULL_DEV) | sed 's|$(GO_DIR)/||')
else
    GO_CLI_TARGETS := $(shell find $(GO_DIR)/$(GO_CLI_BUILD_DIR) -name "main.go" -exec dirname {} \; 2>$(NULL_DEV) | sed 's|$(GO_DIR)/||')
    GO_SRV_TARGETS := $(shell find $(GO_DIR)/$(GO_SRV_BUILD_DIR) -name "main.go" -exec dirname {} \; 2>$(NULL_DEV) | sed 's|$(GO_DIR)/||')
endif

# Go command wrapper - runs from correct directory
ifeq ($(GO_DIR),.)
    GO_CMD := 
    GO_MODULE_DIR := .
else
    GO_CMD := cd $(GO_DIR) &&
    GO_MODULE_DIR := $(GO_DIR)
endif

# Calculate output paths
ifeq ($(GO_DIR),.)
    # Single-language repo
    define make_go_target
        $(1)/$(notdir $(1))$(EXECUTABLE_EXT)
    endef
    define make_go_output
        $(1)
    endef
else
    # Multi-language repo
    define make_go_target
        $(GO_DIR)/$(1)/$(notdir $(1))$(EXECUTABLE_EXT)
    endef
    define make_go_output
        ../$(GO_DIR)/$(1)/$(notdir $(1))$(EXECUTABLE_EXT)
    endef
endif

# Create target lists
GO_CLI_BINARIES := $(foreach target,$(GO_CLI_TARGETS),$(call make_go_target,$(target)))
GO_SRV_BINARIES := $(foreach target,$(GO_SRV_TARGETS),$(call make_go_target,$(target)))
ALL_GO_BINARIES := $(GO_CLI_BINARIES) $(GO_SRV_BINARIES)

# Go-specific targets
.PHONY: go-setup go-deps go-tidy go-build go-build-cli go-build-srv go-test go-lint go-format go-clean go-check go-info go-mod-verify go-protoc go-proto-gen

## Go Setup - Initialize Go development environment
go-setup: go-deps go-tidy
	@echo "✓ Go development environment ready"

## Go Dependencies - Download and verify Go modules
go-deps:
	@echo "Downloading Go dependencies..."
	$(GO_CMD) go mod download
	$(GO_CMD) go mod verify

## Go Tidy - Clean up go.mod and go.sum
go-tidy:
	@echo "Tidying Go modules..."
	$(GO_CMD) go mod tidy

## Go Build - Build all Go components
go-build: go-build-cli go-build-srv

## Go Build CLI - Build all CLI executables
go-build-cli:
	@echo "Building Go CLI components..."
ifneq ($(GO_CLI_TARGETS),)
	@$(foreach target,$(GO_CLI_TARGETS),$(call build_go_binary,$(target)))
else
	@echo "No CLI targets found in $(GO_DIR)/$(GO_CLI_BUILD_DIR)"
endif

## Go Build Server - Build all server executables  
go-build-srv:
	@echo "Building Go server components..."
ifneq ($(GO_SRV_TARGETS),)
	@$(foreach target,$(GO_SRV_TARGETS),$(call build_go_binary,$(target)))
else
	@echo "No server targets found in $(GO_DIR)/$(GO_SRV_BUILD_DIR)"
endif

# Function to build a Go binary
define build_go_binary
	@echo "  Building $(1)..."
	@$(MKDIR) $(dir $(call make_go_target,$(1))) 2>$(NULL_DEV) || true
	$(GO_CMD) go build $(GO_BUILD_FLAGS) -o $(call make_go_output,$(1)) ./$(1)/
	@echo "  ✓ Built: $(call make_go_target,$(1))"

endef

## Go Test - Run all Go tests
go-test:
	@echo "Running Go tests..."
	$(GO_CMD) go test $(GO_MOD_FLAGS) $(GO_TEST_FLAGS) -timeout $(GO_TEST_TIMEOUT) ./...

## Go Lint - Run Go linting and static analysis
go-lint:
	@echo "Running Go linting..."
	$(GO_CMD) go vet ./...
ifeq ($(shell which golangci-lint 2>$(NULL_DEV)),)
	@echo "  Note: golangci-lint not found, using basic go vet only"
else
	$(GO_CMD) golangci-lint run ./...
endif

## Go Format - Format Go code
go-format:
	@echo "Formatting Go code..."
	$(GO_CMD) go fmt ./...
ifeq ($(shell which goimports 2>$(NULL_DEV)),)
	@echo "  Note: goimports not found, using basic go fmt only"
else
	$(GO_CMD) goimports -w .
endif

## Go Clean - Clean Go build artifacts
go-clean:
	@echo "Cleaning Go artifacts..."
	$(GO_CMD) go clean ./...
	@$(foreach binary,$(ALL_GO_BINARIES),$(RM) $(binary) 2>$(NULL_DEV) || true;)
	# Clean generated protobuf files
ifeq ($(GO_DIR),.)
	-$(RM_RF) "$(GO_LIB_BUILD_DIR)/generated" 2>$(NULL_DEV) || true
else
	-$(RM_RF) "$(GO_DIR)/$(GO_LIB_BUILD_DIR)/generated" 2>$(NULL_DEV) || true
endif
	@echo "✓ Go clean complete"

## Go Check - Run all Go checks (lint + test)
go-check: go-lint go-test

## Go Module Verify - Verify module dependencies
go-mod-verify:
	@echo "Verifying Go modules..."
	$(GO_CMD) go mod verify
	$(GO_CMD) go list -m all

## Go Info - Show Go environment information
go-info:
	@echo "=== Go Environment Information ==="
	@echo "Go Directory: $(GO_DIR)"
	@echo "Go Module Directory: $(GO_MODULE_DIR)"
	@echo "Go Module Name: $(GO_MODULE_NAME)"
	@echo "Go Version: $$($(GO_CMD) go version)"
	@echo "Go Root: $$($(GO_CMD) go env GOROOT)"
	@echo "Go Path: $$($(GO_CMD) go env GOPATH)"
	@echo "Go Module Cache: $$($(GO_CMD) go env GOCACHE)"
	@echo ""
	@echo "Build Configuration:"
	@echo "Build Flags: $(GO_BUILD_FLAGS)"
	@echo "Test Flags: $(GO_TEST_FLAGS)"
	@echo "Test Timeout: $(GO_TEST_TIMEOUT)"
	@echo ""
	@echo "Domain Configuration:"
	@echo "Domain Import Prefix: $(DOMAIN_IMPORT_PREFIX)"
	@echo "Domain Replace Target: $(DOMAIN_REPLACE_TARGET)"
	@echo ""
	@echo "Detected Targets:"
	@echo "CLI Targets: $(GO_CLI_TARGETS)"
	@echo "Server Targets: $(GO_SRV_TARGETS)"
	@echo "CLI Binaries: $(GO_CLI_BINARIES)"
	@echo "Server Binaries: $(GO_SRV_BINARIES)"

# Diagnostic targets for domain-based imports
.PHONY: go-lint-diag go-lint-fix

## Go Lint Diagnostics - Diagnose domain import issues
go-lint-diag:
	@echo "=== Go Domain Import Diagnostics ==="
	@echo "Repository Type: $(REPO_TYPE)"
	@echo "Go Module Directory: $(GO_MODULE_DIR)"
	@echo ""
	@echo "1. Current go.mod content:"
	@cat $(GO_MODULE_DIR)/go.mod
	@echo ""
	@echo "2. Go module list:"
	@$(GO_CMD) go list -m all
	@echo ""
	@echo "3. Checking replace directive effectiveness:"
	@echo "   Trying to resolve domain module..."
	@$(GO_CMD) go list -m $(DOMAIN_IMPORT_PREFIX) || echo "   ❌ Domain module not found"
	@echo ""
	@echo "4. Import analysis:"
	@echo "   Checking domain imports in source files..."
ifeq ($(SHELL_TYPE),MSYS)
	@bash -c "find $(GO_DIR) -name '*.go' -exec grep -l '$(DOMAIN_IMPORT_PREFIX)' {} \;" | head -5 | sed 's/^/   /' || echo "   No domain imports found"
else
	@find $(GO_DIR) -name "*.go" -exec grep -l "$(DOMAIN_IMPORT_PREFIX)" {} \; | head -5 | sed 's/^/   /' || echo "   No domain imports found"
endif
	@echo ""
	@echo "5. Suggested fixes:"
	@echo "   - Ensure go.mod has: replace $(DOMAIN_IMPORT_PREFIX) => $(DOMAIN_REPLACE_TARGET)"
	@echo "   - Try: make go-lint-fix"
	@echo "   - For VSCode: restart Go language server"

# =============================================================================
# Protocol Buffer Generation Support
# =============================================================================

# Protocol Buffer Configuration
PROTO_API_DIR ?= api/proto

# Handle single-language vs multi-language directory structure
ifeq ($(GO_DIR),.)
    # Single-language repo: output directly to pkg/generated/api/proto
    GO_PROTO_OUTPUT_DIR := $(GO_LIB_BUILD_DIR)/generated/api/proto
else
    # Multi-language repo: output to go/pkg/generated/api/proto
    GO_PROTO_OUTPUT_DIR := $(GO_DIR)/$(GO_LIB_BUILD_DIR)/generated/api/proto
endif

## Go Protocol Buffer Generation - Generate Go code from .proto files
go-protoc go-proto-gen:
	@echo "Generating Go protobuf code..."
	@if [ ! -d "$(PROTO_API_DIR)" ]; then \
		echo "❌ Proto directory not found: $(PROTO_API_DIR)"; \
		exit 1; \
	fi
	@if [ ! -f "$(PROTO_API_DIR)"/*.proto ]; then \
		echo "❌ No .proto files found in $(PROTO_API_DIR)"; \
		exit 1; \
	fi
	@if ! which protoc >/dev/null 2>&1; then \
		echo "❌ protoc not found. Install Protocol Buffers compiler"; \
		echo "   See: https://grpc.io/docs/protoc-installation/"; \
		exit 1; \
	fi
	@$(MKDIR) "$(GO_PROTO_OUTPUT_DIR)" 2>$(NULL_DEV) || true
	@echo "📁 Output directory: $(GO_PROTO_OUTPUT_DIR)"
	
	# Copy infrastructure files from make/ to generated folder
	@if [ -f "$(INCLUDE_PREFIX)version.go.template" ]; then \
		echo "📋 Copying version.go infrastructure file..."; \
		$(COPY) "$(INCLUDE_PREFIX)version.go.template" "$(GO_PROTO_OUTPUT_DIR)/version.go"; \
	fi
	
	@for proto_file in $(PROTO_API_DIR)/*.proto; do \
		echo "🔄 Processing: $$proto_file"; \
		protoc \
			--go_out="$(GO_PROTO_OUTPUT_DIR)" \
			--go_opt=paths=source_relative \
			--go-grpc_out="$(GO_PROTO_OUTPUT_DIR)" \
			--go-grpc_opt=paths=source_relative \
			-I="$(PROTO_API_DIR)" \
			"$$proto_file"; \
	done
	@echo "✅ Go protobuf generation complete"
	@echo "📂 Generated files in: $(GO_PROTO_OUTPUT_DIR)" 

## Go Lint Fix - Attempt common fixes for linter issues
go-lint-fix:
	@echo "=== Attempting Go Linter Fixes ==="
	@echo "1. Cleaning module cache..."
	$(GO_CMD) go clean -modcache
	@echo "2. Tidying modules..."
	$(GO_CMD) go mod tidy
	@echo "3. Downloading fresh dependencies..."
	$(GO_CMD) go mod download
	@echo "4. Verifying modules..."
	$(GO_CMD) go mod verify
	@echo ""
	@echo "✓ Basic fixes applied. Try running your editor/linter again."
	@echo "  For persistent issues, check domain import configuration." 
