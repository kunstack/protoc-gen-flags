package types

import (
	"io"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/kunstack/protoc-gen-flags/utils"
)

var _ pflag.Value = (*DoubleSliceValue)(nil)

type DoubleSliceValue struct {
	value   *[]*wrapperspb.DoubleValue
	changed bool
}

// Set converts, and assigns, the comma-separated double argument string representation as the []*wrapperspb.DoubleValue value of this flag.
// If Set is called on a flag that already has a []*wrapperspb.DoubleValue assigned, the newly converted values will be appended.
func (s *DoubleSliceValue) Set(val string) error {
	// remove all quote characters
	rmQuote := strings.NewReplacer(`"`, "", `'`, "", "`", "")

	// read flag arguments with CSV parser
	doubleStrSlice, err := utils.ReadAsCSV(rmQuote.Replace(val))
	if err != nil && err != io.EOF {
		return err
	}

	// parse double values into slice
	out := make([]*wrapperspb.DoubleValue, 0, len(doubleStrSlice))
	for _, doubleStr := range doubleStrSlice {
		d, err := strconv.ParseFloat(strings.TrimSpace(doubleStr), 64)
		if err != nil {
			return err
		}
		out = append(out, wrapperspb.Double(d))
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
func (s *DoubleSliceValue) Type() string {
	return "doubleSliceValue"
}

// String defines a "native" format for this double slice flag value.
func (s *DoubleSliceValue) String() string {
	doubleStrSlice := make([]string, len(*s.value))
	for i, d := range *s.value {
		doubleStrSlice[i] = strconv.FormatFloat(d.Value, 'g', -1, 64)
	}
	out, _ := utils.WriteAsCSV(doubleStrSlice)

	return "[" + out + "]"
}

func DoubleSlice(v []*wrapperspb.DoubleValue) *DoubleSliceValue {
	return &DoubleSliceValue{value: &v}
}
