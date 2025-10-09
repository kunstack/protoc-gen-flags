package types

import (
	"encoding/base64"
	"testing"

	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestBytesSliceValue_Set(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []*wrapperspb.BytesValue
		wantErr bool
	}{
		{
			name:  "single base64 string",
			input: "SGVsbG8gV29ybGQ=", // "Hello World" in base64
			want:  []*wrapperspb.BytesValue{wrapperspb.Bytes([]byte("Hello World"))},
		},
		{
			name:  "multiple base64 strings",
			input: "SGVsbG8=,V29ybGQ=", // "Hello" and "World" in base64
			want:  []*wrapperspb.BytesValue{wrapperspb.Bytes([]byte("Hello")), wrapperspb.Bytes([]byte("World"))},
		},
		{
			name:  "with spaces",
			input: " SGVsbG8= , V29ybGQ= ",
			want:  []*wrapperspb.BytesValue{wrapperspb.Bytes([]byte("Hello")), wrapperspb.Bytes([]byte("World"))},
		},
		{
			name:    "invalid base64 string",
			input:   "SGVsbG8=,invalid-base64!",
			wantErr: true,
		},
		{
			name:  "empty string",
			input: "",
			want:  []*wrapperspb.BytesValue{},
		},
		{
			name:  "empty base64",
			input: "",
			want:  []*wrapperspb.BytesValue{},
		},
		{
			name:  "binary data",
			input: "AAECAwQFBgcICQoL", // binary data in base64
			want:  []*wrapperspb.BytesValue{wrapperspb.Bytes([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11})},
		},
		{
			name:  "quoted strings",
			input: `"SGVsbG8=","V29ybGQ="`,
			want:  []*wrapperspb.BytesValue{wrapperspb.Bytes([]byte("Hello")), wrapperspb.Bytes([]byte("World"))},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytesSlice := []*wrapperspb.BytesValue{}
			bs := BytesSliceValue{value: &bytesSlice}
			err := bs.Set(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("BytesSliceValue.Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(*bs.value) != len(tt.want) {
					t.Errorf("BytesSliceValue.Set() length = %v, want %v", len(*bs.value), len(tt.want))
					return
				}

				for i, wantVal := range tt.want {
					if string((*bs.value)[i].Value) != string(wantVal.Value) {
						t.Errorf("BytesSliceValue.Set()[%d] = %v, want %v", i, (*bs.value)[i].Value, wantVal.Value)
					}
				}
			}
		})
	}
}

