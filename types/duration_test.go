package types

import (
	"testing"
	"time"

	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/durationpb"
)

func TestDuration_String(t *testing.T) {
	duration := Duration(*durationpb.New(5 * time.Minute))
	expected := "5m0s"
	if got := duration.String(); got != expected {
		t.Errorf("Duration.String() = %v, want %v", got, expected)
	}
}

func TestDuration_Set(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    time.Duration
		wantErr bool
	}{
		{
			name:  "valid duration",
			input: "5m30s",
			want:  5*time.Minute + 30*time.Second,
		},
		{
			name:  "hours",
			input: "2h15m",
			want:  2*time.Hour + 15*time.Minute,
		},
		{
			name:    "invalid duration",
			input:   "invalid",
			wantErr: true,
		},
		{
			name:  "milliseconds",
			input: "500ms",
			want:  500 * time.Millisecond,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d Duration
			err := d.Set(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Duration.Set() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				pbDuration := durationpb.Duration(d)
				got := pbDuration.AsDuration()
				if got != tt.want {
					t.Errorf("Duration.Set() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestDuration_Type(t *testing.T) {
	var d Duration
	if got := d.Type(); got != "pbduration" {
		t.Errorf("Duration.Type() = %v, want %v", got, "pbduration")
	}
}

func TestDuration_pflagInterface(t *testing.T) {
	// Test that Duration implements pflag.Value interface
	var _ pflag.Value = (*Duration)(nil)

	// Test with actual pflag usage
	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	var d Duration

	fs.Var(&d, "duration", "Test duration")

	// Test parsing
	err := fs.Parse([]string{"--duration", "1h30m"})
	if err != nil {
		t.Fatalf("Failed to parse duration: %v", err)
	}

	expected := 1*time.Hour + 30*time.Minute
	pbDuration := durationpb.Duration(d)
	got := pbDuration.AsDuration()

	if got != expected {
		t.Errorf("Parsed duration = %v, want %v", got, expected)
	}
}

func TestDuration_EmptyString(t *testing.T) {
	var d Duration
	err := d.Set("")
	if err == nil {
		t.Error("Duration.Set(\"\") should return an error")
	}
}

func TestDuration_NilPointer(t *testing.T) {
	var d *Duration = nil

	// Test String() on nil pointer
	if got := d.String(); got != "0s" {
		t.Errorf("Duration.String() on nil = %v, want %v", got, "0s")
	}

	// Test Set() on nil pointer
	err := d.Set("5m")
	if err == nil {
		t.Error("Duration.Set() on nil should return an error")
	}

	// Test Type() on nil pointer
	if got := d.Type(); got != "pbduration" {
		t.Errorf("Duration.Type() on nil = %v, want %v", got, "pbduration")
	}
}
