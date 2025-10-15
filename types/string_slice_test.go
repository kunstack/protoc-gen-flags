package types_test

import (
	"testing"

	"github.com/kunstack/protoc-gen-flags/types"
	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestStringSliceValue_Set(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []*wrapperspb.StringValue
	}{
		{
			name:  "single string",
			input: "hello",
			want:  []*wrapperspb.StringValue{wrapperspb.String("hello")},
		},
		{
			name:  "multiple strings",
			input: "hello,world,test",
			want:  []*wrapperspb.StringValue{wrapperspb.String("hello"), wrapperspb.String("world"), wrapperspb.String("test")},
		},
		{
			name:  "with spaces",
			input: " hello , world ",
			want:  []*wrapperspb.StringValue{wrapperspb.String(" hello "), wrapperspb.String(" world ")},
		},
		{
			name:  "empty string",
			input: "",
			want:  []*wrapperspb.StringValue{},
		},
		{
			name:  "string with special characters",
			input: "hello@world,test#123",
			want:  []*wrapperspb.StringValue{wrapperspb.String("hello@world"), wrapperspb.String("test#123")},
		},
		{
			name:  "string with unicode",
			input: "こんにちは,world",
			want:  []*wrapperspb.StringValue{wrapperspb.String("こんにちは"), wrapperspb.String("world")},
		},
		{
			name:  "string with numbers",
			input: "123,456,789",
			want:  []*wrapperspb.StringValue{wrapperspb.String("123"), wrapperspb.String("456"), wrapperspb.String("789")},
		},
		{
			name:  "string with commas in values",
			input: "hello,world,with,commas",
			want:  []*wrapperspb.StringValue{wrapperspb.String("hello"), wrapperspb.String("world"), wrapperspb.String("with"), wrapperspb.String("commas")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags := make([]*wrapperspb.StringValue, 0)
			ss := types.StringSlice(&flags)
			err := ss.Set(tt.input)

			if err != nil {
				t.Errorf("StringSliceValue.Set() unexpected error = %v", err)
				return
			}

			if len(flags) != len(tt.want) {
				t.Errorf("StringSliceValue.Set() length = %v, want %v", len(flags), len(tt.want))
				return
			}

			for i, wantVal := range tt.want {
				if (flags)[i].Value != wantVal.Value {
					t.Errorf("StringSliceValue.Set()[%d] = %v, want %v", i, (flags)[i].Value, wantVal.Value)
				}
			}
		})
	}
}

func TestStringSliceValue_Set_Append(t *testing.T) {
	// Test that subsequent Set calls append to the slice
	flags := make([]*wrapperspb.StringValue, 0)
	ss := types.StringSlice(&flags)

	// First Set call
	err := ss.Set("hello,world")
	if err != nil {
		t.Fatalf("First Set() failed: %v", err)
	}

	if len(flags) != 2 {
		t.Errorf("After first Set() length = %v, want %v", len(flags), 2)
	}

	// Second Set call should append
	err = ss.Set("test,another")
	if err != nil {
		t.Fatalf("Second Set() failed: %v", err)
	}

	if len(flags) != 4 {
		t.Errorf("After second Set() length = %v, want %v", len(flags), 4)
	}

	expected := []string{"hello", "world", "test", "another"}
	for i, exp := range expected {
		if (flags)[i].Value != exp {
			t.Errorf("After multiple Set() calls [%d] = %v, want %v", i, (flags)[i].Value, exp)
		}
	}
}

func TestStringSliceValue_Set_InitialEmpty(t *testing.T) {
	// Test Set behavior on initially empty slice
	flags := []*wrapperspb.StringValue{}
	ss := types.StringSlice(&flags)

	err := ss.Set("first")
	if err != nil {
		t.Fatalf("Set() failed on empty slice: %v", err)
	}

	if len(flags) != 1 {
		t.Errorf("After Set() on empty slice length = %v, want %v", len(flags), 1)
	}

	if (flags)[0].Value != "first" {
		t.Errorf("After Set() on empty slice value = %v, want %v", (flags)[0].Value, "first")
	}
}

