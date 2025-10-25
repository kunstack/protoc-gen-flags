# protoc-gen-flags

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.23-blue.svg)](https://golang.org/doc/go1.23)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)
[![GoDoc](https://godoc.org/github.com/kunstack/protoc-gen-flags?status.svg)](https://godoc.org/github.com/kunstack/protoc-gen-flags)

**protoc-gen-flags** is a powerful Protocol Buffer compiler plugin that automatically generates command-line flag bindings for protobuf messages. It creates `AddFlags` methods that integrate seamlessly with the `spf13/pflag` library to provide POSIX/GNU-style command-line flag parsing.

## ✨ Features

- **🚀 Automatic Flag Generation**: Generate CLI flags from protobuf message definitions
- **🎯 Comprehensive Type Support**: All protobuf scalar types, enums, repeated fields, maps, and well-known types
- **🔧 Flexible Configuration**: Extensive customization via protobuf options
- **🏗️ Nested Message Support**: Hierarchical flag organization with prefix support
- **📦 Well-Known Types**: Duration, Timestamp, and wrapper types with format options
- **🔒 Secure**: Bytes field encoding options (base64, hex)
- **⚡ High Performance**: Efficient flag parsing with minimal runtime overhead
- **🎨 Modern Go**: Built with Go 1.23+ and latest protobuf libraries

## 🚀 Quick Start

### Installation

```bash
# Install the plugin
go install github.com/kunstack/protoc-gen-flags@latest

# Or using go get
go get github.com/kunstack/protoc-gen-flags
```

### Basic Usage

1. **Define your protobuf message with flag options:**

```protobuf
syntax = "proto3";

package example;

import "flags/flags.proto";

option go_package = "github.com/example/project;example";

message Config {
    string host = 1 [(flags.value).string = {
        name: "host"
        short: "H"
        usage: "Server hostname"
    }];

    int32 port = 2 [(flags.value).int32 = {
        name: "port"
        short: "p"
        usage: "Server port"
        default: "8080"
    }];

    bool verbose = 3 [(flags.value).bool = {
        name: "verbose"
        short: "v"
        usage: "Enable verbose logging"
    }];
}
```

2. **Generate the code:**

```bash
protoc -I. -I flags --go_out=paths=source_relative:. --flags_out=paths=source_relative:. config.proto
```

3. **Use in your application:**

```go
package main

import (
    "fmt"
    "os"

    "github.com/spf13/pflag"
    pb "github.com/example/project"
)

func main() {
    var config pb.Config

    // Create flag set and add flags
    fs := pflag.NewFlagSet("myapp", pflag.ExitOnError)
    if err := config.AddFlags(fs, ""); err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }

    // Parse flags
    fs.Parse(os.Args[1:])

    // Use configuration
    fmt.Printf("Server: %s:%d (verbose: %v)\n",
        config.GetHost(), config.GetPort(), config.GetVerbose())
}
```

## 📋 Requirements

- **Go**: 1.23 or higher
- **Protocol Buffers**: protoc 3.0+
- **protoc-gen-go**: Go protobuf plugin

## 🔧 Advanced Configuration

### Message-Level Options

Control flag generation at the message level:

```protobuf
message MyConfig {
    // Disable flag generation for this message
    option (flags.disabled) = true;

    // Generate unexported AddFlags method
    option (flags.unexported) = true;

    // Allow flag generation without field configurations
    option (flags.allow_empty) = true;

    string field = 1;
}
```

### Field Configuration Options

All field types support these common options:

| Option | Type | Description | Example |
|--------|------|-------------|---------|
| `name` | string | Custom flag name (kebab-case) | `"db-host"` |
| `short` | string | Short flag alias | `"H"` |
| `usage` | string | Help text | `"Database hostname"` |
| `hidden` | bool | Hide from help output | `true` |
| `deprecated` | bool | Mark as deprecated | `true` |
| `deprecated_usage` | string | Deprecation message | `"Use --new-flag instead"` |

### Supported Types

#### Scalar Types
- **Numeric**: `int32`, `int64`, `uint32`, `uint64`, `sint32`, `sint64`, `fixed32`, `fixed64`, `sfixed32`, `sfixed64`, `float`, `double`
- **Basic**: `bool`, `string`
- **Bytes**: `bytes` with encoding options (base64, hex)

#### Well-Known Types
- **Duration**: `google.protobuf.Duration` with flexible parsing
- **Timestamp**: `google.protobuf.Timestamp` with multiple format support
- **Wrappers**: `google.protobuf.*Value` types (StringValue, Int32Value, etc.)

#### Complex Types
- **Enums**: Protocol buffer enum types
- **Repeated**: All scalar types support repeated fields
- **Maps**: Map fields with JSON and native format support
- **Nested Messages**: Hierarchical flag organization
- **Oneof**: Fields within oneof blocks

### Bytes Encoding

Bytes fields support multiple encoding formats:

```protobuf
bytes data = 1 [(flags.value).bytes = {
    name: "data"
    encoding: BYTES_ENCODING_TYPE_HEX  // or BYTES_ENCODING_TYPE_BASE64
}];
```

- `BYTES_ENCODING_TYPE_UNSPECIFIED`: Default base64 encoding
- `BYTES_ENCODING_TYPE_BASE64`: Standard base64 encoding
- `BYTES_ENCODING_TYPE_HEX`: Hexadecimal encoding

### Timestamp Formats

Timestamp fields support multiple time formats:

```protobuf
google.protobuf.Timestamp created = 1 [(flags.value).timestamp = {
    name: "created"
    formats: ["RFC3339", "ISO8601Time", "2006-01-02 15:04:05"]
}];
```

Supported formats include RFC3339, ISO8601, Kitchen, Stamp, and custom Go time layout strings.

### Nested Messages

Generate hierarchical flags for nested messages:

```protobuf
message DatabaseConfig {
    string host = 1 [(flags.value).string = {name: "host" usage: "DB host"}];
    int32 port = 2 [(flags.value).int32 = {name: "port" usage: "DB port"}];
}

message AppConfig {
    DatabaseConfig database = 1 [(flags.value).message = {
        name: "db"
        nested: true  // Generate --db.host and --db.port flags
    }];
}

## 📖 Complete Examples

### Basic Configuration

```protobuf
syntax = "proto3";

package example;

import "flags/flags.proto";
import "google/protobuf/timestamp.proto";
import "google/protobuf/duration.proto";

option go_package = "github.com/example/project;example";

message Config {
    option (flags.allow_empty) = true;

    // Basic types
    string host = 1 [(flags.value).string = {
        name: "host"
        short: "H"
        usage: "Server hostname"
    }];

    int32 port = 2 [(flags.value).int32 = {
        name: "port"
        short: "p"
        usage: "Server port"
    }];

    bool verbose = 3 [(flags.value).bool = {
        name: "verbose"
        short: "v"
        usage: "Enable verbose logging"
    }];

    // Well-known types
    google.protobuf.Duration timeout = 4 [(flags.value).duration = {
        name: "timeout"
        short: "t"
        usage: "Connection timeout"
    }];

    google.protobuf.Timestamp start_time = 5 [(flags.value).timestamp = {
        name: "start-time"
        usage: "Start time for the operation"
        formats: ["2006-01-02T15:04:05"]
    }];

    // Bytes with encoding
    bytes secret = 6 [(flags.value).bytes = {
        name: "secret"
        usage: "Secret key (hex encoded)"
        encoding: BYTES_ENCODING_TYPE_HEX
    }];

    // Nested configuration
    DatabaseConfig database = 7 [(flags.value).message = {
        name: "database"
        nested: true
    }];

    // Oneof selection
    oneof auth_mode {
        string token = 8 [(flags.value).string = {
            name: "token"
            usage: "Authentication token"
        }];

        string api_key = 9 [(flags.value).string = {
            name: "api-key"
            usage: "API key for authentication"
        }];
    }
}

message DatabaseConfig {
    string url = 1 [(flags.value).string = {
        name: "url"
        usage: "Database connection URL"
    }];

    int32 max_connections = 2 [(flags.value).int32 = {
        name: "max-connections"
        usage: "Maximum number of database connections"
    }];
}
```

### Generated Code Usage

```go
package main

import (
    "fmt"
    "os"

    "github.com/spf13/pflag"
    pb "github.com/example/project"
)

func main() {
    var config pb.Config

    // Add flags to the flag set
    fs := pflag.NewFlagSet("config", pflag.ExitOnError)
    if err := config.AddFlags(fs, ""); err != nil {
        fmt.Fprintf(os.Stderr, "Error adding flags: %v\n", err)
        os.Exit(1)
    }

    // Parse command line flags
    fs.Parse(os.Args[1:])

    // Use the configuration
    fmt.Printf("Host: %s, Port: %d\n", config.GetHost(), config.GetPort())
}
```

### Command Line Examples

```bash
# Basic usage
./myapp --host localhost --port 8080 --verbose

# Short flags
./myapp -H localhost -p 8080 -v

# Duration and timestamp
./myapp --timeout 30s --start-time "2024-01-01T12:00:00"

# Nested flags
./myapp --database.url "postgres://localhost/mydb" --database.max-connections 100

# Oneof selection (use one of)
./myapp --token "my-secret-token"
# or
./myapp --api-key "my-api-key"

# Bytes with hex encoding
./myapp --secret "48656c6c6f20576f726c64"

# Multiple formats for timestamps
./myapp --start-time "2024-01-01 12:00:00"
./myapp --start-time "Jan 1, 2024"
```

## 🛠️ Development

### Building from Source

```bash
# Clone the repository
git clone https://github.com/kunstack/protoc-gen-flags.git
cd protoc-gen-flags

# Install dependencies
make deps

# Build the plugin
make build

# Run tests
make test
```

### Development Commands

```bash
# Format Go code and protobuf files
make fmt

# Run static analysis
make vet

# Run comprehensive linting
make lint

# Clean up module dependencies
make tidy

# Clean build artifacts
make clean

# Generate Go code from protobuf definitions
make generate
```

### Project Structure

```
protoc-gen-flags/
├── main.go              # Plugin entry point
├── module/              # Core generation logic
│   ├── module.go        # Main module implementation
│   ├── common.go        # Common flag generation utilities
│   ├── defaults.go      # Default value generation
│   └── [type].go        # Type-specific generators
├── flags/               # Protobuf extension definitions
│   ├── flags.proto      # Extension definitions
│   └── flags.go         # Interface definitions
├── types/               # Custom pflag type implementations
│   ├── bytes.go         # Bytes encoding types
│   ├── duration.go      # Duration parser
│   ├── timestamp.go     # Timestamp parser
│   └── [type]_slice.go  # Slice type implementations
├── tests/               # Test files and examples
├── Makefile             # Build automation
├── buf.yaml            # Buf configuration
└── buf.gen.yaml        # Code generation configuration
```

## 🔍 Troubleshooting

### Common Issues

#### Plugin Not Found

If you get `protoc-gen-flags: program not found or is not executable`:

```bash
# Ensure the plugin is in your PATH
echo $PATH
which protoc-gen-flags

# If not found, install it
go install github.com/kunstack/protoc-gen-flags@latest
```

#### Import Path Issues

Make sure to include the flags proto import path:

```bash
# Correct - include the flags directory
protoc -I. -I flags --go_out=paths=source_relative:. --flags_out=paths=source_relative:. config.proto

# Incorrect - missing flags import
protoc -I. --go_out=paths=source_relative:. --flags_out=paths=source_relative:. config.proto
```

#### Generated Code Errors

If you encounter compilation errors in generated code:

1. Check that you're using compatible versions of protoc-gen-go and protoc-gen-flags
2. Ensure your protobuf definitions are valid
3. Verify that all required imports are present

### Debug Mode

Enable verbose output to debug generation issues:

```bash
protoc --flags_out=paths=source_relative,debug=true:. config.proto
```

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Development Setup

1. Fork the repository
2. Create a feature branch: `git checkout -b feature-name`
3. Make your changes and add tests
4. Run the test suite: `make test`
5. Commit your changes: `git commit -am 'Add new feature'`
6. Push to the branch: `git push origin feature-name`
7. Submit a pull request

### Code Style

- Follow standard Go conventions
- Run `make fmt` before committing
- Ensure `make lint` passes
- Add tests for new functionality
- Update documentation as needed

## 📚 Advanced Topics

### Custom Type Support

The plugin supports custom types through the `pflag.Value` interface. See the `types/` directory for examples.

### Performance Considerations

- Generated code is optimized for minimal runtime overhead
- Flag parsing uses efficient string-to-type conversions
- Memory allocation is minimized through careful design

### Integration with Other Tools

**protoc-gen-flags** works well with:
- **protoc-gen-go**: Standard Go protobuf generation
- **protoc-gen-gogo**: Alternative Go protobuf implementation
- **buf**: Modern protobuf build tool
- **golangci-lint**: Go linting tool

## 📄 License

This project is licensed under the Apache License, Version 2.0. See the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [protoc-gen-star](https://github.com/lyft/protoc-gen-star) - Protocol buffer code generation framework
- [spf13/pflag](https://github.com/spf13/pflag) - POSIX/GNU-style command-line flag parsing
- [protocolbuffers/protobuf](https://github.com/protocolbuffers/protobuf) - Protocol Buffers

## 📞 Support

- **Issues**: [GitHub Issues](https://github.com/kunstack/protoc-gen-flags/issues)
- **Discussions**: [GitHub Discussions](https://github.com/kunstack/protoc-gen-flags/discussions)
- **Documentation**: [Wiki](https://github.com/kunstack/protoc-gen-flags/wiki)

---

**Made with ❤️ by the protoc-gen-flags team**

## 🧪 Testing

### Running Tests

```bash
# Run all tests
make test

# Run specific test categories
go test ./types/...      # Type-specific tests
go test ./module/...     # Module tests
go test ./tests/...      # Integration tests
```

### Test Coverage

```bash
# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## 📊 Benchmarks

Performance benchmarks are available for critical components:

```bash
# Run benchmarks
go test -bench=. ./types/...
go test -bench=. ./module/...
```

## 🔗 Related Projects

- [protoc-gen-go](https://github.com/protocolbuffers/protobuf-go) - Go protobuf compiler plugin
- [protoc-gen-gogo](https://github.com/gogo/protobuf) - Alternative Go protobuf implementation
- [protoc-gen-validate](https://github.com/envoyproxy/protoc-gen-validate) - Protobuf validation
- [protoc-gen-doc](https://github.com/pseudomuto/protoc-gen-doc) - Documentation generation
- [buf](https://github.com/bufbuild/buf) - Modern protobuf tooling

## 🗺️ Roadmap

### Upcoming Features

- [ ] **Default Value Support**: Enhanced default value handling
- [ ] **Validation Integration**: Built-in validation rules
- [ ] **Custom Type Plugins**: Extensible type system
- [ ] **Configuration Files**: Support for config file generation
- [ ] **Environment Variables**: Automatic env var binding
- [ ] **Web UI**: Optional web interface for configuration

## ⚖️ Comparison with Alternatives

| Feature | protoc-gen-flags | Manual Flag Binding | Other Code Generators |
|---------|------------------|-------------------|----------------------|
| **Type Safety** | ✅ Strong typing | ⚠️ Manual validation | ✅ Varies |
| **Protobuf Integration** | ✅ Native | ❌ Manual mapping | ⚠️ Limited |
| **Code Generation** | ✅ Automatic | ❌ Manual | ✅ Varies |
| **Well-Known Types** | ✅ Full support | ❌ Manual handling | ⚠️ Limited |
| **Nested Messages** | ✅ Hierarchical | ❌ Complex | ⚠️ Limited |
| **Maintenance** | ✅ Low effort | ❌ High effort | ⚠️ Varies |

## 📈 Changelog

See [CHANGELOG.md](CHANGELOG.md) for a detailed history of changes and improvements.

## 🏷️ Versioning

This project follows [Semantic Versioning](https://semver.org/). For the versions available, see the [tags on this repository](https://github.com/kunstack/protoc-gen-flags/tags).

