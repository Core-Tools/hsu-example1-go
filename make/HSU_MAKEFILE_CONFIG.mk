# HSU Makefile Configuration
# Template configuration file for HSU Repository Portability Framework

# Project Information
PROJECT_NAME := hsu-example-template
PROJECT_DOMAIN := example
PROJECT_VERSION := 1.0.0

# Include Path Configuration
# Must be either empty or single path segment ending with /
# Examples: INCLUDE_PREFIX := make/  or  INCLUDE_PREFIX := tools/  or  INCLUDE_PREFIX := 
INCLUDE_PREFIX :=  # Default: empty (root-level includes) | Common: make/ build/ scripts/

# Repository Structure (auto-detected if not set)
# REPO_TYPE := single-language | multi-language | multi-repo
# GO_DIR := . | go | path/to/go/module
# PYTHON_DIR := . | python | path/to/python/module

# Language Support
ENABLE_GO := yes
ENABLE_PYTHON := no
ENABLE_CPP := no
ENABLE_RUST := no

# Go Configuration
GO_MODULE_NAME := github.com/core-tools/$(PROJECT_NAME)
GO_BUILD_FLAGS := -v
GO_MOD_FLAGS := -mod=readonly

# Python Configuration  
PYTHON_PACKAGE_NAME := $(subst -,_,$(PROJECT_NAME))
PYTHON_BUILD_TOOL := setuptools  # setuptools | poetry | pdm

# Build Targets
DEFAULT_PORT := 50055
BUILD_CLI := yes
BUILD_SRV := yes
BUILD_LIB := yes

# Domain-specific settings
DOMAIN_IMPORT_PREFIX := github.com/core-tools/hsu-$(PROJECT_DOMAIN)
DOMAIN_REPLACE_TARGET := .

# Language-specific build directories (relative to language root)
# Go directories (for Go or multi-language projects)
GO_CLI_BUILD_DIR := cmd/cli
GO_SRV_BUILD_DIR := cmd/srv
GO_LIB_BUILD_DIR := pkg

# Python directories (for Python or multi-language projects)  
PYTHON_CLI_BUILD_DIR := cli
PYTHON_SRV_BUILD_DIR := srv
PYTHON_LIB_BUILD_DIR := lib

# Testing
TEST_TIMEOUT := 10m
TEST_VERBOSE := yes

# Development
ENABLE_LINTING := yes
ENABLE_FORMATTING := yes
ENABLE_BENCHMARKS := yes

# Nuitka Build Configuration (Python binary compilation)
ENABLE_NUITKA := no
NUITKA_OUTPUT_NAME := server
NUITKA_ENTRY_POINT := srv/run_server.py
NUITKA_EXCLUDES_FILE := nuitka_excludes.txt
NUITKA_EXTRA_MODULES := 
NUITKA_EXTRA_PACKAGES := 
NUITKA_EXTRA_FOLLOW_IMPORTS := 
NUITKA_BUILD_MODE := onefile

# Platform
TARGET_PLATFORMS := linux/amd64 windows/amd64 darwin/amd64

# Docker support
ENABLE_DOCKER := no
DOCKER_REGISTRY := ghcr.io/core-tools
DOCKER_TAG := $(PROJECT_VERSION)

# Documentation
ENABLE_DOCS := yes
DOCS_FORMAT := markdown  # markdown | sphinx | godoc 