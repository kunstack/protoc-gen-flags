package types

import (
	"testing"

	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestFloatSliceValue_Set(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []float32
	}{
		{
			name:     "single value",
			input:    "3.14",
			expected: []float32{3.14},
		},
		{
			name:     "multiple values",
			input:    "1.23,4.56,7.89",
			expected: []float32{1.23, 4.56, 7.89},
		},
		{
			name:     "values with spaces",
			input:    " 1.0 , 2.0 , 3.0 ",
			expected: []float32{1.0, 2.0, 3.0},
		},
		{
			name:     "scientific notation",
			input:    "1.23e-4,5.67e+2",
			expected: []float32{1.23e-4, 5.67e+2},
		},
		{
			name:     "quoted values",
			input:    `"1.23","4.56"`,
			expected: []float32{1.23, 4.56},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var slice []*wrapperspb.FloatValue
			fsv := FloatSlice(&slice)

			err := fsv.Set(tt.input)
			if err != nil {
				t.Errorf("FloatSliceValue.Set() error = %v", err)
				return
			}

			if len(slice) != len(tt.expected) {
				t.Errorf("FloatSliceValue.Set() length = %v, want %v", len(slice), len(tt.expected))
				return
			}

			for i, expected := range tt.expected {
				if slice[i].Value != expected {
					t.Errorf("FloatSliceValue.Set()[%d] = %v, want %v", i, slice[i].Value, expected)
				}
			}
		})
	}
}

func TestFloatSliceValue_String(t *testing.T) {
	tests := []struct {
		name     string
		input    []*wrapperspb.FloatValue
		expected string
	}{
		{
			name:     "empty slice",
			input:    []*wrapperspb.FloatValue{},
			expected: "[]",
		},
		{
			name:     "single value",
			input:    []*wrapperspb.FloatValue{wrapperspb.Float(3.14)},
			expected: "[3.14]",
		},
		{
			name:     "multiple values",
			input:    []*wrapperspb.FloatValue{wrapperspb.Float(1.23), wrapperspb.Float(4.56), wrapperspb.Float(7.89)},
			expected: "[1.23,4.56,7.89]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsv := FloatSlice(&tt.input)
			result := fsv.String()
			if result != tt.expected {
				t.Errorf("FloatSliceValue.String() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFloatSliceValue_Type(t *testing.T) {
	var slice []*wrapperspb.FloatValue
	fsv := FloatSlice(&slice)
	if fsv.Type() != "floatSliceValue" {
		t.Errorf("FloatSliceValue.Type() = %v, want %v", fsv.Type(), "floatSliceValue")
	}
}

func TestFloatSliceValue_Append(t *testing.T) {
	var slice []*wrapperspb.FloatValue
	fsv := FloatSlice(&slice)

	// Set initial values
	err := fsv.Set("1.23,4.56")
	if err != nil {
		t.Errorf("FloatSliceValue.Set() error = %v", err)
		return
	}

	// Append more values
	err = fsv.Set("7.89,10.11")
	if err != nil {
		t.Errorf("FloatSliceValue.Set() error = %v", err)
		return
	}

	expected := []float32{1.23, 4.56, 7.89, 10.11}
	if len(slice) != len(expected) {
		t.Errorf("FloatSliceValue append length = %v, want %v", len(slice), len(expected))
		return
	}

	for i, expected := range expected {
		if slice[i].Value != expected {
			t.Errorf("FloatSliceValue append[%d] = %v, want %v", i, slice[i].Value, expected)
		}
	}
}
