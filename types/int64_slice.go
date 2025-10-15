package types

import (
	"io"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/kunstack/protoc-gen-flags/utils"
)

var _ pflag.Value = (*Int64SliceValue)(nil)

type Int64SliceValue struct {
	value   *[]*wrapperspb.Int64Value
	changed bool
}

// Set converts, and assigns, the comma-separated int64 argument string representation as the []*wrapperspb.Int64Value value of this flag.
// If Set is called on a flag that already has a []*wrapperspb.Int64Value assigned, the newly converted values will be appended.
func (s *Int64SliceValue) Set(val string) error {
	// remove all quote characters
	rmQuote := strings.NewReplacer(`"`, "", "'", "", "`", "")

	// read flag arguments with CSV parser
	int64StrSlice, err := utils.ReadAsCSV(rmQuote.Replace(val))
	if err != nil && err != io.EOF {
		return err
	}

	// parse int64 values into slice
	out := make([]*wrapperspb.Int64Value, 0, len(int64StrSlice))
	for _, int64Str := range int64StrSlice {
		v, err := strconv.ParseInt(strings.TrimSpace(int64Str), 10, 64)
		if err != nil {
			return err
		}
		out = append(out, wrapperspb.Int64(v))
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
func (s *Int64SliceValue) Type() string {
	return "int64SliceValue"
}

// String defines a "native" format for this int64 slice flag value.
func (s *Int64SliceValue) String() string {
	int64StrSlice := make([]string, len(*s.value))
	for i, v := range *s.value {
		int64StrSlice[i] = strconv.FormatInt(v.Value, 10)
	}
	out, _ := utils.WriteAsCSV(int64StrSlice)

	return "[" + out + "]"
}

func Int64Slice(v *[]*wrapperspb.Int64Value) *Int64SliceValue {
	return &Int64SliceValue{value: v}
}