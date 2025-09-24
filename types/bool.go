package types

import (
	"strconv"

	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var _ pflag.Value = (*BoolValue)(nil)

type BoolValue wrapperspb.BoolValue

func (b *BoolValue) String() string {
	return strconv.FormatBool(b.Value)
}

func (b *BoolValue) Set(s string) error {
	v, err := strconv.ParseBool(s)
	b.Value = v
	return err
}

func (b *BoolValue) Type() string {
	return "boolValue"
}

func Bool(v *wrapperspb.BoolValue) *BoolValue {
	return (*BoolValue)(v)
}
