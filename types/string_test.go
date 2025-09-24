package types

import (
	"testing"

	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestStringValue_String(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		expected string
	}{
		{
			name:     "empty string",
			value:    "",
			expected: "",
		},
		{
			name:     "simple string",
			value:    "hello",
			expected: "hello",
		},
		{
			name:     "string with spaces",
			value:    "hello world",
			expected: "hello world",
		},
		{
			name:     "string with special characters",
			value:    "hello@world#123",
			expected: "hello@world#123",
		},
		{
			name:     "string with unicode",
			value:    "こんにちは",
			expected: "こんにちは",
		},
		{
			name:     "string with numbers",
			value:    "12345",
			expected: "12345",
		},
		{
			name:     "long string",
			value:    "this is a very long string that contains multiple words and should be handled correctly by the StringValue implementation",
			expected: "this is a very long string that contains multiple words and should be handled correctly by the StringValue implementation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StringValue{Value: tt.value}
			if got := s.String(); got != tt.expected {
				t.Errorf("StringValue.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestStringValue_String_Nil(t *testing.T) {
	var s *StringValue
	got := s.String()
	expected := ""
	if got != expected {
		t.Errorf("StringValue.String() nil = %v, want %v", got, expected)
	}
}

func TestStringValue_Set(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "simple string",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "string with spaces",
			input:    "hello world",
			expected: "hello world",
		},
		{
			name:     "string with special characters",
			input:    "hello@world#123",
			expected: "hello@world#123",
		},
		{
			name:     "string with unicode",
			input:    "こんにちは",
			expected: "こんにちは",
		},
		{
			name:     "string with numbers",
			input:    "12345",
			expected: "12345",
		},
		{
			name:     "long string",
			input:    "this is a very long string that contains multiple words and should be handled correctly by the StringValue implementation",
			expected: "this is a very long string that contains multiple words and should be handled correctly by the StringValue implementation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StringValue{}
			err := s.Set(tt.input)

			if err != nil {
				t.Errorf("StringValue.Set() unexpected error for input %q: %v", tt.input, err)
				return
			}

			if s.Value != tt.expected {
				t.Errorf("StringValue.Set() = %v, want %v", s.Value, tt.expected)
			}
		})
	}
}

func TestStringValue_Type(t *testing.T) {
	s := &StringValue{}
	if got := s.Type(); got != "stringValue" {
		t.Errorf("StringValue.Type() = %v, want %v", got, "stringValue")
	}
}

func TestStringValue_ImplementsPflagValue(t *testing.T) {
	var _ pflag.Value = (*StringValue)(nil)
}

func TestStringValue_PflagIntegration(t *testing.T) {
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	s := &StringValue{}

	fs.Var(s, "string", "Test string value")

	// Test valid input
	err := fs.Parse([]string{"--string", "hello world"})
	if err != nil {
		t.Fatalf("Failed to parse string: %v", err)
	}

	expected := "hello world"
	if s.Value != expected {
		t.Errorf("Parsed string = %v, want %v", s.Value, expected)
	}

	// Test String() method after Set()
	if got := s.String(); got != "hello world" {
		t.Errorf("StringValue.String() after Set() = %v, want %v", got, "hello world")
	}
}

func TestStringValue_PflagIntegrationEmpty(t *testing.T) {
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	s := &StringValue{}

	fs.Var(s, "string", "Test string value")

	// Test with no argument (should use default)
	err := fs.Parse([]string{})
	if err != nil {
		t.Fatalf("Failed to parse empty args: %v", err)
	}

	expected := ""
	if s.Value != expected {
		t.Errorf("Default string = %v, want %v", s.Value, expected)
	}
}

func TestStringValue_PflagIntegrationSpecialChars(t *testing.T) {
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	s := &StringValue{}

	fs.Var(s, "string", "Test string value")

	// Test input with special characters
	err := fs.Parse([]string{"--string", "hello@world#123"})
	if err != nil {
		t.Fatalf("Failed to parse string with special chars: %v", err)
	}

	expected := "hello@world#123"
	if s.Value != expected {
		t.Errorf("Parsed string with special chars = %v, want %v", s.Value, expected)
	}
}

func TestString_Constructor(t *testing.T) {
	wrap := &wrapperspb.StringValue{Value: "hello world"}
	s := String(wrap)

	// Test that the returned StringValue has correct fields
	if (*wrapperspb.StringValue)(s) != wrap {
		t.Errorf("String() wrapper = %p, want %p", (*wrapperspb.StringValue)(s), wrap)
	}

	// Test that the value is correct
	if s.Value != "hello world" {
		t.Errorf("String().Value = %v, want %v", s.Value, "hello world")
	}

	// Test String() method on constructed string
	expected := "hello world"
	if got := s.String(); got != expected {
		t.Errorf("String().String() = %v, want %v", got, expected)
	}
}

func TestString_ConstructorNil(t *testing.T) {
	s := String(nil)

	// Test that it works with nil wrapper
	if s != nil {
		t.Errorf("String(nil) = %p, want nil", s)
	}
}

func TestStringValue_StringAfterSet(t *testing.T) {
	s := &StringValue{}

	// Set a value
	err := s.Set("test string value")
	if err != nil {
		t.Fatalf("Failed to set string value: %v", err)
	}

	// Test String() method
	if got := s.String(); got != "test string value" {
		t.Errorf("StringValue.String() after Set() = %v, want %v", got, "test string value")
	}
}

func TestStringValue_MultipleSets(t *testing.T) {
	s := &StringValue{}

	// Set first value
	err := s.Set("first value")
	if err != nil {
		t.Fatalf("Failed to set first string value: %v", err)
	}
	if s.Value != "first value" {
		t.Errorf("First set failed: got %v, want %v", s.Value, "first value")
	}

	// Set second value
	err = s.Set("second value")
	if err != nil {
		t.Fatalf("Failed to set second string value: %v", err)
	}
	if s.Value != "second value" {
		t.Errorf("Second set failed: got %v, want %v", s.Value, "second value")
	}

	// Test final String() value
	if got := s.String(); got != "second value" {
		t.Errorf("Final String() = %v, want %v", got, "second value")
	}
}

func TestStringValue_EdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty", "", ""},
		{"single char", "a", "a"},
		{"whitespace", "   ", "   "},
		{"tabs", "\t\t", "\t\t"},
		{"newlines", "hello\nworld", "hello\nworld"},
		{"unicode", "こんにちは", "こんにちは"},
		{"emoji", "👋", "👋"},
		{"numbers", "12345", "12345"},
		{"mixed", "hello 123 🌍", "hello 123 🌍"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &StringValue{}
			err := s.Set(tt.input)
			if err != nil {
				t.Fatalf("Failed to set string value %q: %v", tt.input, err)
			}
			if s.Value != tt.want {
				t.Errorf("StringValue edge case: input=%q, got=%v, want=%v", tt.input, s.Value, tt.want)
			}
		})
	}
}

