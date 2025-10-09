package types

import (
	"fmt"
	"io"
	"strings"

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
	// remove all quote characters
	rmQuote := strings.NewReplacer(`"`, "", `'`, "", "`", "")

	// read flag arguments with CSV parser
	enumStrSlice, err := utils.ReadAsCSV(rmQuote.Replace(val))
	if err != nil && err != io.EOF {
		return err
	}

	// parse enum values into slice
	out := make([]protoreflect.Enum, 0, len(enumStrSlice))
	for _, enumStr := range enumStrSlice {
		enumStr = strings.TrimSpace(enumStr)

		// Try to find the enum value by name
		if ev := s.descriptors.ByName(protoreflect.Name(enumStr)); ev != nil {
			// Create a new enum value and set it
			// We need to create the appropriate enum type - this is a simplified approach
			// In practice, the caller should provide the concrete enum type
			out = append(out, ev)
		} else {
			return fmt.Errorf("invalid enum value %q, allowed values are: %s", enumStr, strings.Join(s.allowedTypes, ", "))
		}
	}

	if !s.changed {
		*s.value = out
	} else {
		*s.value = append(*s.value, out...)
	}

	s.changed = true

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