package types

import (
	"math"
	"strconv"
	"testing"

	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestDoubleSliceValue_Set(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []*wrapperspb.DoubleValue
		wantErr bool
	}{
		{
			name:  "single double",
			input: "3.14159",
			want:  []*wrapperspb.DoubleValue{wrapperspb.Double(3.14159)},
		},
		{
			name:  "multiple doubles",
			input: "1.5,2.7,3.14",
			want:  []*wrapperspb.DoubleValue{wrapperspb.Double(1.5), wrapperspb.Double(2.7), wrapperspb.Double(3.14)},
		},
		{
			name:    "with spaces",
			input:   " 1.5 , 2.7 , 3.14 ",
			wantErr: true, // strconv.ParseFloat doesn't handle spaces
		},
		{
			name:    "invalid double",
			input:   "1.5,invalid,3.14",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true, // strconv.ParseFloat doesn't handle empty string
		},
		{
			name:  "zero value",
			input: "0",
			want:  []*wrapperspb.DoubleValue{wrapperspb.Double(0)},
		},
		{
			name:  "negative values",
			input: "-1.5,-2.7,-3.14",
			want:  []*wrapperspb.DoubleValue{wrapperspb.Double(-1.5), wrapperspb.Double(-2.7), wrapperspb.Double(-3.14)},
		},
		{
			name:  "scientific notation",
			input: "1.23e-4,5.67e+8",
			want:  []*wrapperspb.DoubleValue{wrapperspb.Double(1.23e-4), wrapperspb.Double(5.67e+8)},
		},
		{
			name:    "quoted strings",
			input:   `"1.5","2.7","3.14"`,
			wantErr: true, // strconv.ParseFloat doesn't handle quotes
		},
		{
			name:  "integer input",
			input: "1,2,3",
			want:  []*wrapperspb.DoubleValue{wrapperspb.Double(1), wrapperspb.Double(2), wrapperspb.Double(3)},
		},
		{
			name:  "very large numbers",
			input: "1.7976931348623157e+308,-1.7976931348623157e+308",
			want:  []*wrapperspb.DoubleValue{wrapperspb.Double(1.7976931348623157e+308), wrapperspb.Double(-1.7976931348623157e+308)},
		},
		{
			name:  "very small numbers",
			input: "5e-324,-5e-324",
			want:  []*wrapperspb.DoubleValue{wrapperspb.Double(5e-324), wrapperspb.Double(-5e-324)},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doubleSlice := []*wrapperspb.DoubleValue{}
			ds := DoubleSliceValue{value: &doubleSlice}
			err := ds.Set(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("DoubleSliceValue.Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(*ds.value) != len(tt.want) {
					t.Errorf("DoubleSliceValue.Set() length = %v, want %v", len(*ds.value), len(tt.want))
					return
				}

				for i, wantVal := range tt.want {
					if (*ds.value)[i].Value != wantVal.Value {
						t.Errorf("DoubleSliceValue.Set()[%d] = %v, want %v", i, (*ds.value)[i].Value, wantVal.Value)
					}
				}
			}
		})
	}
}

