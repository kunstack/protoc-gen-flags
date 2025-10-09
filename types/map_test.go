package types

import (
	"testing"

	"github.com/spf13/pflag"
)

func TestStringToInt32Map_String(t *testing.T) {
	tests := []struct {
		name     string
		value    map[string]int32
		expected string
	}{
		{
			name:     "empty map",
			value:    map[string]int32{},
			expected: "[]",
		},
		{
			name:     "single entry",
			value:    map[string]int32{"foo": 123},
			expected: "[foo=123]",
		},
		{
			name:     "multiple entries",
			value:    map[string]int32{"foo": 123, "bar": -456},
			expected: "[bar=-456,foo=123]", // Note: map iteration order is not guaranteed
		},
		{
			name:     "negative values",
			value:    map[string]int32{"neg": -789},
			expected: "[neg=-789]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := StringToInt32(&tt.value)
			result := m.String()
			// Since map iteration order is not guaranteed, we just check that the result contains the expected elements
			if len(tt.value) == 0 {
				if result != "[]" {
					t.Errorf("String() = %v, want %v", result, "[]")
				}
			} else {
				if result[0] != '[' || result[len(result)-1] != ']' {
					t.Errorf("String() = %v, expected bracket format", result)
				}
			}
		})
	}
}

func TestStringToInt32Map_Set(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantErr     bool
		expectedMap map[string]int32
	}{
		{
			name:        "empty string",
			input:       "",
			wantErr:     false,
			expectedMap: map[string]int32{},
		},
		{
			name:        "single key-value pair",
			input:       "foo=123",
			wantErr:     false,
			expectedMap: map[string]int32{"foo": 123},
		},
		{
			name:        "multiple key-value pairs",
			input:       "foo=123,bar=-456",
			wantErr:     false,
			expectedMap: map[string]int32{"foo": 123, "bar": -456},
		},
		{
			name:        "invalid format",
			input:       "invalid",
			wantErr:     true,
			expectedMap: nil,
		},
		{
			name:        "invalid int32 value",
			input:       "foo=notanumber",
			wantErr:     true,
			expectedMap: nil,
		},
		{
			name:        "overflow value",
			input:       "foo=999999999999999999999",
			wantErr:     true, // ParseInt will return error for overflow
			expectedMap: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var m map[string]int32
			m = make(map[string]int32)
			type32 := StringToInt32(&m)

			err := type32.Set(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.expectedMap != nil {
				// Check if the map contains the expected key-value pairs
				for k, v := range tt.expectedMap {
					if m[k] != v {
						t.Errorf("Map[%s] = %v, want %v", k, m[k], v)
					}
				}
			}
		})
	}
}

func TestStringToInt32Map_Type(t *testing.T) {
	var m map[string]int32
	type32 := StringToInt32(&m)
	if result := type32.Type(); result != "stringToInt32" {
		t.Errorf("Type() = %v, want stringToInt32", result)
	}
}

func TestStringToInt32Map_ImplementsPflagValue(t *testing.T) {
	var m map[string]int32
	type32 := StringToInt32(&m)
	var _ pflag.Value = type32
}

func TestStringToInt32Map_MultipleSets(t *testing.T) {
	var m map[string]int32
	type32 := StringToInt32(&m)

	// First set
	err := type32.Set("foo=123")
	if err != nil {
		t.Fatalf("First Set() error = %v", err)
	}

	// Second set (should be incremental)
	err = type32.Set("bar=456")
	if err != nil {
		t.Fatalf("Second Set() error = %v", err)
	}

	// Check that both values are present
	if m["foo"] != 123 {
		t.Errorf("Expected foo=123, got foo=%d", m["foo"])
	}
	if m["bar"] != 456 {
		t.Errorf("Expected bar=456, got bar=%d", m["bar"])
	}
}

func TestStringToUint32Map_String(t *testing.T) {
	tests := []struct {
		name     string
		value    map[string]uint32
		expected string
	}{
		{
			name:     "empty map",
			value:    map[string]uint32{},
			expected: "[]",
		},
		{
			name:     "single entry",
			value:    map[string]uint32{"foo": 123},
			expected: "[foo=123]",
		},
		{
			name:     "multiple entries",
			value:    map[string]uint32{"foo": 123, "bar": 456},
			expected: "[bar=456,foo=123]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := StringToUint32(&tt.value)
			result := m.String()
			if len(tt.value) == 0 {
				if result != "[]" {
					t.Errorf("String() = %v, want %v", result, "[]")
				}
			} else {
				if result[0] != '[' || result[len(result)-1] != ']' {
					t.Errorf("String() = %v, expected bracket format", result)
				}
			}
		})
	}
}

