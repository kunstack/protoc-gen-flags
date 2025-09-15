package types

import (
	"errors"
	"fmt"
	"time"

	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ pflag.Value = (*Timestamp)(nil)

type Timestamp timestamppb.Timestamp

func parseTime(s string) (time.Time, error) {
	for _, format := range []string{
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
	} {
		t, err := time.Parse(format, s)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, errors.New("cannot parse timestamp, timestamp supported format: RFC822 / RFC822Z / RFC850 / RFC1123 / RFC1123Z / RFC3339")
}

func (t *Timestamp) String() string {
	if t == nil {
		return time.Time{}.Format(time.RFC3339)
	}
	timestamp := timestamppb.Timestamp(*t)
	return timestamp.AsTime().Format(time.RFC3339)
}

func (t *Timestamp) Set(s string) error {
	if t == nil {
		return fmt.Errorf("cannot set nil Timestamp")
	}
	timestamp, err := parseTime(s)
	if err != nil {
		return err
	}
	pbTimestamp := timestamppb.New(timestamp)
	*t = Timestamp(*pbTimestamp)
	return nil
}

func (t *Timestamp) Type() string {
	return "pbtimestamp"
}
