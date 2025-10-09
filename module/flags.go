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

func (m *Module) processPrimitiveFlag(f pgs.Field, name pgs.Name, flag commonFlag, wk pgs.WellKnownType, varFunc, varPFunc, typesFunc, fieldType string) (string, bool) {
	if flag.GetDisabled() {
		return fmt.Sprint("\n// ", name, ": flags disabled by disabled=true"), true
	}

	var decl string

	flagName := flag.GetName()

	if flagName == "" {
		flagName = strings.ToLower(name.String())
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
		`, name, fieldType, flagName, flag.GetShort(), flag.GetUsage(), flag.GetHidden(), flag.GetDeprecated(), flag.GetDeprecatedUsage())

	if wk != "" && wk != pgs.UnknownWKT {
		decl += fmt.Sprintf("fs.VarP(types.%s(x.%s), utils.BuildFlagName(prefix,%q), %q, %q)\n", typesFunc, name, flagName, flag.GetShort(), flag.GetUsage())
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
	return decl, true
}

func (m *Module) processBytesFlag(_ pgs.Field, name pgs.Name, flag *flags.BytesFlag, wk pgs.WellKnownType, isRepeatedFlag bool) (string, bool) {
	var (
		decl      string
		typesFunc = "Bytes"
		varFunc   = "BytesBase64VarP"
	)

	if isRepeatedFlag {
		varFunc = ""
		typesFunc = "BytesSlice"
	}

	if flag.GetName() == "" {
		flag.Name = strings.ToLower(name.String())
	}

	if flag.GetDisabled() {
		return fmt.Sprintf("\n// %s: flags disabled by disabled=true", name), true
	}

	if flag.GetEncoding() == flags.BytesEncodingType_BYTES_ENCODING_TYPE_HEX {
		if isRepeatedFlag {
			varFunc = ""
			typesFunc = "BytesHexSlice"
		} else {
			typesFunc = "BytesHex"
			varFunc = "BytesHexVarP"
		}
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
		`, name, flag.GetName(), flag.GetShort(), flag.GetUsage(), flag.GetHidden(), flag.GetDeprecated(), flag.GetDeprecatedUsage(), flag.GetEncoding())

	if (wk != "" && wk != pgs.UnknownWKT) || isRepeatedFlag {
		decl += fmt.Sprintf("fs.VarP(types.%s(x.%s), utils.BuildFlagName(prefix,%q), %q, %q)\n", typesFunc, name, flag.GetName(), flag.GetShort(), flag.GetUsage())
	} else {
		decl += fmt.Sprintf("fs.%s(&x.%s, utils.BuildFlagName(prefix, %q), %q, x.%s, %q)\n", varFunc, name, flag.GetName(), flag.GetShort(), name, flag.GetUsage())
	}

	if flag.GetHidden() {
		decl += fmt.Sprintf("fs.MarkHidden(%q)\n", flag.GetName())
	}
	if flag.GetDeprecated() {
		decl += fmt.Sprintf("fs.MarkDeprecated(%q, %q)\n", flag.GetName(), flag.GetDeprecatedUsage())
	}

	return decl, true
}

