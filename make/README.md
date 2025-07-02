# HSU Universal Makefile System

Core makefile templates and build automation for the HSU (Highly Structured Universal) Repository Portability Framework.

## Overview

This directory contains the HSU makefiles that provide cross-platform, multi-language build automation for HSU-based repositories.

## Key Features

- **Universal Commands**: `make proto`, `make build`, `make test`, `make lint`
- **Multi-Language Support**: Go, Python, C++, Rust
- **Cross-Platform**: Windows, macOS, Linux (with shell auto-detection)
- **Repository Portability**: Supports all 3 HSU approaches
- **Binary Compilation**: Nuitka support for Python projects

## Quick Start

Copy the makefile templates to your project:

```bash
# Copy templates to your project's make/ directory
cp HSU_MAKEFILE_ROOT.mk /your-project/make/
cp HSU_MAKE_CONFIG_TMPL.mk /your-project/Makefile.config

# Build protobuf code and project
make proto && make build

# Run tests and create binaries
make test && make nuitka
```

## Important Files

- **HSU_MAKEFILE_ROOT.mk** - Main coordinator makefile (include this in your project Makefile)
- **HSU_MAKE_CONFIG_TMPL.mk** - Configuration template (copy and customize for your project)
- **HSU_MAKEFILE_GO.mk** - Go-specific build rules
- **HSU_MAKEFILE_PYTHON.mk** - Python-specific build rules (includes Nuitka support)

## Documentation

For comprehensive guides, setup instructions, and advanced configuration:
https://github.com/Core-Tools/docs/makefile_guide/index.md
