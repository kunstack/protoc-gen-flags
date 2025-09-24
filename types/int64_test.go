package types

import (
	"testing"

	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestInt64Value_String(t *testing.T) {
	tests := []struct {
		name     string
		value    int64
		expected string
	}{
		{
			name:     "zero value",
			value:    0,
			expected: "0",
		},
		{
			name:     "positive integer",
			value:    42,
			expected: "42",
		},
		{
			name:     "negative integer",
			value:    -17,
			expected: "-17",
		},
		{
			name:     "maximum int64",
			value:    9223372036854775807,
			expected: "9223372036854775807",
		},
		{
			name:     "minimum int64",
			value:    -9223372036854775808,
			expected: "-9223372036854775808",
		},
		{
			name:     "large positive",
			value:    123456789012345,
			expected: "123456789012345",
		},
		{
			name:     "large negative",
			value:    -987654321098765,
			expected: "-987654321098765",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Int64Value{Value: tt.value}
			if got := i.String(); got != tt.expected {
				t.Errorf("Int64Value.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestInt64Value_String_Nil(t *testing.T) {
	var i *Int64Value
	got := i.String()
	expected := "0"
	if got != expected {
		t.Errorf("Int64Value.String() nil = %v, want %v", got, expected)
	}
}

func TestInt64Value_Set(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    int64
		expectError bool
	}{
		{
			name:     "zero",
			input:    "0",
			expected: 0,
		},
		{
			name:     "positive integer",
			input:    "42",
			expected: 42,
		},
		{
			name:     "negative integer",
			input:    "-17",
			expected: -17,
		},
		{
			name:     "maximum int64",
			input:    "9223372036854775807",
			expected: 9223372036854775807,
		},
		{
			name:     "minimum int64",
			input:    "-9223372036854775808",
			expected: -9223372036854775808,
		},
		{
			name:     "large positive",
			input:    "123456789012345",
			expected: 123456789012345,
		},
		{
			name:     "large negative",
			input:    "-987654321098765",
			expected: -987654321098765,
		},
		{
			name:        "invalid input",
			input:       "not a number",
			expectError: true,
		},
		{
			name:        "empty string",
			input:       "",
			expectError: true,
		},
		{
			name:        "overflow",
			input:       "9223372036854775808",
			expectError: true,
		},
		{
			name:        "underflow",
			input:       "-9223372036854775809",
			expectError: true,
		},
		{
			name:        "float input",
			input:       "3.14",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Int64Value{}
			err := i.Set(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("Int64Value.Set() expected error for input %q, got nil", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("Int64Value.Set() unexpected error for input %q: %v", tt.input, err)
				return
			}

			if i.Value != tt.expected {
				t.Errorf("Int64Value.Set() = %v, want %v", i.Value, tt.expected)
			}
		})
	}
}

func TestInt64Value_Type(t *testing.T) {
	i := &Int64Value{}
	if got := i.Type(); got != "int64Value" {
		t.Errorf("Int64Value.Type() = %v, want %v", got, "int64Value")
	}
}

func TestInt64Value_ImplementsPflagValue(t *testing.T) {
	var _ pflag.Value = (*Int64Value)(nil)
}

func TestInt64Value_PflagIntegration(t *testing.T) {
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	i := &Int64Value{}

	fs.Var(i, "int64", "Test int64 value")

	// Test valid input
	err := fs.Parse([]string{"--int64", "123456789012345"})
	if err != nil {
		t.Fatalf("Failed to parse int64: %v", err)
	}

	expected := int64(123456789012345)
	if i.Value != expected {
		t.Errorf("Parsed int64 = %v, want %v", i.Value, expected)
	}

	// Test String() method after Set()
	if got := i.String(); got != "123456789012345" {
		t.Errorf("Int64Value.String() after Set() = %v, want %v", got, "123456789012345")
	}
}

func TestInt64Value_PflagIntegrationError(t *testing.T) {
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	i := &Int64Value{}

	fs.Var(i, "int64", "Test int64 value")

	// Test invalid input
	err := fs.Parse([]string{"--int64", "invalid"})
	if err == nil {
		t.Fatalf("Expected error for invalid int64 input, got nil")
	}
}

func TestInt64_Constructor(t *testing.T) {
	wrap := &wrapperspb.Int64Value{Value: 123456789012345}
	i := Int64(wrap)

	// Test that the returned Int64Value has correct fields
	if (*wrapperspb.Int64Value)(i) != wrap {
		t.Errorf("Int64() wrapper = %p, want %p", (*wrapperspb.Int64Value)(i), wrap)
	}

	// Test that the value is correct
	if i.Value != 123456789012345 {
		t.Errorf("Int64().Value = %v, want %v", i.Value, 123456789012345)
	}

	// Test String() method on constructed int64
	expected := "123456789012345"
	if got := i.String(); got != expected {
		t.Errorf("Int64().String() = %v, want %v", got, expected)
	}
}

func TestInt64_ConstructorNil(t *testing.T) {
	i := Int64(nil)

	// Test that it works with nil wrapper
	if i != nil {
		t.Errorf("Int64(nil) = %p, want nil", i)
	}
}

func TestInt64Value_StringAfterSet(t *testing.T) {
	i := &Int64Value{}

	// Set a value
	err := i.Set("-987654321098765")
	if err != nil {
		t.Fatalf("Failed to set int64 value: %v", err)
	}

	// Test String() method
	if got := i.String(); got != "-987654321098765" {
		t.Errorf("Int64Value.String() after Set() = %v, want %v", got, "-987654321098765")
	}
}

func TestInt64Value_MultipleSets(t *testing.T) {
	i := &Int64Value{}

	// Set first value
	err := i.Set("123456789012345")
	if err != nil {
		t.Fatalf("Failed to set first int64 value: %v", err)
	}
	if i.Value != 123456789012345 {
		t.Errorf("First set failed: got %v, want %v", i.Value, 123456789012345)
	}

	// Set second value
	err = i.Set("-987654321098765")
	if err != nil {
		t.Fatalf("Failed to set second int64 value: %v", err)
	}
	if i.Value != -987654321098765 {
		t.Errorf("Second set failed: got %v, want %v", i.Value, -987654321098765)
	}

	// Test final String() value
	if got := i.String(); got != "-987654321098765" {
		t.Errorf("Final String() = %v, want %v", got, "-987654321098765")
	}
}

func TestInt64Value_EdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int64
	}{
		{"zero", "0", 0},
		{"positive one", "1", 1},
		{"negative one", "-1", -1},
		{"max int64", "9223372036854775807", 9223372036854775807},
		{"min int64", "-9223372036854775808", -9223372036854775808},
		{"large positive", "123456789012345678", 123456789012345678},
		{"large negative", "-123456789012345678", -123456789012345678},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Int64Value{}
			err := i.Set(tt.input)
			if err != nil {
				t.Fatalf("Failed to set int64 value %q: %v", tt.input, err)
			}
			if i.Value != tt.want {
				t.Errorf("Int64Value edge case: input=%q, got=%v, want=%v", tt.input, i.Value, tt.want)
			}
		})
	}
}

func TestInt64Value_InvalidInputs(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty string", ""},
		{"text", "abc"},
		{"float", "3.14"},
		{"hex", "0x123"},
		{"comma", "1,234"},
		{"space", "123 456"},
		{"double negative", "--123"},
		{"overflow", "9223372036854775808"},
		{"underflow", "-9223372036854775809"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Int64Value{}
			err := i.Set(tt.input)
			if err == nil {
				t.Errorf("Int64Value.Set() expected error for invalid input %q, got nil", tt.input)
			}
		})
	}
}
