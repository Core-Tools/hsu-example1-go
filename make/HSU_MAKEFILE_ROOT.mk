# HSU Universal Makefile - Root Coordinator
# Supports all 3 approaches in HSU Repository Portability Framework
# HSU Makefile System Version: 1.0.0

# Include project configuration (with defaults)
-include HSU_MAKE_CONFIG_TMPL.mk
-include Makefile.config

# Auto-detect OS and Shell Environment
ifeq ($(OS),Windows_NT)
    ifneq ($(findstring msys,$(shell uname -s 2>/dev/null | tr A-Z a-z)),)
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

# Default project settings (if not configured)
PROJECT_NAME ?= hsu-project
PROJECT_DOMAIN ?= $(shell basename $(PWD) | sed 's/^hsu-//')
PROJECT_VERSION ?= 1.0.0
DEFAULT_PORT ?= 50055

# Auto-detect repository structure
ifndef REPO_TYPE
    ifneq ($(wildcard go/go.mod),)
        ifneq ($(wildcard python/pyproject.toml),)
            REPO_TYPE := multi-language
        else
            REPO_TYPE := multi-language
        endif
    else ifneq ($(wildcard go.mod),)
        REPO_TYPE := single-language-go
    else ifneq ($(wildcard pyproject.toml),)
        REPO_TYPE := single-language-python
    else ifneq ($(wildcard requirements.txt),)
        REPO_TYPE := single-language-python
    else
        REPO_TYPE := unknown
    endif
endif

# Auto-detect language directories
ifndef GO_DIR
    ifneq ($(wildcard go/go.mod),)
        GO_DIR := go
    else ifneq ($(wildcard go.mod),)
        GO_DIR := .
    else
        GO_DIR := 
    endif
endif

ifndef PYTHON_DIR
    ifneq ($(wildcard python/pyproject.toml),)
        PYTHON_DIR := python
    else ifneq ($(wildcard python/requirements.txt),)
        PYTHON_DIR := python
    else ifneq ($(wildcard pyproject.toml),)
        PYTHON_DIR := .
    else ifneq ($(wildcard requirements.txt),)
        PYTHON_DIR := .
    else
        PYTHON_DIR := 
    endif
endif

# Auto-detect enabled languages
ifndef ENABLE_GO
    ifneq ($(GO_DIR),)
        ENABLE_GO := yes
    else
        ENABLE_GO := no
    endif
endif

ifndef ENABLE_PYTHON
    ifneq ($(PYTHON_DIR),)
        ENABLE_PYTHON := yes
    else
        ENABLE_PYTHON := no
    endif
endif

# Language-specific includes  
ifeq ($(ENABLE_GO),yes)
    ifneq ($(GO_DIR),)
        include $(INCLUDE_PREFIX)HSU_MAKEFILE_GO.mk
    endif
endif

ifeq ($(ENABLE_PYTHON),yes)
    ifneq ($(PYTHON_DIR),)
        include $(INCLUDE_PREFIX)HSU_MAKEFILE_PYTHON.mk
    endif
endif

# Default target
.DEFAULT_GOAL := help

# Global targets
.PHONY: help info clean build test setup lint format check all proto protoc proto-gen

## Help - Show available targets
help:
	@echo "HSU Universal Makefile ($(REPO_TYPE)) - Available Targets:"
	@echo ""
	@echo "Universal Commands:"
	@echo "  help         - Show this help message"
	@echo "  info         - Show build environment information"
	@echo "  setup        - Initialize development environment"
	@echo "  clean        - Clean all build artifacts"
	@echo "  build        - Build all enabled languages"
	@echo "  test         - Run tests for all enabled languages"
	@echo "  lint         - Run linting for all enabled languages"
	@echo "  format       - Format code for all enabled languages"
	@echo "  check        - Run all checks for all enabled languages"
	@echo "  proto        - Generate protobuf code for all enabled languages"
	@echo ""
ifeq ($(ENABLE_GO),yes)
	@echo "Go Commands:"
	@echo "  go-build     - Build Go components"
	@echo "  go-test      - Run Go tests"
	@echo "  go-lint      - Lint Go code"
	@echo "  go-format    - Format Go code" 
	@echo "  go-clean     - Clean Go artifacts"
	@echo "  go-protoc    - Generate Go protobuf code"
	@echo ""
endif
ifeq ($(ENABLE_PYTHON),yes)
	@echo "Python Commands:"
	@echo "  py-build     - Build Python components"
	@echo "  py-test      - Run Python tests"
	@echo "  py-lint      - Lint Python code"
	@echo "  py-format    - Format Python code"
	@echo "  py-clean     - Clean Python artifacts"
	@echo "  py-protoc    - Generate Python protobuf code"
ifeq ($(ENABLE_NUITKA),yes)
	@echo "  py-nuitka    - Build Python binary with Nuitka"
	@echo "  nuitka       - Alias for py-nuitka"
	@echo "  nuitka-info  - Show Nuitka configuration"
	@echo "  nuitka-clean - Clean Nuitka artifacts"
endif
	@echo ""