func TestStringValue_SetNeverFails(t *testing.T) {
	// StringValue.Set() should never fail since any string input is valid
	tests := []string{
		"",
		"hello",
		"123",
		"true",
		"false",
		"null",
		"undefined",
		"special chars: !@#$%^&*()",
		"unicode: こんにちは",
		"emoji: 👋",
		"very long string: " + string(make([]byte, 10000)),
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			s := &StringValue{}
			err := s.Set(input)
			if err != nil {
				t.Errorf("StringValue.Set() should never fail, but got error for input %q: %v", input, err)
			}
			if s.Value != input {
				t.Errorf("StringValue.Set() value mismatch for input %q: got %v, want %v", input, s.Value, input)
			}
		})
	}
}

func TestStringValue_DefaultValue(t *testing.T) {
	var s *StringValue

	// Test that nil StringValue returns empty string
	if got := s.String(); got != "" {
		t.Errorf("Default nil StringValue.String() = %v, want %v", got, "")
	}

	// Test that new StringValue has empty string by default
	s = &StringValue{}
	if got := s.String(); got != "" {
		t.Errorf("Default StringValue.String() = %v, want %v", got, "")
	}
}

func TestStringValue_LongString(t *testing.T) {
	s := &StringValue{}

	// Test with a very long string
	longString := string(make([]byte, 10000))
	longBytes := []byte(longString)
	for i := range longBytes {
		longBytes[i] = byte(i % 256) // Fill with various byte values
	}
	longString = string(longBytes)

	err := s.Set(longString)
	if err != nil {
		t.Fatalf("Failed to set very long string: %v", err)
	}

	if s.Value != longString {
		t.Errorf("Long string value mismatch")
	}

	if got := s.String(); got != longString {
		t.Errorf("Long string String() mismatch")
	}
}
