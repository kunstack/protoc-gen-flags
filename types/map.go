// Package types provides implementations for converting protobuf field types to pflag.Value types.
// This package enables command-line flag generation for protobuf messages by implementing
// the pflag.Value interface for various protobuf field types.
package types

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
)

var _ pflag.Value = (*StringToInt32Map)(nil)
var _ pflag.Value = (*StringToUint32Map)(nil)
var _ pflag.Value = (*StringToUint64Map)(nil)

// StringToInt32Map implements pflag.Value for map[string]int32
type StringToInt32Map struct {
	value   *map[string]int32
	changed bool
}

// String returns the string representation of the map
func (s *StringToInt32Map) String() string {
	if s.value == nil || *s.value == nil {
		return ""
	}
	var buf bytes.Buffer
	buf.WriteRune('[')
	i := 0
	for k, v := range *s.value {
		if i > 0 {
			buf.WriteRune(',')
		}
		buf.WriteString(k)
		buf.WriteRune('=')
		buf.WriteString(strconv.FormatInt(int64(v), 10))
		i++
	}
	buf.WriteRune(']')
	return buf.String()
}

// Set parses a comma-separated list of key=value pairs
func (s *StringToInt32Map) Set(val string) error {
	// Handle empty input
	if val == "" {
		if !s.changed {
			*s.value = make(map[string]int32)
		}
		s.changed = true
		return nil
	}

	// Parse the comma-separated key=value pairs
	ss := strings.Split(val, ",")
	out := make(map[string]int32, len(ss))
	for _, pair := range ss {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			return fmt.Errorf("%s must be formatted as key=value", pair)
		}
		var err error
		parsedVal, err := strconv.ParseInt(kv[1], 10, 32)
		if err != nil {
			return err
		}
		out[kv[0]] = int32(parsedVal)
	}

	// Apply the changes
	if !s.changed {
		*s.value = out
	} else {
		for k, v := range out {
			(*s.value)[k] = v
		}
	}
	s.changed = true
	return nil
}

// Type returns the type name for help text
func (s *StringToInt32Map) Type() string {
	return "stringToInt32"
}

// StringToInt32 Constructor functions
func StringToInt32(v interface{}) *StringToInt32Map {
	if v == nil {
		return &StringToInt32Map{value: nil, changed: false}
	}
	// Try to cast to the expected type
	if m, ok := v.(*map[string]int32); ok {
		return &StringToInt32Map{value: m, changed: false}
	}
	return &StringToInt32Map{value: nil, changed: false}
}

func StringToUint32(v interface{}) *StringToUint32Map {
	if v == nil {
		return &StringToUint32Map{value: nil, changed: false}
	}
	// Try to cast to the expected type
	if m, ok := v.(*map[string]uint32); ok {
		return &StringToUint32Map{value: m, changed: false}
	}
	return &StringToUint32Map{value: nil, changed: false}
}

func StringToUint64(v interface{}) *StringToUint64Map {
	if v == nil {
		return &StringToUint64Map{value: nil, changed: false}
	}
	// Try to cast to the expected type
	if m, ok := v.(*map[string]uint64); ok {
		return &StringToUint64Map{value: m, changed: false}
	}
	return &StringToUint64Map{value: nil, changed: false}
}

// StringToUint32Map implements pflag.Value for map[string]uint32
type StringToUint32Map struct {
	value   *map[string]uint32
	changed bool
}

// String returns the string representation of the map
func (s *StringToUint32Map) String() string {
	if s.value == nil || *s.value == nil {
		return ""
	}
	var buf bytes.Buffer
	buf.WriteRune('[')
	i := 0
	for k, v := range *s.value {
		if i > 0 {
			buf.WriteRune(',')
		}
		buf.WriteString(k)
		buf.WriteRune('=')
		buf.WriteString(strconv.FormatUint(uint64(v), 10))
		i++
	}
	buf.WriteRune(']')
	return buf.String()
}

// Set parses a comma-separated list of key=value pairs
func (s *StringToUint32Map) Set(val string) error {
	// Handle empty input
	if val == "" {
		if !s.changed {
			*s.value = make(map[string]uint32)
		}
		s.changed = true
		return nil
	}

	// Parse the comma-separated key=value pairs
	ss := strings.Split(val, ",")
	out := make(map[string]uint32, len(ss))
	for _, pair := range ss {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			return fmt.Errorf("%s must be formatted as key=value", pair)
		}
		parsedVal, err := strconv.ParseUint(kv[1], 10, 32)
		if err != nil {
			return err
		}
		out[kv[0]] = uint32(parsedVal)
	}

	// Apply the changes
	if !s.changed {
		*s.value = out
	} else {
		for k, v := range out {
			(*s.value)[k] = v
		}
	}
	s.changed = true
	return nil
}

// Type returns the type name for help text
func (s *StringToUint32Map) Type() string {
	return "stringToUint32"
}

// StringToUint64Map implements pflag.Value for map[string]uint64
type StringToUint64Map struct {
	value   *map[string]uint64
	changed bool
}

// String returns the string representation of the map
func (s *StringToUint64Map) String() string {
	if s.value == nil || *s.value == nil {
		return ""
	}
	var buf bytes.Buffer
	buf.WriteRune('[')
	i := 0
	for k, v := range *s.value {
		if i > 0 {
			buf.WriteRune(',')
		}
		buf.WriteString(k)
		buf.WriteRune('=')
		buf.WriteString(strconv.FormatUint(v, 10))
		i++
	}
	buf.WriteRune(']')
	return buf.String()
}

// Set parses a comma-separated list of key=value pairs
func (s *StringToUint64Map) Set(val string) error {
	// Handle empty input
	if val == "" {
		if !s.changed {
			*s.value = make(map[string]uint64)
		}
		s.changed = true
		return nil
	}

	// Parse the comma-separated key=value pairs
	ss := strings.Split(val, ",")
	out := make(map[string]uint64, len(ss))
	for _, pair := range ss {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			return fmt.Errorf("%s must be formatted as key=value", pair)
		}
		var err error
		out[kv[0]], err = strconv.ParseUint(kv[1], 10, 64)
		if err != nil {
			return err
		}
	}

	// Apply the changes
	if !s.changed {
		*s.value = out
	} else {
		for k, v := range out {
			(*s.value)[k] = v
		}
	}
	s.changed = true
	return nil
}

// Type returns the type name for help text
func (s *StringToUint64Map) Type() string {
	return "stringToUint64"
}
