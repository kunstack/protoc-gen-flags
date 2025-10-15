package module

import (
	"fmt"
	"strings"

	pgs "github.com/lyft/protoc-gen-star"
)

func (m *Module) processPrimitiveFlag(f pgs.Field, name pgs.Name, flag commonFlag, wk pgs.WellKnownType, varFunc, varPFunc, typesFunc, _ string) string {
	if flag.GetDisabled() {
		return fmt.Sprint("\n// ", name, ": flags disabled by disabled=true\n")
	}

	var decl string

	flagName := flag.GetName()

	if flagName == "" {
		flagName = strings.ToLower(name.String())
	}

	if wk != "" && wk != pgs.UnknownWKT {
		decl += fmt.Sprintf("fs.VarP(types.%s(&x.%s), utils.BuildFlagName(prefix,%q), %q, %q)\n", typesFunc, name, flagName, flag.GetShort(), flag.GetUsage())
	} else if f.HasOptionalKeyword() {
		decl += fmt.Sprintf("fs.%s(x.%s, utils.BuildFlagName(prefix, %q), %q, *(x.%s), %q)\n", varFunc, name, flagName, flag.GetShort(), name, flag.GetUsage())
	} else {
		decl += fmt.Sprintf("fs.%s(&x.%s, utils.BuildFlagName(prefix, %q), %q, x.%s, %q)\n", varPFunc, name, flagName, flag.GetShort(), name, flag.GetUsage())
	}

	if flag.GetHidden() {
		decl += fmt.Sprintf("fs.MarkHidden(%q)\n", flagName)
	}
	if flag.GetDeprecated() {
		decl += fmt.Sprintf("fs.MarkDeprecated(%q, %q)\n", flagName, flag.GetDeprecatedUsage())
	}
	return decl
}

func (m *Module) processPrimitive(f pgs.Field, name pgs.Name, flag commonFlag, wk pgs.WellKnownType, wrapper, nativeWrapper string) string {
	if flag.GetDisabled() {
		return fmt.Sprint("\n// ", name, ": flags disabled by disabled=true\n")
	}

	flagName := flag.GetName()

	if flagName == "" {
		flagName = strings.ToLower(name.String())
	}

	if wk != "" && wk != pgs.UnknownWKT {
		decl += fmt.Sprintf("fs.VarP(types.%s(&x.%s), utils.BuildFlagName(prefix,%q), %q, %q)\n", typesFunc, name, flagName, flag.GetShort(), flag.GetUsage())
	} else if f.HasOptionalKeyword() {
		decl += fmt.Sprintf("fs.%s(x.%s, utils.BuildFlagName(prefix, %q), %q, *(x.%s), %q)\n", varFunc, name, flagName, flag.GetShort(), name, flag.GetUsage())
	} else {
		decl += fmt.Sprintf("fs.%s(&x.%s, utils.BuildFlagName(prefix, %q), %q, x.%s, %q)\n", varPFunc, name, flagName, flag.GetShort(), name, flag.GetUsage())
	}

	if flag.GetHidden() {
		decl += fmt.Sprintf("fs.MarkHidden(%q)\n", flagName)
	}
	if flag.GetDeprecated() {
		decl += fmt.Sprintf("fs.MarkDeprecated(%q, %q)\n", flagName, flag.GetDeprecatedUsage())
	}
	return decl
}

func (m *Module) processPrimitiveSlice(f pgs.Field, name pgs.Name, flag commonFlag, wk pgs.WellKnownType, varFunc, varPFunc, typesFunc string) string {
	if flag.GetDisabled() {
		return fmt.Sprint("\n// ", name, ": flags disabled by disabled=true\n")
	}

	var decl string

	flagName := flag.GetName()

	if flagName == "" {
		flagName = strings.ToLower(name.String())
	}

	if wk != "" && wk != pgs.UnknownWKT {
		decl += fmt.Sprintf("fs.VarP(types.%s(&x.%s), utils.BuildFlagName(prefix,%q), %q, %q)\n", typesFunc, name, flagName, flag.GetShort(), flag.GetUsage())
	} else if f.HasOptionalKeyword() {
		decl += fmt.Sprintf("fs.%s(x.%s, utils.BuildFlagName(prefix, %q), %q, *(x.%s), %q)\n", varFunc, name, flagName, flag.GetShort(), name, flag.GetUsage())
	} else {
		decl += fmt.Sprintf("fs.%s(&x.%s, utils.BuildFlagName(prefix, %q), %q, x.%s, %q)\n", varPFunc, name, flagName, flag.GetShort(), name, flag.GetUsage())
	}

	if flag.GetHidden() {
		decl += fmt.Sprintf("fs.MarkHidden(%q)\n", flagName)
	}
	if flag.GetDeprecated() {
		decl += fmt.Sprintf("fs.MarkDeprecated(%q, %q)\n", flagName, flag.GetDeprecatedUsage())
	}
	return decl
}
