package types

import (
	"strconv"

	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var _ pflag.Value = (*Int32Value)(nil)

type Int32Value wrapperspb.Int32Value

func (i *Int32Value) String() string {
	if i == nil {
		return "0"
	}
	return strconv.FormatInt(int64(i.Value), 10)
}

func (i *Int32Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return err
	}
	i.Value = int32(v)
	return nil
}

func (i *Int32Value) Type() string {
	return "int32Value"
}

func Int32(v *wrapperspb.Int32Value) *Int32Value {
	return (*Int32Value)(v)
}
