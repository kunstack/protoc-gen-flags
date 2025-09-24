package types

import (
	"strconv"

	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var _ pflag.Value = (*UInt64Value)(nil)

type UInt64Value wrapperspb.UInt64Value

func (u *UInt64Value) String() string {
	if u == nil {
		return "0"
	}
	return strconv.FormatUint(u.Value, 10)
}

func (u *UInt64Value) Set(s string) error {
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return err
	}
	u.Value = v
	return nil
}

func (u *UInt64Value) Type() string {
	return "uint64Value"
}

func UInt64(v *wrapperspb.UInt64Value) *UInt64Value {
	return (*UInt64Value)(v)
}
