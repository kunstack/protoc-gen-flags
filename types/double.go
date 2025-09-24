package types

import (
	"strconv"

	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var _ pflag.Value = (*DoubleValue)(nil)

type DoubleValue wrapperspb.DoubleValue

func (d *DoubleValue) String() string {
	if d == nil {
		return "0"
	}
	return strconv.FormatFloat(d.Value, 'g', -1, 64)
}

func (d *DoubleValue) Set(s string) error {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return err
	}
	d.Value = v
	return nil
}

func (d *DoubleValue) Type() string {
	return "doubleValue"
}

func Double(v *wrapperspb.DoubleValue) *DoubleValue {
	return (*DoubleValue)(v)
}
