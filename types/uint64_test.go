package types

import (
	"testing"

	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestUInt64Value_String(t *testing.T) {
	tests := []struct {
		name     string
		value    uint64
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
			name:     "maximum uint64",
			value:    18446744073709551615,
			expected: "18446744073709551615",
		},
		{
			name:     "large positive",
			value:    123456789012345,
			expected: "123456789012345",
		},
		{
			name:     "very large positive",
			value:    18446744073709551614,
			expected: "18446744073709551614",
		},
		{
			name:     "one",
			value:    1,
			expected: "1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UInt64Value{Value: tt.value}
			if got := u.String(); got != tt.expected {
				t.Errorf("UInt64Value.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestUInt64Value_String_Nil(t *testing.T) {
	var u *UInt64Value
	got := u.String()
	expected := "0"
	if got != expected {
		t.Errorf("UInt64Value.String() nil = %v, want %v", got, expected)
	}
}

func TestUInt64Value_Set(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    uint64
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
			name:     "maximum uint64",
			input:    "18446744073709551615",
			expected: 18446744073709551615,
		},
		{
			name:     "large positive",
			input:    "123456789012345",
			expected: 123456789012345,
		},
		{
			name:     "very large positive",
			input:    "18446744073709551614",
			expected: 18446744073709551614,
		},
		{
			name:     "one",
			input:    "1",
			expected: 1,
		},
		{
			name:        "negative integer",
			input:       "-17",
			expectError: true,
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
			input:       "18446744073709551616",
			expectError: true,
		},
		{
			name:        "float input",
			input:       "3.14",
			expectError: true,
		},
		{
			name:        "hex input",
			input:       "0xff",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UInt64Value{}
			err := u.Set(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("UInt64Value.Set() expected error for input %q, got nil", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("UInt64Value.Set() unexpected error for input %q: %v", tt.input, err)
				return
			}

			if u.Value != tt.expected {
				t.Errorf("UInt64Value.Set() = %v, want %v", u.Value, tt.expected)
			}
		})
	}
}

func TestUInt64Value_Type(t *testing.T) {
	u := &UInt64Value{}
	if got := u.Type(); got != "uint64Value" {
		t.Errorf("UInt64Value.Type() = %v, want %v", got, "uint64Value")
	}
}

func TestUInt64Value_ImplementsPflagValue(t *testing.T) {
	var _ pflag.Value = (*UInt64Value)(nil)
}

func TestUInt64Value_PflagIntegration(t *testing.T) {
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	u := &UInt64Value{}

	fs.Var(u, "uint64", "Test uint64 value")

	// Test valid input
	err := fs.Parse([]string{"--uint64", "123456789012345"})
	if err != nil {
		t.Fatalf("Failed to parse uint64: %v", err)
	}

	expected := uint64(123456789012345)
	if u.Value != expected {
		t.Errorf("Parsed uint64 = %v, want %v", u.Value, expected)
	}

	// Test String() method after Set()
	if got := u.String(); got != "123456789012345" {
		t.Errorf("UInt64Value.String() after Set() = %v, want %v", got, "123456789012345")
	}
}

func TestUInt64Value_PflagIntegrationError(t *testing.T) {
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	u := &UInt64Value{}

	fs.Var(u, "uint64", "Test uint64 value")

	// Test invalid input
	err := fs.Parse([]string{"--uint64", "invalid"})
	if err == nil {
		t.Fatalf("Expected error for invalid uint64 input, got nil")
	}
}

func TestUInt64_Constructor(t *testing.T) {
	wrap := &wrapperspb.UInt64Value{Value: 123456789012345}
	u := UInt64(wrap)

	// Test that the returned UInt64Value has correct fields
	if (*wrapperspb.UInt64Value)(u) != wrap {
		t.Errorf("UInt64() wrapper = %p, want %p", (*wrapperspb.UInt64Value)(u), wrap)
	}

	// Test that the value is correct
	if u.Value != 123456789012345 {
		t.Errorf("UInt64().Value = %v, want %v", u.Value, 123456789012345)
	}

	// Test String() method on constructed uint64
	expected := "123456789012345"
	if got := u.String(); got != expected {
		t.Errorf("UInt64().String() = %v, want %v", got, expected)
	}
}

func TestUInt64_ConstructorNil(t *testing.T) {
	u := UInt64(nil)

	// Test that it works with nil wrapper
	if u != nil {
		t.Errorf("UInt64(nil) = %p, want nil", u)
	}
}

func TestUInt64Value_StringAfterSet(t *testing.T) {
	u := &UInt64Value{}

	// Set a value
	err := u.Set("18446744073709551615")
	if err != nil {
		t.Fatalf("Failed to set uint64 value: %v", err)
	}

	// Test String() method
	if got := u.String(); got != "18446744073709551615" {
		t.Errorf("UInt64Value.String() after Set() = %v, want %v", got, uint64(18446744073709551615))
	}
}

func TestUInt64Value_MultipleSets(t *testing.T) {
	u := &UInt64Value{}

	// Set first value
	err := u.Set("123456789012345")
	if err != nil {
		t.Fatalf("Failed to set first uint64 value: %v", err)
	}
	if u.Value != 123456789012345 {
		t.Errorf("First set failed: got %v, want %v", u.Value, 123456789012345)
	}

	// Set second value
	err = u.Set("18446744073709551615")
	if err != nil {
		t.Fatalf("Failed to set second uint64 value: %v", err)
	}
	if u.Value != 18446744073709551615 {
		t.Errorf("Second set failed: got %v, want %v", u.Value, uint64(18446744073709551615))
	}

	// Test final String() value
	if got := u.String(); got != "18446744073709551615" {
		t.Errorf("Final String() = %v, want %v", got, uint64(18446744073709551615))
	}
}

func TestUInt64Value_EdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  uint64
	}{
		{"zero", "0", 0},
		{"one", "1", 1},
		{"max uint64", "18446744073709551615", 18446744073709551615},
		{"max uint64 minus one", "18446744073709551614", 18446744073709551614},
		{"large positive", "123456789012345678", 123456789012345678},
		{"small positive", "42", 42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UInt64Value{}
			err := u.Set(tt.input)
			if err != nil {
				t.Fatalf("Failed to set uint64 value %q: %v", tt.input, err)
			}
			if u.Value != tt.want {
				t.Errorf("UInt64Value edge case: input=%q, got=%v, want=%v", tt.input, u.Value, tt.want)
			}
		})
	}
}

