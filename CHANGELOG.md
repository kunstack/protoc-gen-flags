# Changelog

All notable changes to the protoc-gen-flags project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Enhanced enum slice implementation with improved protobuf definitions
- Comprehensive interface documentation for better developer experience
- Generic type support for bytes slice operations
- Int32SliceValue and Int64SliceValue implementations with full test coverage
- Modular refactoring with timestamp support and comprehensive type handling

### Changed
- Improved flag generation system with optimized modular design
- Enhanced timestamp format support with additional parsing options
- Updated bytes slice implementation for better type safety
- Refined protobuf flag generation with better validation

### Fixed
- Resolved issues with enum slice handling
- Fixed timestamp parsing edge cases
- Improved error handling in type conversions

## [0.2.0] - 2025-10-15

### Added
- **Complete slice type support** for all scalar types (int32, int64, uint32, uint64, float, double, bool, string, bytes)
- **CSV parsing** for repeated fields with proper escaping and validation
- **Generic type implementations** using Go 1.18+ generics for type-safe operations
- **Enhanced bytes encoding** with both base64 and hexadecimal support
- **Comprehensive test coverage** for all slice value types
- **Modular type system** with dedicated handlers for each protobuf type

### Changed
- **Major architectural refactoring** with clean separation of concerns
- **Improved template system** for more maintainable code generation
- **Enhanced validation** with better error messages and type checking
- **Optimized performance** through efficient flag parsing algorithms

### Fixed
- Resolved memory allocation issues in slice operations
- Fixed parsing edge cases for numeric types
- Improved handling of empty and nil values

## [0.1.5] - 2025-10-09

### Added
- **Map field support** with JSON-based configuration
- **Enhanced duration parsing** with multiple format support
- **Improved timestamp handling** with extensive format options
- **Better nested message support** with hierarchical flag organization
- **Additional time formats** including RFC3339, ISO8601, Kitchen, and custom layouts

### Changed
- **Streamlined flag generation** with reduced code duplication
- **Improved error reporting** with detailed validation messages
- **Enhanced protobuf extension** with more configuration options

### Fixed
- Corrected timestamp parsing in edge cases
- Fixed map field generation for complex types
- Resolved issues with nested message prefixes

## [0.1.4] - 2025-09-24

### Added
- **Enhanced timestamp format support** with 15+ predefined formats
- **Custom time layout support** for specialized parsing needs
- **Improved documentation** with comprehensive examples
- **Better format validation** for timestamp inputs

### Changed
- **Optimized timestamp parsing** for better performance
- **Improved format detection** with automatic fallback
- **Enhanced error messages** for parsing failures

### Fixed
- Resolved timezone handling issues
- Fixed format string validation
- Improved parsing reliability for edge cases

## [0.1.3] - 2025-09-15

### Added
- **Comprehensive well-known types support** for google.protobuf.Duration and google.protobuf.Timestamp
- **Wrapper type integration** for optional fields (StringValue, Int32Value, etc.)
- **Enhanced protobuf options** with more configuration flexibility
- **Improved flag validation** with better type checking

### Changed
- **Refactored flag generation system** with optimized design patterns
- **Better integration** with protoc-gen-star framework
- **Improved code organization** with modular architecture

### Fixed
- Enhanced compatibility with different protobuf versions
- Fixed wrapper type initialization issues
- Improved error handling in flag generation

## [0.1.2] - 2025-09-12

### Added
- **Enum type support** with integer representation
- **Oneof field handling** for complex message structures
- **Enhanced field options** with deprecation and hiding support
- **Better message-level configuration** with disable and unexported options

### Changed
- **Improved code generation** with better template organization
- **Enhanced validation** with comprehensive error checking
- **Better performance** through optimized flag registration

### Fixed
- Fixed enum value parsing issues
- Resolved oneof field generation problems
- Improved handling of deprecated fields

## [0.1.1] - 2025-09-11

### Added
- **Basic scalar type support** (int32, int64, uint32, uint64, float, double, bool, string, bytes)
- **Message-level options** (disabled, unexported, allow_empty)
- **Field-level configuration** (name, short, usage, hidden, deprecated)
- **Bytes encoding options** (base64, hex)
- **Nested message support** with hierarchical flags
- **Basic repeated field support** for scalar types

### Changed
- **Initial architecture** with template-based code generation
- **Basic protobuf extension** for flag configuration
- **Core module implementation** with protoc-gen-star integration

### Fixed
- Initial setup and configuration issues
- Basic type mapping problems
- Template generation edge cases

## [0.1.0] - 2025-09-11

### Added
- **Initial project release** with core functionality
- **protoc-gen-star integration** for code generation
- **Basic AddFlags method generation** for protobuf messages
- **spf13/pflag integration** for POSIX/GNU-style flag parsing
- **Foundation architecture** with extensible design

---

## Release Notes Template

When creating a new release, use this template:

```markdown
## [X.Y.Z] - YYYY-MM-DD

### Added
- New features

### Changed
- Changes in existing functionality

### Deprecated
- Soon-to-be removed features

### Removed
- Now removed features

### Fixed
- Bug fixes

### Security
- Security improvements
```

## Version History

- **0.2.x series**: Focus on comprehensive type support and architectural improvements
- **0.1.x series**: Initial development and core feature implementation

## Contributing

Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on how to contribute to this project.

## License

This project is licensed under the Apache License, Version 2.0. See [LICENSE](LICENSE) for details.