func TestDoubleSliceValue_String(t *testing.T) {
	tests := []struct {
		name  string
		input []*wrapperspb.DoubleValue
		want  string
	}{
		{
			name:  "empty slice",
			input: []*wrapperspb.DoubleValue{},
			want:  "[]",
		},
		{
			name:  "single double",
			input: []*wrapperspb.DoubleValue{wrapperspb.Double(3.14159)},
			want:  "[3.14159]",
		},
		{
			name:  "multiple values",
			input: []*wrapperspb.DoubleValue{wrapperspb.Double(1.5), wrapperspb.Double(2.7), wrapperspb.Double(3.14)},
			want:  "[1.5,2.7,3.14]",
		},
		{
			name:  "zero values",
			input: []*wrapperspb.DoubleValue{wrapperspb.Double(0), wrapperspb.Double(0)},
			want:  "[0,0]",
		},
		{
			name:  "negative values",
			input: []*wrapperspb.DoubleValue{wrapperspb.Double(-1.5), wrapperspb.Double(-2.7)},
			want:  "[-1.5,-2.7]",
		},
		{
			name:  "scientific notation",
			input: []*wrapperspb.DoubleValue{wrapperspb.Double(1.23e-4), wrapperspb.Double(5.67e+8)},
			want:  "[0.000123,5.67e+08]",
		},
		{
			name:  "large values",
			input: []*wrapperspb.DoubleValue{wrapperspb.Double(1.7976931348623157e+308)},
			want:  "[1.7976931348623157e+308]",
		},
		{
			name:  "small values",
			input: []*wrapperspb.DoubleValue{wrapperspb.Double(5e-324)},
			want:  "[5e-324]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds := DoubleSliceValue{value: &tt.input}
			if got := ds.String(); got != tt.want {
				t.Errorf("DoubleSliceValue.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDoubleSliceValue_Type(t *testing.T) {
	doubleSlice := []*wrapperspb.DoubleValue{}
	ds := DoubleSliceValue{value: &doubleSlice}
	if got := ds.Type(); got != "doubleSliceValue" {
		t.Errorf("DoubleSliceValue.Type() = %v, want doubleSliceValue", got)
	}
}

func TestDoubleSliceValue_ImplementsPflagValue(t *testing.T) {
	doubleSlice := []*wrapperspb.DoubleValue{}
	ds := DoubleSliceValue{value: &doubleSlice}
	var _ pflag.Value = &ds
}

func TestDoubleSliceValue_MultipleSets(t *testing.T) {
	doubleSlice := []*wrapperspb.DoubleValue{}
	ds := DoubleSliceValue{value: &doubleSlice}

	// First set
	err := ds.Set("1.5,2.7,3.14")
	if err != nil {
		t.Fatalf("First Set() error = %v", err)
	}

	// Second set (should append)
	err = ds.Set("4.2,5.6,6.9")
	if err != nil {
		t.Fatalf("Second Set() error = %v", err)
	}

	// Verify all values are present
	expected := []float64{1.5, 2.7, 3.14, 4.2, 5.6, 6.9}
	if len(*ds.value) != len(expected) {
		t.Fatalf("Expected %d values, got %d", len(expected), len(*ds.value))
	}

	for i, expectedVal := range expected {
		if (*ds.value)[i].Value != expectedVal {
			t.Errorf("Expected %v at index %d, got %v", expectedVal, i, (*ds.value)[i].Value)
		}
	}
}

func TestDoubleSliceValue_EmptyString(t *testing.T) {
	doubleSlice := []*wrapperspb.DoubleValue{}
	ds := DoubleSliceValue{value: &doubleSlice}

	err := ds.Set("")
	if err == nil {
		t.Fatalf("Expected error for empty string, but got none")
	}

	// Verify slice remains empty after failed Set
	if len(*ds.value) != 0 {
		t.Errorf("Expected empty slice after failed Set, got %d elements", len(*ds.value))
	}
}

func TestDoubleSliceValue_SpecialValues(t *testing.T) {
	// Test with special float values
	specialValues := []float64{math.Inf(1), math.Inf(-1), math.NaN()}
	encoded := make([]string, len(specialValues))
	for i, val := range specialValues {
		encoded[i] = strconv.FormatFloat(val, 'g', -1, 64)
	}

	doubleSlice := []*wrapperspb.DoubleValue{}
	ds := DoubleSliceValue{value: &doubleSlice}

	input := encoded[0] + "," + encoded[1] + "," + encoded[2]
	err := ds.Set(input)
	if err != nil {
		t.Fatalf("Set() with special values error = %v", err)
	}

	if len(*ds.value) != 3 {
		t.Fatalf("Expected 3 values, got %d", len(*ds.value))
	}

	// Check that String() method encodes back correctly
	got := ds.String()
	if got != "["+input+"]" {
		t.Errorf("String() = %v, want %v", got, "["+input+"]")
	}
}

func TestDoubleSlice_Constructor(t *testing.T) {
	// Test with non-nil slice
	initialDoubles := []*wrapperspb.DoubleValue{wrapperspb.Double(3.14)}
	result := DoubleSlice(&initialDoubles)

	if result == nil {
		t.Fatal("DoubleSlice() returned nil")
	}

	if result.value == nil {
		t.Error("DoubleSlice() did not set the value pointer")
	}

	if result.changed {
		t.Error("DoubleSlice() should initialize with changed=false")
	}
}

func TestDoubleSlice_ConstructorNil(t *testing.T) {
	// Test with nil slice
	var nilDoubles []*wrapperspb.DoubleValue
	result := DoubleSlice(&nilDoubles)

	if result == nil {
		t.Fatal("DoubleSlice() returned nil")
	}

	if result.value == nil {
		t.Error("DoubleSlice() did not set the value pointer for nil input")
	}

	if result.changed {
		t.Error("DoubleSlice() should initialize with changed=false for nil input")
	}
}

func TestDoubleSliceValue_PflagInterface(t *testing.T) {
	doubleSlice := make([]*wrapperspb.DoubleValue, 0)
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)

	// Test that DoubleSliceValue implements pflag.Value interface
	ds := DoubleSlice(&doubleSlice)
	fs.Var(ds, "doubles", "test double slice")

	// Test setting flags
	args := []string{"--doubles", "1.5,2.7,3.14"}
	if err := fs.Parse(args); err != nil {
		t.Errorf("Failed to parse flags: %v", err)
	}

	// Verify the parsed values
	if len(doubleSlice) != 3 {
		t.Errorf("Expected 3 values, got %d", len(doubleSlice))
	}

	if (doubleSlice)[0].Value != 1.5 {
		t.Errorf("Expected 1.5 at index 0, got %v", (doubleSlice)[0].Value)
	}

	if (doubleSlice)[1].Value != 2.7 {
		t.Errorf("Expected 2.7 at index 1, got %v", (doubleSlice)[1].Value)
	}

	if (doubleSlice)[2].Value != 3.14 {
		t.Errorf("Expected 3.14 at index 2, got %v", (doubleSlice)[2].Value)
	}
}