func TestBytesSliceValue_String(t *testing.T) {
	tests := []struct {
		name  string
		input []*wrapperspb.BytesValue
		want  string
	}{
		{
			name:  "empty slice",
			input: []*wrapperspb.BytesValue{},
			want:  "[]",
		},
		{
			name:  "single base64 string",
			input: []*wrapperspb.BytesValue{wrapperspb.Bytes([]byte("Hello World"))},
			want:  "[SGVsbG8gV29ybGQ=]",
		},
		{
			name:  "multiple values",
			input: []*wrapperspb.BytesValue{wrapperspb.Bytes([]byte("Hello")), wrapperspb.Bytes([]byte("World"))},
			want:  "[SGVsbG8=,V29ybGQ=]",
		},
		{
			name:  "binary data",
			input: []*wrapperspb.BytesValue{wrapperspb.Bytes([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11})},
			want:  "[AAECAwQFBgcICQoL]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := BytesSliceValue{value: &tt.input}
			if got := bs.String(); got != tt.want {
				t.Errorf("BytesSliceValue.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytesSliceValue_Type(t *testing.T) {
	bytesSlice := []*wrapperspb.BytesValue{}
	bs := BytesSliceValue{value: &bytesSlice}
	if got := bs.Type(); got != "bytesSliceValue" {
		t.Errorf("BytesSliceValue.Type() = %v, want bytesSliceValue", got)
	}
}

func TestBytesSliceValue_ImplementsPflagValue(t *testing.T) {
	bytesSlice := []*wrapperspb.BytesValue{}
	bs := BytesSliceValue{value: &bytesSlice}
	var _ pflag.Value = &bs
}

func TestBytesSliceValue_MultipleSets(t *testing.T) {
	bytesSlice := []*wrapperspb.BytesValue{}
	bs := BytesSliceValue{value: &bytesSlice}

	// First set
	err := bs.Set("SGVsbG8=,V29ybGQ=") // "Hello", "World"
	if err != nil {
		t.Fatalf("First Set() error = %v", err)
	}

	// Second set (should append)
	err = bs.Set("VGhpcyBJcw==,QSBUZXN0") // "This Is", "A Test"
	if err != nil {
		t.Fatalf("Second Set() error = %v", err)
	}

	// Verify all values are present
	expected := []string{"Hello", "World", "This Is", "A Test"}
	if len(*bs.value) != len(expected) {
		t.Fatalf("Expected %d values, got %d", len(expected), len(*bs.value))
	}

	for i, expectedStr := range expected {
		if string((*bs.value)[i].Value) != expectedStr {
			t.Errorf("Expected %s at index %d, got %s", expectedStr, i, (*bs.value)[i].Value)
		}
	}
}

func TestBytesSliceValue_EmptyString(t *testing.T) {
	bytesSlice := []*wrapperspb.BytesValue{}
	bs := BytesSliceValue{value: &bytesSlice}

	err := bs.Set("")
	if err != nil {
		t.Fatalf("Set() with empty string error = %v", err)
	}

	if len(*bs.value) != 0 {
		t.Errorf("Expected empty slice, got %d elements", len(*bs.value))
	}

	if bs.String() != "[]" {
		t.Errorf("Expected '[]', got '%s'", bs.String())
	}
}

func TestBytesSliceValue_SpecialCharacters(t *testing.T) {
	// Test with special characters that need base64 encoding
	specialBytes := []byte{0, 1, 2, 255, 254, 253}
	encoded := base64.StdEncoding.EncodeToString(specialBytes)

	bytesSlice := []*wrapperspb.BytesValue{}
	bs := BytesSliceValue{value: &bytesSlice}

	err := bs.Set(encoded)
	if err != nil {
		t.Fatalf("Set() with special characters error = %v", err)
	}

	if len(*bs.value) != 1 {
		t.Fatalf("Expected 1 value, got %d", len(*bs.value))
	}

	if string((*bs.value)[0].Value) != string(specialBytes) {
		t.Errorf("Expected %v, got %v", specialBytes, (*bs.value)[0].Value)
	}

	// Check that String() method encodes back correctly
	expectedString := "[" + encoded + "]"
	if bs.String() != expectedString {
		t.Errorf("String() = %v, want %v", bs.String(), expectedString)
	}
}

func TestBytesSlice_Constructor(t *testing.T) {
	// Test with non-nil slice
	initialBytes := []*wrapperspb.BytesValue{wrapperspb.Bytes([]byte("test"))}
	result := BytesSlice(initialBytes)

	if result == nil {
		t.Fatal("BytesSlice() returned nil")
	}

	// Type assert to access internal fields
	if bs, ok := result.(*BytesSliceValue); ok {
		if bs.value == nil {
			t.Error("BytesSlice() did not set the value pointer")
		}

		if bs.changed {
			t.Error("BytesSlice() should initialize with changed=false")
		}
	} else {
		t.Error("BytesSlice() did not return *BytesSliceValue")
	}
}

func TestBytesSlice_ConstructorNil(t *testing.T) {
	// Test with nil slice
	var nilBytes []*wrapperspb.BytesValue
	result := BytesSlice(nilBytes)

	if result == nil {
		t.Fatal("BytesSlice() returned nil")
	}

	// Type assert to access internal fields
	if bs, ok := result.(*BytesSliceValue); ok {
		if bs.value == nil {
			t.Error("BytesSlice() did not set the value pointer for nil input")
		}

		if bs.changed {
			t.Error("BytesSlice() should initialize with changed=false for nil input")
		}
	} else {
		t.Error("BytesSlice() did not return *BytesSliceValue for nil input")
	}
}

func TestBytesSliceValue_PflagInterface(t *testing.T) {
	bytesSlice := []*wrapperspb.BytesValue{}
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)

	// Test that BytesSliceValue implements pflag.Value interface
	bs := BytesSlice(bytesSlice)
	fs.Var(bs, "bytes", "test bytes slice")

	// Test setting flags
	args := []string{"--bytes", "SGVsbG8=,V29ybGQ="} // "Hello", "World" in base64
	if err := fs.Parse(args); err != nil {
		t.Errorf("Failed to parse flags: %v", err)
	}

	// Type assert to access internal values
	if bsSlice, ok := bs.(*BytesSliceValue); ok {
		// Verify the parsed values
		if len(*bsSlice.value) != 2 {
			t.Errorf("Expected 2 values, got %d", len(*bsSlice.value))
		}

		if string((*bsSlice.value)[0].Value) != "Hello" {
			t.Errorf("Expected 'Hello', got '%s'", (*bsSlice.value)[0].Value)
		}

		if string((*bsSlice.value)[1].Value) != "World" {
			t.Errorf("Expected 'World', got '%s'", (*bsSlice.value)[1].Value)
		}
	} else {
		t.Error("Expected *BytesSliceValue type")
	}
}
