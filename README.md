# protoc-gen-flags

*This project is currently in **alpha**. The API should be considered unstable and likely to change*

**protoc-gen-flags** is a protoc plugin generating the implementation of an `AddFlags` method on protobuf messages that integrates with the `spf13/pflag` library to provide POSIX/GNU-style command-line flag parsing.

```go
type Interface interface {
	AddFlags(fs *pflag.FlagSet, prefix string) error
}
```

## Installation

```bash
go get github.com/kunstack/protoc-gen-flags
```

## Usage

### Overview

**protoc-gen-flags** makes use of **Protobuf** options to define command-line flag bindings for protobuf message fields.

### Generation

**protoc-gen-flags** works the same way as other **protoc** plugins.

Example:
```bash
protoc -I. -I flags --go_out=paths=source_relative:. --flags_out=paths=source_relative:. types.proto
```

### Disable generation

Flag generation can be disabled with the `(flags.disabled) = true` message option.

```proto
message NoFlags {
    option (flags.disabled) = true;
    string string_field = 1 [(flags.value).string = {
        name: "string-field"
        usage: "This won't be generated"
    }];
}
```

### Unexported generation

Unexported flag methods can be generated with the `(flags.unexported) = true` message option.

```proto
message UnexportedFlags {
    option (flags.unexported) = true;
    string string_field = 1 [(flags.value).string = {
        name: "string-field"
        usage: "Generated method will be unexported"
    }];
}
```

### Allow empty generation

Flag generation can be allowed even without field configurations with the `(flags.allow_empty) = true` message option.

```proto
message EmptyFlags {
    option (flags.allow_empty) = true;
    string string_field = 1; // No flag configuration needed
}
```

### Scalar and Well-Known Types

Each scalar or Well-Known type has its corresponding `(flags.value).[type] = {name: string, short: string, usage: string, ...}` option.

All scalar types support the following options:
- `name`: Custom flag name (defaults to field name converted to kebab-case)
- `short`: Short flag alias (single character)
- `usage`: Description shown in help output
- `hidden`: Hide flag from help output
- `deprecated`: Mark flag as deprecated
- `deprecated_usage`: Additional context for deprecated flags

#### Numeric Types

- **float**:
    ```proto
    float float = 1 [(flags.value).float = {
        name: "float-value"
        short: "f"
        usage: "Float value parameter"
    }];
    ```
- **double**:
    ```proto
    double double = 2 [(flags.value).double = {
        name: "double-value"
        short: "d"
        usage: "Double value parameter"
    }];
    ```
- **int32**:
    ```proto
    int32 int32 = 3 [(flags.value).int32 = {
        name: "int32-value"
        short: "i"
        usage: "32-bit signed integer"
    }];
    ```
- **int64**:
    ```proto
    int64 int64 = 4 [(flags.value).int64 = {
        name: "int64-value"
        short: "I"
        usage: "64-bit signed integer"
    }];
    ```
- **uint32**:
    ```proto
    uint32 uint32 = 5 [(flags.value).uint32 = {
        name: "uint32-value"
        short: "u"
        usage: "32-bit unsigned integer"
    }];
    ```
- **uint64**:
    ```proto
    uint64 uint64 = 6 [(flags.value).uint64 = {
        name: "uint64-value"
        short: "U"
        usage: "64-bit unsigned integer"
    }];
    ```
- **sint32**:
    ```proto
    sint32 sint32 = 7 [(flags.value).sint32 = {
        name: "sint32-value"
        usage: "32-bit signed integer (zigzag encoded)"
    }];
    ```
- **sint64**:
    ```proto
    sint64 sint64 = 8 [(flags.value).sint64 = {
        name: "sint64-value"
        usage: "64-bit signed integer (zigzag encoded)"
    }];
    ```
- **fixed32**:
    ```proto
    fixed32 fixed32 = 9 [(flags.value).fixed32 = {
        name: "fixed32-value"
        usage: "32-bit fixed-point integer"
    }];
    ```
- **fixed64**:
    ```proto
    fixed64 fixed64 = 10 [(flags.value).fixed64 = {
        name: "fixed64-value"
        usage: "64-bit fixed-point integer"
    }];
    ```
- **sfixed32**:
    ```proto
    sfixed32 sfixed32 = 11 [(flags.value).sfixed32 = {
        name: "sfixed32-value"
        usage: "32-bit signed fixed-point integer"
    }];
    ```
- **sfixed64**:
    ```proto
    sfixed64 sfixed64 = 12 [(flags.value).sfixed64 = {
        name: "sfixed64-value"
        usage: "64-bit signed fixed-point integer"
    }];
    ```

#### Boolean and String Types

- **bool**:
    ```proto
    bool bool = 13 [(flags.value).bool = {
        name: "bool-value"
        short: "b"
        usage: "Boolean flag"
    }];
    ```
- **string**:
    ```proto
    string string = 14 [(flags.value).string = {
        name: "string-value"
        short: "s"
        usage: "String value parameter"
    }];
    ```

#### Bytes Type

Bytes fields support additional encoding options:

