package types

import (
	"testing"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestInt32SliceValue_Set(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []int32
	}{
		{
			name:     "single value",
			input:    "42",
			expected: []int32{42},
		},
		{
			name:     "multiple values",
			input:    "1,2,3",
			expected: []int32{1, 2, 3},
		},
		{
			name:     "values with spaces",
			input:    " 1 , 2 , 3 ",
			expected: []int32{1, 2, 3},
		},
		{
			name:     "negative values",
			input:    "-1,-2,-3",
			expected: []int32{-1, -2, -3},
		},
		{
			name:     "zero values",
			input:    "0,0,0",
			expected: []int32{0, 0, 0},
		},
		{
			name:     "max int32 values",
			input:    "2147483647",
			expected: []int32{2147483647},
		},
		{
			name:     "min int32 values",
			input:    "-2147483648",
			expected: []int32{-2147483648},
		},

		{
			name:     "mixed positive and negative",
			input:    "1,-2,3,-4",
			expected: []int32{1, -2, 3, -4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var slice []*wrapperspb.Int32Value
			isv := Int32Slice(&slice)

			err := isv.Set(tt.input)
			if err != nil {
				t.Errorf("Int32SliceValue.Set() error = %v", err)
				return
			}

			if len(slice) != len(tt.expected) {
				t.Errorf("Int32SliceValue.Set() length = %v, want %v", len(slice), len(tt.expected))
				return
			}

			for i, expected := range tt.expected {
				if slice[i].Value != expected {
					t.Errorf("Int32SliceValue.Set()[%d] = %v, want %v", i, slice[i].Value, expected)
				}
			}
		})
	}
}

func TestInt32SliceValue_Set_Error(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
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
			name:    "overflows int32",
			input:   "2147483648",
			wantErr: true,
		},
		{
			name:    "underflows int32",
			input:   "-2147483649",
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
			var slice []*wrapperspb.Int32Value
			isv := Int32Slice(&slice)

			err := isv.Set(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Int32SliceValue.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInt32SliceValue_String(t *testing.T) {
	tests := []struct {
		name     string
		input    []*wrapperspb.Int32Value
		expected string
	}{
		{
			name:     "empty slice",
			input:    []*wrapperspb.Int32Value{},
			expected: "[]",
		},
		{
			name:     "single value",
			input:    []*wrapperspb.Int32Value{wrapperspb.Int32(42)},
			expected: "[42]",
		},
		{
			name:     "multiple values",
			input:    []*wrapperspb.Int32Value{wrapperspb.Int32(1), wrapperspb.Int32(2), wrapperspb.Int32(3)},
			expected: "[1,2,3]",
		},
		{
			name:     "negative values",
			input:    []*wrapperspb.Int32Value{wrapperspb.Int32(-1), wrapperspb.Int32(-2)},
			expected: "[-1,-2]",
		},
		{
			name:     "zero values",
			input:    []*wrapperspb.Int32Value{wrapperspb.Int32(0), wrapperspb.Int32(0)},
			expected: "[0,0]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isv := Int32Slice(&tt.input)
			result := isv.String()
			if result != tt.expected {
				t.Errorf("Int32SliceValue.String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestInt32SliceValue_Type(t *testing.T) {
	var slice []*wrapperspb.Int32Value
	isv := Int32Slice(&slice)
	if isv.Type() != "int32SliceValue" {
		t.Errorf("Int32SliceValue.Type() = %v, want %v", isv.Type(), "int32SliceValue")
	}
}

func TestInt32SliceValue_Append(t *testing.T) {
	var slice []*wrapperspb.Int32Value
	isv := Int32Slice(&slice)

	// Set initial values
	err := isv.Set("1,2")
	if err != nil {
		t.Errorf("Int32SliceValue.Set() error = %v", err)
		return
	}

	// Append more values
	err = isv.Set("3,4")
	if err != nil {
		t.Errorf("Int32SliceValue.Set() error = %v", err)
		return
	}

	expected := []int32{1, 2, 3, 4}
	if len(slice) != len(expected) {
		t.Errorf("Int32SliceValue append length = %v, want %v", len(slice), len(expected))
		return
	}

	for i, expected := range expected {
		if slice[i].Value != expected {
			t.Errorf("Int32SliceValue append[%d] = %v, want %v", i, slice[i].Value, expected)
		}
	}
}

func TestInt32Slice_Constructor(t *testing.T) {
	var slice []*wrapperspb.Int32Value
	isv := Int32Slice(&slice)

	// Test that the returned Int32SliceValue has correct fields
	if isv.value != &slice {
		t.Errorf("Int32Slice() value field not set correctly")
	}
	if isv.changed {
		t.Errorf("Int32Slice() changed field should be false initially")
	}

	// Test Type() method on constructed int32 slice
	if isv.Type() != "int32SliceValue" {
		t.Errorf("Int32Slice() Type() = %v, want %v", isv.Type(), "int32SliceValue")
	}
}

func TestInt32Slice_ConstructorNil(t *testing.T) {
	var slice []*wrapperspb.Int32Value
	isv := Int32Slice(&slice)

	// Test that it works with nil value
	if isv.value != &slice {
		t.Errorf("Int32Slice() with nil value not set correctly")
	}
	if isv.changed {
		t.Errorf("Int32Slice() with nil value changed field should be false initially")
	}
}

func TestInt32SliceValue_PflagInterface(t *testing.T) {
	var slice []*wrapperspb.Int32Value
	isv := Int32Slice(&slice)

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

	if isv.Type() != "int32SliceValue" {
		t.Errorf("pflag.Value.Type() = %v, want int32SliceValue", isv.Type())
	}
}
