# Nested Configuration Example

This example demonstrates how to use nested messages to create hierarchical configuration structures.

## Overview

This example shows how to:
- Define nested protobuf messages
- Configure nested message flag generation
- Use hierarchical flag naming
- Organize complex configurations into logical sections

## Project Structure

```
nested/
├── go.mod
├── main.go
├── buf.yaml
├── buf.gen.yaml
├── Makefile
├── README.md
└── proto/
    └── config.proto    # Hierarchical configuration with nested messages
```

## Configuration Structure

The example defines an `AppConfig` message with three nested sections:

### ServerConfig
- `host` - Server host address
- `port` - Server port number
- `https` - Enable HTTPS

### DatabaseConfig
- `url` - Database connection URL
- `max_connections` - Maximum connection pool size

### LoggingConfig
- `level` - Log level
- `format` - Log format (json/text)

## Usage

### Build and Run

```bash
# Generate code and build
make build

# Run with defaults
./bin/app

# Run with custom configuration
./bin/app \
  --server.host 0.0.0.0 \
  --server.port 9000 \
  --server.https \
  --db.database-url postgres://localhost/mydb \
  --db.database-max-connections 50 \
  --log.logging-level debug \
  --log.logging-format json

# Show help to see all available flags
./bin/app --help
```

### Expected Output

With defaults:
```
Application Configuration:
=========================

Server:
  Host: localhost
  Port: 8080
  HTTPS: false

Database:
  URL: postgres://localhost:5432/defaultdb
  Max Connections: 10

Logging:
  Level: info
  Format: text

Configuration loaded successfully!
```

## Key Concepts

### 1. Nested Message Definition

In the protobuf file, nested messages are defined as separate message types:

```protobuf
message ServerConfig {
  string host = 1 [(flags.value).string = { ... }];
  // ...
}

message AppConfig {
  ServerConfig server = 1 [(flags.value).message = {
    name: "server"
    nested: true  // Enable nested flag generation
  }];
}
```

### 2. Hierarchical Flag Names

When `nested: true` is set, flags are generated with prefixes:
- `--server.host` (from server.host field)
- `--db.database-url` (from database.url field)
- `--log.logging-level` (from logging.level field)

The prefix comes from the field name in the parent message.

### 3. Organizing Complex Configuration

Nested messages help organize related configuration into logical sections:
- Server settings
- Database settings
- Logging settings
- Authentication settings
- etc.

This makes configuration more maintainable and easier to understand.

## Benefits

1. **Clear Organization**: Related configuration grouped together
2. **Namespace Isolation**: Prevents flag name conflicts
3. **Reusability**: Nested messages can be reused across different configurations
4. **Type Safety**: Each section maintains strong typing

## Next Steps

After understanding nested configuration, explore:
- **advanced/** - All supported protobuf types
- **hierarchical/** - Custom prefixes and delimiters for even more control