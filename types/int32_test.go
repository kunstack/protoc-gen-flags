package types

import (
	"testing"

	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestInt32Value_String(t *testing.T) {
	tests := []struct {
		name     string
		value    int32
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
			name:     "maximum int32",
			value:    2147483647,
			expected: "2147483647",
		},
		{
			name:     "minimum int32",
			value:    -2147483648,
			expected: "-2147483648",
		},
		{
			name:     "large positive",
			value:    123456789,
			expected: "123456789",
		},
		{
			name:     "large negative",
			value:    -987654321,
			expected: "-987654321",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Int32Value{Value: tt.value}
			if got := i.String(); got != tt.expected {
				t.Errorf("Int32Value.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestInt32Value_String_Nil(t *testing.T) {
	var i *Int32Value
	got := i.String()
	expected := "0"
	if got != expected {
		t.Errorf("Int32Value.String() nil = %v, want %v", got, expected)
	}
}

func TestInt32Value_Set(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    int32
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
			name:     "maximum int32",
			input:    "2147483647",
			expected: 2147483647,
		},
		{
			name:     "minimum int32",
			input:    "-2147483648",
			expected: -2147483648,
		},
		{
			name:     "large positive",
			input:    "123456789",
			expected: 123456789,
		},
		{
			name:     "large negative",
			input:    "-987654321",
			expected: -987654321,
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
			input:       "2147483648",
			expectError: true,
		},
		{
			name:        "underflow",
			input:       "-2147483649",
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
			i := &Int32Value{}
			err := i.Set(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("Int32Value.Set() expected error for input %q, got nil", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("Int32Value.Set() unexpected error for input %q: %v", tt.input, err)
				return
			}

			if i.Value != tt.expected {
				t.Errorf("Int32Value.Set() = %v, want %v", i.Value, tt.expected)
			}
		})
	}
}

func TestInt32Value_Type(t *testing.T) {
	i := &Int32Value{}
	if got := i.Type(); got != "int32Value" {
		t.Errorf("Int32Value.Type() = %v, want %v", got, "int32Value")
	}
}

func TestInt32Value_ImplementsPflagValue(t *testing.T) {
	var _ pflag.Value = (*Int32Value)(nil)
}

func TestInt32Value_PflagIntegration(t *testing.T) {
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	i := &Int32Value{}

	fs.Var(i, "int32", "Test int32 value")

	// Test valid input
	err := fs.Parse([]string{"--int32", "123456789"})
	if err != nil {
		t.Fatalf("Failed to parse int32: %v", err)
	}

	expected := int32(123456789)
	if i.Value != expected {
		t.Errorf("Parsed int32 = %v, want %v", i.Value, expected)
	}

	// Test String() method after Set()
	if got := i.String(); got != "123456789" {
		t.Errorf("Int32Value.String() after Set() = %v, want %v", got, "123456789")
	}
}

func TestInt32Value_PflagIntegrationError(t *testing.T) {
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	i := &Int32Value{}

	fs.Var(i, "int32", "Test int32 value")

	// Test invalid input
	err := fs.Parse([]string{"--int32", "invalid"})
	if err == nil {
		t.Fatalf("Expected error for invalid int32 input, got nil")
	}
}

func TestInt32_Constructor(t *testing.T) {
	wrap := &wrapperspb.Int32Value{Value: 123456789}
	i := Int32(wrap)

	// Test that the returned Int32Value has correct fields
	if (*wrapperspb.Int32Value)(i) != wrap {
		t.Errorf("Int32() wrapper = %p, want %p", (*wrapperspb.Int32Value)(i), wrap)
	}

	// Test that the value is correct
	if i.Value != 123456789 {
		t.Errorf("Int32().Value = %v, want %v", i.Value, 123456789)
	}

	// Test String() method on constructed int32
	expected := "123456789"
	if got := i.String(); got != expected {
		t.Errorf("Int32().String() = %v, want %v", got, expected)
	}
}

func TestInt32_ConstructorNil(t *testing.T) {
	i := Int32(nil)

	// Test that it works with nil wrapper
	if i != nil {
		t.Errorf("Int32(nil) = %p, want nil", i)
	}
}

func TestInt32Value_StringAfterSet(t *testing.T) {
	i := &Int32Value{}

	// Set a value
	err := i.Set("-987654321")
	if err != nil {
		t.Fatalf("Failed to set int32 value: %v", err)
	}

	// Test String() method
	if got := i.String(); got != "-987654321" {
		t.Errorf("Int32Value.String() after Set() = %v, want %v", got, "-987654321")
	}
}

func TestInt32Value_MultipleSets(t *testing.T) {
	i := &Int32Value{}

	// Set first value
	err := i.Set("123456789")
	if err != nil {
		t.Fatalf("Failed to set first int32 value: %v", err)
	}
	if i.Value != 123456789 {
		t.Errorf("First set failed: got %v, want %v", i.Value, 123456789)
	}

	// Set second value
	err = i.Set("-987654321")
	if err != nil {
		t.Fatalf("Failed to set second int32 value: %v", err)
	}
	if i.Value != -987654321 {
		t.Errorf("Second set failed: got %v, want %v", i.Value, -987654321)
	}

	// Test final String() value
	if got := i.String(); got != "-987654321" {
		t.Errorf("Final String() = %v, want %v", got, "-987654321")
	}
}

func TestInt32Value_EdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int32
	}{
		{"zero", "0", 0},
		{"positive one", "1", 1},
		{"negative one", "-1", -1},
		{"max int32", "2147483647", 2147483647},
		{"min int32", "-2147483648", -2147483648},
		{"large positive", "123456789", 123456789},
		{"large negative", "-123456789", -123456789},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Int32Value{}
			err := i.Set(tt.input)
			if err != nil {
				t.Fatalf("Failed to set int32 value %q: %v", tt.input, err)
			}
			if i.Value != tt.want {
				t.Errorf("Int32Value edge case: input=%q, got=%v, want=%v", tt.input, i.Value, tt.want)
			}
		})
	}
}

func TestInt32Value_InvalidInputs(t *testing.T) {
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
		{"overflow", "2147483648"},
		{"underflow", "-2147483649"},
		{"very large overflow", "999999999999999999999999999999"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Int32Value{}
			err := i.Set(tt.input)
			if err == nil {
				t.Errorf("Int32Value.Set() expected error for invalid input %q, got nil", tt.input)
			}
		})
	}
}

func TestInt32Value_BoundaryValues(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int32
	}{
		{"zero", "0", 0},
		{"one", "1", 1},
		{"max int32", "2147483647", 2147483647},
		{"min int32", "-2147483648", -2147483648},
		{"max int32 minus one", "2147483646", 2147483646},
		{"min int32 plus one", "-2147483647", -2147483647},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &Int32Value{}
			err := i.Set(tt.input)
			if err != nil {
				t.Fatalf("Failed to set int32 value %q: %v", tt.input, err)
			}
			if i.Value != tt.want {
				t.Errorf("Int32Value boundary test: input=%q, got=%v, want=%v", tt.input, i.Value, tt.want)
			}
		})
	}
}
