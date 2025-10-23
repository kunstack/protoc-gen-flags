package module

import (
	"fmt"
	"strings"

	"github.com/kunstack/protoc-gen-flags/flags"
	pgs "github.com/lyft/protoc-gen-star"
)

func (m *Module) genDuration(f pgs.Field, name pgs.Name, flag *flags.DurationFlag, wk pgs.WellKnownType) string {
	var declBuilder = &strings.Builder{}

	if flag.GetDisabled() {
		return fmt.Sprintf("// %s: flags disabled by disabled=true\n", name)
	}

	if flag.GetName() == "" {
		flag.Name = strings.ToLower(name.String())
	}

	_, _ = fmt.Fprintf(declBuilder, `
			if x.%s  == nil {
				x.%s = new(%s)
			}
		`,
		name, name, m.ctx.Type(f).Value(),
	)

	_, _ = fmt.Fprintf(declBuilder, `
			fs.VarP(types.Duration(x.%s), utils.BuildFlagName(prefix, %q), %q, %q)
		`,
		name, flag.Name, flag.GetShort(), flag.GetUsage(),
	)

	// 添加可选的 flag 配置
	_, _ = declBuilder.WriteString(m.genMark(flag))
	return declBuilder.String()
}

func (m *Module) genDurationSlice(f pgs.Field, name pgs.Name, flag *flags.RepeatedDurationFlag, wk pgs.WellKnownType) string {
	var declBuilder = &strings.Builder{}

	if flag.GetDisabled() {
		return fmt.Sprintf("// %s: flags disabled by disabled=true\n", name)
	}

	if flag.GetName() == "" {
		flag.Name = strings.ToLower(name.String())
	}

	_, _ = fmt.Fprintf(declBuilder, `
			fs.VarP(types.DurationSlice(&x.%s), utils.BuildFlagName(prefix, %q), %q, %q)
		`,
		name, flag.Name, flag.GetShort(), flag.GetUsage(),
	)

	// 添加可选的 flag 配置
	_, _ = declBuilder.WriteString(m.genMark(flag))
	return declBuilder.String()
}