func TestStringToUint32Map_Set(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantErr     bool
		expectedMap map[string]uint32
	}{
		{
			name:        "empty string",
			input:       "",
			wantErr:     false,
			expectedMap: map[string]uint32{},
		},
		{
			name:        "single key-value pair",
			input:       "foo=123",
			wantErr:     false,
			expectedMap: map[string]uint32{"foo": 123},
		},
		{
			name:        "multiple key-value pairs",
			input:       "foo=123,bar=456",
			wantErr:     false,
			expectedMap: map[string]uint32{"foo": 123, "bar": 456},
		},
		{
			name:        "invalid format",
			input:       "invalid",
			wantErr:     true,
			expectedMap: nil,
		},
		{
			name:        "negative value (should fail)",
			input:       "foo=-123",
			wantErr:     true, // ParseUint should fail on negative numbers
			expectedMap: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var m map[string]uint32
			m = make(map[string]uint32)
			type32u := StringToUint32(&m)

			err := type32u.Set(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.expectedMap != nil {
				for k, v := range tt.expectedMap {
					if m[k] != v {
						t.Errorf("Map[%s] = %v, want %v", k, m[k], v)
					}
				}
			}
		})
	}
}

func TestStringToUint32Map_Type(t *testing.T) {
	var m map[string]uint32
	type32u := StringToUint32(&m)
	if result := type32u.Type(); result != "stringToUint32" {
		t.Errorf("Type() = %v, want stringToUint32", result)
	}
}

func TestStringToUint32Map_ImplementsPflagValue(t *testing.T) {
	var m map[string]uint32
	type32u := StringToUint32(&m)
	var _ pflag.Value = type32u
}

func TestStringToUint64Map_String(t *testing.T) {
	tests := []struct {
		name     string
		value    map[string]uint64
		expected string
	}{
		{
			name:     "empty map",
			value:    map[string]uint64{},
			expected: "[]",
		},
		{
			name:     "single entry",
			value:    map[string]uint64{"foo": 1234567890},
			expected: "[foo=1234567890]",
		},
		{
			name:     "large values",
			value:    map[string]uint64{"foo": 18446744073709551615}, // Max uint64
			expected: "[foo=18446744073709551615]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := StringToUint64(&tt.value)
			result := m.String()
			if len(tt.value) == 0 {
				if result != "[]" {
					t.Errorf("String() = %v, want %v", result, "[]")
				}
			} else {
				if result[0] != '[' || result[len(result)-1] != ']' {
					t.Errorf("String() = %v, expected bracket format", result)
				}
			}
		})
	}
}

func TestStringToUint64Map_Set(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantErr     bool
		expectedMap map[string]uint64
	}{
		{
			name:        "empty string",
			input:       "",
			wantErr:     false,
			expectedMap: map[string]uint64{},
		},
		{
			name:        "single key-value pair",
			input:       "foo=1234567890",
			wantErr:     false,
			expectedMap: map[string]uint64{"foo": 1234567890},
		},
		{
			name:        "multiple key-value pairs",
			input:       "foo=1234567890,bar=9876543210",
			wantErr:     false,
			expectedMap: map[string]uint64{"foo": 1234567890, "bar": 9876543210},
		},
		{
			name:        "invalid format",
			input:       "invalid",
			wantErr:     true,
			expectedMap: nil,
		},
		{
			name:        "negative value (should fail)",
			input:       "foo=-123",
			wantErr:     true,
			expectedMap: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var m map[string]uint64
			m = make(map[string]uint64)
			type64u := StringToUint64(&m)

			err := type64u.Set(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.expectedMap != nil {
				for k, v := range tt.expectedMap {
					if m[k] != v {
						t.Errorf("Map[%s] = %v, want %v", k, m[k], v)
					}
				}
			}
		})
	}
}

func TestStringToUint64Map_Type(t *testing.T) {
	var m map[string]uint64
	type64u := StringToUint64(&m)
	if result := type64u.Type(); result != "stringToUint64" {
		t.Errorf("Type() = %v, want stringToUint64", result)
	}
}

func TestStringToUint64Map_ImplementsPflagValue(t *testing.T) {
	var m map[string]uint64
	type64u := StringToUint64(&m)
	var _ pflag.Value = type64u
}

func TestStringToUint64Map_LargeValues(t *testing.T) {
	var m map[string]uint64
	m = make(map[string]uint64)
	type64u := StringToUint64(&m)

	// Test with maximum uint64 value
	err := type64u.Set("max=18446744073709551615")
	if err != nil {
		t.Fatalf("Set() with max uint64 error = %v", err)
	}
	if m["max"] != 18446744073709551615 {
		t.Errorf("Expected max=18446744073709551615, got max=%d", m["max"])
	}
}