func (m *Module) processEnumFlag(f pgs.Field, name pgs.Name, flag *flags.EnumFlag, wk pgs.WellKnownType, genOneOfField ...bool) (string, bool) {
	if flag.GetDisabled() {
		return fmt.Sprint("\n// ", name, ": flags disabled by disabled=true"), true
	}

	var decl string

	if flag.GetName() == "" {
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
		`, name.String(), flag.GetName(), flag.GetShort(), flag.GetUsage(), flag.GetHidden(), flag.GetDeprecated(), flag.GetDeprecatedUsage())

	// Generate flag binding code for enum
	if isOk(genOneOfField) || f.HasOptionalKeyword() {
		// For oneof fields or optional fields, use direct value (no pointer)
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

	return decl, true
}

func (m *Module) processDurationFlag(f pgs.Field, name pgs.Name, flag *flags.DurationFlag, wk pgs.WellKnownType, isRepeatedFlag bool) (string, bool) {
	if flag.GetDisabled() {
		return fmt.Sprint("\n// ", name, ": flags disabled by disabled=true"), true
	}

	var decl string

	if flag.GetName() == "" {
		flag.Name = strings.ToLower(name.String())
	}

	fieldType := "duration"
	if isRepeatedFlag {
		fieldType = "repeated.duration"
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
		`, name.String(), fieldType, flag.GetName(), flag.GetShort(), flag.GetUsage(), flag.GetHidden(), flag.GetDeprecated(), flag.GetDeprecatedUsage())

	if isRepeatedFlag {
		// For oneof fields, use direct value (no pointer)
		decl += fmt.Sprintf("fs.VarP(types.DurationSlice(x.%s), utils.BuildFlagName(prefix, %q), %q, %q)\n", name, flag.GetName(), flag.GetShort(), flag.GetUsage())
	} else {
		// For oneof fields, use direct value (no pointer)
		decl += fmt.Sprintf("fs.VarP(types.Duration(x.%s), utils.BuildFlagName(prefix, %q), %q, %q)\n", name, flag.GetName(), flag.GetShort(), flag.GetUsage())
	}

	if flag.GetHidden() {
		decl += fmt.Sprintf("fs.MarkHidden(%q)\n", flag.GetName())
	}
	if flag.GetDeprecated() {
		decl += fmt.Sprintf("fs.MarkDeprecated(%q, %q)\n", flag.GetName(), flag.GetDeprecatedUsage())
	}

	return decl, true
}

func (m *Module) processRepeatedFlag(f pgs.Field, name pgs.Name, repeated *flags.RepeatedFlags) (string, bool) {
	if repeated == nil {
		return "", false
	}
	wk := pgs.UnknownWKT
	if emb := f.Type().Element().Embed(); emb != nil {
		wk = emb.WellKnownType()
	}

	switch r := repeated.Type.(type) {
	case *flags.RepeatedFlags_Float:
		return m.processPrimitiveFlag(f, name, r.Float, wk, "Float32SliceVarP", "Float32SliceVarP", "Float32Slice", "float")
	case *flags.RepeatedFlags_Double:
		return m.processPrimitiveFlag(f, name, r.Double, wk, "Float64SliceVarP", "Float64SliceVarP", "DoubleSlice", "double")
	case *flags.RepeatedFlags_Int32:
		return m.processPrimitiveFlag(f, name, r.Int32, wk, "Int32SliceVarP", "Int32SliceVarP", "Int32Slice", "int32")
	case *flags.RepeatedFlags_Int64:
		return m.processPrimitiveFlag(f, name, r.Int64, wk, "Int64SliceVarP", "Int64SliceVarP", "Int64Slice", "int64")
	case *flags.RepeatedFlags_Uint32:
		return m.processPrimitiveFlag(f, name, r.Uint32, wk, "Uint32SliceVarP", "Uint32SliceVarP", "UInt32Slice", "uint32")
	case *flags.RepeatedFlags_Uint64:
		return m.processPrimitiveFlag(f, name, r.Uint64, wk, "Uint64SliceVarP", "Uint64SliceVarP", "UInt64Slice", "uint64")
	case *flags.RepeatedFlags_Sint32:
		return m.processPrimitiveFlag(f, name, r.Sint32, wk, "Int32SliceVarP", "Int32SliceVarP", "Int32Slice", "sint32")
	case *flags.RepeatedFlags_Sint64:
		return m.processPrimitiveFlag(f, name, r.Sint64, wk, "Int64SliceVarP", "Int64SliceVarP", "Int64Slice", "sint64")
	case *flags.RepeatedFlags_Fixed32:
		return m.processPrimitiveFlag(f, name, r.Fixed32, wk, "Uint32SliceVarP", "Uint32SliceVarP", "UInt32Slice", "fixed32")
	case *flags.RepeatedFlags_Fixed64:
		return m.processPrimitiveFlag(f, name, r.Fixed64, wk, "Uint64SliceVarP", "Uint64SliceVarP", "UInt32Slice", "fixed64")
	case *flags.RepeatedFlags_Sfixed32:
		return m.processPrimitiveFlag(f, name, r.Sfixed32, wk, "Int32SliceVarP", "Int32SliceVarP", "Int32Slice", "sfixed32")
	case *flags.RepeatedFlags_Sfixed64:
		return m.processPrimitiveFlag(f, name, r.Sfixed64, wk, "Int64SliceVarP", "Int64SliceVarP", "Int64Slice", "sfixed64")
	case *flags.RepeatedFlags_Bool:
		return m.processPrimitiveFlag(f, name, r.Bool, wk, "BoolSliceVarP", "BoolSliceVarP", "BoolSlice", "bool")
	case *flags.RepeatedFlags_String_:
		return m.processPrimitiveFlag(f, name, r.String_, wk, "StringSliceVarP", "StringSliceVarP", "StringSlice", "string")
	case *flags.RepeatedFlags_Bytes:
		return m.processBytesFlag(f, name, r.Bytes, wk, true)
	case *flags.RepeatedFlags_Enum:
		return m.processEnumFlag(f, name, r.Enum, wk)
	case *flags.RepeatedFlags_Duration:
		return m.processDurationFlag(f, name, r.Duration, wk, true)
	case *flags.RepeatedFlags_Timestamp:
		return m.processTimestampFlag(f, name, r.Timestamp, wk)
	default:
		return "", false
	}
}

