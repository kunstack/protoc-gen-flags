package types

import (
	"encoding/hex"
	"io"
	"strings"

	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/kunstack/protoc-gen-flags/utils"
)

var (
	_ pflag.Value = (*BytesHexSliceValue)(nil)
	_ pflag.Value = (*NativeBytesHexSliceValue)(nil)
)

type BytesHexSliceValue struct {
	value   *[]*wrapperspb.BytesValue
	changed bool
}

// Set converts, and assigns, the comma-separated hexadecimal-encoded bytes argument string representation as the []*wrapperspb.BytesValue value of this flag.
// If Set is called on a flag that already has a []*wrapperspb.BytesValue assigned, the newly converted values will be appended.
func (s *BytesHexSliceValue) Set(val string) error {
	// remove all quote characters
	rmQuote := strings.NewReplacer(`"`, "", `'`, "", "`", "")

	// read flag arguments with CSV parser
	hexStrSlice, err := utils.ReadAsCSV(rmQuote.Replace(val))
	if err != nil && err != io.EOF {
		return err
	}

	// parse hexadecimal values into slice
	out := make([]*wrapperspb.BytesValue, 0, len(hexStrSlice))
	for _, hexStr := range hexStrSlice {
		decodedBytes, err := hex.DecodeString(strings.TrimSpace(hexStr))
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
func (s *BytesHexSliceValue) Type() string {
	return "bytesHexSliceValue"
}

// String defines a "native" format for this hexadecimal slice flag value.
func (s *BytesHexSliceValue) String() string {
	hexStrSlice := make([]string, len(*s.value))
	for i, b := range *s.value {
		hexStrSlice[i] = strings.ToUpper(hex.EncodeToString(b.Value))
	}
	out, _ := utils.WriteAsCSV(hexStrSlice)

	return "[" + out + "]"
}

type NativeBytesHexSliceValue struct {
	value   *[][]byte
	changed bool
}

// Set converts, and assigns, the comma-separated hex-encoded bytes argument string representation as the [][]byte value of this flag.
// If Set is called on a flag that already has a [][]byte assigned, the newly converted values will be appended.
func (s *NativeBytesHexSliceValue) Set(val string) error {
	// remove all quote characters
	rmQuote := strings.NewReplacer(`"`, "", `'`, "", "`", "")

	// read flag arguments with CSV parser
	bytesStrSlice, err := utils.ReadAsCSV(rmQuote.Replace(val))
	if err != nil && err != io.EOF {
		return err
	}

	// parse hex values into slice
	out := make([][]byte, 0, len(bytesStrSlice))
	for _, bytesStr := range bytesStrSlice {
		decodedBytes, err := hex.DecodeString(strings.TrimSpace(bytesStr))
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
func (s *NativeBytesHexSliceValue) Type() string {
	return "nativeBytesHexSliceValue"
}

// String defines a "native" format for this bytes slice flag value.
func (s *NativeBytesHexSliceValue) String() string {
	bytesStrSlice := make([]string, len(*s.value))
	for i, b := range *s.value {
		bytesStrSlice[i] = hex.EncodeToString(b)
	}
	out, _ := utils.WriteAsCSV(bytesStrSlice)

	return "[" + out + "]"
}

func BytesHexSlice[T []*wrapperspb.BytesValue | [][]byte](v T) pflag.Value {
	switch wrap := any(v).(type) {
	case []*wrapperspb.BytesValue:
		return &BytesHexSliceValue{value: &wrap}
	case [][]byte:
		return &NativeBytesHexSliceValue{value: &wrap}
	default:
		// This should never happen due to type constraints
		panic("BytesSlice: unsupported type")
	}
}
