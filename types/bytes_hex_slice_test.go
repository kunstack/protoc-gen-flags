package types

import (
	"encoding/hex"
	"strings"
	"testing"

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
			want:  []*wrapperspb.BytesValue{},
		},
		{
			name:  "empty hex",
			input: "",
			want:  []*wrapperspb.BytesValue{},
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
			bytesSlice := []*wrapperspb.BytesValue{}
			bs := BytesHexSliceValue{value: &bytesSlice}
			err := bs.Set(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("BytesHexSliceValue.Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(*bs.value) != len(tt.want) {
					t.Errorf("BytesHexSliceValue.Set() length = %v, want %v", len(*bs.value), len(tt.want))
					return
				}

				for i, wantVal := range tt.want {
					if string((*bs.value)[i].Value) != string(wantVal.Value) {
						t.Errorf("BytesHexSliceValue.Set()[%d] = %v, want %v", i, (*bs.value)[i].Value, wantVal.Value)
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
			input: []*wrapperspb.BytesValue{},
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
			bs := BytesHexSliceValue{value: &tt.input}
			if got := bs.String(); got != tt.want {
				t.Errorf("BytesHexSliceValue.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBytesHexSliceValue_Type(t *testing.T) {
	bytesSlice := []*wrapperspb.BytesValue{}
	bs := BytesHexSliceValue{value: &bytesSlice}
	if got := bs.Type(); got != "bytesHexSliceValue" {
		t.Errorf("BytesHexSliceValue.Type() = %v, want bytesHexSliceValue", got)
	}
}

func TestBytesHexSliceValue_ImplementsPflagValue(t *testing.T) {
	bytesSlice := []*wrapperspb.BytesValue{}
	bs := BytesHexSliceValue{value: &bytesSlice}
	var _ pflag.Value = &bs
}

func TestBytesHexSliceValue_MultipleSets(t *testing.T) {
	bytesSlice := []*wrapperspb.BytesValue{}
	bs := BytesHexSliceValue{value: &bytesSlice}

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
	if len(*bs.value) != len(expected) {
		t.Fatalf("Expected %d values, got %d", len(expected), len(*bs.value))
	}

	for i, expectedStr := range expected {
		if string((*bs.value)[i].Value) != expectedStr {
			t.Errorf("Expected %s at index %d, got %s", expectedStr, i, (*bs.value)[i].Value)
		}
	}
}

func TestBytesHexSliceValue_EmptyString(t *testing.T) {
	bytesSlice := []*wrapperspb.BytesValue{}
	bs := BytesHexSliceValue{value: &bytesSlice}

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

func TestBytesHexSliceValue_SpecialCharacters(t *testing.T) {
	// Test with special characters that need hex encoding
	specialBytes := []byte{0, 1, 2, 255, 254, 253}
	encoded := hex.EncodeToString(specialBytes)

	bytesSlice := []*wrapperspb.BytesValue{}
	bs := BytesHexSliceValue{value: &bytesSlice}

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

	// Check that String() method encodes back correctly (should be uppercase)
	expectedString := "[" + strings.ToUpper(encoded) + "]"
	if bs.String() != expectedString {
		t.Errorf("String() = %v, want %v", bs.String(), expectedString)
	}
}

func TestBytesHexSlice_Constructor(t *testing.T) {
	// Test with non-nil slice
	initialBytes := []*wrapperspb.BytesValue{wrapperspb.Bytes([]byte("test"))}
	result := BytesHexSlice(initialBytes)

	if result == nil {
		t.Fatal("BytesHexSlice() returned nil")
	}

	// Type assert to access internal fields
	if bs, ok := result.(*BytesHexSliceValue); ok {
		if bs.value == nil {
			t.Error("BytesHexSlice() did not set the value pointer")
		}

		if bs.changed {
			t.Error("BytesHexSlice() should initialize with changed=false")
		}
	} else {
		t.Error("BytesHexSlice() did not return *BytesHexSliceValue")
	}
}

func TestBytesHexSlice_ConstructorNil(t *testing.T) {
	// Test with nil slice
	var nilBytes []*wrapperspb.BytesValue
	result := BytesHexSlice(nilBytes)

	if result == nil {
		t.Fatal("BytesHexSlice() returned nil")
	}

	// Type assert to access internal fields
	if bs, ok := result.(*BytesHexSliceValue); ok {
		if bs.value == nil {
			t.Error("BytesHexSlice() did not set the value pointer for nil input")
		}

		if bs.changed {
			t.Error("BytesHexSlice() should initialize with changed=false for nil input")
		}
	} else {
		t.Error("BytesHexSlice() did not return *BytesHexSliceValue for nil input")
	}
}

func TestBytesHexSliceValue_PflagInterface(t *testing.T) {
	bytesSlice := []*wrapperspb.BytesValue{}
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)

	// Test that BytesHexSliceValue implements pflag.Value interface
	bs := BytesHexSlice(bytesSlice)
	fs.Var(bs, "bytes", "test bytes hex slice")

	// Test setting flags
	args := []string{"--bytes", "48656C6C6F,576F726C64"} // "Hello", "World" in hex
	if err := fs.Parse(args); err != nil {
		t.Errorf("Failed to parse flags: %v", err)
	}

	// Type assert to access internal values
	if bsHex, ok := bs.(*BytesHexSliceValue); ok {
		// Verify the parsed values
		if len(*bsHex.value) != 2 {
			t.Errorf("Expected 2 values, got %d", len(*bsHex.value))
		}

		if string((*bsHex.value)[0].Value) != "Hello" {
			t.Errorf("Expected 'Hello', got '%s'", (*bsHex.value)[0].Value)
		}

		if string((*bsHex.value)[1].Value) != "World" {
			t.Errorf("Expected 'World', got '%s'", (*bsHex.value)[1].Value)
		}
	} else {
		t.Error("Expected *BytesHexSliceValue type")
	}
}

func TestNativeBytesHexSliceValue_Set(t *testing.T) {
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
			nbhs := NativeBytesHexSliceValue{value: &bytesSlice}
			err := nbhs.Set(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("NativeBytesHexSliceValue.Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(*nbhs.value) != len(tt.want) {
					t.Errorf("NativeBytesHexSliceValue.Set() length = %v, want %v", len(*nbhs.value), len(tt.want))
					return
				}

				for i, wantVal := range tt.want {
					if string((*nbhs.value)[i]) != string(wantVal) {
						t.Errorf("NativeBytesHexSliceValue.Set()[%d] = %v, want %v", i, (*nbhs.value)[i], wantVal)
					}
				}
			}
		})
	}
}

func TestNativeBytesHexSliceValue_String(t *testing.T) {
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
			nbhs := NativeBytesHexSliceValue{value: &tt.input}
			if got := nbhs.String(); got != tt.want {
				t.Errorf("NativeBytesHexSliceValue.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNativeBytesHexSliceValue_Type(t *testing.T) {
	var bytesSlice [][]byte
	nbhs := NativeBytesHexSliceValue{value: &bytesSlice}
	if got := nbhs.Type(); got != "nativeBytesHexSliceValue" {
		t.Errorf("NativeBytesHexSliceValue.Type() = %v, want nativeBytesHexSliceValue", got)
	}
}

func TestNativeBytesHexSliceValue_ImplementsPflagValue(t *testing.T) {
	var bytesSlice [][]byte
	nbhs := NativeBytesHexSliceValue{value: &bytesSlice}
	var _ pflag.Value = &nbhs
}

func TestNativeBytesHexSliceValue_MultipleSets(t *testing.T) {
	var bytesSlice [][]byte
	nbhs := NativeBytesHexSliceValue{value: &bytesSlice}

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
	if len(*nbhs.value) != len(expected) {
		t.Fatalf("Expected %d values, got %d", len(expected), len(*nbhs.value))
	}

	for i, expectedVal := range expected {
		if string((*nbhs.value)[i]) != expectedVal {
			t.Errorf("Expected %v at index %d, got %v", expectedVal, i, (*nbhs.value)[i])
		}
	}
}

func TestNativeBytesHexSliceValue_PflagInterface(t *testing.T) {
	var bytesSlice [][]byte
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)

	// Test that NativeBytesHexSliceValue implements pflag.Value interface
	nbhs := &NativeBytesHexSliceValue{value: &bytesSlice}
	fs.Var(nbhs, "hex-bytes", "test native hex bytes slice")

	// Test setting flags
	args := []string{"--hex-bytes", "48656c6c6f,576f726c64"}
	if err := fs.Parse(args); err != nil {
		t.Errorf("Failed to parse flags: %v", err)
	}

	// Verify the parsed values
	if len(*nbhs.value) != 2 {
		t.Errorf("Expected 2 values, got %d", len(*nbhs.value))
	}

	if string((*nbhs.value)[0]) != "Hello" {
		t.Errorf("Expected 'Hello' at index 0, got %v", (*nbhs.value)[0])
	}

	if string((*nbhs.value)[1]) != "World" {
		t.Errorf("Expected 'World' at index 1, got %v", (*nbhs.value)[1])
	}
}

func TestNativeBytesHexSlice_Constructor(t *testing.T) {
	// Test with non-nil slice - using direct struct initialization instead of deprecated NativeBytesHexSlice
	initialBytes := [][]byte{[]byte("test")}
	result := &NativeBytesHexSliceValue{value: &initialBytes}

	if result == nil {
		t.Fatal("NativeBytesHexSliceValue constructor returned nil")
	}

	if result.value == nil {
		t.Error("NativeBytesHexSliceValue did not set the value pointer")
	}

	if result.changed {
		t.Error("NativeBytesHexSliceValue should initialize with changed=false")
	}
}

func TestNativeBytesHexSlice_ConstructorNil(t *testing.T) {
	// Test with nil slice - using direct struct initialization instead of deprecated NativeBytesHexSlice
	var nilBytes [][]byte
	result := &NativeBytesHexSliceValue{value: &nilBytes}

	if result == nil {
		t.Fatal("NativeBytesHexSliceValue constructor returned nil")
	}

	if result.value == nil {
		t.Error("NativeBytesHexSliceValue did not set the value pointer for nil input")
	}

	if result.changed {
		t.Error("NativeBytesHexSliceValue should initialize with changed=false for nil input")
	}
}

func TestNativeBytesHexSliceValue_RoundTrip(t *testing.T) {
	// Test that String() and Set() are consistent
	originalData := [][]byte{
		[]byte("Hello World"),
		[]byte("Test Data"),
		[]byte("Binary\x00Data"),
	}

	var bytesSlice [][]byte
	nbhs := NativeBytesHexSliceValue{value: &bytesSlice}

	// Set the data
	err := nbhs.Set("48656c6c6f20576f726c64,546573742044617461,42696e6172790044617461")
	if err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	// Get the string representation
	str := nbhs.String()

	// Create a new slice and set from the string
	var newBytesSlice [][]byte
	newNbhs := NativeBytesHexSliceValue{value: &newBytesSlice}
	err = newNbhs.Set(str[1 : len(str)-1]) // Remove brackets
	if err != nil {
		t.Fatalf("Set() from string error = %v", err)
	}

	// Verify the data matches
	if len(*newNbhs.value) != len(originalData) {
		t.Errorf("Expected %d values, got %d", len(originalData), len(*newNbhs.value))
	}

	for i, expected := range originalData {
		if string((*newNbhs.value)[i]) != string(expected) {
			t.Errorf("Round trip failed at index %d: got %v, want %v", i, (*newNbhs.value)[i], expected)
		}
	}
}
