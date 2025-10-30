package utils

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"google.golang.org/protobuf/types/known/durationpb"
)

// ParseTimeFormat converts a format string name to the actual time format
func ParseTimeFormat(format string) string {
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

// Units are required to go in order from biggest to smallest.
// This guards against confusion from "1m1d" being 1 minute + 1 day, not 1 month + 1 day.
var unitMap = map[string]struct {
	pos  int
	mult uint64
}{
	"ns": {8, uint64(time.Nanosecond)},
	"ms": {7, uint64(time.Millisecond)},
	"s":  {6, uint64(time.Second)},
	"m":  {5, uint64(time.Minute)},
	"h":  {4, uint64(time.Hour)},
	"d":  {3, uint64(24 * time.Hour)},
	"w":  {2, uint64(7 * 24 * time.Hour)},
	"y":  {1, uint64(365 * 24 * time.Hour)},
}

func digit(c byte) bool { return c >= '0' && c <= '9' }

func ParseDuration(s string) (*durationpb.Duration, error) {
	zero := &durationpb.Duration{}

	switch s {
	case "0":
		// Allow 0 without a unit.
		return zero, nil
	case "":
		return zero, errors.New("empty duration string")
	}

	orig := s
	var dur uint64
	lastUnitPos := 0

	for s != "" {
		if !digit(s[0]) {
			return zero, fmt.Errorf("not a valid duration string: %q", orig)
		}
		// Consume [0-9]*
		i := 0
		for ; i < len(s) && digit(s[i]); i++ {
		}
		v, err := strconv.ParseUint(s[:i], 10, 0)
		if err != nil {
			return zero, fmt.Errorf("not a valid duration string: %q", orig)
		}
		s = s[i:]

		// Consume unit.
		for i = 0; i < len(s) && !digit(s[i]); i++ {
		}
		if i == 0 {
			return zero, fmt.Errorf("not a valid duration string: %q", orig)
		}
		u := s[:i]
		s = s[i:]
		unit, ok := unitMap[u]
		if !ok {
			return zero, fmt.Errorf("unknown unit %q in duration %q", u, orig)
		}
		if unit.pos <= lastUnitPos { // Units must go in order from biggest to smallest.
			return zero, fmt.Errorf("not a valid duration string: %q", orig)
		}
		lastUnitPos = unit.pos
		// Check if the provided duration overflows time.Duration (> ~ 290years).
		if v > 1<<63/unit.mult {
			return zero, errors.New("duration out of range")
		}
		dur += v * unit.mult
		if dur > 1<<63-1 {
			return zero, errors.New("duration out of range")
		}
	}

	return durationpb.New(time.Duration(dur)), nil
}
