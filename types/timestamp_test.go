package types

import (
	"testing"
	"time"

	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestTimestamp_String(t *testing.T) {
	timestamp := Timestamp(*timestamppb.New(time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)))
	expected := "2023-12-25T10:30:00Z"
	if got := timestamp.String(); got != expected {
		t.Errorf("Timestamp.String() = %v, want %v", got, expected)
	}
}

func TestTimestamp_Set(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    time.Time
		wantErr bool
	}{
		{
			name:  "valid RFC3339",
			input: "2023-12-25T10:30:00Z",
			want:  time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC),
		},
		{
			name:  "valid RFC3339 with timezone",
			input: "2023-12-25T10:30:00+02:00",
			want:  time.Date(2023, 12, 25, 8, 30, 0, 0, time.UTC),
		},
		{
			name:  "valid RFC822",
			input: "25 Dec 23 10:30 UTC",
			want:  time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC),
		},
		{
			name:  "valid RFC822Z",
			input: "25 Dec 23 10:30 +0000",
			want:  time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC),
		},
		{
			name:  "valid RFC850",
			input: "Monday, 25-Dec-23 10:30:00 UTC",
			want:  time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC),
		},
		{
			name:  "valid RFC1123",
			input: "Mon, 25 Dec 2023 10:30:00 UTC",
			want:  time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC),
		},
		{
			name:  "valid RFC1123Z",
			input: "Mon, 25 Dec 2023 10:30:00 +0000",
			want:  time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC),
		},
		{
			name:    "invalid format",
			input:   "invalid-timestamp",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ts Timestamp
			err := ts.Set(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Timestamp.Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				pbTimestamp := timestamppb.Timestamp(ts)
				got := pbTimestamp.AsTime().UTC()
				if !got.Equal(tt.want.UTC()) {
					t.Errorf("Timestamp.Set() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestTimestamp_Type(t *testing.T) {
	var ts Timestamp
	if got := ts.Type(); got != "pbtimestamp" {
		t.Errorf("Timestamp.Type() = %v, want %v", got, "pbtimestamp")
	}
}

func TestTimestamp_pflagInterface(t *testing.T) {
	// Test that Timestamp implements pflag.Value interface
	var _ pflag.Value = (*Timestamp)(nil)

	// Test with actual pflag usage
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	var ts Timestamp

	fs.Var(&ts, "timestamp", "Test timestamp")

	// Test parsing
	err := fs.Parse([]string{"--timestamp", "2023-12-25T10:30:00Z"})
	if err != nil {
		t.Fatalf("Failed to parse timestamp: %v", err)
	}

	expected := time.Date(2023, 12, 25, 10, 30, 0, 0, time.UTC)
	pbTimestamp := timestamppb.Timestamp(ts)
	got := pbTimestamp.AsTime()

	if !got.Equal(expected) {
		t.Errorf("Parsed timestamp = %v, want %v", got, expected)
	}
}

func TestTimestamp_NilPointer(t *testing.T) {
	var ts *Timestamp = nil

	// Test String() on nil pointer
	if got := ts.String(); got != "0001-01-01T00:00:00Z" {
		t.Errorf("Timestamp.String() on nil = %v, want %v", got, "0001-01-01T00:00:00Z")
	}

	// Test Set() on nil pointer
	err := ts.Set("2023-12-25T10:30:00Z")
	if err == nil {
		t.Error("Timestamp.Set() on nil should return an error")
	}

	// Test Type() on nil pointer
	if got := ts.Type(); got != "pbtimestamp" {
		t.Errorf("Timestamp.Type() on nil = %v, want %v", got, "pbtimestamp")
	}
}
