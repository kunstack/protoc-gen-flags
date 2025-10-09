package types

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ pflag.Value = (*TimestampValue)(nil)

type TimestampValue struct {
	layouts []string
	wrap    *timestamppb.Timestamp
}

func (t *TimestampValue) String() string {
	return t.wrap.AsTime().Format(time.RFC3339Nano)
}

// parseTimeFormat converts a format string name to the actual time format
func parseTimeFormat(format string) string {
	switch format {
	case "RFC339":
		return time.RFC3339
	case "RFC3339":
		return time.RFC3339
	case "RFC3339Nano":
		return time.RFC3339Nano
	case "RFC822":
		return time.RFC822
	case "RFC822Z":
		return time.RFC822Z
	case "RFC850":
		return time.RFC850
	case "RFC1123":
		return time.RFC1123
	case "RFC1123Z":
		return time.RFC1123Z
	case "ISO8601":
		return "2006-01-02"
	case "ISO8601Time":
		return "2006-01-02T15:04:05"
	case "Kitchen":
		return time.Kitchen
	case "Stamp":
		return time.Stamp
	case "StampMilli":
		return time.StampMilli
	case "StampMicro":
		return time.StampMicro
	case "StampNano":
		return time.StampNano
	case "DateTime":
		return "2006-01-02 15:04:05"
	case "DateOnly":
		return "2006-01-02"
	case "TimeOnly":
		return "15:04:05"
	default:
		return format
	}
}

func (t *TimestampValue) Set(s string) error {
	for _, format := range t.layouts {
		actualFormat := parseTimeFormat(format)
		value, err := time.Parse(actualFormat, s)
		if err == nil {
			t.wrap.Seconds = value.Unix()
			t.wrap.Nanos = int32(value.Nanosecond())
			return nil
		}
	}
	return fmt.Errorf("invalid time format `%s` must be one of: %s", s, strings.Join(t.layouts, ","))
}

func (t *TimestampValue) Type() string {
	return "timestampValue"
}

func Timestamp(wrap *timestamppb.Timestamp, formats []string) *TimestampValue {
	return &TimestampValue{wrap: wrap, layouts: formats}
}
