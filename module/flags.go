// Copyright 2021 Linka Cloud  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package module

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kunstack/protoc-gen-flags/flags"
	pgs "github.com/lyft/protoc-gen-star"
)

func isOk(v []bool) bool {
	return len(v) > 0 && v[0]
}

func (m *Module) processPrimitiveFlag(f pgs.Field, name pgs.Name, flag *flags.PrimitiveFlag, wk pgs.WellKnownType, varFunc, varPFunc, typesFunc, fieldType string) (string, bool) {
	if flag.Disabled {
		return fmt.Sprint("\n// ", name, ": flags disabled by disabled=true"), true
	}

	var decl string

	if flag.Name == "" {
		flag.Name = strings.ToLower(name.String())
	}

	decl += fmt.Sprintf(`
		// %s flag generated for [(flags.value).%s = {
		//     disabled: false,
		//     name: %q,
		//     short: %q,
		//     usage: %q,
		//     hidden: %v,
		//     deprecated: %v,
		//     deprecated_usage: %q,
		// }]
		`, name, fieldType, flag.Name, flag.Short, flag.Usage, flag.Hidden, flag.Deprecated, flag.DeprecatedUsage)

	if wk != "" && wk != pgs.UnknownWKT {
		decl += fmt.Sprintf("fs.VarP(types.%s(x.%s), utils.BuildFlagName(prefix,%q), %q, %q)\n", typesFunc, name, flag.Name, flag.Short, flag.Usage)
	} else if f.HasOptionalKeyword() {
		decl += fmt.Sprintf("fs.%s(x.%s, utils.BuildFlagName(prefix, %q), %q, *(x.%s), %q)\n", varFunc, name, flag.Name, flag.Short, name, flag.Usage)
	} else {
		decl += fmt.Sprintf("fs.%s(&x.%s, utils.BuildFlagName(prefix, %q), %q, x.%s, %q)\n", varPFunc, name, flag.Name, flag.Short, name, flag.Usage)
	}

	if flag.Hidden {
		decl += fmt.Sprintf("fs.MarkHidden(%q)\n", flag.Name)
	}
	if flag.Deprecated {
		decl += fmt.Sprintf("fs.MarkDeprecated(%q, %q)\n", flag.Name, flag.DeprecatedUsage)
	}
	return decl, true
}

func (m *Module) processBytesFlag(f pgs.Field, name pgs.Name, flag *flags.BytesFlag, wk pgs.WellKnownType) (string, bool) {
	var (
		decl      string
		varFunc   = "BytesBase64VarP"
		typesFunc = "Bytes"
	)

	if flag.Name == "" {
		flag.Name = strings.ToLower(name.String())
	}

	if flag.Disabled {
		return fmt.Sprintf("\n// %s: flags disabled by disabled=true", name), true
	}

	if flag.GetEncoding() == flags.BytesEncodingType_BYTES_ENCODING_TYPE_HEX {
		typesFunc = "BytesHex"
		varFunc = "BytesHexVarP"
	}

	decl += fmt.Sprintf(`
		// %s flag generated for (flags.value).bytes = {
		//  disabled: false,
		//     name: %q,
		//     short: %q,
		//     usage: %q,
		//     hidden: %v,
		//     deprecated: %v,
		//     deprecated_usage: %q,
		//     encoding: %s
		// }
		`, name, flag.Name, flag.Short, flag.Usage, flag.Hidden, flag.Deprecated, flag.DeprecatedUsage, flag.Encoding)

	if wk != "" && wk != pgs.UnknownWKT {
		decl += fmt.Sprintf("fs.VarP(types.%s(x.%s), utils.BuildFlagName(prefix,%q), %q, %q)\n", typesFunc, name, flag.Name, flag.Short, flag.Usage)
	} else {
		decl += fmt.Sprintf("fs.%s(&x.%s, utils.BuildFlagName(prefix, %q), %q, x.%s, %q)\n", varFunc, name, flag.Name, flag.Short, name, flag.Usage)
	}

	if flag.Hidden {
		decl += fmt.Sprintf("fs.MarkHidden(%q)\n", flag.Name)
	}
	if flag.Deprecated {
		decl += fmt.Sprintf("fs.MarkDeprecated(%q, %q)\n", flag.Name, flag.DeprecatedUsage)
	}

	return decl, true
}

