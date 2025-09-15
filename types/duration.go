package types

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/durationpb"
)

var _ pflag.Value = (*Duration)(nil)

type Duration durationpb.Duration

func (d *Duration) String() string {
	if d == nil {
		return "0s"
	}
	duration := durationpb.Duration(*d)
	return duration.AsDuration().String()
}

func (d *Duration) Set(s string) error {
	if d == nil {
		return fmt.Errorf("cannot set nil Duration")
	}
	duration, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	pbDuration := durationpb.New(duration)
	*d = Duration(*pbDuration)
	return nil
}

func (d *Duration) Type() string {
	return "pbduration"
}