func TestUInt64Value_InvalidInputs(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty string", ""},
		{"text", "abc"},
		{"float", "3.14"},
		{"hex", "0x123"},
		{"negative", "-123"},
		{"negative zero", "-0"},
		{"comma", "1,234"},
		{"space", "123 456"},
		{"plus prefix", "+123"},
		{"double negative", "--123"},
		{"overflow", "18446744073709551616"},
		{"very large overflow", "999999999999999999999999999999"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UInt64Value{}
			err := u.Set(tt.input)
			if err == nil {
				t.Errorf("UInt64Value.Set() expected error for invalid input %q, got nil", tt.input)
			}
		})
	}
}

func TestUInt64Value_BoundaryValues(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  uint64
	}{
		{"zero", "0", 0},
		{"one", "1", 1},
		{"max uint64", "18446744073709551615", 18446744073709551615},
		{"max uint64 minus one", "18446744073709551614", 18446744073709551614},
		{"large but valid", "9223372036854775808", 9223372036854775808}, // > max int64
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UInt64Value{}
			err := u.Set(tt.input)
			if err != nil {
				t.Fatalf("Failed to set uint64 value %q: %v", tt.input, err)
			}
			if u.Value != tt.want {
				t.Errorf("UInt64Value boundary test: input=%q, got=%v, want=%v", tt.input, u.Value, tt.want)
			}
		})
	}
}
