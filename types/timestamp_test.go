package types

import (
	"testing"
	"time"

	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestTimestampValue_String(t *testing.T) {
	tests := []struct {
		name     string
		format   string
		warpTime time.Time
		expected string
	}{
		{
			name:     "RFC3339Nano format (default)",
			format:   time.RFC3339,
			warpTime: time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC),
			expected: "2023-12-25T10:30:00Z",
		},
		{
			name:     "RFC1123 format",
			format:   time.RFC1123,
			warpTime: time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC),
			expected: "2023-12-25T10:30:00Z",
		},
		{
			name:     "RFC822 format",
			format:   time.RFC822,
			warpTime: time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC),
			expected: "2023-12-25T10:30:00Z",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := &TimestampValue{
				layouts: []string{tt.format},
				wrap:    timestamppb.New(tt.warpTime),
			}
			if got := ts.String(); got != tt.expected {
				t.Errorf("TimestampValue.String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestTimestampValue_Set(t *testing.T) {
	tests := []struct {
		name         string
		formats      []string
		input        string
		wantTime     time.Time
		wantErr      bool
		errorMessage string
	}{
		{
			name:     "valid RFC3339",
			formats:  []string{time.RFC3339},
			input:    "2023-12-25T10:30:00Z",
			wantTime: time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC),
		},
		{
			name:     "valid RFC822Z",
			formats:  []string{time.RFC822Z},
			input:    "25 Dec 23 10:30 +0000",
			wantTime: time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC),
		},
		{
			name:     "valid RFC822",
			formats:  []string{time.RFC822},
			input:    "25 Dec 23 10:30 UTC",
			wantTime: time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC),
		},
		{
			name:     "valid RFC850",
			formats:  []string{time.RFC850},
			input:    "Monday, 25-Dec-23 10:30:00 UTC",
			wantTime: time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC),
		},
		{
			name:     "valid RFC1123",
			formats:  []string{time.RFC1123},
			input:    "Mon, 25 Dec 2023 10:30:00 UTC",
			wantTime: time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC),
		},
		{
			name:     "valid RFC1123Z",
			formats:  []string{time.RFC1123Z},
			input:    "Mon, 25 Dec 2023 10:30:00 +0000",
			wantTime: time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC),
		},
		{
			name:         "invalid format",
			formats:      []string{time.RFC3339},
			input:        "invalid-timestamp",
			wantErr:      true,
			errorMessage: "invalid time format `invalid-timestamp` must be one of: 2006-01-02T15:04:05Z07:00",
		},
		{
			name:         "empty string",
			formats:      []string{time.RFC3339},
			input:        "",
			wantErr:      true,
			errorMessage: "invalid time format `` must be one of: 2006-01-02T15:04:05Z07:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := &TimestampValue{
				wrap:    &timestamppb.Timestamp{},
				layouts: tt.formats,
			}
			err := ts.Set(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("TimestampValue.Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errorMessage {
				t.Errorf("TimestampValue.Set() error message = %v, want %v", err.Error(), tt.errorMessage)
			}
			if !tt.wantErr {
				got := ts.wrap.AsTime()
				if !got.Equal(tt.wantTime) {
					t.Errorf("TimestampValue.Set() time = %v, want %v", got, tt.wantTime)
				}
			}
		})
	}
}

func TestTimestampValue_Type(t *testing.T) {
	ts := &TimestampValue{}
	if got := ts.Type(); got != "timestampValue" {
		t.Errorf("TimestampValue.Type() = %v, want %v", got, "timestampValue")
	}
}

func TestTimestampValue_ImplementsPflagValue(t *testing.T) {
	var _ pflag.Value = (*TimestampValue)(nil)
}

func TestTimestampValue_PflagIntegration(t *testing.T) {
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	ts := &TimestampValue{
		wrap:    &timestamppb.Timestamp{},
		layouts: []string{time.RFC3339},
	}

	fs.Var(ts, "timestamp", "Test timestamp")

	err := fs.Parse([]string{"--timestamp", "2023-12-25T10:30:00Z"})
	if err != nil {
		t.Fatalf("Failed to parse timestamp: %v", err)
	}

	expected := time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)
	got := ts.wrap.AsTime()

	if !got.Equal(expected) {
		t.Errorf("Parsed timestamp = %v, want %v", got, expected)
	}
}

func TestTimestampValue_EmptyFormat(t *testing.T) {
	ts := &TimestampValue{
		wrap:    timestamppb.New(time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)),
		layouts: []string{},
	}

	// Test String() with empty format - should use RFC3339Nano
	got := ts.String()
	expected := "2023-12-25T10:30:00Z"
	if got != expected {
		t.Errorf("TimestampValue.String() with empty format = %v, want %v", got, expected)
	}
}

func TestTimestampValue_MultipleFormatsPriority(t *testing.T) {
	ts := &TimestampValue{
		wrap:    &timestamppb.Timestamp{},
		layouts: []string{time.RFC3339, time.RFC1123},
	}

	// Test that it tries RFC3339 first (which matches this format)
	err := ts.Set("2023-12-25T10:30:00Z")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)
	got := ts.wrap.AsTime()

	if !got.Equal(expected) {
		t.Errorf("TimestampValue.Set() time = %v, want %v", got, expected)
	}
}

func TestTimestampValue_TimezoneHandling(t *testing.T) {
	ts := &TimestampValue{
		wrap:    &timestamppb.Timestamp{},
		layouts: []string{time.RFC1123Z},
	}

	// Test timezone offset handling
	err := ts.Set("Mon, 25 Dec 2023 10:30:00 +0200")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := time.Date(2023, 12, 25, 8, 30, 0, 0, time.UTC) // 10:30 +02:00 = 8:30 UTC
	got := ts.wrap.AsTime()

	if !got.Equal(expected) {
		t.Errorf("Timezone conversion failed: got %v, want %v", got, expected)
	}
}

func TestTimestamp_Constructor(t *testing.T) {
	wrap := timestamppb.New(time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC))

	ts := Timestamp(wrap, []string{time.RFC3339})

	// Test that the returned TimestampValue has correct fields
	if ts.wrap != wrap {
		t.Errorf("Timestamp() wrap = %p, want %p", ts.wrap, wrap)
	}

	if len(ts.layouts) == 0 || ts.layouts[0] != time.RFC3339 {
		t.Errorf("Timestamp() format = %v, want %v", ts.layouts, time.RFC3339)
	}

	// Test String() method on constructed timestamp
	expected := "2023-12-25T10:30:00Z"
	if got := ts.String(); got != expected {
		t.Errorf("Timestamp().String() = %v, want %v", got, expected)
	}
}

func TestTimestamp_ConstructorNil(t *testing.T) {
	ts := Timestamp(nil, []string{time.RFC3339})

	// Test that it works with nil wrap
	if ts.wrap != nil {
		t.Errorf("Timestamp(nil).wrap = %p, want nil", ts.wrap)
	}

	if len(ts.layouts) == 0 || ts.layouts[0] != time.RFC3339 {
		t.Errorf("Timestamp(nil).layouts = %v, want %v", ts.layouts, time.RFC3339)
	}
}