```proto
bytes bytes = 15 [(flags.value).bytes = {
    name: "bytes-value"
    short: "B"
    usage: "Bytes data"
    encoding: BYTES_ENCODING_TYPE_HEX  // or BYTES_ENCODING_TYPE_BASE64
}];
```

Supported encoding types:
- `BYTES_ENCODING_TYPE_UNSPECIFIED`: Default base64 encoding
- `BYTES_ENCODING_TYPE_BASE64`: Standard base64 encoding
- `BYTES_ENCODING_TYPE_HEX`: Hexadecimal encoding

#### Well-Known Types

- **google.protobuf.Duration**:
    ```proto
    google.protobuf.Duration duration = 16 [(flags.value).duration = {
        name: "duration"
        short: "D"
        usage: "Duration value (e.g., 30s, 5m, 1h)"
    }];
    ```
- **google.protobuf.Timestamp**:
    ```proto
    google.protobuf.Timestamp timestamp = 17 [(flags.value).timestamp = {
        name: "timestamp"
        short: "T"
        usage: "Timestamp value"
        formats: ["2006-01-02", "RFC3339", "ISO8601Time"]
    }];
    ```

    Timestamp fields support multiple time formats through the `formats` array. The parser will try each format in order until one successfully parses the input string.

    Supported format names:
    - **RFC339**: RFC3339 format (alias for RFC3339)
    - **RFC3339**: Standard RFC3339 format
    - **RFC3339Nano**: RFC3339 with nanosecond precision
    - **RFC822**: RFC822 format
    - **RFC822Z**: RFC822 with timezone
    - **RFC850**: RFC850 format
    - **RFC1123**: RFC1123 format
    - **RFC1123Z**: RFC1123 with timezone
    - **ISO8601**: ISO8601 date format (2006-01-02)
    - **ISO8601Time**: ISO8601 datetime format (2006-01-02T15:04:05)
    - **Kitchen**: Kitchen time format (3:04PM)
    - **Stamp**: Timestamp format (Jan _2 15:04:05)
    - **StampMilli**: Timestamp with milliseconds
    - **StampMicro**: Timestamp with microseconds
    - **StampNano**: Timestamp with nanoseconds
    - **DateTime**: Custom datetime format (2006-01-02 15:04:05)
    - **DateOnly**: Date only format (2006-01-02)
    - **TimeOnly**: Time only format (15:04:05)

    Custom Go time layout strings are also supported.
- **Wrapper types** (google.protobuf.*Value):
    ```proto
    google.protobuf.StringValue string_value = 18 [(flags.value).string = {
        name: "string-value"
        usage: "Optional string value"
    }];
    ```

### Enum Types

```proto
enum TestEnum {
    UNKNOWN = 0;
    VALUE1 = 1;
    VALUE2 = 2;
}

TestEnum enum = 19 [(flags.value).enum = {
    name: "enum-value"
    short: "e"
    usage: "Enum value parameter"
}];
```

### Message Types

Message fields can be configured to generate nested flags:

```proto
Message nested = 20 [(flags.value).message = {
    name: "nested"
    nested: true  // Generate flags for nested message fields
}];
```

When `nested: true`, the generated code will call `AddFlags` on the nested message if it implements the flags interface. The `name` field provides a prefix for nested flags (e.g., `--nested.field-name`).

### Repeated Types

Repeated fields are supported with type-specific configuration:

```proto
repeated string string_list = 21 [(flags.value).repeated.string = {
    name: "string-list"
    usage: "List of string values (can be specified multiple times)"
}];
```

All scalar types support repeated field configuration.

### oneof Fields

Fields within oneof blocks are fully supported:

```proto
oneof choice {
    string string_value = 22 [(flags.value).string = {
        name: "string-value"
        short: "s"
        usage: "String value in oneof"
    }];

    int32 int_value = 23 [(flags.value).int32 = {
        name: "int-value"
        short: "i"
        usage: "Integer value in oneof"
    }];

    google.protobuf.Duration duration_value = 24 [(flags.value).duration = {
        name: "duration-value"
        short: "d"
        usage: "Duration value in oneof"
    }];
}
```

## Complete Example

```proto
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

## Generated Code Usage

The plugin generates an `AddFlags(fs *pflag.FlagSet, prefix string) error` method for each configured message:

```go
package main

import (
    "flag"
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
    pflag.Parse()

    // Use the configuration
    fmt.Printf("Host: %s, Port: %d\n", config.GetHost(), config.GetPort())
}
```

## Development

### Building

```bash
make build
```

### Testing

```bash
make test
```

### Code Quality

```bash
make fmt      # Format Go code and protobuf files
make vet      # Run static analysis
make lint     # Run comprehensive linting
make tidy     # Clean up module dependencies
```

## Dependencies

- **github.com/lyft/protoc-gen-star**: Protocol buffer code generation framework
- **github.com/spf13/pflag**: POSIX/GNU-style command-line flag parsing
- **google.golang.org/protobuf**: Google's protocol buffer implementation

## License

This project is licensed under the Apache License, Version 2.0. See the LICENSE file for details.