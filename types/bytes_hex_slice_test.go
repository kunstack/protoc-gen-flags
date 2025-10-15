package types_test

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/kunstack/protoc-gen-flags/types"
	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestBytesHexSliceValue_Set(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []*wrapperspb.BytesValue
		wantErr bool
	}{
		{
			name:  "single hex string",
			input: "48656C6C6F20576F726C64", // "Hello World" in hex
			want:  []*wrapperspb.BytesValue{wrapperspb.Bytes([]byte("Hello World"))},
		},
		{
			name:  "multiple hex strings",
			input: "48656C6C6F,576F726C64", // "Hello" and "World" in hex
			want:  []*wrapperspb.BytesValue{wrapperspb.Bytes([]byte("Hello")), wrapperspb.Bytes([]byte("World"))},
		},
		{
			name:  "with spaces",
			input: " 48656C6C6F , 576F726C64 ",
			want:  []*wrapperspb.BytesValue{wrapperspb.Bytes([]byte("Hello")), wrapperspb.Bytes([]byte("World"))},
		},
		{
			name:    "invalid hex string",
			input:   "48656C6C6F,invalid-hex!",
			wantErr: true,
		},
		{
			name:  "empty string",
			input: "",
			want:  make([]*wrapperspb.BytesValue, 0),
		},
		{
			name:  "empty hex",
			input: "",
			want:  make([]*wrapperspb.BytesValue, 0),
		},
		{
			name:  "binary data",
			input: "000102030405060708090A0B", // binary data in hex
			want:  []*wrapperspb.BytesValue{wrapperspb.Bytes([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11})},
		},
		{
			name:  "quoted strings",
			input: `"48656C6C6F","576F726C64"`,
			want:  []*wrapperspb.BytesValue{wrapperspb.Bytes([]byte("Hello")), wrapperspb.Bytes([]byte("World"))},
		},
		{
			name:  "lowercase hex",
			input: "48656c6c6f,576f726c64",
			want:  []*wrapperspb.BytesValue{wrapperspb.Bytes([]byte("Hello")), wrapperspb.Bytes([]byte("World"))},
		},
		{
			name:  "mixed case hex",
			input: "48656C6c6F,576F726c64",
			want:  []*wrapperspb.BytesValue{wrapperspb.Bytes([]byte("Hello")), wrapperspb.Bytes([]byte("World"))},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytesSlice := make([]*wrapperspb.BytesValue, 0)
			bs := types.BytesHexSlice(&bytesSlice)
			err := bs.Set(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("BytesHexSliceValue.Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(bytesSlice) != len(tt.want) {
					t.Errorf("BytesHexSliceValue.Set() length = %v, want %v", len(bytesSlice), len(tt.want))
					return
				}

				for i, wantVal := range tt.want {
					if string((bytesSlice)[i].Value) != string(wantVal.Value) {
						t.Errorf("BytesHexSliceValue.Set()[%d] = %v, want %v", i, (bytesSlice)[i].Value, wantVal.Value)
					}
				}
			}
		})
	}
}

