package types

import (
	"encoding/base64"
	"strings"

	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var _ pflag.Value = (*BytesValue)(nil)

type BytesValue wrapperspb.BytesValue

func (b *BytesValue) String() string {
	if b == nil {
		return ""
	}
	return base64.StdEncoding.EncodeToString(b.Value)
}

func (b *BytesValue) Set(value string) error {
	bin, err := base64.StdEncoding.DecodeString(strings.TrimSpace(value))
	if err != nil {
		return err
	}
	b.Value = bin
	return nil
}

func (b *BytesValue) Type() string {
	return "bytesBase64"
}

func Bytes(v *wrapperspb.BytesValue) *BytesValue {
	return (*BytesValue)(v)
}
