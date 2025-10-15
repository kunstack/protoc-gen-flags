package types

import (
	"encoding/base64"
	"io"
	"strings"

	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/kunstack/protoc-gen-flags/utils"
)

var _ pflag.Value = (*BytesSliceValue)(nil)

type BytesSliceValue struct {
	value   *[]*wrapperspb.BytesValue
	changed bool
}

// Set converts, and assigns, the comma-separated base64-encoded bytes argument string representation as the []*wrapperspb.BytesValue value of this flag.
// If Set is called on a flag that already has a []*wrapperspb.BytesValue assigned, the newly converted values will be appended.
func (s *BytesSliceValue) Set(val string) error {
	// remove all quote characters
	rmQuote := strings.NewReplacer(`"`, "", `'`, "", "`", "")

	// read flag arguments with CSV parser
	bytesStrSlice, err := utils.ReadAsCSV(rmQuote.Replace(val))
	if err != nil && err != io.EOF {
		return err
	}

	// parse base64 values into slice
	out := make([]*wrapperspb.BytesValue, 0, len(bytesStrSlice))
	for _, bytesStr := range bytesStrSlice {
		decodedBytes, err := base64.StdEncoding.DecodeString(strings.TrimSpace(bytesStr))
		if err != nil {
			return err
		}
		out = append(out, wrapperspb.Bytes(decodedBytes))
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
func (s *BytesSliceValue) Type() string {
	return "bytesSliceValue"
}

// String defines a "native" format for this bytes slice flag value.
func (s *BytesSliceValue) String() string {
	bytesStrSlice := make([]string, len(*s.value))
	for i, b := range *s.value {
		bytesStrSlice[i] = base64.StdEncoding.EncodeToString(b.Value)
	}
	out, _ := utils.WriteAsCSV(bytesStrSlice)
	return "[" + out + "]"
}

// NativeBytesSliceValue is a pflag.Value implementation for handling native [][]byte fields.
// It supports base64 encoding for command-line flag parsing.
type NativeBytesSliceValue struct {
	value   *[][]byte
	changed bool
}

// Set converts, and assigns, the comma-separated base64-encoded bytes argument string representation as the [][]byte value of this flag.
// If Set is called on a flag that already has a [][]byte assigned, the newly converted values will be appended.
func (s *NativeBytesSliceValue) Set(val string) error {
	// remove all quote characters
	rmQuote := strings.NewReplacer(`"`, "", `'`, "", "`", "")

	// read flag arguments with CSV parser
	bytesStrSlice, err := utils.ReadAsCSV(rmQuote.Replace(val))
	if err != nil && err != io.EOF {
		return err
	}

	// parse base64 values into slice
	out := make([][]byte, 0, len(bytesStrSlice))
	for _, bytesStr := range bytesStrSlice {
		decodedBytes, err := base64.StdEncoding.DecodeString(strings.TrimSpace(bytesStr))
		if err != nil {
			return err
		}
		out = append(out, decodedBytes)
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
func (s *NativeBytesSliceValue) Type() string {
	return "nativeBytesSliceValue"
}

// String defines a "native" format for this bytes slice flag value.
func (s *NativeBytesSliceValue) String() string {
	bytesStrSlice := make([]string, len(*s.value))
	for i, b := range *s.value {
		bytesStrSlice[i] = base64.StdEncoding.EncodeToString(b)
	}
	out, _ := utils.WriteAsCSV(bytesStrSlice)

	return "[" + out + "]"
}

func BytesSlice[T []*wrapperspb.BytesValue | [][]byte](v *T) pflag.Value {
	switch wrap := any(v).(type) {
	case *[]*wrapperspb.BytesValue:
		return &BytesSliceValue{value: wrap}
	case *[][]byte:
		return &NativeBytesSliceValue{value: wrap}
	default:
		// This should never happen due to type constraints
		panic("BytesSlice: unsupported type")
	}
}
