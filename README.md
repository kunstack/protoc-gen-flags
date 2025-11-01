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
- **🎯 Default Values**: Comprehensive default value support for all types
- **📊 Repeated Fields**: Full slice support with default values
- **🗺️ Map Fields**: JSON and native format support for map fields
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
- **Enums**: Protocol buffer enum types with default value support
- **Repeated**: All scalar types support repeated fields with slice defaults
- **Maps**: Map fields with JSON and native format support and defaults
- **Nested Messages**: Hierarchical flag organization with prefix support
- **Oneof**: Fields within oneof blocks

### Complete Type Reference

#### Numeric Types
All numeric types support default values and repeated variants:

| Proto Type | Go Type | Default Support | Repeated Support | Example |
|------------|---------|----------------|------------------|---------|
| `float` | `float32` | ✅ | ✅ | `3.14159` |
| `double` | `float64` | ✅ | ✅ | `2.71828` |
| `int32` | `int32` | ✅ | ✅ | `42` |
| `int64` | `int64` | ✅ | ✅ | `9223372036854775807` |
| `uint32` | `uint32` | ✅ | ✅ | `1000` |
| `uint64` | `uint64` | ✅ | ✅ | `18446744073709551615` |
| `sint32` | `int32` | ✅ | ✅ | `-42` |
| `sint64` | `int64` | ✅ | ✅ | `-9223372036854775808` |
| `fixed32` | `uint32` | ✅ | ✅ | `8080` |
| `fixed64` | `uint64` | ✅ | ✅ | `3000000000` |
| `sfixed32` | `int32` | ✅ | ✅ | `-1000` |
| `sfixed64` | `int64` | ✅ | ✅ | `-3000000000` |

#### Special Types
| Proto Type | Go Type | Features | Example |
|------------|---------|----------|---------|
| `bool` | `bool` | Default values, repeated | `true`, `false` |
| `string` | `string` | Default values, repeated | `"hello world"` |
| `bytes` | `[]byte` | Base64/hex encoding, defaults, repeated | `"aGVsbG8="` (base64) |
| `enum` | Enum type | Default values, repeated | `1` (enum value) |

#### Well-Known Types
| Proto Type | Go Type | Features | Example |
|------------|---------|----------|---------|
| `google.protobuf.Duration` | `*durationpb.Duration` | Default values, repeated | `"30s"`, `"1h"` |
| `google.protobuf.Timestamp` | `*timestamppb.Timestamp` | Multiple formats, defaults, repeated | `"2024-01-01T00:00:00Z"` |
| `google.protobuf.StringValue` | `*wrapperspb.StringValue` | Default values, repeated | `"wrapper"` |
| `google.protobuf.Int32Value` | `*wrapperspb.Int32Value` | Default values, repeated | `42` |
| `google.protobuf.Int64Value` | `*wrapperspb.Int64Value` | Default values, repeated | `9223372036854775807` |
| `google.protobuf.UInt32Value` | `*wrapperspb.UInt32Value` | Default values, repeated | `1000` |
| `google.protobuf.UInt64Value` | `*wrapperspb.UInt64Value` | Default values, repeated | `18446744073709551615` |
| `google.protobuf.FloatValue` | `*wrapperspb.FloatValue` | Default values, repeated | `3.14159` |
| `google.protobuf.DoubleValue` | `*wrapperspb.DoubleValue` | Default values, repeated | `2.71828` |
| `google.protobuf.BoolValue` | `*wrapperspb.BoolValue` | Default values, repeated | `true` |
| `google.protobuf.BytesValue` | `*wrapperspb.BytesValue` | Base64/hex encoding, defaults, repeated | `"aGVsbG8="` |

