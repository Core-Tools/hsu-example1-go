# HSU Makefile Configuration
# Template configuration file for HSU Repository Portability Framework

# Project Information
PROJECT_NAME := hsu-example-template
PROJECT_DOMAIN := example
PROJECT_VERSION := 1.0.0

# Include Path Configuration
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

# Custom build directories (relative to language root)
CLI_BUILD_DIR := cmd/cli
SRV_BUILD_DIR := cmd/srv
LIB_BUILD_DIR := pkg

# Testing
TEST_TIMEOUT := 10m
TEST_VERBOSE := yes

# Development
ENABLE_LINTING := yes
ENABLE_FORMATTING := yes
ENABLE_BENCHMARKS := yes

# Platform
TARGET_PLATFORMS := linux/amd64 windows/amd64 darwin/amd64

# Docker support
ENABLE_DOCKER := no
DOCKER_REGISTRY := ghcr.io/core-tools
DOCKER_TAG := $(PROJECT_VERSION)

# Documentation
ENABLE_DOCS := yes
DOCS_FORMAT := markdown  # markdown | sphinx | godoc 