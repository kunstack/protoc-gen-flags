// Copyright 2021 Aapeli <aapeli.nian@gmail.com> All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package flags

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDelimiterConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{
			name:     "DelimiterDot",
			constant: DelimiterDot,
			expected: ".",
		},
		{
			name:     "DelimiterDash",
			constant: DelimiterDash,
			expected: "-",
		},
		{
			name:     "DelimiterUnderscore",
			constant: DelimiterUnderscore,
			expected: "_",
		},
		{
			name:     "DelimiterColon",
			constant: DelimiterColon,
			expected: ":",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.constant)
		})
	}
}

func TestWithDelimiter(t *testing.T) {
	t.Run("default delimiter", func(t *testing.T) {
		opts := &Options{}
		WithDelimiter(DelimiterDot)(opts)
		assert.Equal(t, DelimiterDot, opts.Delimiter)
	})

	t.Run("custom delimiter dash", func(t *testing.T) {
		opts := &Options{}
		WithDelimiter(DelimiterDash)(opts)
		assert.Equal(t, DelimiterDash, opts.Delimiter)
	})

	t.Run("custom delimiter underscore", func(t *testing.T) {
		opts := &Options{}
		WithDelimiter(DelimiterUnderscore)(opts)
		assert.Equal(t, DelimiterUnderscore, opts.Delimiter)
	})

	t.Run("custom delimiter colon", func(t *testing.T) {
		opts := &Options{}
		WithDelimiter(DelimiterColon)(opts)
		assert.Equal(t, DelimiterColon, opts.Delimiter)
	})

	t.Run("custom delimiter custom value", func(t *testing.T) {
		opts := &Options{}
		WithDelimiter("::")(opts)
		assert.Equal(t, "::", opts.Delimiter)
	})
}

func TestWithRenamer(t *testing.T) {
	t.Run("identity renamer", func(t *testing.T) {
		opts := &Options{}
		renamer := func(s string) string { return s }
		WithRenamer(renamer)(opts)
		assert.NotNil(t, opts.Renamer)
	})

	t.Run("lowercase renamer", func(t *testing.T) {
		opts := &Options{}
		WithRenamer(strings.ToLower)(opts)
		assert.NotNil(t, opts.Renamer)
		result := opts.Renamer("TEST")
		assert.Equal(t, "test", result)
	})

	t.Run("uppercase renamer", func(t *testing.T) {
		opts := &Options{}
		WithRenamer(strings.ToUpper)(opts)
		assert.NotNil(t, opts.Renamer)
		result := opts.Renamer("test")
		assert.Equal(t, "TEST", result)
	})

	t.Run("custom renamer", func(t *testing.T) {
		opts := &Options{}
		customRenamer := func(s string) string { return "prefix_" + s }
		WithRenamer(customRenamer)(opts)
		assert.NotNil(t, opts.Renamer)
		result := opts.Renamer("value")
		assert.Equal(t, "prefix_value", result)
	})
}

func TestWithPrefix(t *testing.T) {
	t.Run("single prefix", func(t *testing.T) {
		opts := &Options{}
		WithPrefix("server")(opts)
		assert.Equal(t, []string{"server"}, opts.Prefix)
	})

	t.Run("multiple prefixes", func(t *testing.T) {
		opts := &Options{}
		WithPrefix("server", "database")(opts)
		assert.Equal(t, []string{"server", "database"}, opts.Prefix)
	})

	t.Run("prefix with dots preserved until Build", func(t *testing.T) {
		// WithPrefix no longer trims - it's deferred to Build()
		opts := &Options{}
		WithPrefix(".server.", ".database.")(opts)
		// Dots are preserved in Options
		assert.Equal(t, []string{".server.", ".database."}, opts.Prefix)

		// But Build() will trim them based on delimiter
		builder := NewNameBuilder(WithPrefix(".server."), WithDelimiter(DelimiterDot))
		result := builder.Build("port")
		assert.Equal(t, "server.port", result)
	})

	t.Run("empty prefix removed", func(t *testing.T) {
		opts := &Options{}
		WithPrefix("", "server", "", "database", "")(opts)
		assert.Equal(t, []string{"server", "database"}, opts.Prefix)
	})

	t.Run("only empty prefixes", func(t *testing.T) {
		opts := &Options{}
		WithPrefix("", "", "")(opts)
		assert.Empty(t, opts.Prefix)
	})

	t.Run("nested prefix", func(t *testing.T) {
		opts := &Options{}
		WithPrefix("server.database.pool")(opts)
		assert.Equal(t, []string{"server.database.pool"}, opts.Prefix)
	})
}