#### Map Types
| Map Type | Format Support | Default Values | Example |
|----------|---------------|----------------|---------|
| `map<string, string>` | JSON, native | ✅ | `{"key": "value"}` or `key=value` |
| `map<string, int32>` | JSON, native int | ✅ | `{"key": 123}` or `key=123` |
| `map<string, int64>` | JSON, native int | ✅ | `{"key": 456}` or `key=456` |
| `map<string, uint32>` | JSON, native int | ✅ | `{"key": 789}` or `key=789` |
| `map<string, uint64>` | JSON, native int | ✅ | `{"key": 1000}` or `key=1000` |
| `map<string, float>` | JSON | ✅ | `{"key": 3.14}` |
| `map<string, double>` | JSON | ✅ | `{"key": 2.718}` |
| `map<string, bool>` | JSON | ✅ | `{"key": true}` |

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

    // Basic types with default values
    string host = 1 [(flags.value).string = {
        name: "host"
        short: "H"
        usage: "Server hostname"
        default: "localhost"
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
        default: "false"
    }];

    // Well-known types with defaults
    google.protobuf.Duration timeout = 4 [(flags.value).duration = {
        name: "timeout"
        short: "t"
        usage: "Connection timeout"
        default: "30s"
    }];

    google.protobuf.Timestamp start_time = 5 [(flags.value).timestamp = {
        name: "start-time"
        usage: "Start time for the operation"
        formats: ["RFC3339", "ISO8601"]
        default: "2024-01-01T00:00:00Z"
    }];

    // Bytes with encoding and defaults
    bytes api_key = 6 [(flags.value).bytes = {
        name: "api-key"
        short: "k"
        usage: "API key in hex format"
        encoding: BYTES_ENCODING_TYPE_HEX
        default: "48656c6c6f20576f726c64"
    }];

    bytes config_data = 7 [(flags.value).bytes = {
        name: "config-data"
        usage: "Configuration data in base64 format"
        encoding: BYTES_ENCODING_TYPE_BASE64
        default: "aGVsbG8gd29ybGQ="
    }];

    // Repeated fields with defaults
    repeated string servers = 8 [(flags.value).repeated.string = {
        name: "servers"
        short: "s"
        usage: "Server addresses"
        default: ["localhost:8080", "localhost:8081"]
    }];

    repeated int32 allowed_ports = 9 [(flags.value).repeated.int32 = {
        name: "allowed-ports"
        usage: "Allowed port numbers"
        default: [80, 443, 8080]
    }];

    repeated bytes certificates = 10 [(flags.value).repeated.bytes = {
        name: "certificates"
        usage: "SSL certificates in base64 format"
        encoding: BYTES_ENCODING_TYPE_BASE64
        default: ["Y2VydDE=", "Y2VydDI="]
    }];

    repeated google.protobuf.Duration retry_intervals = 11 [(flags.value).repeated.duration = {
        name: "retry-intervals"
        usage: "Retry intervals for failed operations"
        default: ["1s", "2s", "5s", "10s"]
    }];

    // Map fields
    map<string, string> labels = 12 [(flags.value).map = {
        name: "labels"
        short: "l"
        usage: "Resource labels"
        format: MAP_FORMAT_TYPE_STRING_TO_STRING
        default: "{\"env\": \"production\", \"team\": \"backend\"}"
    }];

    map<string, int32> quotas = 13 [(flags.value).map = {
        name: "quotas"
        usage: "Resource quotas"
        format: MAP_FORMAT_TYPE_STRING_TO_INT
        default: "{\"requests\": 1000, \"connections\": 100}"
    }];

    // Nested configuration
    DatabaseConfig database = 14 [(flags.value).message = {
        name: "database"
        nested: true
    }];

    // Enum with default
    LogLevel log_level = 15 [(flags.value).enum = {
        name: "log-level"
        usage: "Logging level"
        default: 1  // INFO
    }];
}

enum LogLevel {
    LOG_LEVEL_UNSPECIFIED = 0;
    LOG_LEVEL_DEBUG = 1;
    LOG_LEVEL_INFO = 2;
    LOG_LEVEL_WARN = 3;
    LOG_LEVEL_ERROR = 4;
}

