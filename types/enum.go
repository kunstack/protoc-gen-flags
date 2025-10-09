// Package types provides implementations for converting protobuf field types to pflag.Value types.
// This package enables command-line flag generation for protobuf messages by implementing
// the pflag.Value interface for various protobuf field types.
package types

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var _ pflag.Value = (*EnumValue)(nil)

// EnumValue implements the pflag.Value interface for protobuf enum types.
// It provides string representation and parsing capabilities for enum values,
// allowing them to be used as command-line flags.
type EnumValue struct {
	allowedTypes []string                          // List of valid enum value names
	wrap         protoreflect.Enum                 // Pointer to the actual enum value (int32)
	descriptors  protoreflect.EnumValueDescriptors // Enum value descriptors for validation
}

// String returns the string representation of the enum value.
// If the enum value has a corresponding name in the descriptors, it returns the name.
// Otherwise, it returns the numeric value as a string.
func (e *EnumValue) String() string {
	n := e.wrap.Number()
	ev := e.descriptors.ByNumber(n)
	if ev != nil {
		return string(ev.Name())
	}
	return strconv.Itoa(int(n))
}

// Set parses the string and updates the enum value accordingly.
// It validates that the string corresponds to a valid enum value name.
// Returns an error if the string is not a valid enum value or if the enum cannot be set.
func (e *EnumValue) Set(s string) error {
	if e == nil {
		return fmt.Errorf("enum value is nil or invalid")
	}

	v := reflect.ValueOf(e.wrap)

	if v.Kind() != reflect.Ptr {
		return errors.New("argument must be a pointer to enum")
	}

	elem := v.Elem()
	if !elem.CanSet() {
		return errors.New("cannot set enum value (ensure pointer is to a settable field)")
	}

	if val := e.descriptors.ByName(protoreflect.Name(s)); val != nil {
		elem.Set(reflect.ValueOf(val))
		return nil
	}

	return fmt.Errorf("invalid enum value %q, allowed values are: %s", s, strings.Join(e.allowedTypes, ", "))
}

// Type returns the data type name of the enum value.
// This method implements the pflag.Value interface.
func (e *EnumValue) Type() string {
	return "enumValue"
}

// Enum creates a new EnumValue wrapper for the given protobuf enum.
// It extracts the allowed enum values from the enum descriptor and
// returns an EnumValue that can be used as a pflag.Value.
func Enum(wrap protoreflect.Enum) *EnumValue {
	values := wrap.Descriptor().Values()
	allowed := make([]string, 0, values.Len())

	for i := 0; i < values.Len(); i++ {
		allowed = append(allowed, string(values.Get(i).Name()))
	}

	return &EnumValue{
		wrap:         wrap,
		descriptors:  values,
		allowedTypes: allowed,
	}
}
