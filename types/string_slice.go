package types

import (
	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/kunstack/protoc-gen-flags/utils"
)

var _ pflag.Value = (*StringSliceValue)(nil)

type StringSliceValue struct {
	value   *[]*wrapperspb.StringValue
	changed bool
}

// Set converts, and assigns, the comma-separated string argument string representation as the []*wrapperspb.StringValue value of this flag.
// If Set is called on a flag that already has a []*wrapperspb.StringValue assigned, the newly converted values will be appended.
func (s *StringSliceValue) Set(val string) error {
	values, err := utils.ReadAsCSV(val)
	if err != nil {
		return err
	}
	if !s.changed {
		*s.value = make([]*wrapperspb.StringValue, 0, len(values))
		for _, v := range values {
			*s.value = append(*s.value, wrapperspb.String(v))
		}
	} else {
		for _, v := range values {
			*s.value = append(*s.value, wrapperspb.String(v))
		}
	}
	s.changed = true
	return nil
}

// Type returns a string that uniquely represents this flag's type.
func (s *StringSliceValue) Type() string {
	return "stringSliceValue"
}

// String defines a "native" format for this string slice flag value.
func (s *StringSliceValue) String() string {
	stringStrSlice := make([]string, len(*s.value))
	for i, str := range *s.value {
		stringStrSlice[i] = str.Value
	}
	out, _ := utils.WriteAsCSV(stringStrSlice)

	return "[" + out + "]"
}

func StringSlice(v *[]*wrapperspb.StringValue) *StringSliceValue {
	return &StringSliceValue{value: v}
}