message DatabaseConfig {
    string host = 1 [(flags.value).string = {
        name: "host"
        usage: "Database host"
        default: "localhost"
    }];

    int32 port = 2 [(flags.value).int32 = {
        name: "port"
        usage: "Database port"
        default: 5432
    }];

    string username = 3 [(flags.value).string = {
        name: "username"
        short: "u"
        usage: "Database username"
        default: "admin"
    }];

    bytes password = 4 [(flags.value).bytes = {
        name: "password"
        short: "p"
        usage: "Database password (base64 encoded)"
        encoding: BYTES_ENCODING_TYPE_BASE64
        hidden: true
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

### Advanced Type Examples

#### Default Values and Wrapper Types

```protobuf
syntax = "proto3";

package example;

import "flags/flags.proto";
import "google/protobuf/wrappers.proto";

option go_package = "github.com/example/project;example";

message AdvancedConfig {
    option (flags.allow_empty) = true;

    // All numeric types with defaults
    float pi = 1 [(flags.value).float = {
        name: "pi"
        usage: "Pi constant"
        default: 3.14159
    }];

    double euler = 2 [(flags.value).double = {
        name: "euler"
        usage: "Euler's number"
        default: 2.71828
    }];

    int32 max_retries = 3 [(flags.value).int32 = {
        name: "max-retries"
        usage: "Maximum retry attempts"
        default: 3
    }];

    uint64 memory_limit = 4 [(flags.value).uint64 = {
        name: "memory-limit"
        usage: "Memory limit in bytes"
        default: 1073741824
    }];

    // Wrapper types with defaults
    optional google.protobuf.StringValue api_key = 5 [(flags.value).string = {
        name: "api-key"
        usage: "API key"
        default: "default-key"
    }];

    optional google.protobuf.Int32Value timeout = 6 [(flags.value).int32 = {
        name: "timeout"
        usage: "Request timeout"
        default: 30
    }];

    optional google.protobuf.BoolValue debug = 7 [(flags.value).bool = {
        name: "debug"
        usage: "Enable debug mode"
        default: "true"
    }];

    // Repeated wrapper types
    repeated google.protobuf.StringValue tags = 8 [(flags.value).repeated.string = {
        name: "tags"
        usage: "Resource tags"
        default: ["production", "api", "v1"]
    }];

    repeated google.protobuf.DoubleValue metrics = 9 [(flags.value).repeated.double = {
        name: "metrics"
        usage: "Performance metrics"
        default: [0.95, 0.99, 0.999]
    }];
}
```

#### Complex Map Examples

```protobuf
message MapExamples {
    option (flags.allow_empty) = true;

    // JSON format (default)
    map<string, string> metadata = 1 [(flags.value).map = {
        name: "metadata"
        usage: "Resource metadata"
        default: "{\"owner\": \"team-a\", \"environment\": \"prod\"}"
    }];

    // Native string-to-string format
    map<string, string> features = 2 [(flags.value).map = {
        name: "features"
        usage: "Feature flags"
        format: MAP_FORMAT_TYPE_STRING_TO_STRING
        default: "feature-a=true,feature-b=false"
    }];

    // Native string-to-int format
    map<string, int32> limits = 3 [(flags.value).map = {
        name: "limits"
        usage: "Resource limits"
        format: MAP_FORMAT_TYPE_STRING_TO_INT
        default: "cpu=1000,memory=2048,storage=10240"
    }];

    // Complex nested JSON
    map<string, string> complex_config = 4 [(flags.value).map = {
        name: "complex-config"
        usage: "Complex configuration"
        default: "{\"database\": {\"host\": \"localhost\", \"port\": 5432}, \"cache\": {\"ttl\": 300}}"
    }];
}
```

### Command Line Examples

```bash
# Basic usage with defaults
./myapp
# Uses: --host localhost --port 8080 --verbose=false

# Override specific values
./myapp --host example.com --port 9000 --verbose

# Short flags
./myapp -H example.com -p 9000 -v

# Duration and timestamp with multiple formats
./myapp --timeout 45s --start-time "2024-01-01T12:00:00"
./myapp --start-time "Jan 1, 2024 at 12:00"

# Nested flags
./myapp --database.host db.example.com --database.port 5432 --database.username admin

# Bytes with different encodings
./myapp --api-key "48656c6c6f576f726c64"  # hex
./myapp --config-data "SGVsbG8gV29ybGQ="    # base64

# Repeated fields
./myapp --servers server1.example.com --servers server2.example.com
./myapp --allowed-ports 80 --allowed-ports 443 --allowed-ports 8080
./myapp --certificates "Y2VydDE=" --certificates "Y2VydDI="

# Duration slices
./myapp --retry-intervals 1s --retry-intervals 2s --retry-intervals 5s --retry-intervals 10s

# Map fields - JSON format
./myapp --labels '{"env": "staging", "team": "frontend"}'
./myapp --metadata '{"version": "1.2.3", "build": "12345"}'

# Map fields - Native format
./myapp --features "feature-a=true,feature-b=false"
./myapp --limits "cpu=2000,memory=4096"

# Wrapper types
./myapp --api-key "custom-api-key" --timeout 60 --debug=false

# Repeated wrapper types
./myapp --tags "staging" --tags "api" --tags "v2"
./myapp --metrics 0.99 --metrics 0.999 --metrics 0.9999

# Enum values
./myapp --log-level 3  # LOG_LEVEL_WARN

# Deprecated flags with warnings
./myapp --old-flag value
# Warning: --old-flag is deprecated, use --new-flag instead

# Hidden flags (only shown with --help-all)
./myapp --help-all
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
├── main.go              # Plugin entry point and initialization
├── module/              # Core generation logic
│   ├── module.go        # Main module implementation with protoc-gen-star interface
│   ├── common.go        # Common flag generation utilities and helpers
│   ├── defaults.go      # Default value generation for all supported types
│   ├── flags.go         # Field-level flag processing and dispatch
│   ├── checker.go       # Validation and checking logic for flag configurations
│   ├── bytes.go         # Bytes field handling with encoding support
│   ├── duration.go      # Duration field parsing and generation
│   ├── timestamp.go     # Timestamp field parsing with multiple formats
│   ├── enum.go          # Enum field processing
│   ├── message.go       # Nested message field handling
│   └── map.go           # Map field generation with format support
├── flags/               # Protobuf extension definitions
│   ├── flags.proto      # Complete extension definitions for all flag types
│   ├── flags.pb.go      # Generated protobuf Go code
│   ├── flags.go         # Interface definitions and types
│   └── builder.go       # Flag name building and configuration utilities
├── types/               # Custom pflag type implementations
│   ├── bytes.go         # Base64 bytes encoding type
│   ├── bytes_hex.go     # Hexadecimal bytes encoding type
│   ├── bytes_slice.go   # Base64 bytes slice type
│   ├── bytes_hex_slice.go # Hexadecimal bytes slice type
│   ├── duration.go      # Duration parsing with flexible format support
│   ├── duration_slice.go # Duration slice type
│   ├── timestamp.go     # Timestamp parsing with multiple formats
│   ├── timestamp_slice.go # Timestamp slice type
│   ├── json.go          # JSON map type
│   ├── enum.go          # Enum slice type
│   ├── map.go           # Map type with native format support
│   ├── string.go        # String slice type
│   ├── [numeric].go     # All numeric types (int32, int64, uint32, uint64, float, double)
│   └── [type]_slice.go  # Slice implementations for all types
├── utils/               # Utility functions and helpers
│   ├── strings.go       # String manipulation utilities
│   └── time.go          # Time parsing utilities
├── tests/               # Test files and comprehensive examples
│   ├── test.proto       # Comprehensive test protobuf definitions
│   ├── test.pb.go       # Generated protobuf Go code
│   └── test.pb.flags.go # Generated flag bindings for testing
├── Makefile             # Build automation with all development commands
├── buf.yaml            # Buf configuration for protobuf linting and building
├── buf.gen.yaml        # Code generation configuration
├── CHANGELOG.md        # Detailed changelog of all versions
└── README.md           # This comprehensive documentation
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

### Recently Completed Features ✅

- [x] **Default Value Support**: Comprehensive default value handling for all types
- [x] **Bytes Encoding**: Multiple encoding formats (base64, hex) for bytes fields
- [x] **Repeated Field Support**: Full slice support with default values
- [x] **Wrapper Type Support**: Complete google.protobuf.*Value type support
- [x] **Map Field Enhancements**: JSON and native format support
- [x] **Duration and Timestamp Slices**: Repeated well-known type support

### Upcoming Features

- [ ] **Validation Integration**: Built-in validation rules and constraints
- [ ] **Custom Type Plugins**: Extensible type system for custom types
- [ ] **Configuration Files**: Support for config file generation (YAML, JSON, TOML)
- [ ] **Environment Variables**: Automatic env var binding with fallback
- [ ] **Flag Groups**: Organize flags into logical groups for help output
- [ ] **Completion Scripts**: Generate shell completion scripts (bash, zsh, fish)
- [ ] **Web UI**: Optional web interface for configuration management
- [ ] **Flag Constraints**: Add constraints between flags (e.g., mutually exclusive flags)

## ⚖️ Comparison with Alternatives

| Feature | protoc-gen-flags | Manual Flag Binding | Other Code Generators |
|---------|------------------|-------------------|----------------------|
| **Type Safety** | ✅ Strong typing | ⚠️ Manual validation | ✅ Varies |
| **Protobuf Integration** | ✅ Native | ❌ Manual mapping | ⚠️ Limited |
| **Code Generation** | ✅ Automatic | ❌ Manual | ✅ Varies |
| **Well-Known Types** | ✅ Full support (Duration, Timestamp, Wrappers) | ❌ Manual handling | ⚠️ Limited |
| **Nested Messages** | ✅ Hierarchical with prefixes | ❌ Complex | ⚠️ Limited |
| **Default Values** | ✅ Comprehensive for all types | ⚠️ Manual implementation | ❌ Limited |
| **Repeated Fields** | ✅ Full slice support with defaults | ❌ Manual parsing | ⚠️ Limited |
| **Bytes Encoding** | ✅ Base64/Hex with validation | ❌ Manual encoding | ❌ No support |
| **Map Fields** | ✅ JSON and native formats | ❌ Manual parsing | ⚠️ Limited |
| **Wrapper Types** | ✅ All protobuf wrappers | ❌ Manual handling | ❌ No support |
| **Maintenance** | ✅ Low effort | ❌ High effort | ⚠️ Varies |
| **Documentation** | ✅ Auto-generated help | ⚠️ Manual updates | ⚠️ Varies |

## 📈 Changelog

See [CHANGELOG.md](CHANGELOG.md) for a detailed history of changes and improvements.

## 🏷️ Versioning

This project follows [Semantic Versioning](https://semver.org/). For the versions available, see the [tags on this repository](https://github.com/kunstack/protoc-gen-flags/tags).