func TestStringSliceValue_String(t *testing.T) {
	tests := []struct {
		name  string
		input []*wrapperspb.StringValue
		want  string
	}{
		{
			name:  "empty slice",
			input: []*wrapperspb.StringValue{},
			want:  "[]",
		},
		{
			name:  "single string",
			input: []*wrapperspb.StringValue{wrapperspb.String("hello")},
			want:  "[hello]",
		},
		{
			name:  "multiple strings",
			input: []*wrapperspb.StringValue{wrapperspb.String("hello"), wrapperspb.String("world"), wrapperspb.String("test")},
			want:  "[hello,world,test]",
		},
		{
			name:  "strings with special characters",
			input: []*wrapperspb.StringValue{wrapperspb.String("hello@world"), wrapperspb.String("test#123")},
			want:  "[hello@world,test#123]",
		},
		{
			name:  "strings with unicode",
			input: []*wrapperspb.StringValue{wrapperspb.String("こんにちは"), wrapperspb.String("world")},
			want:  "[こんにちは,world]",
		},
		{
			name:  "empty strings",
			input: []*wrapperspb.StringValue{wrapperspb.String(""), wrapperspb.String("")},
			want:  "[,]",
		},
		{
			name:  "mixed content",
			input: []*wrapperspb.StringValue{wrapperspb.String("hello123"), wrapperspb.String("world!@#"), wrapperspb.String("test")},
			want:  "[hello123,world!@#,test]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ss := types.StringSlice(&tt.input)
			if got := ss.String(); got != tt.want {
				t.Errorf("StringSliceValue.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringSliceValue_Type(t *testing.T) {
	flags := []*wrapperspb.StringValue{}
	ss := types.StringSlice(&flags)
	if got := ss.Type(); got != "stringSliceValue" {
		t.Errorf("StringSliceValue.Type() = %v, want %v", got, "stringSliceValue")
	}
}

func TestStringSliceValue_PflagInterface(t *testing.T) {
	flags := make([]*wrapperspb.StringValue, 0)
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)

	// Test that StringSliceValue implements pflag.Value interface
	ss := types.StringSlice(&flags)
	fs.Var(ss, "strings", "test string slice")

	// Test setting flags
	args := []string{"--strings", "hello,world,test"}
	if err := fs.Parse(args); err != nil {
		t.Errorf("Failed to parse flags: %v", err)
	}

	// Check the values through the StringSliceValue interface
	if len(flags) != 3 {
		t.Errorf("Expected 3 values, got %d", len(flags))
	}

	expected := []string{"hello", "world", "test"}
	for i, exp := range expected {
		if (flags)[i].Value != exp {
			t.Errorf("Value[%d] = %v, want %v", i, (flags)[i].Value, exp)
		}
	}
}

func TestStringSliceValue_PflagMultipleCalls(t *testing.T) {
	flags := make([]*wrapperspb.StringValue, 0)
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)

	ss := types.StringSlice(&flags)
	fs.Var(ss, "strings", "test string slice")

	// Test multiple calls to the same flag (should append)
	args := []string{"--strings", "first,second", "--strings", "third,fourth"}
	if err := fs.Parse(args); err != nil {
		t.Errorf("Failed to parse multiple flag calls: %v", err)
	}

	// Should have all 4 values
	if len(flags) != 4 {
		t.Errorf("Expected 4 values after multiple calls, got %d", len(flags))
	}

	expected := []string{"first", "second", "third", "fourth"}
	for i, exp := range expected {
		if (flags)[i].Value != exp {
			t.Errorf("Value after multiple calls[%d] = %v, want %v", i, (flags)[i].Value, exp)
		}
	}
}

func TestStringSliceFunction(t *testing.T) {
	flags := []*wrapperspb.StringValue{wrapperspb.String("hello"), wrapperspb.String("world")}
	ss := types.StringSlice(&flags)

	if ss == nil {
		t.Error("StringSlice() returned nil")
	}

	if len(flags) != 2 {
		t.Errorf("StringSlice() length = %v, want %v", len(flags), 2)
	}

	if (flags)[0].Value != "hello" || (flags)[1].Value != "world" {
		t.Errorf("Values not as expected: %v, %v", (flags)[0].Value, (flags)[1].Value)
	}
}

func TestStringSliceFunction_Nil(t *testing.T) {
	var flags []*wrapperspb.StringValue
	ss := types.StringSlice(&flags)

	if ss == nil {
		t.Error("StringSlice() returned nil for nil slice")
	}

	// Should be able to call Set on nil slice
	err := ss.Set("test")
	if err != nil {
		t.Errorf("StringSlice().Set() failed on nil slice: %v", err)
	}

	if len(flags) != 1 || (flags)[0].Value != "test" {
		t.Errorf("StringSlice() failed to initialize nil slice properly")
	}
}

func TestStringSliceValue_EdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{"empty", "", []string{}},
		{"single char", "a", []string{"a"}},
		{"whitespace", "   ", []string{"   "}},
		{"tabs", "\t\t", []string{"\t\t"}},
		{"newlines", "\"hello\nworld\"", []string{"hello\nworld"}},
		{"unicode", "こんにちは", []string{"こんにちは"}},
		{"emoji", "👋", []string{"👋"}},
		{"numbers", "12345", []string{"12345"}},
		{"mixed", "hello 123 🌍", []string{"hello 123 🌍"}},
		{"quotes", "\"hello\"", []string{"hello"}},
		{"brackets", "[hello]", []string{"[hello]"}},
		{"url", "https://example.com", []string{"https://example.com"}},
		{"email", "test@example.com", []string{"test@example.com"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags := make([]*wrapperspb.StringValue, 0)
			ss := types.StringSlice(&flags)
			err := ss.Set(tt.input)
			if err != nil {
				t.Fatalf("Failed to set string slice value %q: %v", tt.input, err)
			}

			if len(flags) != len(tt.want) {
				t.Errorf("StringSlice edge case: input=%q, got length=%d, want length=%d",
					tt.input, len(flags), len(tt.want))
				return
			}

			for i, want := range tt.want {
				if (flags)[i].Value != want {
					t.Errorf("StringSlice edge case: input=%q, got[%d]=%v, want=%v",
						tt.input, i, (flags)[i].Value, want)
				}
			}
		})
	}
}

