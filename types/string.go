package types

import (
	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var _ pflag.Value = (*StringValue)(nil)

type StringValue wrapperspb.StringValue

func (s *StringValue) String() string {
	if s == nil {
		return ""
	}
	return s.Value
}

func (s *StringValue) Set(value string) error {
	s.Value = value
	return nil
}

func (s *StringValue) Type() string {
	return "stringValue"
}

func String(v *wrapperspb.StringValue) *StringValue {
	return (*StringValue)(v)
}
