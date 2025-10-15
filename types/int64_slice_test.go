package types

import (
	"testing"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestInt64SliceValue_Set(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []int64
	}{
		{
			name:     "single value",
			input:    "42",
			expected: []int64{42},
		},
		{
			name:     "multiple values",
			input:    "1,2,3",
			expected: []int64{1, 2, 3},
		},
		{
			name:     "values with spaces",
			input:    " 1 , 2 , 3 ",
			expected: []int64{1, 2, 3},
		},
		{
			name:     "negative values",
			input:    "-1,-2,-3",
			expected: []int64{-1, -2, -3},
		},
		{
			name:     "zero values",
			input:    "0,0,0",
			expected: []int64{0, 0, 0},
		},
		{
			name:     "max int64 values",
			input:    "9223372036854775807",
			expected: []int64{9223372036854775807},
		},
		{
			name:     "min int64 values",
			input:    "-9223372036854775808",
			expected: []int64{-9223372036854775808},
		},
		{
			name:     "quoted values",
			input:    `"1","2","3"`,
			expected: []int64{1, 2, 3},
		},
		{
			name:     "mixed positive and negative",
			input:    "1,-2,3,-4",
			expected: []int64{1, -2, 3, -4},
		},
		{
			name:     "large numbers",
			input:    "1000000000,2000000000,3000000000",
			expected: []int64{1000000000, 2000000000, 3000000000},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var slice []*wrapperspb.Int64Value
			isv := Int64Slice(&slice)

			err := isv.Set(tt.input)
			if err != nil {
				t.Errorf("Int64SliceValue.Set() error = %v", err)
				return
			}

			if len(slice) != len(tt.expected) {
				t.Errorf("Int64SliceValue.Set() length = %v, want %v", len(slice), len(tt.expected))
				return
			}

			for i, expected := range tt.expected {
				if slice[i].Value != expected {
					t.Errorf("Int64SliceValue.Set()[%d] = %v, want %v", i, slice[i].Value, expected)
				}
			}
		})
	}
}

