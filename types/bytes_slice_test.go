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
			name:    "invalid base64 string",
			input:   "invalid-base64!",
			wantErr: true,
		},
		{
			name:  "binary data",
			input: "AAECAwQFBgcICQoL", // binary data in base64
			want:  []*wrapperspb.BytesValue{wrapperspb.Bytes([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11})},
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
	bs := BytesSliceValue{value: &[]*wrapperspb.BytesValue{}}
	if got := bs.Type(); got != "bytesSliceValue" {
		t.Errorf("BytesSliceValue.Type() = %v, want bytesSliceValue", got)
	}
}

func TestBytesSliceValue_ImplementsPflagValue(t *testing.T) {
	bs := BytesSliceValue{value: &[]*wrapperspb.BytesValue{}}
	var _ pflag.Value = &bs
}

func TestBytesSliceValue_MultipleSets(t *testing.T) {
	bs := BytesSliceValue{value: &[]*wrapperspb.BytesValue{}}

	// First set - single value
	err := bs.Set("SGVsbG8=") // "Hello"
	if err != nil {
		t.Fatalf("First Set() error = %v", err)
	}

	// Second set - another single value (should append)
	err = bs.Set("V29ybGQ=") // "World"
	if err != nil {
		t.Fatalf("Second Set() error = %v", err)
	}

	// Third set - another single value (should append)
	err = bs.Set("VGhpcyBJcw==") // "This Is"
	if err != nil {
		t.Fatalf("Third Set() error = %v", err)
	}

	// Fourth set - another single value (should append)
	err = bs.Set("QSBUZXN0") // "A Test"
	if err != nil {
		t.Fatalf("Fourth Set() error = %v", err)
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

func TestBytesSliceValue_Append(t *testing.T) {
	tests := []struct {
		name    string
		initial string
		append  string
		want    []*wrapperspb.BytesValue
		wantErr bool
	}{
		{
			name:    "append to empty slice",
			initial: "",
			append:  "SGVsbG8=", // "Hello" in base64
			want:    []*wrapperspb.BytesValue{wrapperspb.Bytes([]byte("Hello"))},
			wantErr: false,
		},
		{
			name:    "append to existing slice",
			initial: "SGVsbG8=", // "Hello" in base64
			append:  "V29ybGQ=", // "World" in base64
			want: []*wrapperspb.BytesValue{
				wrapperspb.Bytes([]byte("Hello")),
				wrapperspb.Bytes([]byte("World")),
			},
			wantErr: false,
		},
		{
			name:    "append invalid base64",
			initial: "SGVsbG8=", // "Hello" in base64
			append:  "invalid-base64!",
			want: []*wrapperspb.BytesValue{
				wrapperspb.Bytes([]byte("Hello")),
			},
			wantErr: true,
		},
		{
			name:    "append with whitespace",
			initial: "SGVsbG8=",     // "Hello" in base64
			append:  "  V29ybGQ=  ", // "World" in base64 with spaces
			want: []*wrapperspb.BytesValue{
				wrapperspb.Bytes([]byte("Hello")),
				wrapperspb.Bytes([]byte("World")),
			},
			wantErr: false,
		},
		{
			name:    "append binary data",
			initial: "SGVsbG8=",         // "Hello" in base64
			append:  "AAECAwQFBgcICQoL", // binary data in base64
			want: []*wrapperspb.BytesValue{
				wrapperspb.Bytes([]byte("Hello")),
				wrapperspb.Bytes([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup initial state
			bytesSlice := []*wrapperspb.BytesValue{}
			bs := BytesSliceValue{value: &bytesSlice}

			// Set initial value if provided
			if tt.initial != "" {
				if err := bs.Set(tt.initial); err != nil {
					t.Fatalf("Failed to set initial value: %v", err)
				}
			}

			// Test Append
			err := bs.Append(tt.append)
			if (err != nil) != tt.wantErr {
				t.Errorf("Append() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(*bs.value) != len(tt.want) {
					t.Errorf("Append() length = %v, want %v", len(*bs.value), len(tt.want))
					return
				}

				for i, wantVal := range tt.want {
					if string((*bs.value)[i].Value) != string(wantVal.Value) {
						t.Errorf("Append()[%d] = %v, want %v", i, (*bs.value)[i].Value, wantVal.Value)
					}
				}
			}
		})
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
	result := BytesSlice(&initialBytes)

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
	result := BytesSlice(&nilBytes)

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

func TestBytesSliceValue_BoundaryConditions(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "single character base64",
			input:   "QQ==", // "A" in base64
			wantErr: false,
		},
		{
			name:    "empty base64 result",
			input:   "", // empty string decodes to empty bytes
			wantErr: false,
		},
		{
			name:    "whitespace only",
			input:   "   ", // whitespace only
			wantErr: false,
		},
		{
			name:    "maximum single byte",
			input:   "/w==", // 0xFF in base64
			wantErr: false,
		},
		{
			name:    "minimum single byte",
			input:   "AA==", // 0x00 in base64
			wantErr: false,
		},
		{
			name:    "invalid base64 padding",
			input:   "SGVsbG8", // missing padding
			wantErr: true,
		},
		{
			name:    "invalid base64 character",
			input:   "SGVsbG8!", // invalid character
			wantErr: true,
		},
		{
			name:    "base64 with special chars",
			input:   "SGVsbG8gV29ybGQh", // "Hello World!" with special char
			wantErr: false,
		},
		{
			name:    "unicode string base64",
			input:   "44GT44KM44Gv44GP", // "こんにちは" in base64
			wantErr: false,
		},
		{
			name:    "very long valid base64",
			input:   "VGhpcyBpcyBhIHZlcnkgbG9uZyBzdHJpbmcgdGhhdCBzaG91bGQgd29yayBmaW5lIHdpdGggYmFzZTY0IGVuY29kaW5n", // long string
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytesSlice := []*wrapperspb.BytesValue{}
			bs := BytesSliceValue{value: &bytesSlice}
			err := bs.Set(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBytesSliceValue_ErrorRecovery(t *testing.T) {
	bytesSlice := []*wrapperspb.BytesValue{}
	bs := BytesSliceValue{value: &bytesSlice}

	// Set initial values
	if err := bs.Set("SGVsbG8="); err != nil { // "Hello"
		t.Fatalf("Failed to set initial value: %v", err)
	}
	if err := bs.Set("V29ybGQ="); err != nil { // "World"
		t.Fatalf("Failed to set second value: %v", err)
	}

	// Try to set an invalid value
	if err := bs.Set("invalid-base64!"); err == nil {
		t.Error("Expected error for invalid base64, but got none")
	}

	// Verify the original values are still there
	if len(bytesSlice) != 2 {
		t.Errorf("Expected 2 elements after failed Set, got %d", len(bytesSlice))
	}
	if string(bytesSlice[0].Value) != "Hello" || string(bytesSlice[1].Value) != "World" {
		t.Errorf("Values not as expected after failed Set")
	}

	// Test Replace with partial failure
	if err := bs.Replace([]string{"VGVzdA==", "invalid!"}); err == nil {
		t.Error("Expected error for invalid base64 in Replace, but got none")
	}

	// Verify state is consistent after failed Replace
	if len(bytesSlice) != 2 {
		t.Errorf("Expected 2 elements after failed Replace, got %d", len(bytesSlice))
	}
}

func TestBytesSliceValue_LargeData(t *testing.T) {
	// Test with large data
	largeData := make([]byte, 1024*1024) // 1MB
	for i := range largeData {
		largeData[i] = byte(i % 256)
	}

	hexStr := base64.StdEncoding.EncodeToString(largeData)
	bytesSlice := []*wrapperspb.BytesValue{}
	bs := BytesSliceValue{value: &bytesSlice}

	err := bs.Set(hexStr)
	if err != nil {
		t.Fatalf("Set() with large data failed: %v", err)
	}

	if len(bytesSlice) != 1 {
		t.Fatalf("Expected 1 element, got %d", len(bytesSlice))
	}

	if len(bytesSlice[0].Value) != len(largeData) {
		t.Errorf("Data length mismatch: got %d, want %d", len(bytesSlice[0].Value), len(largeData))
	}

	// Verify data integrity
	for i := 0; i < len(largeData); i++ {
		if bytesSlice[0].Value[i] != largeData[i] {
			t.Errorf("Data mismatch at index %d: got %v, want %v", i, bytesSlice[0].Value[i], largeData[i])
			break
		}
	}
}

func TestBytesSliceValue_MassiveOperations(t *testing.T) {
	bytesSlice := []*wrapperspb.BytesValue{}
	bs := BytesSliceValue{value: &bytesSlice}

	// Add 1000 small strings
	for i := 0; i < 1000; i++ {
		value := "string_" + string(rune('0'+i%10)) + "_" + string(rune('a'+i%26))
		if err := bs.Set(base64.StdEncoding.EncodeToString([]byte(value))); err != nil {
			t.Fatalf("Failed to set string %d: %v", i, err)
		}
	}

	if len(bytesSlice) != 1000 {
		t.Errorf("Expected 1000 elements, got %d", len(bytesSlice))
	}

	// Test String() with massive data
	result := bs.String()
	if len(result) == 0 {
		t.Error("String() returned empty result for massive data")
	}

	// Test GetSlice() with massive data
	slice := bs.GetSlice()
	if len(slice) != 1000 {
		t.Errorf("GetSlice() returned %d elements, want 1000", len(slice))
	}
}

func TestBytesSliceValue_Replace(t *testing.T) {
	tests := []struct {
		name    string
		initial []string
		replace []string
		want    []*wrapperspb.BytesValue
		wantErr bool
	}{
		{
			name:    "replace empty with values",
			initial: []string{},
			replace: []string{"SGVsbG8=", "V29ybGQ="}, // "Hello", "World"
			want: []*wrapperspb.BytesValue{
				wrapperspb.Bytes([]byte("Hello")),
				wrapperspb.Bytes([]byte("World")),
			},
			wantErr: false,
		},
		{
			name:    "replace existing values",
			initial: []string{"SGVsbG8="},                 // "Hello"
			replace: []string{"V29ybGQ=", "VGhpcyBJcw=="}, // "World", "This Is"
			want: []*wrapperspb.BytesValue{
				wrapperspb.Bytes([]byte("World")),
				wrapperspb.Bytes([]byte("This Is")),
			},
			wantErr: false,
		},
		{
			name:    "replace with invalid base64",
			initial: []string{"SGVsbG8="}, // "Hello"
			replace: []string{"invalid-base64!"},
			want: []*wrapperspb.BytesValue{
				wrapperspb.Bytes([]byte("Hello")), // Should remain unchanged on error
			},
			wantErr: true,
		},
		{
			name:    "replace with whitespace",
			initial: []string{"SGVsbG8="},                         // "Hello"
			replace: []string{"  V29ybGQ=  ", "  VGhpcyBJcw==  "}, // "World", "This Is" with spaces
			want: []*wrapperspb.BytesValue{
				wrapperspb.Bytes([]byte("World")),
				wrapperspb.Bytes([]byte("This Is")),
			},
			wantErr: false,
		},
		{
			name:    "replace with empty slice",
			initial: []string{"SGVsbG8="}, // "Hello"
			replace: []string{},
			want:    []*wrapperspb.BytesValue{},
			wantErr: false,
		},
		{
			name:    "replace with binary data",
			initial: []string{"SGVsbG8="},                     // "Hello"
			replace: []string{"AAECAwQFBgcICQoL", "//79/f39"}, // binary data
			want: []*wrapperspb.BytesValue{
				wrapperspb.Bytes([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}),
				wrapperspb.Bytes([]byte{255, 254, 253, 253, 253, 253}),
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup initial state
			bytesSlice := make([]*wrapperspb.BytesValue, len(tt.initial))
			for i, v := range tt.initial {
				decoded, _ := base64.StdEncoding.DecodeString(v) // Assume initial values are valid
				bytesSlice[i] = wrapperspb.Bytes(decoded)
			}
			bs := BytesSliceValue{value: &bytesSlice}

			// Test Replace
			err := bs.Replace(tt.replace)
			if (err != nil) != tt.wantErr {
				t.Errorf("Replace() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(bytesSlice) != len(tt.want) {
					t.Errorf("Replace() length = %v, want %v", len(bytesSlice), len(tt.want))
					return
				}

				for i, wantVal := range tt.want {
					if string(bytesSlice[i].Value) != string(wantVal.Value) {
						t.Errorf("Replace()[%d] = %v, want %v", i, bytesSlice[i].Value, wantVal.Value)
					}
				}
			}
		})
	}
}

func TestBytesSliceValue_GetSlice(t *testing.T) {
	tests := []struct {
		name  string
		setup func(*BytesSliceValue)
		want  []string
	}{
		{
			name: "get from empty slice",
			setup: func(bs *BytesSliceValue) {
				// Leave empty
			},
			want: []string{},
		},
		{
			name: "get from single value slice",
			setup: func(bs *BytesSliceValue) {
				bs.Set("SGVsbG8=") // "Hello"
			},
			want: []string{"SGVSBG8="},
		},
		{
			name: "get from multiple values slice",
			setup: func(bs *BytesSliceValue) {
				bs.Set("SGVsbG8=")     // "Hello"
				bs.Set("V29ybGQ=")     // "World"
				bs.Set("VGhpcyBJcw==") // "This Is"
			},
			want: []string{"SGVSBG8=", "V29YBGQ=", "VGHPCYBJCW=="},
		},
		{
			name: "get with binary data",
			setup: func(bs *BytesSliceValue) {
				bs.Set("AAECAwQ=") // binary data
				bs.Set("//79/f39") // high bytes
			},
			want: []string{"AAECAWQ=", "//79/F39"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytesSlice := []*wrapperspb.BytesValue{}
			bs := BytesSliceValue{value: &bytesSlice}

			// Setup the test scenario
			tt.setup(&bs)

			// Test GetSlice
			got := bs.GetSlice()

			if len(got) != len(tt.want) {
				t.Errorf("GetSlice() length = %v, want %v", len(got), len(tt.want))
				return
			}

			for i, want := range tt.want {
				if got[i] != want {
					t.Errorf("GetSlice()[%d] = %v, want %v", i, got[i], want)
				}
			}
		})
	}
}

func TestBytesSliceValue_PflagInterface(t *testing.T) {
	bytesSlice := []*wrapperspb.BytesValue{}
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)

	// Test that BytesSliceValue implements pflag.Value interface
	bs := BytesSlice(&bytesSlice)
	fs.Var(bs, "bytes", "test bytes slice")

	// Test setting flags - first value
	args := []string{"--bytes", "SGVsbG8="} // "Hello" in base64
	if err := fs.Parse(args); err != nil {
		t.Errorf("Failed to parse flags: %v", err)
	}

	// Test setting flags - second value (should append)
	args = []string{"--bytes", "V29ybGQ="} // "World" in base64
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
