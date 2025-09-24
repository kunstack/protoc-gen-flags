package types

import (
	"testing"

	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestUInt32Value_String(t *testing.T) {
	tests := []struct {
		name     string
		value    uint32
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
			name:     "maximum uint32",
			value:    4294967295,
			expected: "4294967295",
		},
		{
			name:     "large positive",
			value:    3147483647,
			expected: "3147483647",
		},
		{
			name:     "one",
			value:    1,
			expected: "1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UInt32Value{Value: tt.value}
			if got := u.String(); got != tt.expected {
				t.Errorf("UInt32Value.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestUInt32Value_String_Nil(t *testing.T) {
	var u *UInt32Value
	got := u.String()
	expected := "0"
	if got != expected {
		t.Errorf("UInt32Value.String() nil = %v, want %v", got, expected)
	}
}

func TestUInt32Value_Set(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    uint32
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
			name:     "maximum uint32",
			input:    "4294967295",
			expected: 4294967295,
		},
		{
			name:     "large positive",
			input:    "3147483647",
			expected: 3147483647,
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
			input:       "4294967296",
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
			u := &UInt32Value{}
			err := u.Set(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("UInt32Value.Set() expected error for input %q, got nil", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("UInt32Value.Set() unexpected error for input %q: %v", tt.input, err)
				return
			}

			if u.Value != tt.expected {
				t.Errorf("UInt32Value.Set() = %v, want %v", u.Value, tt.expected)
			}
		})
	}
}

func TestUInt32Value_Type(t *testing.T) {
	u := &UInt32Value{}
	if got := u.Type(); got != "uint32Value" {
		t.Errorf("UInt32Value.Type() = %v, want %v", got, "uint32Value")
	}
}

func TestUInt32Value_ImplementsPflagValue(t *testing.T) {
	var _ pflag.Value = (*UInt32Value)(nil)
}

func TestUInt32Value_PflagIntegration(t *testing.T) {
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	u := &UInt32Value{}

	fs.Var(u, "uint32", "Test uint32 value")

	// Test valid input
	err := fs.Parse([]string{"--uint32", "123456789"})
	if err != nil {
		t.Fatalf("Failed to parse uint32: %v", err)
	}

	expected := uint32(123456789)
	if u.Value != expected {
		t.Errorf("Parsed uint32 = %v, want %v", u.Value, expected)
	}

	// Test String() method after Set()
	if got := u.String(); got != "123456789" {
		t.Errorf("UInt32Value.String() after Set() = %v, want %v", got, "123456789")
	}
}

func TestUInt32Value_PflagIntegrationError(t *testing.T) {
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	u := &UInt32Value{}

	fs.Var(u, "uint32", "Test uint32 value")

	// Test invalid input
	err := fs.Parse([]string{"--uint32", "invalid"})
	if err == nil {
		t.Fatalf("Expected error for invalid uint32 input, got nil")
	}
}

func TestUInt32_Constructor(t *testing.T) {
	wrap := &wrapperspb.UInt32Value{Value: 123456789}
	u := UInt32(wrap)

	// Test that the returned UInt32Value has correct fields
	if (*wrapperspb.UInt32Value)(u) != wrap {
		t.Errorf("UInt32() wrapper = %p, want %p", (*wrapperspb.UInt32Value)(u), wrap)
	}

	// Test that the value is correct
	if u.Value != uint32(123456789) {
		t.Errorf("UInt32().Value = %v, want %v", u.Value, uint32(123456789))
	}

	// Test String() method on constructed uint32
	expected := "123456789"
	if got := u.String(); got != expected {
		t.Errorf("UInt32().String() = %v, want %v", got, expected)
	}
}

func TestUInt32_ConstructorNil(t *testing.T) {
	u := UInt32(nil)

	// Test that it works with nil wrapper
	if u != nil {
		t.Errorf("UInt32(nil) = %p, want nil", u)
	}
}

func TestUInt32Value_StringAfterSet(t *testing.T) {
	u := &UInt32Value{}

	// Set a value
	err := u.Set("4294967295")
	if err != nil {
		t.Fatalf("Failed to set uint32 value: %v", err)
	}

	// Test String() method
	if got := u.String(); got != "4294967295" {
		t.Errorf("UInt32Value.String() after Set() = %v, want %v", got, "4294967295")
	}
}

func TestUInt32Value_MultipleSets(t *testing.T) {
	u := &UInt32Value{}

	// Set first value
	err := u.Set("123456789")
	if err != nil {
		t.Fatalf("Failed to set first uint32 value: %v", err)
	}
	if u.Value != uint32(123456789) {
		t.Errorf("First set failed: got %v, want %v", u.Value, uint32(123456789))
	}

	// Set second value
	err = u.Set("4294967295")
	if err != nil {
		t.Fatalf("Failed to set second uint32 value: %v", err)
	}
	if u.Value != uint32(4294967295) {
		t.Errorf("Second set failed: got %v, want %v", u.Value, uint32(4294967295))
	}

	// Test final String() value
	if got := u.String(); got != "4294967295" {
		t.Errorf("Final String() = %v, want %v", got, "4294967295")
	}
}

func TestUInt32Value_EdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  uint32
	}{
		{"zero", "0", 0},
		{"one", "1", 1},
		{"max uint32", "4294967295", 4294967295},
		{"max uint32 minus one", "4294967294", 4294967294},
		{"large positive", "2147483648", 2147483648},
		{"small positive", "42", 42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UInt32Value{}
			err := u.Set(tt.input)
			if err != nil {
				t.Fatalf("Failed to set uint32 value %q: %v", tt.input, err)
			}
			if u.Value != tt.want {
				t.Errorf("UInt32Value edge case: input=%q, got=%v, want=%v", tt.input, u.Value, tt.want)
			}
		})
	}
}

func TestUInt32Value_InvalidInputs(t *testing.T) {
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
		{"overflow", "4294967296"},
		{"very large overflow", "999999999999999999999999999999"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UInt32Value{}
			err := u.Set(tt.input)
			if err == nil {
				t.Errorf("UInt32Value.Set() expected error for invalid input %q, got nil", tt.input)
			}
		})
	}
}

func TestUInt32Value_BoundaryValues(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  uint32
	}{
		{"zero", "0", 0},
		{"one", "1", 1},
		{"max uint32", "4294967295", 4294967295},
		{"max uint32 minus one", "4294967294", 4294967294},
		{"large but valid", "2147483648", 2147483648},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &UInt32Value{}
			err := u.Set(tt.input)
			if err != nil {
				t.Fatalf("Failed to set uint32 value %q: %v", tt.input, err)
			}
			if u.Value != tt.want {
				t.Errorf("UInt32Value boundary test: input=%q, got=%v, want=%v", tt.input, u.Value, tt.want)
			}
		})
	}
}
