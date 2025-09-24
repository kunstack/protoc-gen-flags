package types

import (
	"strconv"

	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var _ pflag.Value = (*FloatValue)(nil)

type FloatValue wrapperspb.FloatValue

func (b *FloatValue) String() string {
	if b == nil {
		return "0"
	}
	return strconv.FormatFloat(float64(b.Value), 'g', -1, 32)
}

func (b *FloatValue) Set(s string) error {
	v, err := strconv.ParseFloat(s, 32)
	if err != nil {
		return err
	}
	b.Value = float32(v)
	return nil
}

func (b *FloatValue) Type() string {
	return "floatValue"
}

func Float(v *wrapperspb.FloatValue) *FloatValue {
	return (*FloatValue)(v)
}