func (m *Module) processTimestampFlag(_ pgs.Field, name pgs.Name, flag *flags.TimestampFlag, wk pgs.WellKnownType) (string, bool) {
	if flag.GetDisabled() {
		return fmt.Sprint("\n// ", name, ": flags disabled by disabled=true"), true
	}

	var decl string

	if flag.GetName() == "" {
		flag.Name = strings.ToLower(name.String())
	}

	decl += fmt.Sprintf(`
		// %s flag generated for [(flags.value).timestamp = {
		//     disabled: false,
		//     name: %q,
		//     short: %q,
		//     usage: %q,
		//     hidden: %v,
		//     deprecated: %v,
		//     deprecated_usage: %q,
		// }]
		`, name.String(), flag.GetName(), flag.GetShort(), flag.GetUsage(), flag.GetHidden(), flag.GetDeprecated(), flag.GetDeprecatedUsage())

	// Generate flag binding code for timestamp
	if wk != "" && wk != pgs.UnknownWKT {
		decl += fmt.Sprintf("fs.VarP(types.Timestamp(x.%s), utils.BuildFlagName(prefix,%q), %q, %q)\n", name, flag.GetName(), flag.GetShort(), flag.GetUsage())
	} else {
		decl += fmt.Sprintf("fs.VarP(types.Timestamp(&x.%s), utils.BuildFlagName(prefix, %q), %q, %q)\n", name, flag.GetName(), flag.GetShort(), flag.GetUsage())
	}

	if flag.GetHidden() {
		decl += fmt.Sprintf("fs.MarkHidden(%q)\n", flag.GetName())
	}
	if flag.GetDeprecated() {
		decl += fmt.Sprintf("fs.MarkDeprecated(%q, %q)\n", flag.GetName(), flag.GetDeprecatedUsage())
	}

	return decl, true
}

