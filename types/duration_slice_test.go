package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/durationpb"
)

func TestDurationSliceValue_Set(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []*durationpb.Duration
		wantErr  bool
	}{
		{
			name:     "single duration",
			input:    "5s",
			expected: []*durationpb.Duration{durationpb.New(5 * time.Second)},
			wantErr:  false,
		},
		{
			name:  "multiple durations",
			input: "1s,2m,3h",
			expected: []*durationpb.Duration{
				durationpb.New(1 * time.Second),
				durationpb.New(2 * time.Minute),
				durationpb.New(3 * time.Hour),
			},
			wantErr: false,
		},
		{
			name:  "quoted durations",
			input: `"1s","2m"`,
			expected: []*durationpb.Duration{
				durationpb.New(1 * time.Second),
				durationpb.New(2 * time.Minute),
			},
			wantErr: false, // Quoted durations are supported with CSV parser
		},
		{
			name:    "invalid duration",
			input:   "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value := []*durationpb.Duration{}
			d := DurationSliceValue{value: &value}

			err := d.Set(tt.input)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, value)
			}
		})
	}
}

func TestDurationSliceValue_String(t *testing.T) {
	tests := []struct {
		name     string
		value    []*durationpb.Duration
		expected string
	}{
		{
			name:     "empty slice",
			value:    []*durationpb.Duration{},
			expected: "[]",
		},
		{
			name: "single duration",
			value: []*durationpb.Duration{
				durationpb.New(5 * time.Second),
			},
			expected: "[5s]",
		},
		{
			name: "multiple durations",
			value: []*durationpb.Duration{
				durationpb.New(1 * time.Second),
				durationpb.New(2 * time.Minute),
				durationpb.New(3 * time.Hour),
			},
			expected: "[1s,2m0s,3h0m0s]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := DurationSlice(tt.value)
			assert.Equal(t, tt.expected, d.String())
		})
	}
}

func TestDurationSliceValue_Type(t *testing.T) {
	d := DurationSlice([]*durationpb.Duration{})
	assert.Equal(t, "durationSliceValue", d.Type())
}

func TestDurationSliceValue_Append(t *testing.T) {
	value := []*durationpb.Duration{}
	d := DurationSliceValue{value: &value}

	// First set
	err := d.Set("1s,2m")
	assert.NoError(t, err)
	assert.Len(t, value, 2)

	// Second set should append
	err = d.Set("3h")
	assert.NoError(t, err)
	assert.Len(t, value, 3)
	assert.Equal(t, "1s", value[0].AsDuration().String())
	assert.Equal(t, "2m0s", value[1].AsDuration().String())
	assert.Equal(t, "3h0m0s", value[2].AsDuration().String())
}
