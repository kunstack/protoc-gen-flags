package types

import (
	"io"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/kunstack/protoc-gen-flags/utils"
)

var _ pflag.Value = (*Int32SliceValue)(nil)

type Int32SliceValue struct {
	value   *[]*wrapperspb.Int32Value
	changed bool
}

// Set converts, and assigns, the comma-separated int32 argument string representation as the []*wrapperspb.Int32Value value of this flag.
// If Set is called on a flag that already has a []*wrapperspb.Int32Value assigned, the newly converted values will be appended.
func (s *Int32SliceValue) Set(val string) error {
	// remove all quote characters
	rmQuote := strings.NewReplacer(`"`, "", "'", "", "`", "")

	// read flag arguments with CSV parser
	int32StrSlice, err := utils.ReadAsCSV(rmQuote.Replace(val))
	if err != nil && err != io.EOF {
		return err
	}

	// parse int32 values into slice
	out := make([]*wrapperspb.Int32Value, 0, len(int32StrSlice))
	for _, int32Str := range int32StrSlice {
		v, err := strconv.ParseInt(strings.TrimSpace(int32Str), 10, 32)
		if err != nil {
			return err
		}
		out = append(out, wrapperspb.Int32(int32(v)))
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
func (s *Int32SliceValue) Type() string {
	return "int32SliceValue"
}

// String defines a "native" format for this int32 slice flag value.
func (s *Int32SliceValue) String() string {
	int32StrSlice := make([]string, len(*s.value))
	for i, v := range *s.value {
		int32StrSlice[i] = strconv.FormatInt(int64(v.Value), 10)
	}
	out, _ := utils.WriteAsCSV(int32StrSlice)

	return "[" + out + "]"
}

func Int32Slice(v *[]*wrapperspb.Int32Value) *Int32SliceValue {
	return &Int32SliceValue{value: v}
}