func (m *Module) processMapFlag(f pgs.Field, name pgs.Name, flag *flags.MapFlag, wk pgs.WellKnownType) (string, bool) {
	if flag.GetDisabled() {
		return fmt.Sprint("\n// ", name, ": flags disabled by disabled=true"), true
	}

	var decl string

	if flag.GetName() == "" {
		flag.Name = strings.ToLower(name.String())
	}

	// Add format information to the comment
	decl += fmt.Sprintf(`
		// %s flag generated for [(flags.value).map = {
		//     disabled: false,
		//     name: %q,
		//     short: %q,
		//     usage: %q,
		//     hidden: %v,
		//     deprecated: %v,
		//     deprecated_usage: %q,
		//     format: %s
		// }]
		`, name.String(), flag.GetName(), flag.GetShort(), flag.GetUsage(), flag.GetHidden(), flag.GetDeprecated(), flag.GetDeprecatedUsage(), flag.GetFormat())

	// Determine the format to use
	mapFormat := flag.GetFormat()

	// If unspecified, default to JSON format for backward compatibility
	if mapFormat == flags.MapFormatType_MAP_FORMAT_TYPE_UNSPECIFIED {
		mapFormat = flags.MapFormatType_MAP_FORMAT_TYPE_JSON
	}

	// Validate the field type matches the specified format
	if !f.Type().IsMap() {
		m.Logf("Warning: field %s is not a map type but map format is specified", name)
	} else {
		keyType := f.Type().Key()
		valueType := f.Type().Element()

		switch mapFormat {
		case flags.MapFormatType_MAP_FORMAT_TYPE_STRING_TO_STRING:
			if keyType.ProtoType() != pgs.StringT {
				m.Logf("Warning: field %s key type is not string for STRING_TO_STRING format", name)
			}
			if valueType.ProtoType() != pgs.StringT {
				m.Logf("Warning: field %s value type is not string for STRING_TO_STRING format", name)
			}
		case flags.MapFormatType_MAP_FORMAT_TYPE_STRING_TO_INT:
			if keyType.ProtoType() != pgs.StringT {
				m.Logf("Warning: field %s key type is not string for STRING_TO_INT format", name)
			}
			// Support all integer types for STRING_TO_INT format
			switch valueType.ProtoType() {
			case pgs.Int32T, pgs.SInt32, pgs.SFixed32,
				pgs.Int64T, pgs.SInt64, pgs.SFixed64,
				pgs.UInt32T, pgs.Fixed32T,
				pgs.UInt64T, pgs.Fixed64T:
				// These are all valid integer types
			default:
				m.Logf("Warning: field %s value type is not a valid integer type for STRING_TO_INT format", name)
			}
		}
	}

	// Generate flag binding based on format
	switch mapFormat {
	case flags.MapFormatType_MAP_FORMAT_TYPE_STRING_TO_STRING:
		// For string-to-string maps, use native pflag StringToStringVarP
		if f.HasOptionalKeyword() {
			decl += fmt.Sprintf("fs.StringToStringVarP(x.%s, utils.BuildFlagName(prefix, %q), %q, x.%s, %q)\n", name, flag.GetName(), flag.GetShort(), name, flag.GetUsage())
		} else {
			decl += fmt.Sprintf("fs.StringToStringVarP(&x.%s, utils.BuildFlagName(prefix, %q), %q, x.%s, %q)\n", name, flag.GetName(), flag.GetShort(), name, flag.GetUsage())
		}

	case flags.MapFormatType_MAP_FORMAT_TYPE_STRING_TO_INT:
		// For string-to-int maps, determine the specific int type based on the field
		valueType := f.Type().Element()

		switch valueType.ProtoType() {
		case pgs.Int64T, pgs.SInt64, pgs.SFixed64:
			// Use native pflag StringToInt64 for int64 types
			if f.HasOptionalKeyword() {
				decl += fmt.Sprintf("fs.StringToInt64VarP(x.%s, utils.BuildFlagName(prefix, %q), %q, x.%s, %q)\n", name, flag.GetName(), flag.GetShort(), name, flag.GetUsage())
			} else {
				decl += fmt.Sprintf("fs.StringToInt64VarP(&x.%s, utils.BuildFlagName(prefix, %q), %q, x.%s, %q)\n", name, flag.GetName(), flag.GetShort(), name, flag.GetUsage())
			}
		case pgs.Int32T, pgs.SInt32, pgs.SFixed32:
			// Use custom type for int32 since pflag's StringToIntVarP expects *map[string]int, not *map[string]int32
			if f.HasOptionalKeyword() {
				decl += fmt.Sprintf("fs.VarP(types.StringToInt32(x.%s), utils.BuildFlagName(prefix, %q), %q, %q)\n", name, flag.GetName(), flag.GetShort(), flag.GetUsage())
			} else {
				decl += fmt.Sprintf("fs.VarP(types.StringToInt32(&x.%s), utils.BuildFlagName(prefix, %q), %q, %q)\n", name, flag.GetName(), flag.GetShort(), flag.GetUsage())
			}
		case pgs.UInt32T, pgs.Fixed32T:
			// Use custom type for uint32 since pflag doesn't have native support
			if f.HasOptionalKeyword() {
				decl += fmt.Sprintf("fs.VarP(types.StringToUint32(x.%s), utils.BuildFlagName(prefix, %q), %q, %q)\n", name, flag.GetName(), flag.GetShort(), flag.GetUsage())
			} else {
				decl += fmt.Sprintf("fs.VarP(types.StringToUint32(&x.%s), utils.BuildFlagName(prefix, %q), %q, %q)\n", name, flag.GetName(), flag.GetShort(), flag.GetUsage())
			}
		case pgs.UInt64T, pgs.Fixed64T:
			// Use custom type for uint64 since pflag doesn't have native support
			if f.HasOptionalKeyword() {
				decl += fmt.Sprintf("fs.VarP(types.StringToUint64(x.%s), utils.BuildFlagName(prefix, %q), %q, %q)\n", name, flag.GetName(), flag.GetShort(), flag.GetUsage())
			} else {
				decl += fmt.Sprintf("fs.VarP(types.StringToUint64(&x.%s), utils.BuildFlagName(prefix, %q), %q, %q)\n", name, flag.GetName(), flag.GetShort(), flag.GetUsage())
			}
		default:
			// Default to int32 for unknown types
			if f.HasOptionalKeyword() {
				decl += fmt.Sprintf("fs.VarP(types.StringToInt32(x.%s), utils.BuildFlagName(prefix, %q), %q, %q)\n", name, flag.GetName(), flag.GetShort(), flag.GetUsage())
			} else {
				decl += fmt.Sprintf("fs.VarP(types.StringToInt32(&x.%s), utils.BuildFlagName(prefix, %q), %q, %q)\n", name, flag.GetName(), flag.GetShort(), flag.GetUsage())
			}
		}

	case flags.MapFormatType_MAP_FORMAT_TYPE_JSON:
		// For JSON format, use the existing JSON handling
		if f.HasOptionalKeyword() {
			decl += fmt.Sprintf("fs.VarP(types.JSON(x.%s), utils.BuildFlagName(prefix, %q), %q, %q)\n", name, flag.GetName(), flag.GetShort(), flag.GetUsage())
		} else {
			decl += fmt.Sprintf("fs.VarP(types.JSON(&x.%s), utils.BuildFlagName(prefix, %q), %q, %q)\n", name, flag.GetName(), flag.GetShort(), flag.GetUsage())
		}
	}

	if flag.GetHidden() {
		decl += fmt.Sprintf("fs.MarkHidden(%q)\n", flag.GetName())
	}
	if flag.GetDeprecated() {
		decl += fmt.Sprintf("fs.MarkDeprecated(%q, %q)\n", flag.GetName(), flag.GetDeprecatedUsage())
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
		return m.processBytesFlag(f, name, r.Bytes, wk, false)
	case *flags.FieldFlags_Enum:
		return m.processEnumFlag(f, name, r.Enum, wk, genOneOfField...)
	case *flags.FieldFlags_Duration:
		return m.processDurationFlag(f, name, r.Duration, wk, false)
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
	case *flags.FieldFlags_Map:
		return m.processMapFlag(f, name, r.Map, wk)
	case *flags.FieldFlags_Repeated:
		return m.processRepeatedFlag(f, name, r.Repeated)
	case nil: // noop
	default:
		m.Failf("unknown rule type (%T)", field.Type)

	}
	return fmt.Sprint("\n// ", f.Name()), true
}