func TestNameBuilder_Build(t *testing.T) {
	t.Run("simple name without prefix", func(t *testing.T) {
		builder := NewNameBuilder()
		result := builder.Build("port")
		assert.Equal(t, "port", result)
	})

	t.Run("name with single prefix", func(t *testing.T) {
		builder := NewNameBuilder(WithPrefix("server"))
		result := builder.Build("port")
		assert.Equal(t, "server.port", result)
	})

	t.Run("name with multiple prefixes", func(t *testing.T) {
		builder := NewNameBuilder(WithPrefix("server", "database"))
		result := builder.Build("port")
		assert.Equal(t, "server.database.port", result)
	})

	t.Run("name with custom delimiter", func(t *testing.T) {
		builder := NewNameBuilder(WithDelimiter(DelimiterDash))
		result := builder.Build("port")
		assert.Equal(t, "port", result)
	})

	t.Run("name with prefix and custom delimiter", func(t *testing.T) {
		builder := NewNameBuilder(
			WithPrefix("server"),
			WithDelimiter(DelimiterDash),
		)
		result := builder.Build("port")
		assert.Equal(t, "server-port", result)
	})

	t.Run("name with prefix and underscore delimiter", func(t *testing.T) {
		builder := NewNameBuilder(
			WithPrefix("server", "database"),
			WithDelimiter(DelimiterUnderscore),
		)
		result := builder.Build("port")
		assert.Equal(t, "server_database_port", result)
	})

	t.Run("name with prefix and colon delimiter", func(t *testing.T) {
		builder := NewNameBuilder(
			WithPrefix("server", "database"),
			WithDelimiter(DelimiterColon),
		)
		result := builder.Build("port")
		assert.Equal(t, "server:database:port", result)
	})

	t.Run("name with renamer lowercase", func(t *testing.T) {
		builder := NewNameBuilder(
			WithPrefix("Server"),
			WithRenamer(strings.ToLower),
		)
		result := builder.Build("Port")
		assert.Equal(t, "server.port", result)
	})

	t.Run("name with renamer uppercase", func(t *testing.T) {
		builder := NewNameBuilder(
			WithPrefix("server"),
			WithRenamer(strings.ToUpper),
		)
		result := builder.Build("port")
		assert.Equal(t, "SERVER.PORT", result)
	})

	t.Run("name with prefix, delimiter, and renamer", func(t *testing.T) {
		builder := NewNameBuilder(
			WithPrefix("Server", "Database"),
			WithDelimiter(DelimiterDash),
			WithRenamer(strings.ToLower),
		)
		result := builder.Build("Port")
		assert.Equal(t, "server-database-port", result)
	})

	t.Run("name with custom renamer function", func(t *testing.T) {
		builder := NewNameBuilder(
			WithPrefix("api"),
			WithRenamer(func(s string) string { return "v1_" + s }),
		)
		result := builder.Build("endpoint")
		assert.Equal(t, "v1_api.endpoint", result)
	})

	t.Run("empty name", func(t *testing.T) {
		builder := NewNameBuilder(WithPrefix("server"))
		result := builder.Build("")
		assert.Equal(t, "server.", result)
	})

	t.Run("name with special characters", func(t *testing.T) {
		builder := NewNameBuilder(WithPrefix("api"))
		result := builder.Build("endpoint-v2")
		assert.Equal(t, "api.endpoint-v2", result)
	})
}

func TestNewNameBuilder(t *testing.T) {
	t.Run("no options uses defaults", func(t *testing.T) {
		builder := NewNameBuilder()
		assert.Equal(t, DelimiterDot, builder.options.Delimiter)
		assert.NotNil(t, builder.options.Renamer)
		assert.Empty(t, builder.options.Prefix)
	})

	t.Run("single option", func(t *testing.T) {
		builder := NewNameBuilder(WithDelimiter(DelimiterDash))
		assert.Equal(t, DelimiterDash, builder.options.Delimiter)
	})

	t.Run("multiple options", func(t *testing.T) {
		builder := NewNameBuilder(
			WithPrefix("server"),
			WithDelimiter(DelimiterDash),
			WithRenamer(strings.ToLower),
		)
		assert.Equal(t, DelimiterDash, builder.options.Delimiter)
		assert.Equal(t, []string{"server"}, builder.options.Prefix)
		assert.NotNil(t, builder.options.Renamer)
		// Test that renamer works
		result := builder.Build("Test")
		assert.Equal(t, "server-test", result)
	})

	t.Run("chained prefix calls", func(t *testing.T) {
		builder := NewNameBuilder(
			WithPrefix("server"),
			WithPrefix("database"),
		)
		assert.Equal(t, []string{"server", "database"}, builder.options.Prefix)
	})

	t.Run("renamer is identity by default", func(t *testing.T) {
		builder := NewNameBuilder()
		result := builder.Build("TestName")
		assert.Equal(t, "TestName", result)
	})

	t.Run("options are shared across builders if not careful", func(t *testing.T) {
		builder1 := NewNameBuilder(WithPrefix("server1"))
		builder2 := NewNameBuilder(WithPrefix("server2"))

		// Note: This test demonstrates that options are copied by reference
		// In practice, each NewNameBuilder call creates a new Options instance
		assert.NotEqual(t, builder1.options, builder2.options)
		assert.Equal(t, []string{"server1"}, builder1.options.Prefix)
		assert.Equal(t, []string{"server2"}, builder2.options.Prefix)
	})
}