func TestBytesHexSliceValue_String(t *testing.T) {
	tests := []struct {
		name  string
		input []*wrapperspb.BytesValue
		want  string
	}{
		{
			name:  "empty slice",
			input: make([]*wrapperspb.BytesValue, 0),
			want:  "[]",
		},
		{
			name:  "single hex string",
			input: []*wrapperspb.BytesValue{wrapperspb.Bytes([]byte("Hello World"))},
			want:  "[48656C6C6F20576F726C64]",
		},
		{
			name:  "multiple values",
			input: []*wrapperspb.BytesValue{wrapperspb.Bytes([]byte("Hello")), wrapperspb.Bytes([]byte("World"))},
			want:  "[48656C6C6F,576F726C64]",
		},
		{
			name:  "binary data",
			input: []*wrapperspb.BytesValue{wrapperspb.Bytes([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11})},
			want:  "[000102030405060708090A0B]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bs := types.BytesHexSlice(&tt.input)
			if got := bs.String(); got != tt.want {
				t.Errorf("BytesHexSliceValue.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytesHexSliceValue_Type(t *testing.T) {
	bytesSlice := make([]*wrapperspb.BytesValue, 0)
	bs := types.BytesHexSlice(&bytesSlice)
	if got := bs.Type(); got != "bytesHexSliceValue" {
		t.Errorf("BytesHexSliceValue.Type() = %v, want bytesHexSliceValue", got)
	}
}

func TestBytesHexSliceValue_ImplementsPflagValue(t *testing.T) {
	bytesSlice := make([]*wrapperspb.BytesValue, 0)
	bs := types.BytesHexSlice(&bytesSlice)
	var _ pflag.Value = bs
}

func TestBytesHexSliceValue_MultipleSets(t *testing.T) {
	bytesSlice := make([]*wrapperspb.BytesValue, 0)
	bs := types.BytesHexSlice(&bytesSlice)

	// First set
	err := bs.Set("48656C6C6F,576F726C64") // "Hello", "World" in hex
	if err != nil {
		t.Fatalf("First Set() error = %v", err)
	}

	// Second set (should append)
	err = bs.Set("54686973204973,412054657374") // "This Is", "A Test" in hex
	if err != nil {
		t.Fatalf("Second Set() error = %v", err)
	}

	// Verify all values are present
	expected := []string{"Hello", "World", "This Is", "A Test"}
	if len(bytesSlice) != len(expected) {
		t.Fatalf("Expected %d values, got %d", len(expected), len(bytesSlice))
	}

	for i, expectedStr := range expected {
		if string(bytesSlice[i].Value) != expectedStr {
			t.Errorf("Expected %s at index %d, got %s", expectedStr, i, bytesSlice[i].Value)
		}
	}
}

func TestBytesHexSliceValue_EmptyString(t *testing.T) {
	bytesSlice := make([]*wrapperspb.BytesValue, 0)
	bs := types.BytesHexSlice(&bytesSlice)

	err := bs.Set("")
	if err != nil {
		t.Fatalf("Set() with empty string error = %v", err)
	}

	if len(bytesSlice) != 0 {
		t.Errorf("Expected empty slice, got %d elements", len(bytesSlice))
	}

	if bs.String() != "[]" {
		t.Errorf("Expected '[]', got '%s'", bs.String())
	}
}

func TestBytesHexSliceValue_SpecialCharacters(t *testing.T) {
	// Test with special characters that need hex encoding
	specialBytes := []byte{0, 1, 2, 255, 254, 253}
	encoded := hex.EncodeToString(specialBytes)

	bytesSlice := make([]*wrapperspb.BytesValue, 0)
	bs := types.BytesHexSlice(&bytesSlice)

	err := bs.Set(encoded)
	if err != nil {
		t.Fatalf("Set() with special characters error = %v", err)
	}

	if len(bytesSlice) != 1 {
		t.Fatalf("Expected 1 value, got %d", len(bytesSlice))
	}

	if string(bytesSlice[0].Value) != string(specialBytes) {
		t.Errorf("Expected %v, got %v", specialBytes, bytesSlice[0].Value)
	}

	// Check that String() method encodes back correctly (should be uppercase)
	expectedString := "[" + strings.ToUpper(encoded) + "]"
	if bs.String() != expectedString {
		t.Errorf("String() = %v, want %v", bs.String(), expectedString)
	}
}

func TestBytesHexSlice_Constructor(t *testing.T) {
	// Test with non-nil slice
	initialBytes := []*wrapperspb.BytesValue{wrapperspb.Bytes([]byte("test"))}
	result := types.BytesHexSlice(&initialBytes)

	if result == nil {
		t.Fatal("BytesHexSlice() returned nil")
	}

	// Test that it implements pflag.Value
	var _ pflag.Value = result

	// Test that the initial value is preserved
	if result.String() != "[74657374]" {
		t.Errorf("Expected [54455354], got %s", result.String())
	}
}

func TestBytesHexSlice_ConstructorNil(t *testing.T) {
	// Test with nil slice
	var nilBytes []*wrapperspb.BytesValue
	result := types.BytesHexSlice(&nilBytes)

	if result == nil {
		t.Fatal("BytesHexSlice() returned nil")
	}

	// Test that it implements pflag.Value
	var _ pflag.Value = result

	// Test that it works with nil input
	if result.String() != "[]" {
		t.Errorf("Expected '[]', got %s", result.String())
	}
}

func TestBytesHexSliceValue_PflagInterface(t *testing.T) {
	bytesSlice := make([]*wrapperspb.BytesValue, 0)
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)

	// Test that BytesHexSliceValue implements pflag.Value interface
	bs := types.BytesHexSlice(&bytesSlice)
	fs.Var(bs, "bytes", "test bytes hex slice")

	// Test setting flags
	args := []string{"--bytes", "48656C6C6F,576F726C64"} // "Hello", "World" in hex
	if err := fs.Parse(args); err != nil {
		t.Errorf("Failed to parse flags: %v", err)
	}

	// Verify the parsed values
	if len(bytesSlice) != 2 {
		t.Errorf("Expected 2 values, got %d", len(bytesSlice))
	}

	if string(bytesSlice[0].Value) != "Hello" {
		t.Errorf("Expected 'Hello', got '%s'", (bytesSlice)[0].Value)
	}

	if string(bytesSlice[1].Value) != "World" {
		t.Errorf("Expected 'World', got '%s'", (bytesSlice)[1].Value)
	}
}

// Tests for [][]byte (native byte slice)
func TestBytesHexNativeSliceValue_Set(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    [][]byte
		wantErr bool
	}{
		{
			name:  "single hex",
			input: "48656c6c6f20576f726c64",
			want:  [][]byte{[]byte("Hello World")},
		},
		{
			name:  "multiple hex",
			input: "48656c6c6f,576f726c64,54657374",
			want:  [][]byte{[]byte("Hello"), []byte("World"), []byte("Test")},
		},
		{
			name:  "with spaces",
			input: " 48656c6c6f , 576f726c64 , 54657374 ",
			want:  [][]byte{[]byte("Hello"), []byte("World"), []byte("Test")},
		},
		{
			name:    "invalid hex",
			input:   "48656c6c6f,invalid,54657374",
			wantErr: true,
		},
		{
			name:  "empty string",
			input: "",
			want:  [][]byte{},
		},
		{
			name:  "quoted strings",
			input: `"48656c6c6f","576f726c64"`,
			want:  [][]byte{[]byte("Hello"), []byte("World")},
		},
		{
			name:  "mixed case hex",
			input: "48656C6C6F,776f726c64",
			want:  [][]byte{[]byte("Hello"), []byte("world")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var bytesSlice [][]byte
			nbhs := types.BytesHexSlice(&bytesSlice)
			err := nbhs.Set(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("BytesHexNativeSliceValue.Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(bytesSlice) != len(tt.want) {
					t.Errorf("BytesHexNativeSliceValue.Set() length = %v, want %v", len(bytesSlice), len(tt.want))
					return
				}

				for i, wantVal := range tt.want {
					if string(bytesSlice[i]) != string(wantVal) {
						t.Errorf("BytesHexNativeSliceValue.Set()[%d] = %v, want %v", i, bytesSlice[i], wantVal)
					}
				}
			}
		})
	}
}

func TestBytesHexNativeSliceValue_String(t *testing.T) {
	tests := []struct {
		name  string
		input [][]byte
		want  string
	}{
		{
			name:  "empty slice",
			input: [][]byte{},
			want:  "[]",
		},
		{
			name:  "single bytes",
			input: [][]byte{[]byte("Hello World")},
			want:  "[48656c6c6f20576f726c64]",
		},
		{
			name:  "multiple values",
			input: [][]byte{[]byte("Hello"), []byte("World"), []byte("Test")},
			want:  "[48656c6c6f,576f726c64,54657374]",
		},
		{
			name:  "empty bytes",
			input: [][]byte{[]byte{}, []byte{}},
			want:  "[,]",
		},
		{
			name:  "binary data",
			input: [][]byte{{0x00, 0x01, 0x02, 0x03}, {0xFF, 0xFE, 0xFD}},
			want:  "[00010203,fffefd]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nbhs := types.BytesHexSlice(&tt.input)
			if got := nbhs.String(); got != tt.want {
				t.Errorf("BytesHexNativeSliceValue.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytesHexNativeSliceValue_Type(t *testing.T) {
	var bytesSlice [][]byte
	nbhs := types.BytesHexSlice(&bytesSlice)
	if got := nbhs.Type(); got != "BytesHexNativeSliceValue" {
		t.Errorf("BytesHexNativeSliceValue.Type() = %v, want BytesHexNativeSliceValue", got)
	}
}

func TestBytesHexNativeSliceValue_ImplementsPflagValue(t *testing.T) {
	var bytesSlice [][]byte
	nbhs := types.BytesHexSlice(&bytesSlice)
	var _ pflag.Value = nbhs
}

func TestBytesHexNativeSliceValue_MultipleSets(t *testing.T) {
	var bytesSlice [][]byte
	nbhs := types.BytesHexSlice(&bytesSlice)

	// First set
	err := nbhs.Set("48656c6c6f,576f726c64")
	if err != nil {
		t.Fatalf("First Set() error = %v", err)
	}

	// Second set (should append)
	err = nbhs.Set("54657374")
	if err != nil {
		t.Fatalf("Second Set() error = %v", err)
	}

	// Verify all values are present
	expected := []string{"Hello", "World", "Test"}
	if len(bytesSlice) != len(expected) {
		t.Fatalf("Expected %d values, got %d", len(expected), len(bytesSlice))
	}

	for i, expectedVal := range expected {
		if string(bytesSlice[i]) != expectedVal {
			t.Errorf("Expected %v at index %d, got %v", expectedVal, i, bytesSlice[i])
		}
	}
}

func TestBytesHexNativeSliceValue_PflagInterface(t *testing.T) {
	var bytesSlice [][]byte
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)

	// Test that BytesHexNativeSliceValue implements pflag.Value interface
	nbhs := types.BytesHexSlice(&bytesSlice)
	fs.Var(nbhs, "hex-bytes", "test native hex bytes slice")

	// Test setting flags
	args := []string{"--hex-bytes", "48656c6c6f,576f726c64"}
	if err := fs.Parse(args); err != nil {
		t.Errorf("Failed to parse flags: %v", err)
	}

	// Verify the parsed values
	if len(bytesSlice) != 2 {
		t.Errorf("Expected 2 values, got %d", len(bytesSlice))
	}

	if string(bytesSlice[0]) != "Hello" {
		t.Errorf("Expected 'Hello' at index 0, got %v", bytesSlice[0])
	}

	if string(bytesSlice[1]) != "World" {
		t.Errorf("Expected 'World' at index 1, got %v", bytesSlice[1])
	}
}
