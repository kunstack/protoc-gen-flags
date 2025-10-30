// Package types provides implementations for converting protobuf field types to pflag.Value types.
// This package enables command-line flag generation for protobuf messages by implementing
// the pflag.Value interface for various protobuf field types.
package types

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/spf13/pflag"
)

var _ pflag.Value = (*JSONValue)(nil)

// JSONValue implements the pflag.Value interface for protobuf map and struct types.
// It provides JSON string representation and parsing capabilities for complex types,
// allowing them to be used as command-line flags.
type JSONValue struct {
	wrap interface{} // Pointer to the actual protobuf message (map or struct)
}

// String returns the JSON string representation of the map/struct value.
func (j *JSONValue) String() string {
	jsonBytes, err := json.Marshal(j.wrap)
	if err != nil {
		return "{}"
	}
	return string(jsonBytes)
}

// Set parses a JSON string and sets the map/struct value.
func (j *JSONValue) Set(value string) error {
	if j.wrap == nil {
		return fmt.Errorf("cannot set value on nil message")
	}
	if err := json.Unmarshal([]byte(value), j.wrap); err != nil {
		return fmt.Errorf("cannot unmarshal JSON: %w", err)
	}
	return nil
}

// Type returns the type name for help text.
func (j *JSONValue) Type() string {
	return "JSONValue"
}

// JSON creates a new JSONValue for the given protobuf message pointer.
func JSON(v interface{}) *JSONValue {
	if v == nil {
		panic("JSON: nil value")
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer || rv.IsNil() {
		panic("JSON: non-nil pointer")
	}
	return &JSONValue{wrap: v}
}
