package types

import (
	"io"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/durationpb"

	"github.com/kunstack/protoc-gen-flags/utils"
)

var _ pflag.Value = (*DurationSliceValue)(nil)

type DurationSliceValue struct {
	value   *[]*durationpb.Duration
	changed bool
}

// Set converts, and assigns, the comma-separated duration argument string representation as the []*durationpb.Duration value of this flag.
// If Set is called on a flag that already has a []*durationpb.Duration assigned, the newly converted values will be appended.
func (s *DurationSliceValue) Set(val string) error {
	// remove all quote characters
	rmQuote := strings.NewReplacer(`"`, "", `'`, "", "`", "")

	// read flag arguments with CSV parser
	ss, err := utils.ReadAsCSV(rmQuote.Replace(val))
	if err != nil && err != io.EOF {
		return err
	}

	// parse duration values into slice
	out := make([]*durationpb.Duration, 0, len(ss))
	for _, durationStr := range ss {
		duration, err := time.ParseDuration(strings.TrimSpace(durationStr))
		if err != nil {
			return err
		}
		out = append(out, durationpb.New(duration))
	}

	if !s.changed {
		*s.value = out
	} else {
		*s.value = append(*s.value, out...)
	}

	s.changed = true
	return nil
}

// Type returns a string that uniquely represents this flag's type.
func (s *DurationSliceValue) Type() string {
	return "durationSliceValue"
}

// String defines a "native" format for this duration slice flag value.
func (s *DurationSliceValue) String() string {
	durationStrSlice := make([]string, len(*s.value))
	for i, d := range *s.value {
		durationStrSlice[i] = d.AsDuration().String()
	}
	out, _ := utils.WriteAsCSV(durationStrSlice)

	return "[" + out + "]"
}

func DurationSlice(v *[]*durationpb.Duration) *DurationSliceValue {
	return &DurationSliceValue{value: v}
}