func (m *Module) processEnumFlag(f pgs.Field, name pgs.Name, flag *flags.PrimitiveFlag, wk pgs.WellKnownType, genOneOfField ...bool) (string, bool) {
	if flag.Disabled {
		return fmt.Sprint("\n// ", name, ": flags disabled by disabled=true"), true
	}

	var decl string

	if flag.Name == "" {
		flag.Name = strings.ToLower(name.String())
	}

	decl += fmt.Sprintf(`
		// %s flag generated for [(flags.value).enum = {
		//     disabled: false,
		//     name: %q,
		//     short: %q,
		//     usage: %q,
		//     hidden: %v,
		//     deprecated: %v,
		//     deprecated_usage: %q,
		// }]
		`, name.String(), flag.Name, flag.Short, flag.Usage, flag.Hidden, flag.Deprecated, flag.DeprecatedUsage)

	// Generate flag binding code for enum
	if isOk(genOneOfField) || f.HasOptionalKeyword() {
		// For oneof fields or optional fields, use direct value (no pointer)
		decl += fmt.Sprintf("fs.VarP(types.Enum(x.%s), utils.BuildFlagName(prefix, %q), %q, %q)\n", name, flag.Name, flag.Short, flag.Usage)
	} else {
		// For regular fields, use a pointer to access the Type() method
		decl += fmt.Sprintf("fs.VarP(types.Enum(&x.%s), utils.BuildFlagName(prefix, %q), %q, %q)\n", name, flag.Name, flag.Short, flag.Usage)
	}

	if flag.Hidden {
		decl += fmt.Sprintf("fs.MarkHidden(%q)\n", flag.Name)
	}
	if flag.Deprecated {
		decl += fmt.Sprintf("fs.MarkDeprecated(%q, %q)\n", flag.Name, flag.DeprecatedUsage)
	}

	return decl, true
}

func (m *Module) processDurationFlag(f pgs.Field, name pgs.Name, flag *flags.PrimitiveFlag, wk pgs.WellKnownType, genOneOfField ...bool) (string, bool) {
	if flag.Disabled {
		return fmt.Sprint("\n// ", name, ": flags disabled by disabled=true"), true
	}

	var decl string

	if flag.Name == "" {
		flag.Name = strings.ToLower(name.String())
	}

	decl += fmt.Sprintf(`
		// %s flag generated for [(flags.value).duration = {
		//     disabled: false,
		//     name: %q,
		//     short: %q,
		//     usage: %q,
		//     hidden: %v,
		//     deprecated: %v,
		//     deprecated_usage: %q,
		// }]
		`, name.String(), flag.Name, flag.Short, flag.Usage, flag.Hidden, flag.Deprecated, flag.DeprecatedUsage)

	// For oneof fields, use direct value (no pointer)
	decl += fmt.Sprintf("fs.VarP(types.Duration(x.%s), utils.BuildFlagName(prefix, %q), %q, %q)\n", name, flag.Name, flag.Short, flag.Usage)

	if flag.Hidden {
		decl += fmt.Sprintf("fs.MarkHidden(%q)\n", flag.Name)
	}
	if flag.Deprecated {
		decl += fmt.Sprintf("fs.MarkDeprecated(%q, %q)\n", flag.Name, flag.DeprecatedUsage)
	}

	return decl, true
}