func TestStringSliceValue_ChangedFlagBehavior(t *testing.T) {
	// Test the changed flag behavior indirectly through Set() method
	flags := make([]*wrapperspb.StringValue, 0)
	ss := types.StringSlice(&flags)

	// First Set call should initialize the slice
	err := ss.Set("hello")
	if err != nil {
		t.Fatalf("First Set() failed: %v", err)
	}

	if len(flags) != 1 {
		t.Errorf("After first Set() expected 1 value, got %d", len(flags))
	}

	// Second Set call should append to the slice
	err = ss.Set("world")
	if err != nil {
		t.Fatalf("Second Set() failed: %v", err)
	}

	if len(flags) != 2 {
		t.Errorf("After second Set() expected 2 values, got %d", len(flags))
	}

	// Verify both values are present
	if (flags)[0].Value != "hello" || (flags)[1].Value != "world" {
		t.Errorf("Values not as expected: got %v and %v, want hello and world",
			(flags)[0].Value, (flags)[1].Value)
	}
}

func TestStringSliceValue_LongStrings(t *testing.T) {
	// Test with very long strings
	longString := string(make([]byte, 1000))
	for i := range longString {
		longString = longString[:i] + "a" + longString[i+1:]
	}

	flags := make([]*wrapperspb.StringValue, 0)
	ss := types.StringSlice(&flags)

	err := ss.Set(longString)
	if err != nil {
		t.Fatalf("Failed to set very long string: %v", err)
	}

	if len(flags) != 1 {
		t.Errorf("Expected 1 value after setting long string, got %d", len(flags))
	}

	if (flags)[0].Value != longString {
		t.Errorf("Long string value mismatch")
	}
}

func TestStringSliceValue_CommaHandling(t *testing.T) {
	// Test various comma scenarios
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{"leading comma", ",hello,world", []string{"", "hello", "world"}},
		{"trailing comma", "hello,world,", []string{"hello", "world", ""}},
		{"multiple commas", "hello,,,world", []string{"hello", "", "", "world"}},
		{"only commas", ",,,", []string{"", "", "", ""}},
		{"no commas", "helloworld", []string{"helloworld"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flags := make([]*wrapperspb.StringValue, 0)
			ss := types.StringSlice(&flags)
			err := ss.Set(tt.input)
			if err != nil {
				t.Fatalf("Failed to set comma test input %q: %v", tt.input, err)
			}

			if len(flags) != len(tt.want) {
				t.Errorf("Comma test: input=%q, got length=%d, want length=%d",
					tt.input, len(flags), len(tt.want))
				return
			}

			for i, want := range tt.want {
				if (flags)[i].Value != want {
					t.Errorf("Comma test: input=%q, got[%d]=%v, want=%v",
						tt.input, i, (flags)[i].Value, want)
				}
			}
		})
	}
}
