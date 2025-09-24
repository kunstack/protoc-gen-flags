package types

import (
	"strconv"

	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var _ pflag.Value = (*UInt32Value)(nil)

type UInt32Value wrapperspb.UInt32Value

func (u *UInt32Value) String() string {
	if u == nil {
		return "0"
	}
	return strconv.FormatUint(uint64(u.Value), 10)
}

func (u *UInt32Value) Set(s string) error {
	v, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		return err
	}
	u.Value = uint32(v)
	return nil
}

func (u *UInt32Value) Type() string {
	return "uint32Value"
}

func UInt32(v *wrapperspb.UInt32Value) *UInt32Value {
	return (*UInt32Value)(v)
}
