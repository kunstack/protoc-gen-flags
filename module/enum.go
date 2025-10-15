package module

import (
	"fmt"
	"strings"

	"github.com/kunstack/protoc-gen-flags/flags"
	pgs "github.com/lyft/protoc-gen-star"
)


func (m *Module) checkEnum(ft FieldType, flag *flags.EnumFlag, pt pgs.ProtoType, wrapper pgs.WellKnownType, isSlice bool) {
	m.checkCommon(ft, flag, pt, wrapper, isSlice)

	_, ok := ft.(interface {
		Enum() pgs.Enum
	})

	if !ok {
		m.Failf("unexpected field type (%T)", ft)
	}
}

func (m *Module) processEnumFlag(f pgs.Field, name pgs.Name, flag *flags.EnumFlag, wk pgs.WellKnownType) string {
	if flag.GetDisabled() {
		return fmt.Sprint("\n// ", name, ": flags disabled by disabled=true\n")
	}

	var decl string

	if flag.GetName() == "" {
		flag.Name = strings.ToLower(name.String())
	}

	// Generate flag binding code for enum
	if f.HasOptionalKeyword() {
		// For optional fields, use direct value (no pointer)
		decl += fmt.Sprintf("fs.VarP(types.Enum(x.%s), utils.BuildFlagName(prefix, %q), %q, %q)\n", name, flag.GetName(), flag.GetShort(), flag.GetUsage())
	} else {
		// For regular fields, use a pointer to access the Type() method
		decl += fmt.Sprintf("fs.VarP(types.Enum(&x.%s), utils.BuildFlagName(prefix, %q), %q, %q)\n", name, flag.GetName(), flag.GetShort(), flag.GetUsage())
	}

	if flag.GetHidden() {
		decl += fmt.Sprintf("fs.MarkHidden(%q)\n", flag.GetName())
	}
	if flag.GetDeprecated() {
		decl += fmt.Sprintf("fs.MarkDeprecated(%q, %q)\n", flag.GetName(), flag.GetDeprecatedUsage())
	}

	return decl
}
