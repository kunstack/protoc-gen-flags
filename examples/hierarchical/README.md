# Hierarchical Flags Example

This example demonstrates how to use prefixes and custom delimiters to organize flags hierarchically.

## Overview

This example shows how to:
- Add prefixes to flag names
- Use different delimiter styles (dot, dash, underscore, colon)
- Create multi-level hierarchies
- Use custom renaming functions
- Organize large flag sets effectively

## Project Structure

```
hierarchical/
├── go.mod
├── main.go
├── buf.yaml
├── buf.gen.yaml
├── Makefile
├── README.md
└── proto/
    └── config.proto
```

## Configuration

The example defines a simple `ServiceConfig` message with basic fields, but demonstrates various ways to organize these flags using prefixes and delimiters.

## Usage

### Build and Run

```bash
# Generate code and build
make build

# Run with different delimiter styles
./bin/hierarchical

# The application demonstrates 5 different styles:
# 1. No prefix (default)
# 2. Dot delimiter (default)
# 3. Dash delimiter
# 4. Underscore delimiter
# 5. Colon delimiter
```

## Hierarchical Organization Styles

### 1. No Prefix (Default)

```go
config.AddFlags(fs)
```

Generates:
```bash
--host
--port
--timeout
```

### 2. Dot Delimiter (Default with Prefix)

```go
config.AddFlags(fs, flags.WithPrefix("service"))
```

Generates:
```bash
--service.host
--service.port
--service.timeout
```

### 3. Dash Delimiter

```go
config.AddFlags(fs,
    flags.WithPrefix("service"),
    flags.WithDelimiter("-"))
```

Generates:
```bash
--service-host
--service-port
--service-timeout
```

### 4. Underscore Delimiter

```go
config.AddFlags(fs,
    flags.WithPrefix("service"),
    flags.WithDelimiter("_"))
```

Generates:
```bash
--service_host
--service_port
--service_timeout
```

### 5. Colon Delimiter

```go
config.AddFlags(fs,
    flags.WithPrefix("service"),
    flags.WithDelimiter(":"))
```

Generates:
```bash
--service:host
--service:port
--service:timeout
```

## Multi-Level Prefixes

You can create multi-level hierarchies:

```go
config.AddFlags(fs,
    flags.WithPrefix("myapp", "backend", "service"))
```

Generates (with dot delimiter):
```bash
--myapp.backend.service.host
--myapp.backend.service.port
--myapp.backend.service.timeout
```

## Custom Renaming

Use custom renaming functions:

```go
config.AddFlags(fs,
    flags.WithPrefix("Service"),
    flags.WithRenamer(strings.ToLower))
```

Generates:
```bash
--service.host
--service.port
--service.timeout
```

## Use Cases

### 1. Microservices Architecture

Organize flags by service:

```go
// API service
apiConfig.AddFlags(fs, flags.WithPrefix("api"))
// --api.host, --api.port

// Database service
dbConfig.AddFlags(fs, flags.WithPrefix("db"))
// --db.host, --db.port

// Cache service
cacheConfig.AddFlags(fs, flags.WithPrefix("cache"))
// --cache.host, --cache.port
```

### 2. Environment-Specific Configuration

```go
// Production
prodConfig.AddFlags(fs, flags.WithPrefix("prod"))
// --prod.host, --prod.port

// Staging
stagingConfig.AddFlags(fs, flags.WithPrefix("staging"))
// --staging.host, --staging.port

// Development
devConfig.AddFlags(fs, flags.WithPrefix("dev"))
// --dev.host, --dev.port
```

### 3. Component-Based Organization

```go
// Frontend component
frontendConfig.AddFlags(fs,
    flags.WithPrefix("frontend"),
    flags.WithDelimiter("-"))
// --frontend-host, --frontend-port

// Backend component
backendConfig.AddFlags(fs,
    flags.WithPrefix("backend"),
    flags.WithDelimiter("-"))
// --backend-host, --backend-port

// Database component
databaseConfig.AddFlags(fs,
    flags.WithPrefix("database"),
    flags.WithDelimiter("-"))
// --database-host, --database-port
```

## Available Delimiters

The `flags` package provides these delimiter constants:

- `flags.DelimiterDot` - "." (default)
- `flags.DelimiterDash` - "-"
- `flags.DelimiterUnderscore` - "_"
- `flags.DelimiterColon` - ":"

You can also use custom delimiters:

```go
config.AddFlags(fs,
    flags.WithPrefix("service"),
    flags.WithDelimiter("::"))
```

## Benefits

1. **Namespace Isolation**: Prevent flag name conflicts
2. **Clear Organization**: Group related flags together
3. **Consistency**: Maintain consistent naming patterns
4. **Scalability**: Easy to add new flag groups
5. **Flexibility**: Choose delimiter style that matches your conventions

## Best Practices

1. **Choose Consistent Delimiters**: Use the same delimiter throughout your application
2. **Keep Prefixes Short**: Use concise, meaningful prefixes
3. **Use Dot for Multi-Level**: Dots work well for deeply nested hierarchies
4. **Use Dash for Kubernetes**: Dashes are common in Kubernetes configurations
5. **Document Your Style**: Make your delimiter choice clear in documentation

## Integration with Environment Variables

When using hierarchical flags with viper, you can map them to environment variables:

```go
import "github.com/spf13/viper"

// Use uppercase with underscores for env vars
viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
viper.AutomaticEnv()

// Maps: --service.host -> SERVICE_HOST
```

## Next Steps

Combine hierarchical organization with:
- Nested messages for complex configurations
- Configuration files (YAML, JSON, TOML)
- Environment variables via viper
- Feature flags and runtime configuration
