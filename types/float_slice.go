package types

import (
	"io"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/kunstack/protoc-gen-flags/utils"
)

var _ pflag.Value = (*FloatSliceValue)(nil)

type FloatSliceValue struct {
	value   *[]*wrapperspb.FloatValue
	changed bool
}

// Set converts, and assigns, the comma-separated float argument string representation as the []*wrapperspb.FloatValue value of this flag.
// If Set is called on a flag that already has a []*wrapperspb.FloatValue assigned, the newly converted values will be appended.
func (s *FloatSliceValue) Set(val string) error {
	// remove all quote characters
	rmQuote := strings.NewReplacer(`"`, "", `'`, "", "`", "")

	// read flag arguments with CSV parser
	floatStrSlice, err := utils.ReadAsCSV(rmQuote.Replace(val))
	if err != nil && err != io.EOF {
		return err
	}

	// parse float values into slice
	out := make([]*wrapperspb.FloatValue, 0, len(floatStrSlice))
	for _, floatStr := range floatStrSlice {
		f, err := strconv.ParseFloat(strings.TrimSpace(floatStr), 32)
		if err != nil {
			return err
		}
		out = append(out, wrapperspb.Float(float32(f)))
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
func (s *FloatSliceValue) Type() string {
	return "floatSliceValue"
}

// String defines a "native" format for this float slice flag value.
func (s *FloatSliceValue) String() string {
	floatStrSlice := make([]string, len(*s.value))
	for i, f := range *s.value {
		floatStrSlice[i] = strconv.FormatFloat(float64(f.Value), 'g', -1, 32)
	}
	out, _ := utils.WriteAsCSV(floatStrSlice)

	return "[" + out + "]"
}

func FloatSlice(v *[]*wrapperspb.FloatValue) *FloatSliceValue {
	return &FloatSliceValue{value: v}
}
