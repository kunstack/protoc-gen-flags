package utils

import (
	"encoding/csv"
	"strings"
)

// BuildFlagName builds a hierarchical configuration flag name from prefix and suffix parts.
// This is a convenience function for the common case of building configuration flag names.
//
// Examples:
//
//	prefix := []string{"app", "config"}
//	BuildFlagName(prefix, "timeout") // returns "app.config.timeout"
//	BuildFlagName(nil, "timeout") // returns "timeout"
func BuildFlagName(prefix []string, parts ...string) string {
	allParts := make([]string, 0, len(prefix)+len(parts))
	allParts = append(allParts, prefix...)
	allParts = append(allParts, parts...)
	return DotJoin(allParts...)
}

// DotJoin builds a command-line flag name by combining parts with dots.
// Accepts variadic parameters to build hierarchical flag names.
// Empty parts are ignored, and proper dot separation is ensured.
//
// Examples:
//
//	DotJoin("hello", "world") // returns "hello.world"
//	DotJoin("", "world") // returns "world"
//	DotJoin("app", "config", "timeout") // returns "app.config.timeout"
//	DotJoin("app", "", "timeout") // returns "app.timeout"
//	DotJoin("service") // returns "service"
func DotJoin(parts ...string) string {
	// Filter out empty parts and remove leading/trailing dots
	var nonEmptyParts []string
	for _, part := range parts {
		if part != "" {
			// Remove leading and trailing dots
			trimmed := strings.Trim(part, ".")
			if trimmed != "" {
				nonEmptyParts = append(nonEmptyParts, trimmed)
			}
		}
	}

	return strings.Join(nonEmptyParts, ".")
}

func ReadAsCSV(val string) ([]string, error) {
	if val == "" {
		return []string{}, nil
	}
	stringReader := strings.NewReader(val)
	csvReader := csv.NewReader(stringReader)
	return csvReader.Read()
}

// WriteAsCSV writes a slice of strings as a CSV string.
func WriteAsCSV(vals []string) (string, error) {
	b := &strings.Builder{}
	w := csv.NewWriter(b)
	err := w.Write(vals)
	if err != nil {
		return "", err
	}
	w.Flush()
	return strings.TrimSuffix(b.String(), "\n"), nil
}