func TestNameBuilderIntegration(t *testing.T) {
	t.Run("realistic server config scenario", func(t *testing.T) {
		// Simulating a server configuration with database settings
		serverBuilder := NewNameBuilder(
			WithPrefix("server"),
			WithDelimiter(DelimiterDot),
		)

		assert.Equal(t, "server.host", serverBuilder.Build("host"))
		assert.Equal(t, "server.port", serverBuilder.Build("port"))
		assert.Equal(t, "server.ssl.enabled", serverBuilder.Build("ssl.enabled"))
	})

	t.Run("microservice scenario with kebab-case", func(t *testing.T) {
		// Microservice environment often uses kebab-case
		serviceBuilder := NewNameBuilder(
			WithPrefix("payment-service"),
			WithDelimiter(DelimiterDash),
			WithRenamer(strings.ToLower),
		)

		assert.Equal(t, "payment-service-timeout", serviceBuilder.Build("Timeout"))
		assert.Equal(t, "payment-service-retrycount", serviceBuilder.Build("RetryCount"))
	})

	t.Run("snake_case configuration scenario", func(t *testing.T) {
		// Legacy systems might prefer snake_case
		legacyBuilder := NewNameBuilder(
			WithPrefix("legacy_system"),
			WithDelimiter(DelimiterUnderscore),
			WithRenamer(strings.ToLower),
		)

		assert.Equal(t, "legacy_system_maxconnections", legacyBuilder.Build("MaxConnections"))
		assert.Equal(t, "legacy_system_timeoutseconds", legacyBuilder.Build("TimeoutSeconds"))
	})

	t.Run("namespace scenario with colon delimiter", func(t *testing.T) {
		// Kubernetes-style namespacing
		nsBuilder := NewNameBuilder(
			WithPrefix("production"),
			WithDelimiter(DelimiterColon),
		)

		assert.Equal(t, "production:api:endpoint", nsBuilder.Build("api:endpoint"))
		assert.Equal(t, "production:config:file", nsBuilder.Build("config:file"))
	})
}

func TestEdgeCases(t *testing.T) {
	t.Run("builder with very long prefix chain", func(t *testing.T) {
		builder := NewNameBuilder(
			WithPrefix("a", "b", "c", "d", "e"),
			WithDelimiter(DelimiterDot),
		)
		result := builder.Build("value")
		assert.Equal(t, "a.b.c.d.e.value", result)
	})

	t.Run("builder with unicode in prefix", func(t *testing.T) {
		builder := NewNameBuilder(
			WithPrefix("服务"),
			WithDelimiter(DelimiterDot),
		)
		result := builder.Build("端口")
		assert.Equal(t, "服务.端口", result)
	})

	t.Run("empty string delimiter", func(t *testing.T) {
		builder := NewNameBuilder(
			WithPrefix("server"),
			WithDelimiter(""),
		)
		result := builder.Build("port")
		assert.Equal(t, "serverport", result)
	})

	t.Run("multi-character delimiter", func(t *testing.T) {
		builder := NewNameBuilder(
			WithPrefix("server"),
			WithDelimiter(" -> "),
		)
		result := builder.Build("port")
		assert.Equal(t, "server -> port", result)
	})

	t.Run("renamer that changes delimiter", func(t *testing.T) {
		builder := NewNameBuilder(
			WithPrefix("server"),
			WithDelimiter(DelimiterDot),
			WithRenamer(func(s string) string {
				return strings.ReplaceAll(s, ".", "-")
			}),
		)
		result := builder.Build("port")
		assert.Equal(t, "server-port", result)
	})
}
