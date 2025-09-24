package types

import (
	"encoding/hex"
	"strings"

	"github.com/spf13/pflag"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var _ pflag.Value = (*BytesHexValue)(nil)

type BytesHexValue wrapperspb.BytesValue

func (b *BytesHexValue) String() string {
	if b == nil {
		return ""
	}
	return strings.ToUpper(hex.EncodeToString(b.Value))
}

func (b *BytesHexValue) Set(value string) error {
	bin, err := hex.DecodeString(strings.TrimSpace(value))
	if err != nil {
		return err
	}
	b.Value = bin
	return nil
}

func (b *BytesHexValue) Type() string {
	return "bytesHex"
}

func BytesHex(v *wrapperspb.BytesValue) *BytesHexValue {
	return (*BytesHexValue)(v)
}
