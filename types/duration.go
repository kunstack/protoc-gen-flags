package types

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/durationpb"
)

var _ pflag.Value = (*DurationValue)(nil)

type DurationValue durationpb.Duration

func (d *DurationValue) String() string {
	if d == nil {
		return "0s"
	}
	return fmt.Sprintf("%s", (*durationpb.Duration)(d).AsDuration())
}

func (d *DurationValue) Set(s string) error {
	if d == nil {
		return fmt.Errorf("cannot set nil Duration")
	}
	duration, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	nanos := duration.Nanoseconds()
	secs := nanos / 1e9
	nanos -= secs * 1e9
	d.Seconds = int64(secs)
	d.Nanos = int32(nanos)
	return nil
}

func (d *DurationValue) Type() string {
	return "durationValue"
}

func Duration(value *durationpb.Duration) *DurationValue {
	return (*DurationValue)(value)
}
