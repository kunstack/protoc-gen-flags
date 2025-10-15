package module

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kunstack/protoc-gen-flags/flags"
	pgs "github.com/lyft/protoc-gen-star"
	"github.com/samber/lo"
)


func (m *Module) checkTimestamp(ft FieldType, r *flags.TimestampFlag,  isSlice bool) {
	m.checkCommon(ft, r, pgs.MessageT,pgs.TimestampWKT, isSlice)

	if len(r.Formats) == 0 {
		m.Failf("at least one format must be specified for timestamp flag")
	}

	// Validate that formats are not empty strings
	for i, format := range r.Formats {
		if format == "" {
			m.Failf("timestamp format at index %d is empty", i)
		}
	}

	// Check for duplicate formats
	seenFormats := make(map[string]bool)
	for i, format := range r.Formats {
		if seenFormats[format] {
			m.Failf("timestamp format '%s' at index %d is duplicated", format, i)
		}
		seenFormats[format] = true
	}
}


func (m *Module) genTimestamp(f pgs.Field, name pgs.Name, flag *flags.TimestampFlag) string {
	var (
		declBuilder    = &strings.Builder{}
		formatsBuilder = &strings.Builder{}
	)
	if flag.GetDisabled() {
		return fmt.Sprint("\n// ", name, ": flags disabled by disabled=true\n")
	}
	if flag.GetName() == "" {
		flag.Name = strings.ToLower(name.String())
	}
	_, _ = fmt.Fprint(formatsBuilder,
		`[]string{`,
		strings.Join(
			lo.Map(
				flag.GetFormats(),
				func(item string, _ int) string {
					return strconv.Quote(item)
				},
			),
			",",
		),
		`}`,
	)

	_, _ = fmt.Fprintf(declBuilder, `
			if x.%s  == nil {
				x.%s = new(%s)
			}
		`,
		name, name, m.ctx.Type(f).Value(),
	)

	_, _ = fmt.Fprintf(declBuilder, `
		fs.VarP(types.Timestamp(x.%s, %s), utils.BuildFlagName(prefix,%q), %q, %q)
	`,
		name, formatsBuilder.String(), flag.GetName(), flag.GetShort(), flag.GetUsage(),
	)

	_, _ = declBuilder.WriteString(m.genMark(flag))
	return declBuilder.String()
}

func (m *Module) genTimestampSlice(f pgs.Field, name pgs.Name, flag *flags.TimestampFlag) string {
	var (
		declBuilder    = &strings.Builder{}
		formatsBuilder = &strings.Builder{}
	)
	if flag.GetDisabled() {
		return fmt.Sprint("\n// ", name, ": flags disabled by disabled=true\n")
	}
	if flag.GetName() == "" {
		flag.Name = strings.ToLower(name.String())
	}

	_, _ = fmt.Fprint(formatsBuilder,
		`[]string{`,
		strings.Join(
			lo.Map(
				flag.GetFormats(),
				func(item string, _ int) string {
					return strconv.Quote(item)
				},
			),
			",",
		),
		`}`,
	)

	_, _ = fmt.Fprintf(declBuilder, `
		fs.VarP(types.Timestamp(x.%s, %s), utils.BuildFlagName(prefix,%q), %q, %q)
	`,
		name, formatsBuilder.String(), flag.GetName(), flag.GetShort(), flag.GetUsage(),
	)

	_, _ = declBuilder.WriteString(m.genMark(flag))
	return declBuilder.String()
}