endif
	@echo "Examples:"
	@echo "  make setup && make build"
	@echo "  make test"
	@echo "  make format && make lint"

## Info - Show build environment and detected configuration
info:
	@echo "=== HSU Universal Makefile - Build Environment ==="
	@echo ""
	@echo "Environment:"
	@echo "  Detected OS: $(DETECTED_OS)"
	@echo "  Shell Type: $(SHELL_TYPE)"
	@echo "  Native OS: $(OS)"
	@echo "  Path Separator: $(PATH_SEP)"
	@echo "  Executable Extension: $(EXECUTABLE_EXT)"
	@echo ""
	@echo "Project Configuration:"
	@echo "  Project Name: $(PROJECT_NAME)"
	@echo "  Project Domain: $(PROJECT_DOMAIN)"
	@echo "  Project Version: $(PROJECT_VERSION)"
	@echo "  Repository Type: $(REPO_TYPE)"
	@echo ""
	@echo "Language Support:"
	@echo "  Go Enabled: $(ENABLE_GO)"
ifeq ($(ENABLE_GO),yes)
	@echo "    Go Directory: $(GO_DIR)"
	@echo "    Go Module: $(GO_MODULE_NAME)"
endif
	@echo "  Python Enabled: $(ENABLE_PYTHON)"
ifeq ($(ENABLE_PYTHON),yes)
	@echo "    Python Directory: $(PYTHON_DIR)"
	@echo "    Python Package: $(PYTHON_PACKAGE_NAME)"
endif
	@echo ""
	@echo "Build Configuration:"
	@echo "  Default Port: $(DEFAULT_PORT)"
	@echo "  Build CLI: $(BUILD_CLI)"
	@echo "  Build SRV: $(BUILD_SRV)"
	@echo "  Build LIB: $(BUILD_LIB)"

# Universal build targets
build: $(if $(filter yes,$(ENABLE_GO)),go-build) $(if $(filter yes,$(ENABLE_PYTHON)),py-build)
test: $(if $(filter yes,$(ENABLE_GO)),go-test) $(if $(filter yes,$(ENABLE_PYTHON)),py-test)  
lint: $(if $(filter yes,$(ENABLE_GO)),go-lint) $(if $(filter yes,$(ENABLE_PYTHON)),py-lint)
format: $(if $(filter yes,$(ENABLE_GO)),go-format) $(if $(filter yes,$(ENABLE_PYTHON)),py-format)
clean: $(if $(filter yes,$(ENABLE_GO)),go-clean) $(if $(filter yes,$(ENABLE_PYTHON)),py-clean)
check: lint test
setup: $(if $(filter yes,$(ENABLE_GO)),go-setup) $(if $(filter yes,$(ENABLE_PYTHON)),py-setup)

# Universal aliases for backward compatibility
all: build
tidy: $(if $(filter yes,$(ENABLE_GO)),go-tidy)
deps: $(if $(filter yes,$(ENABLE_GO)),go-deps) $(if $(filter yes,$(ENABLE_PYTHON)),py-deps)

## Protocol Buffer Generation - Generate protobuf code for all enabled languages
proto protoc proto-gen: $(if $(filter yes,$(ENABLE_GO)),go-protoc) $(if $(filter yes,$(ENABLE_PYTHON)),py-protoc)

# Nuitka universal targets
ifeq ($(ENABLE_NUITKA),yes)
.PHONY: nuitka nuitka-clean nuitka-info
nuitka: py-nuitka
nuitka-clean: py-nuitka-clean
nuitka-info: py-nuitka-info
endif

# Running targets (Go-specific for now, extend for other languages)
ifeq ($(ENABLE_GO),yes)
.PHONY: run-srv run-cli run-cli-port run-cli-srvpath

## Run Server - Start the server on default port
run-srv: go-build-srv
	@echo "Starting server on port $(DEFAULT_PORT)..."
	./cmd/srv/echogrpcsrv/echogrpcsrv$(EXECUTABLE_EXT) --port $(DEFAULT_PORT)

## Run CLI - Run CLI client (connects to localhost:DEFAULT_PORT)
run-cli: go-build-cli
	@echo "Running CLI client (connecting to localhost:$(DEFAULT_PORT))..."
	./cmd/cli/echogrpccli/echogrpccli$(EXECUTABLE_EXT) --port $(DEFAULT_PORT)

## Run CLI Port - Run CLI client with custom port (use PORT=xxxx)
run-cli-port: go-build-cli
	@echo "Running CLI client (connecting to localhost:$(or $(PORT),$(DEFAULT_PORT)))..."
	./cmd/cli/echogrpccli/echogrpccli$(EXECUTABLE_EXT) --port $(or $(PORT),$(DEFAULT_PORT))

## Run CLI Server Path - Run CLI client with server executable path
run-cli-srvpath: build
	@echo "Running CLI client with server path..."
	./cmd/cli/echogrpccli/echogrpccli$(EXECUTABLE_EXT) --server "./cmd/srv/echogrpcsrv/echogrpcsrv$(EXECUTABLE_EXT)"

endif 