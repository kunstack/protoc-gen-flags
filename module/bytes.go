package module

import (
	"fmt"
	"strings"

	"github.com/kunstack/protoc-gen-flags/flags"
	pgs "github.com/lyft/protoc-gen-star"
)

func (m *Module) genBytes(f pgs.Field, name pgs.Name, flag *flags.BytesFlag, wk pgs.WellKnownType) string {
	var (
		wrapper       = "Bytes"
		nativeWrapper = "BytesBase64VarP"
		declBuilder   = &strings.Builder{}
	)
	if flag.GetName() == "" {
		flag.Name = strings.ToLower(name.String())
	}
	if flag.GetDisabled() {
		return fmt.Sprintf("// %s: flags disabled by disabled=true\n", name)
	}
	if flag.GetEncoding() == flags.BytesEncodingType_BYTES_ENCODING_TYPE_HEX {
		wrapper = "BytesHex"
		nativeWrapper = "BytesHexVarP"
	}
	// 如果类型是 google.protobuf.Bytes
	if wk != "" && wk != pgs.UnknownWKT {
		_, _ = fmt.Fprintf(declBuilder, `
			if x.%s  == nil {
				x.%s = new(%s)
			}
		`,
			name, name, m.ctx.Type(f).Value(),
		)
		_, _ = fmt.Fprintf(declBuilder, `
			fs.VarP(types.%s(x.%s), utils.BuildFlagName(prefix,%q), %q, %q)
		`,
			wrapper, name, flag.GetName(), flag.GetShort(), flag.GetUsage(),
		)
	} else {
		_, _ = fmt.Fprintf(declBuilder, `
				fs.%s(&x.%s, utils.BuildFlagName(prefix,%q), %q, x.%s, %q)
			`,
			nativeWrapper, name, flag.GetName(), flag.GetShort(), name, flag.GetUsage())
	}
	_, _ = declBuilder.WriteString(m.genMark(flag))
	return declBuilder.String()
}

func (m *Module) genBytesSlice(f pgs.Field, name pgs.Name, flag *flags.RepeatedBytesFlag, wk pgs.WellKnownType) string {
	var (
		wrapper     = "BytesSlice"
		declBuilder = &strings.Builder{}
	)
	if flag.GetName() == "" {
		flag.Name = strings.ToLower(name.String())
	}

	if flag.GetDisabled() {
		return fmt.Sprintf("// %s: flags disabled by disabled=true\n", name)
	}
	if flag.GetEncoding() == flags.BytesEncodingType_BYTES_ENCODING_TYPE_HEX {
		wrapper = "BytesHexSlice"
	}
	_, _ = fmt.Fprintf(declBuilder, `
			fs.VarP(types.%s(&x.%s), utils.BuildFlagName(prefix,%q), %q, %q)
		`,
		wrapper, name, flag.GetName(), flag.GetShort(), flag.GetUsage())

	_, _ = declBuilder.WriteString(m.genMark(flag))
	return declBuilder.String()
}