func TestInt64SliceValue_Set_Error(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantErr  bool
	}{
		{
			name:    "invalid integer",
			input:   "abc",
			wantErr: true,
		},
		{
			name:    "float value",
			input:   "3.14",
			wantErr: true,
		},
		{
			name:    "overflows int64",
			input:   "9223372036854775808",
			wantErr: true,
		},
		{
			name:    "underflows int64",
			input:   "-9223372036854775809",
			wantErr: true,
		},
		{
			name:    "mixed valid and invalid",
			input:   "1,abc,3",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var slice []*wrapperspb.Int64Value
			isv := Int64Slice(&slice)

			err := isv.Set(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Int64SliceValue.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInt64SliceValue_String(t *testing.T) {
	tests := []struct {
		name     string
		input    []*wrapperspb.Int64Value
		expected string
	}{
		{
			name:     "empty slice",
			input:    []*wrapperspb.Int64Value{},
			expected: "[]",
		},
		{
			name:     "single value",
			input:    []*wrapperspb.Int64Value{wrapperspb.Int64(42)},
			expected: "[42]",
		},
		{
			name:     "multiple values",
			input:    []*wrapperspb.Int64Value{wrapperspb.Int64(1), wrapperspb.Int64(2), wrapperspb.Int64(3)},
			expected: "[1,2,3]",
		},
		{
			name:     "negative values",
			input:    []*wrapperspb.Int64Value{wrapperspb.Int64(-1), wrapperspb.Int64(-2)},
			expected: "[-1,-2]",
		},
		{
			name:     "zero values",
			input:    []*wrapperspb.Int64Value{wrapperspb.Int64(0), wrapperspb.Int64(0)},
			expected: "[0,0]",
		},
		{
			name:     "large values",
			input:    []*wrapperspb.Int64Value{wrapperspb.Int64(9223372036854775807), wrapperspb.Int64(-9223372036854775808)},
			expected: "[9223372036854775807,-9223372036854775808]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isv := Int64Slice(&tt.input)
			result := isv.String()
			if result != tt.expected {
				t.Errorf("Int64SliceValue.String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestInt64SliceValue_Type(t *testing.T) {
	var slice []*wrapperspb.Int64Value
	isv := Int64Slice(&slice)
	if isv.Type() != "int64SliceValue" {
		t.Errorf("Int64SliceValue.Type() = %v, want %v", isv.Type(), "int64SliceValue")
	}
}

func TestInt64SliceValue_Append(t *testing.T) {
	var slice []*wrapperspb.Int64Value
	isv := Int64Slice(&slice)

	// Set initial values
	err := isv.Set("1,2")
	if err != nil {
		t.Errorf("Int64SliceValue.Set() error = %v", err)
		return
	}

	// Append more values
	err = isv.Set("3,4")
	if err != nil {
		t.Errorf("Int64SliceValue.Set() error = %v", err)
		return
	}

	expected := []int64{1, 2, 3, 4}
	if len(slice) != len(expected) {
		t.Errorf("Int64SliceValue append length = %v, want %v", len(slice), len(expected))
		return
	}

	for i, expected := range expected {
		if slice[i].Value != expected {
			t.Errorf("Int64SliceValue append[%d] = %v, want %v", i, slice[i].Value, expected)
		}
	}
}

func TestInt64SliceValue_EmptyString(t *testing.T) {
	var slice []*wrapperspb.Int64Value
	isv := Int64Slice(&slice)

	err := isv.Set("")
	if err != nil {
		t.Errorf("Int64SliceValue.Set() error = %v", err)
		return
	}

	if len(slice) != 0 {
		t.Errorf("Int64SliceValue.Set() empty string length = %v, want 0", len(slice))
	}
}

func TestInt64SliceValue_SingleQuotes(t *testing.T) {
	var slice []*wrapperspb.Int64Value
	isv := Int64Slice(&slice)

	err := isv.Set(`'1','2','3'`)
	if err != nil {
		t.Errorf("Int64SliceValue.Set() error = %v", err)
		return
	}

	expected := []int64{1, 2, 3}
	if len(slice) != len(expected) {
		t.Errorf("Int64SliceValue.Set() single quotes length = %v, want %v", len(slice), len(expected))
		return
	}

	for i, expected := range expected {
		if slice[i].Value != expected {
			t.Errorf("Int64SliceValue.Set() single quotes[%d] = %v, want %v", i, slice[i].Value, expected)
		}
	}
}

func TestInt64SliceValue_Backticks(t *testing.T) {
	var slice []*wrapperspb.Int64Value
	isv := Int64Slice(&slice)

	err := isv.Set("`1`,`2`,`3`")
	if err != nil {
		t.Errorf("Int64SliceValue.Set() error = %v", err)
		return
	}

	expected := []int64{1, 2, 3}
	if len(slice) != len(expected) {
		t.Errorf("Int64SliceValue.Set() backticks length = %v, want %v", len(slice), len(expected))
		return
	}

	for i, expected := range expected {
		if slice[i].Value != expected {
			t.Errorf("Int64SliceValue.Set() backticks[%d] = %v, want %v", i, slice[i].Value, expected)
		}
	}
}

func TestInt64Slice_Constructor(t *testing.T) {
	var slice []*wrapperspb.Int64Value
	isv := Int64Slice(&slice)

	// Test that the returned Int64SliceValue has correct fields
	if isv.value != &slice {
		t.Errorf("Int64Slice() value field not set correctly")
	}
	if isv.changed {
		t.Errorf("Int64Slice() changed field should be false initially")
	}

	// Test Type() method on constructed int64 slice
	if isv.Type() != "int64SliceValue" {
		t.Errorf("Int64Slice() Type() = %v, want %v", isv.Type(), "int64SliceValue")
	}
}

func TestInt64Slice_ConstructorNil(t *testing.T) {
	var slice []*wrapperspb.Int64Value
	isv := Int64Slice(&slice)

	// Test that it works with nil value
	if isv.value != &slice {
		t.Errorf("Int64Slice() with nil value not set correctly")
	}
	if isv.changed {
		t.Errorf("Int64Slice() with nil value changed field should be false initially")
	}
}

func TestInt64SliceValue_PflagInterface(t *testing.T) {
	var slice []*wrapperspb.Int64Value
	isv := Int64Slice(&slice)

	// Test that it implements pflag.Value interface
	var _ interface {
		Set(string) error
		String() string
		Type() string
	} = isv

	// Test setting value through interface
	err := isv.Set("42")
	if err != nil {
		t.Errorf("pflag.Value.Set() error = %v", err)
		return
	}

	if isv.String() != "[42]" {
		t.Errorf("pflag.Value.String() = %v, want [42]", isv.String())
	}

	if isv.Type() != "int64SliceValue" {
		t.Errorf("pflag.Value.Type() = %v, want int64SliceValue", isv.Type())
	}
}

func TestInt64SliceValue_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []int64
	}{
		{
			name:     "very large positive",
			input:    "9223372036854775807",
			expected: []int64{9223372036854775807},
		},
		{
			name:     "very large negative",
			input:    "-9223372036854775808",
			expected: []int64{-9223372036854775808},
		},
		{
			name:     "mixed large and small",
			input:    "9223372036854775807,-9223372036854775808,0,1,-1",
			expected: []int64{9223372036854775807, -9223372036854775808, 0, 1, -1},
		},
		{
			name:     "powers of 2",
			input:    "1,2,4,8,16,32,64,128,256,512,1024,2048",
			expected: []int64{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048},
		},
		{
			name:     "powers of 10",
			input:    "10,100,1000,10000,100000,1000000,10000000,100000000,1000000000",
			expected: []int64{10, 100, 1000, 10000, 100000, 1000000, 10000000, 100000000, 1000000000},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var slice []*wrapperspb.Int64Value
			isv := Int64Slice(&slice)

			err := isv.Set(tt.input)
			if err != nil {
				t.Errorf("Int64SliceValue.Set() error = %v", err)
				return
			}

			if len(slice) != len(tt.expected) {
				t.Errorf("Int64SliceValue.Set() length = %v, want %v", len(slice), len(tt.expected))
				return
			}

			for i, expected := range tt.expected {
				if slice[i].Value != expected {
					t.Errorf("Int64SliceValue.Set()[%d] = %v, want %v", i, slice[i].Value, expected)
				}
			}
		})
	}
}