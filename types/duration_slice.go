package types

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/durationpb"
)

var (
	_ pflag.Value      = (*DurationSliceValue)(nil)
	_ pflag.SliceValue = (*DurationSliceValue)(nil)
)

type DurationSliceValue struct {
	value   *[]*durationpb.Duration
	changed bool
}

func (s *DurationSliceValue) Set(val string) error {
	ss := strings.Split(val, ",")
	out := make([]*durationpb.Duration, len(ss))
	for i, d := range ss {
		dur, err := time.ParseDuration(strings.TrimSpace(d))
		if err != nil {
			return err
		}
		out[i] = durationpb.New(dur)
	}
	if !s.changed {
		*s.value = out
	} else {
		*s.value = append(*s.value, out...)
	}
	s.changed = true
	return nil
}

func (s *DurationSliceValue) Append(val string) error {
	i, err := time.ParseDuration(strings.TrimSpace(val))
	if err != nil {
		return err
	}
	*s.value = append(*s.value, durationpb.New(i))
	return nil
}

func (s *DurationSliceValue) Replace(val []string) error {
	out := make([]*durationpb.Duration, len(val))
	for i, d := range val {
		dur, err := time.ParseDuration(strings.TrimSpace(d))
		if err != nil {
			return err
		}
		out[i] = durationpb.New(dur)
	}
	*s.value = out
	return nil
}

func (s *DurationSliceValue) GetSlice() []string {
	out := make([]string, len(*s.value))
	for i, d := range *s.value {
		out[i] = d.AsDuration().String()
	}
	return out
}

// Type returns a string that uniquely represents this flag's type.
func (s *DurationSliceValue) Type() string {
	return "DurationSliceValue"
}

// String defines a "native" format for this duration slice flag value.
func (s *DurationSliceValue) String() string {
	durationStrSlice := make([]string, len(*s.value))
	for i, d := range *s.value {
		durationStrSlice[i] = fmt.Sprintf("%s", d.AsDuration())
	}
	return "[" + strings.Join(durationStrSlice, ",") + "]"
}

func DurationSlice(v *[]*durationpb.Duration) *DurationSliceValue {
	return &DurationSliceValue{value: v}
}
