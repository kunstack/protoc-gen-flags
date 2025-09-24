package types

import (
	"strconv"
	"testing"

	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestFloatValue_String(t *testing.T) {
	tests := []struct {
		name     string
		value    float32
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
			name:     "positive float",
			value:    3.14159,
			expected: "3.14159",
		},
		{
			name:     "negative float",
			value:    -2.71828,
			expected: "-2.71828",
		},
		{
			name:     "very small float",
			value:    0.000001,
			expected: "1e-06",
		},
		{
			name:     "very large float",
			value:    999999.9,
			expected: "999999.9",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FloatValue{Value: tt.value}
			if got := f.String(); got != tt.expected {
				t.Errorf("FloatValue.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestFloatValue_String_Nil(t *testing.T) {
	var f *FloatValue
	got := f.String()
	expected := "0"
	if got != expected {
		t.Errorf("FloatValue.String() nil = %v, want %v", got, expected)
	}
}

func TestFloatValue_Set(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    float32
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
			name:     "positive float",
			input:    "3.14159",
			expected: 3.14159,
		},
		{
			name:     "negative float",
			input:    "-2.71828",
			expected: -2.71828,
		},
		{
			name:     "scientific notation",
			input:    "1.5e10",
			expected: 1.5e10,
		},
		{
			name:     "negative scientific notation",
			input:    "-2.5e-5",
			expected: -2.5e-5,
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
			input:       "1e100",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &FloatValue{}
			err := f.Set(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("FloatValue.Set() expected error for input %q, got nil", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("FloatValue.Set() unexpected error for input %q: %v", tt.input, err)
				return
			}

			// Use a small epsilon for float comparison
			if f.Value != tt.expected {
				t.Errorf("FloatValue.Set() = %v, want %v", f.Value, tt.expected)
			}
		})
	}
}

func TestFloatValue_Type(t *testing.T) {
	f := &FloatValue{}
	if got := f.Type(); got != "floatValue" {
		t.Errorf("FloatValue.Type() = %v, want %v", got, "floatValue")
	}
}

func TestFloatValue_ImplementsPflagValue(t *testing.T) {
	var _ pflag.Value = (*FloatValue)(nil)
}

func TestFloatValue_PflagIntegration(t *testing.T) {
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	f := &FloatValue{}

	fs.Var(f, "float", "Test float value")

	// Test valid input
	err := fs.Parse([]string{"--float", "3.14159"})
	if err != nil {
		t.Fatalf("Failed to parse float: %v", err)
	}

	expected := float32(3.14159)
	if f.Value != expected {
		t.Errorf("Parsed float = %v, want %v", f.Value, expected)
	}

	// Test String() method after Set()
	if got := f.String(); got != "3.14159" {
		t.Errorf("FloatValue.String() after Set() = %v, want %v", got, "3.14159")
	}
}

func TestFloatValue_PflagIntegrationError(t *testing.T) {
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	f := &FloatValue{}

	fs.Var(f, "float", "Test float value")

	// Test invalid input
	err := fs.Parse([]string{"--float", "invalid"})
	if err == nil {
		t.Fatalf("Expected error for invalid float input, got nil")
	}
}

func TestFloat_Constructor(t *testing.T) {
	wrap := &wrapperspb.FloatValue{Value: 3.14159}
	f := Float(wrap)

	// Test that the returned FloatValue has correct fields
	if (*wrapperspb.FloatValue)(f) != wrap {
		t.Errorf("Float() wrapper = %p, want %p", (*wrapperspb.FloatValue)(f), wrap)
	}

	// Test that the value is correct
	if f.Value != 3.14159 {
		t.Errorf("Float().Value = %v, want %v", f.Value, 3.14159)
	}

	// Test String() method on constructed float
	expected := "3.14159"
	if got := f.String(); got != expected {
		t.Errorf("Float().String() = %v, want %v", got, expected)
	}
}

func TestFloat_ConstructorNil(t *testing.T) {
	f := Float(nil)

	// Test that it works with nil wrapper
	if f != nil {
		t.Errorf("Float(nil) = %p, want nil", f)
	}
}

func TestFloatValue_StringAfterSet(t *testing.T) {
	f := &FloatValue{}

	// Set a value
	err := f.Set("-2.71828")
	if err != nil {
		t.Fatalf("Failed to set float value: %v", err)
	}

	// Test String() method
	if got := f.String(); got != "-2.71828" {
		t.Errorf("FloatValue.String() after Set() = %v, want %v", got, "-2.71828")
	}
}

func TestFloatValue_MultipleSets(t *testing.T) {
	f := &FloatValue{}

	// Set first value
	err := f.Set("1.234")
	if err != nil {
		t.Fatalf("Failed to set first float value: %v", err)
	}
	if f.Value != 1.234 {
		t.Errorf("First set failed: got %v, want %v", f.Value, 1.234)
	}

	// Set second value
	err = f.Set("-5.678")
	if err != nil {
		t.Fatalf("Failed to set second float value: %v", err)
	}
	if f.Value != -5.678 {
		t.Errorf("Second set failed: got %v, want %v", f.Value, -5.678)
	}

	// Test final String() value
	if got := f.String(); got != "-5.678" {
		t.Errorf("Final String() = %v, want %v", got, "-5.678")
	}
}

func TestFloatValue_Precision(t *testing.T) {
	// Test that we maintain proper precision for float32
	tests := []struct {
		input    string
		expected float32
	}{
		{"0.1", 0.1},
		{"0.01", 0.01},
		{"0.001", 0.001},
		{"123.456", 123.456},
		{"-0.0001", -0.0001},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			f := &FloatValue{}
			err := f.Set(tt.input)
			if err != nil {
				t.Fatalf("Failed to set float value %q: %v", tt.input, err)
			}

			// Convert back to string to verify precision
			expectedStr := strconv.FormatFloat(float64(tt.expected), 'f', -1, 32)
			gotStr := f.String()

			if gotStr != expectedStr {
				t.Errorf("FloatValue precision test: input=%q, got=%q, want=%q",
					tt.input, gotStr, expectedStr)
			}
		})
	}
}
