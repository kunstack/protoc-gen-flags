package utils

import "strings"

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
