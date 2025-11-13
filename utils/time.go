package utils

import (
	"time"
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
