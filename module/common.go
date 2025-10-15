package module

import (
	"fmt"
	"strings"

	pgs "github.com/lyft/protoc-gen-star"
)

func (m *Module) checkCommon(typ FieldType, r commonFlag, pt pgs.ProtoType, wrapper pgs.WellKnownType, isSlice bool) {
	m.mustType(typ, pt, wrapper)

	if typ, ok := typ.(Repeatable); ok {
		if isSlice && !typ.IsRepeated(){
			m.Fail("repeated fields should use repeated flag")
		}

		if !isSlice && typ.IsRepeated(){
			m.Fail("repeated flag should be used for repeated fields")
		}
	}

	if r.GetUsage() == "" {
		m.Failf("usage is required for flag")
	}
	// Check if deprecated flag has proper deprecation usage message
	if r.GetDeprecated() && r.GetDeprecatedUsage() == "" {
		m.Failf("deprecated flag must provide deprecated_usage message")
	}
}



func (m *Module) genCommon(f pgs.Field, name pgs.Name, flag commonFlag, wk pgs.WellKnownType, wrapper, nativeWrapper string) string {
	var (
		declBuilder = &strings.Builder{}
	)
	if flag.GetDisabled() {
		return fmt.Sprint("\n// ", name, ": flags disabled by disabled=true\n")
	}
	flagName := flag.GetName()
	if flagName == "" {
		flagName = strings.ToLower(name.String())
	}
	if wk != "" && wk != pgs.UnknownWKT {
		_, _ = fmt.Fprintf(declBuilder, `
				if x.%s == nil {
					x.%s = new(%s)
				}`,
			name, name, m.ctx.Type(f).Value(),
		)
		_, _ = fmt.Fprintf(declBuilder, `
				fs.VarP(types.%s(x.%s), utils.BuildFlagName(prefix,%q), %q, %q)
			`,
			wrapper, name, flagName, flag.GetShort(), flag.GetUsage())
	} else if f.HasOptionalKeyword() {
		_, _ = fmt.Fprintf(declBuilder, `
				if x.%s == nil {
					x.%s = new(%s)
				}`,
			name, name, m.ctx.Type(f).Value(),
		)
		_, _ = fmt.Fprintf(declBuilder, `
				fs.%s(x.%s, utils.BuildFlagName(prefix, %q), %q, *(x.%s), %q)
			`,
			nativeWrapper, name, flagName, flag.GetShort(), name, flag.GetUsage())
	} else {
		_, _ = fmt.Fprintf(declBuilder, `
				fs.%s(&x.%s, utils.BuildFlagName(prefix, %q), %q, x.%s, %q)
			`,
			nativeWrapper, name, flagName, flag.GetShort(), name, flag.GetUsage())
	}
	_, _ = declBuilder.WriteString(m.genMark(flag))
	return declBuilder.String()
}

func (m *Module) genCommonSlice(f pgs.Field, name pgs.Name, flag commonFlag, wk pgs.WellKnownType, wrapper, nativeWrapper string) string {
	var (
		declBuilder = &strings.Builder{}
	)

	if flag.GetDisabled() {
		return fmt.Sprint("\n// ", name, ": flags disabled by disabled=true\n")
	}

	flagName := flag.GetName()

	if flagName == "" {
		flagName = strings.ToLower(name.String())
	}

	if wk != "" && wk != pgs.UnknownWKT {
		_, _ = fmt.Fprintf(declBuilder, `
				fs.VarP(types.%s(&x.%s), utils.BuildFlagName(prefix,%q), %q, %q)
			`,
			wrapper, name, flagName, flag.GetShort(), flag.GetUsage())
	} else {
		_, _ = fmt.Fprintf(declBuilder, `
				fs.%s(&x.%s, utils.BuildFlagName(prefix, %q), %q, x.%s, %q)
			`,
			nativeWrapper, name, flagName, flag.GetShort(), name, flag.GetUsage())
	}

	_, _ = declBuilder.WriteString(m.genMark(flag))
	return declBuilder.String()
}
