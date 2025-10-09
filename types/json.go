// Package types provides implementations for converting protobuf field types to pflag.Value types.
// This package enables command-line flag generation for protobuf messages by implementing
// the pflag.Value interface for various protobuf field types.
package types

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/pflag"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/structpb"
)

var _ pflag.Value = (*JSONValue)(nil)

// JSONValue implements the pflag.Value interface for protobuf map and struct types.
// It provides JSON string representation and parsing capabilities for complex types,
// allowing them to be used as command-line flags.
type JSONValue struct {
	wrap protoreflect.Message // Pointer to the actual protobuf message (map or struct)
}

// String returns the JSON string representation of the map/struct value.
func (j *JSONValue) String() string {
	if j.wrap == nil {
		return "{}"
	}

	// Convert protobuf message to JSON
	jsonBytes, err := json.Marshal(j.wrap.Interface())
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

	// Parse JSON string into a generic interface{}
	var jsonData interface{}
	if err := json.Unmarshal([]byte(value), &jsonData); err != nil {
		return fmt.Errorf("invalid JSON: %v", err)
	}

	// Convert to protobuf struct
	structVal, err := structpb.NewStruct(jsonData.(map[string]interface{}))
	if err != nil {
		return fmt.Errorf("cannot convert JSON to protobuf struct: %v", err)
	}

	// Find the appropriate field to set - try common field names
	fields := j.wrap.Descriptor().Fields()

	// For map fields, look for "fields" field
	if field := fields.ByName("fields"); field != nil {
		j.wrap.Set(field, protoreflect.ValueOfMessage(structVal.ProtoReflect()))
		return nil
	}

	// For other struct types, try to unmarshal directly
	// This is a simplified approach - in production you might need more sophisticated handling
	return fmt.Errorf("cannot set JSON value: no suitable field found")
}

// Type returns the type name for help text.
func (j *JSONValue) Type() string {
	return "JSON"
}

// JSON creates a new JSONValue for the given protobuf message pointer.
func JSON(v interface{}) *JSONValue {
	if v == nil {
		return &JSONValue{wrap: nil}
	}

	// Ensure we have a message type
	msg, ok := v.(protoreflect.Message)
	if !ok {
		return &JSONValue{wrap: nil}
	}

	return &JSONValue{wrap: msg}
}
