package types

import (
	"fmt"

	"github.com/spf13/pflag"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/kunstack/protoc-gen-flags/utils"
)

var _ pflag.Value = (*EnumSliceValue)(nil)

type EnumSliceValue struct {
	value        *[]protoreflect.Enum
	changed      bool
	allowedTypes []string                          // List of valid enum value names
	descriptors  protoreflect.EnumValueDescriptors // Enum value descriptors for validation
}

// Set converts, and assigns, the comma-separated enum argument string representation as the []protoreflect.Enum value of this flag.
// If Set is called on a flag that already has a []protoreflect.Enum assigned, the newly converted values will be appended.
func (s *EnumSliceValue) Set(val string) error {

	return nil
}

// Type returns a string that uniquely represents this flag's type.
func (s *EnumSliceValue) Type() string {
	return "enumSliceValue"
}

// String defines a "native" format for this enum slice flag value.
func (s *EnumSliceValue) String() string {
	enumStrSlice := make([]string, len(*s.value))
	for i, e := range *s.value {
		ev := s.descriptors.ByNumber(e.Number())
		if ev != nil {
			enumStrSlice[i] = string(ev.Name())
		} else {
			enumStrSlice[i] = fmt.Sprintf("%d", e.Number())
		}
	}
	out, _ := utils.WriteAsCSV(enumStrSlice)

	return "[" + out + "]"
}

func EnumSlice(v []protoreflect.Enum, descriptors protoreflect.EnumValueDescriptors, allowedTypes []string) *EnumSliceValue {
	result := EnumSliceValue{
		value:        &v,
		descriptors:  descriptors,
		allowedTypes: allowedTypes,
	}
	return &result
}
