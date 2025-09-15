# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

protoc-gen-flags is a Go-based protocol buffer compiler plugin that generates command-line flag bindings for protobuf messages. It creates `AddFlags` methods that integrate with the `spf13/pflag` library to automatically generate CLI flags from protobuf message definitions.

## Development Commands

### Build and Code Quality
```bash
# Format Go code and protobuf files
make fmt

# Run static analysis
make vet

# Run comprehensive linting
make lint

# Clean up module dependencies
make tidy

# Run tests
make test

# Clean build artifacts
make clean
```

### Code Generation
```bash
# Install required tools (buf, protoc-gen-go, golangci-lint)
make deps

# Generate Go code from protobuf definitions
make generate
```

## Architecture

### Core Components

**main.go**: Entry point that initializes the protoc-gen-star plugin framework and registers the flags module.

**modules/modules.go**: Contains the core `Module` struct that implements the protoc-gen-star interface. This module:
- Processes protobuf messages and fields
- Generates `AddFlags` methods using Go templates
- Maps protobuf field types to appropriate pflag methods
- Handles field options and message-level configuration

**flags/**: Defines protobuf extensions for configuring flag generation:
- `flags.proto`: Protocol buffer extensions for field and message options
- `flags.go`: Interface definition for the generated `AddFlags` method

### Key Design Patterns

1. **Template-based Generation**: Uses Go's `text/template` to generate flag binding code
2. **Extension-based Configuration**: Uses protobuf field and message options to control flag generation
3. **Type Mapping**: Automatically maps protobuf types to appropriate pflag methods (BoolVarP, StringVarP, etc.)
4. **Prefix Support**: Allows hierarchical flag organization through prefix parameters

### Extension Usage

The plugin supports configuration via protobuf options:

```protobuf
// Message-level options
message MyMessage {
  option (flags.disabled) = false;
  option (flags.ignored) = false;
  
  string name = 1 [(flags.flag) = {
    enabled: true
    name: "custom-name"
    short: "n"
    usage: "Custom usage text"
  }];
}
```

## Dependencies

- **github.com/lyft/protoc-gen-star**: Protocol buffer code generation framework
- **github.com/spf13/pflag**: POSIX/GNU-style command-line flag parsing
- **google.golang.org/protobuf**: Google's protocol buffer implementation

## Buf Configuration

The project uses buf for protobuf linting and generation:
- `buf.yaml`: Module configuration for buf.build/linka-cloud/protoc-gen-defaults
- `buf.gen.yaml`: Code generation configuration with source-relative paths