func (m *Module) genFieldFlags(f pgs.Field, genOneOfField ...bool) (string, bool) {
	m.Push(f.Name().String())
	defer m.Pop()
	var field flags.FieldFlags
	ok, err := f.Extension(flags.E_Value, &field)
	if err != nil || !ok {
		return "", false
	}
	wk := pgs.UnknownWKT
	if emb := f.Type().Embed(); emb != nil {
		wk = emb.WellKnownType()
	}
	if !isOk(genOneOfField) && f.InRealOneOf() {
		if m.isOneOfDone(f.OneOf()) {
			return "", false
		}
		m.setOneOfDone(f.OneOf())
		var out string
		out += fmt.Sprint(`
			switch x := x.`, m.ctx.Name(f.OneOf()), `.(type) {`)
		for _, f := range f.OneOf().Fields() {
			decl, ok := m.genFieldFlags(f, true)
			if !ok {
				continue
			}
			out += fmt.Sprint(`
				case *`, m.ctx.OneofOption(f), `: `, decl)
		}
		out += `}`
		return out, true
	}
	name := m.ctx.Name(f)
	switch r := field.Type.(type) {
	case *flags.FieldFlags_Float:
		return m.processPrimitiveFlag(f, name, r.Float, wk, "Float32VarP", "Float32VarP", "Float", "float")
	case *flags.FieldFlags_Double:
		return m.processPrimitiveFlag(f, name, r.Double, wk, "Float64VarP", "Float64VarP", "Double", "double")
	case *flags.FieldFlags_Int32:
		return m.processPrimitiveFlag(f, name, r.Int32, wk, "Int32VarP", "Int32VarP", "Int32", "int32")
	case *flags.FieldFlags_Int64:
		return m.processPrimitiveFlag(f, name, r.Int64, wk, "Int64VarP", "Int64VarP", "Int64", "int64")
	case *flags.FieldFlags_Uint32:
		return m.processPrimitiveFlag(f, name, r.Uint32, wk, "Uint32VarP", "Uint32VarP", "UInt32", "uint32")
	case *flags.FieldFlags_Uint64:
		return m.processPrimitiveFlag(f, name, r.Uint64, wk, "Uint64VarP", "Uint64VarP", "UInt64", "uint64")
	case *flags.FieldFlags_Sint32:
		return m.processPrimitiveFlag(f, name, r.Sint32, wk, "Int32VarP", "Int32VarP", "Int32", "sint32")
	case *flags.FieldFlags_Sint64:
		return m.processPrimitiveFlag(f, name, r.Sint64, wk, "Int64VarP", "Int64VarP", "Int64", "sint64")
	case *flags.FieldFlags_Fixed32:
		return m.processPrimitiveFlag(f, name, r.Fixed32, wk, "Uint32VarP", "Uint32VarP", "UInt32", "fixed32")
	case *flags.FieldFlags_Fixed64:
		return m.processPrimitiveFlag(f, name, r.Fixed64, wk, "Uint64VarP", "Uint64VarP", "UInt32", "fixed64")
	case *flags.FieldFlags_Sfixed32:
		return m.processPrimitiveFlag(f, name, r.Sfixed32, wk, "Int32VarP", "Int32VarP", "Int32", "sfixed32")
	case *flags.FieldFlags_Sfixed64:
		return m.processPrimitiveFlag(f, name, r.Sfixed64, wk, "Int64VarP", "Int64VarP", "Int64", "sfixed64")
	case *flags.FieldFlags_Bool:
		return m.processPrimitiveFlag(f, name, r.Bool, wk, "BoolVarP", "BoolVarP", "Bool", "bool")
	case *flags.FieldFlags_String_:
		return m.processPrimitiveFlag(f, name, r.String_, wk, "StringVarP", "StringVarP", "String", "string")
	case *flags.FieldFlags_Bytes:
		return m.processBytesFlag(f, name, r.Bytes, wk)
	case *flags.FieldFlags_Enum:
		return m.processEnumFlag(f, name, r.Enum, wk, genOneOfField...)
	case *flags.FieldFlags_Duration:
		return m.processDurationFlag(f, name, r.Duration, wk, genOneOfField...)
	case *flags.FieldFlags_Message:
		if field.GetMessage() != nil && !field.GetMessage().Nested {
			return fmt.Sprint("\n// ", name, ": flags disabled by [(flags.value).message = {nested: false}]"), true
		}
		prefix := field.GetMessage().Name
		if prefix == "" {
			// use field name instead, convert to kebab-case
			prefix = strings.ToLower(f.Name().String())
		}
		var decl string
		return decl + fmt.Sprint(`
			if v, ok := interface{}(x.`, name, `).(flags.Interface); ok && x.`, name, ` != nil {
				v.AddFlags(fs,`, strconv.Quote(prefix), ` )
			}`), true
	case nil: // noop
	default:
		_ = r
		//m.Failf("unknown rule type (%T)", field.Type)

	}
	return fmt.Sprint("\n// ", f.Name()), true
}
