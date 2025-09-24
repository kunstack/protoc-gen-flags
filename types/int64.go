package types

import (
	"strconv"

	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var _ pflag.Value = (*Int64Value)(nil)

type Int64Value wrapperspb.Int64Value

func (i *Int64Value) String() string {
	if i == nil {
		return "0"
	}
	return strconv.FormatInt(i.Value, 10)
}

func (i *Int64Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}
	i.Value = v
	return nil
}

func (i *Int64Value) Type() string {
	return "int64Value"
}

func Int64(v *wrapperspb.Int64Value) *Int64Value {
	return (*Int64Value)(v)
}
