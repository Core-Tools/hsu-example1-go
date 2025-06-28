# HSU Example1 Go - HSU Universal Makefile System
# HSU Repository Portability Framework - Approach 1 (Single-Repository + Single-Language)
# Version: 1.0.0

# Include the HSU Universal Makefile System
include make/HSU_MAKEFILE_ROOT.mk

# This Makefile automatically detects:
# - Repository Type: single-language-go
# - Language Support: Go only
# - Build Targets: CLI and Server executables
# - Cross-Platform: Windows, macOS, Linux compatibility
#
# Configuration is in Makefile.config
# All standard targets are provided by the HSU Universal system:
#   make help     - Show all available targets
#   make setup    - Initialize development environment
#   make build    - Build all components
#   make test     - Run tests
#   make clean    - Clean artifacts
#   make run-srv  - Start server
#   make run-cli  - Run client
#
# Advanced targets:
#   make info     - Show build environment
#   make go-lint-diag - Diagnose domain import issues
#   make format   - Format code
#   make lint     - Run linters

# Legacy target aliases for backward compatibility
.PHONY: build-all build-cli build-srv lint-diag lint-fix

build-all: build
build-cli: go-build-cli  
build-srv: go-build-srv
lint-diag: go-lint-diag
lint-fix: go-lint-fix 