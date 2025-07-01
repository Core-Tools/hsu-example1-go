# HSU Example1 Py - HSU Universal Makefile System

# Include the HSU Universal Makefile System
include make/HSU_MAKEFILE_ROOT.mk

# Configuration is in Makefile.config

# Legacy target aliases for backward compatibility
.PHONY: build-all

build-all: build 