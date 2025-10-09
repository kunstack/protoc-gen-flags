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
	"strings"

	"github.com/kunstack/protoc-gen-flags/flags"
	pgs "github.com/lyft/protoc-gen-star"
)

// Heavily taken from https://github.com/envoyproxy/protoc-gen-validate/blob/main/module/checker.go

type FieldType interface {
	ProtoType() pgs.ProtoType
	Embed() pgs.Message
}

type Repeatable interface {
	IsRepeated() bool
}

type Element interface {
	Element() pgs.FieldTypeElem
}

func (m *Module) Check(msg pgs.Message) {
	m.Push("msg: " + msg.Name().String())
	defer m.Pop()

	var disabled bool
	_, err := msg.Extension(flags.E_Disabled, &disabled)
	m.CheckErr(err, "unable to read flags extension from message")

	if disabled {
		m.Debug("flags disabled, skipping checks")
		return
	}

	// Track field names to detect duplicates
	fieldNames := make(map[string]pgs.Field)

	for _, f := range msg.Fields() {
		m.Push(f.Name().String())

		var field flags.FieldFlags
		_, err = f.Extension(flags.E_Value, &field)
		m.CheckErr(err, "unable to read flags from field")

		// Check for duplicate field names
		flagName := m.getFlagName(f, &field)
		if flagName != "" {
			if existingField, exists := fieldNames[flagName]; exists {
				m.Failf("duplicate flag name '%s' detected. Field '%s' conflicts with field '%s'",
					flagName, f.Name().String(), existingField.Name().String())
			}
			fieldNames[flagName] = f
		}

		m.CheckFieldRules(f.Type(), &field)
		m.Pop()
	}
}

// getFlagName extracts the flag name from field configuration
func (m *Module) getFlagName(field pgs.Field, flag *flags.FieldFlags) string {
	if flag == nil {
		return ""
	}

	var flagName string
	switch r := flag.Type.(type) {
	case *flags.FieldFlags_Float:
		if r.Float.Name != "" {
			flagName = r.Float.Name
		}
	case *flags.FieldFlags_Double:
		if r.Double.Name != "" {
			flagName = r.Double.Name
		}
	case *flags.FieldFlags_Int32:
		if r.Int32.Name != "" {
			flagName = r.Int32.Name
		}
	case *flags.FieldFlags_Int64:
		if r.Int64.Name != "" {
			flagName = r.Int64.Name
		}
	case *flags.FieldFlags_Uint32:
		if r.Uint32.Name != "" {
			flagName = r.Uint32.Name
		}
	case *flags.FieldFlags_Uint64:
		if r.Uint64.Name != "" {
			flagName = r.Uint64.Name
		}
	case *flags.FieldFlags_Sint32:
		if r.Sint32.Name != "" {
			flagName = r.Sint32.Name
		}
	case *flags.FieldFlags_Sint64:
		if r.Sint64.Name != "" {
			flagName = r.Sint64.Name
		}
	case *flags.FieldFlags_Fixed32:
		if r.Fixed32.Name != "" {
			flagName = r.Fixed32.Name
		}
	case *flags.FieldFlags_Fixed64:
		if r.Fixed64.Name != "" {
			flagName = r.Fixed64.Name
		}
	case *flags.FieldFlags_Sfixed32:
		if r.Sfixed32.Name != "" {
			flagName = r.Sfixed32.Name
		}
	case *flags.FieldFlags_Sfixed64:
		if r.Sfixed64.Name != "" {
			flagName = r.Sfixed64.Name
		}
	case *flags.FieldFlags_Bool:
		if r.Bool.Name != "" {
			flagName = r.Bool.Name
		}
	case *flags.FieldFlags_String_:
		if r.String_.Name != "" {
			flagName = r.String_.Name
		}
	case *flags.FieldFlags_Bytes:
		if r.Bytes.Name != "" {
			flagName = r.Bytes.Name
		}
	case *flags.FieldFlags_Enum:
		if r.Enum.Name != "" {
			flagName = r.Enum.Name
		}
	case *flags.FieldFlags_Duration:
		if r.Duration.Name != "" {
			flagName = r.Duration.Name
		}
	case *flags.FieldFlags_Timestamp:
		if r.Timestamp.Name != "" {
			flagName = r.Timestamp.Name
		}
	case *flags.FieldFlags_Message:
		if r.Message.Name != "" {
			flagName = r.Message.Name
		}
	case *flags.FieldFlags_Map:
		if r.Map.Name != "" {
			flagName = r.Map.Name
		}
	case *flags.FieldFlags_Repeated:
		// Handle repeated field types
		switch r2 := r.Repeated.Type.(type) {
		case *flags.RepeatedFlags_Float:
			if r2.Float.Name != "" {
				flagName = r2.Float.Name
			}
		case *flags.RepeatedFlags_Double:
			if r2.Double.Name != "" {
				flagName = r2.Double.Name
			}
		case *flags.RepeatedFlags_Int32:
			if r2.Int32.Name != "" {
				flagName = r2.Int32.Name
			}
		case *flags.RepeatedFlags_Int64:
			if r2.Int64.Name != "" {
				flagName = r2.Int64.Name
			}
		case *flags.RepeatedFlags_Uint32:
			if r2.Uint32.Name != "" {
				flagName = r2.Uint32.Name
			}
		case *flags.RepeatedFlags_Uint64:
			if r2.Uint64.Name != "" {
				flagName = r2.Uint64.Name
			}
		case *flags.RepeatedFlags_Sint32:
			if r2.Sint32.Name != "" {
				flagName = r2.Sint32.Name
			}
		case *flags.RepeatedFlags_Sint64:
			if r2.Sint64.Name != "" {
				flagName = r2.Sint64.Name
			}
		case *flags.RepeatedFlags_Fixed32:
			if r2.Fixed32.Name != "" {
				flagName = r2.Fixed32.Name
			}
		case *flags.RepeatedFlags_Fixed64:
			if r2.Fixed64.Name != "" {
				flagName = r2.Fixed64.Name
			}
		case *flags.RepeatedFlags_Sfixed32:
			if r2.Sfixed32.Name != "" {
				flagName = r2.Sfixed32.Name
			}
		case *flags.RepeatedFlags_Sfixed64:
			if r2.Sfixed64.Name != "" {
				flagName = r2.Sfixed64.Name
			}
		case *flags.RepeatedFlags_Bool:
			if r2.Bool.Name != "" {
				flagName = r2.Bool.Name
			}
		case *flags.RepeatedFlags_String_:
			if r2.String_.Name != "" {
				flagName = r2.String_.Name
			}
		case *flags.RepeatedFlags_Bytes:
			if r2.Bytes.Name != "" {
				flagName = r2.Bytes.Name
			}
		case *flags.RepeatedFlags_Enum:
			if r2.Enum.Name != "" {
				flagName = r2.Enum.Name
			}
		case *flags.RepeatedFlags_Duration:
			if r2.Duration.Name != "" {
				flagName = r2.Duration.Name
			}
		case *flags.RepeatedFlags_Timestamp:
			if r2.Timestamp.Name != "" {
				flagName = r2.Timestamp.Name
			}
		}
	}

	// If no custom name is specified, use the default field name conversion
	if flagName == "" {
		flagName = strings.ToLower(field.Name().String())
	}

	return flagName
}

func (m *Module) CheckFieldRules(typ FieldType, field *flags.FieldFlags) {
	if field == nil {
		return
	}

	switch r := field.Type.(type) {
	case *flags.FieldFlags_Float:
		m.MustType(typ, pgs.FloatT, pgs.FloatValueWKT)
		m.CheckPrimitiveFlag(typ, r.Float)
	case *flags.FieldFlags_Double:
		m.MustType(typ, pgs.DoubleT, pgs.DoubleValueWKT)
		m.CheckPrimitiveFlag(typ, r.Double)
	case *flags.FieldFlags_Int32:
		m.MustType(typ, pgs.Int32T, pgs.Int32ValueWKT)
		m.CheckPrimitiveFlag(typ, r.Int32)
	case *flags.FieldFlags_Int64:
		m.MustType(typ, pgs.Int64T, pgs.Int64ValueWKT)
		m.CheckPrimitiveFlag(typ, r.Int64)
	case *flags.FieldFlags_Uint32:
		m.MustType(typ, pgs.UInt32T, pgs.UInt32ValueWKT)
		m.CheckPrimitiveFlag(typ, r.Uint32)
	case *flags.FieldFlags_Uint64:
		m.MustType(typ, pgs.UInt64T, pgs.UInt64ValueWKT)
		m.CheckPrimitiveFlag(typ, r.Uint64)
	case *flags.FieldFlags_Sint32:
		m.MustType(typ, pgs.SInt32, pgs.UnknownWKT)
		m.CheckPrimitiveFlag(typ, r.Sint32)
	case *flags.FieldFlags_Sint64:
		m.MustType(typ, pgs.SInt64, pgs.UnknownWKT)
		m.CheckPrimitiveFlag(typ, r.Sint64)
	case *flags.FieldFlags_Fixed32:
		m.MustType(typ, pgs.Fixed32T, pgs.UnknownWKT)
		m.CheckPrimitiveFlag(typ, r.Fixed32)
	case *flags.FieldFlags_Fixed64:
		m.MustType(typ, pgs.Fixed64T, pgs.UnknownWKT)
		m.CheckPrimitiveFlag(typ, r.Fixed64)
	case *flags.FieldFlags_Sfixed32:
		m.MustType(typ, pgs.SFixed32, pgs.UnknownWKT)
		m.CheckPrimitiveFlag(typ, r.Sfixed32)
	case *flags.FieldFlags_Sfixed64:
		m.MustType(typ, pgs.SFixed64, pgs.UnknownWKT)
		m.CheckPrimitiveFlag(typ, r.Sfixed64)
	case *flags.FieldFlags_Bool:
		m.MustType(typ, pgs.BoolT, pgs.BoolValueWKT)
		m.CheckPrimitiveFlag(typ, r.Bool)
	case *flags.FieldFlags_String_:
		m.MustType(typ, pgs.StringT, pgs.StringValueWKT)
		m.CheckPrimitiveFlag(typ, r.String_)
	case *flags.FieldFlags_Bytes:
		m.MustType(typ, pgs.BytesT, pgs.BytesValueWKT)
		m.CheckBytes(typ, r.Bytes)
	case *flags.FieldFlags_Enum:
		m.MustType(typ, pgs.EnumT, pgs.UnknownWKT)
		m.CheckEnum(typ, r.Enum)
		m.CheckPrimitiveFlag(typ, r.Enum)
	case *flags.FieldFlags_Duration:
		m.CheckDuration(typ, r.Duration)
		m.CheckPrimitiveFlag(typ, r.Duration)
	case *flags.FieldFlags_Timestamp:
		m.CheckTimestamp(typ, r.Timestamp)
	case *flags.FieldFlags_Repeated:
		el, ok := typ.(Element)
		if !ok {
			m.Failf("repeated field does not implement Element interface")
			return
		}
		m.CheckRepeatedFlag(el.Element(), r.Repeated)
	case *flags.FieldFlags_Message:
		m.MustType(typ, pgs.MessageT, pgs.UnknownWKT)
	case *flags.FieldFlags_Map:
		m.CheckMap(typ, r.Map)

	case nil: // noop
	default:
		m.Failf("unknown rule type (%T)", field.Type)
	}
}

func (m *Module) MustType(typ FieldType, pt pgs.ProtoType, wrapper pgs.WellKnownType) {
	if emb := typ.Embed(); emb != nil && emb.IsWellKnown() && emb.WellKnownType() == wrapper {
		m.MustType(emb.Fields()[0].Type(), pt, pgs.UnknownWKT)
		return
	}
	if typ, ok := typ.(Repeatable); ok {
		m.Assert(!typ.IsRepeated(), "repeated flag should be used for repeated fields")
	}

	m.Assert(typ.ProtoType() == pt,
		" expected flags for ",
		typ.ProtoType().Proto(),
		" but got ",
		pt.Proto(),
	)
}

func (m *Module) CheckEnum(ft FieldType, _ *flags.EnumFlag) {
	_, ok := ft.(interface {
		Enum() pgs.Enum
	})

	if !ok {
		m.Failf("unexpected field type (%T)", ft)
	}
}

func (m *Module) CheckMessage(f pgs.Field, flag *flags.FieldFlags) {
	m.Assert(f.Type().IsEmbed(), "field is not embedded but got message flags")
	emb := f.Type().Embed()
	if emb != nil && emb.IsWellKnown() {
		switch emb.WellKnownType() {
		case pgs.DurationWKT:
			m.Failf("Duration value should be used for Duration fields")
		case pgs.TimestampWKT:
			m.Failf("Timestamp value should be used for Timestamp fields")
		}
	}
	if !flag.GetMessage().Nested {
		return
	}
	current := m.ctx.ImportPath(f.Message()).String()
	if i := m.ctx.ImportPath(f.Type().Embed()).String(); i != current {
		m.imports[i] = struct{}{}
	}
}

func (m *Module) CheckDuration(ft FieldType, r *flags.DurationFlag) {
	if embed := ft.Embed(); embed == nil || embed.WellKnownType() != pgs.DurationWKT {
		m.Failf("unexpected field type (%T) for Duration, expected google.protobuf.Duration ", ft)
	}
	m.CheckPrimitiveFlag(ft, r)
}

type commonFlag interface {
	GetDisabled() bool
	GetName() string
	GetUsage() string
	GetDeprecated() bool
	GetDeprecatedUsage() string
	GetHidden() bool
	GetShort() string
}

func (m *Module) CheckPrimitiveFlag(_ FieldType, r commonFlag) {
	if r.GetUsage() == "" {
		m.Failf("usage is required for flag")
	}
	// Check if deprecated flag has proper deprecation usage message
	if r.GetDeprecated() && r.GetDeprecatedUsage() == "" {
		m.Failf("deprecated flag must provide deprecated_usage message")
	}
}

func (m *Module) CheckTimestamp(ft FieldType, r *flags.TimestampFlag) {
	if embed := ft.Embed(); embed == nil || embed.WellKnownType() != pgs.TimestampWKT {
		m.Failf("unexpected field type (%T) for Timestamp, expected google.protobuf.Timestamp ", ft)
	}

	if r.Usage == "" {
		m.Failf("usage is required for flag")
	}
	// Check if deprecated flag has proper deprecation usage message
	if r.Deprecated && r.DeprecatedUsage == "" {
		m.Failf("deprecated flag must provide deprecated_usage message")
	}

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

func (m *Module) CheckBytes(_ FieldType, r *flags.BytesFlag) {
	if r.Usage == "" {
		m.Failf("usage is required for flag")
	}
	// Check if deprecated flag has proper deprecation usage message
	if r.Deprecated && r.DeprecatedUsage == "" {
		m.Failf("deprecated flag must provide deprecated_usage message")
	}
}

func (m *Module) CheckMap(typ FieldType, flag *flags.MapFlag) {
	if flag == nil {
		return
	}

	// Ensure the field is actually a map type
	fieldType := m.mustFieldType(typ)
	m.Assert(fieldType.IsMap(), "map flag should be used for map fields")

	// Check that usage is provided
	if flag.Usage == "" {
		m.Failf("usage is required for map flag")
	}

	// Check if deprecated flag has proper deprecation usage message
	if flag.Deprecated && flag.DeprecatedUsage == "" {
		m.Failf("deprecated map flag must provide deprecated_usage message")
	}

	// Validate that the format matches the actual field types
	keyElem := fieldType.Key()
	valueElem := fieldType.Element()

	switch flag.GetFormat() {
	case flags.MapFormatType_MAP_FORMAT_TYPE_STRING_TO_STRING:
		// Validate key is string
		if keyElem.ProtoType() != pgs.StringT {
			m.Failf("STRING_TO_STRING format requires string keys, but got %v", keyElem.ProtoType())
		}
		// Validate value is string
		if valueElem.ProtoType() != pgs.StringT {
			m.Failf("STRING_TO_STRING format requires string values, but got %v", valueElem.ProtoType())
		}

	case flags.MapFormatType_MAP_FORMAT_TYPE_STRING_TO_INT:
		// Validate key is string
		if keyElem.ProtoType() != pgs.StringT {
			m.Failf("STRING_TO_INT format requires string keys, but got %v", keyElem.ProtoType())
		}
		// Validate value is an integer type (support all integer types)
		switch valueElem.ProtoType() {
		case pgs.Int32T, pgs.SInt32, pgs.SFixed32,
			pgs.Int64T, pgs.SInt64, pgs.SFixed64,
			pgs.UInt32T, pgs.Fixed32T,
			pgs.UInt64T, pgs.Fixed64T:
			// These are all valid integer types
		default:
			m.Failf("STRING_TO_INT format requires integer values, but got %v", valueElem.ProtoType())
		}

	case flags.MapFormatType_MAP_FORMAT_TYPE_JSON:
		// JSON format is flexible, no strict type validation needed
		// Just ensure it's actually a map
		break

	case flags.MapFormatType_MAP_FORMAT_TYPE_UNSPECIFIED:
		// Default to JSON format, no additional validation needed
		break

	default:
		m.Failf("unknown map format type: %v", flag.GetFormat())
	}
}

func (m *Module) CheckRepeatedFlag(typ FieldType, repeated *flags.RepeatedFlags) {
	if repeated == nil {
		return
	}

	if typ, ok := typ.(Repeatable); ok {
		m.Assert(typ.IsRepeated(), "repeated flag should be used for repeated fields")
	}

	switch r := repeated.Type.(type) {
	case *flags.RepeatedFlags_Float:
		m.MustType(typ, pgs.FloatT, pgs.FloatValueWKT)
		m.CheckPrimitiveFlag(typ, r.Float)
	case *flags.RepeatedFlags_Double:
		m.MustType(typ, pgs.DoubleT, pgs.DoubleValueWKT)
		m.CheckPrimitiveFlag(typ, r.Double)
	case *flags.RepeatedFlags_Int32:
		m.MustType(typ, pgs.Int32T, pgs.Int32ValueWKT)
		m.CheckPrimitiveFlag(typ, r.Int32)
	case *flags.RepeatedFlags_Int64:
		m.MustType(typ, pgs.Int64T, pgs.Int64ValueWKT)
		m.CheckPrimitiveFlag(typ, r.Int64)
	case *flags.RepeatedFlags_Uint32:
		m.MustType(typ, pgs.UInt32T, pgs.UInt32ValueWKT)
		m.CheckPrimitiveFlag(typ, r.Uint32)
	case *flags.RepeatedFlags_Uint64:
		m.MustType(typ, pgs.UInt64T, pgs.UInt32ValueWKT)
		m.CheckPrimitiveFlag(typ, r.Uint64)
	case *flags.RepeatedFlags_Sint32:
		m.MustType(typ, pgs.SInt32, pgs.UnknownWKT)
		m.CheckPrimitiveFlag(typ, r.Sint32)
	case *flags.RepeatedFlags_Sint64:
		m.MustType(typ, pgs.SInt64, pgs.UnknownWKT)
		m.CheckPrimitiveFlag(typ, r.Sint64)
	case *flags.RepeatedFlags_Fixed32:
		m.MustType(typ, pgs.Fixed32T, pgs.UnknownWKT)
		m.CheckPrimitiveFlag(typ, r.Fixed32)
	case *flags.RepeatedFlags_Fixed64:
		m.MustType(typ, pgs.Fixed64T, pgs.UnknownWKT)
		m.CheckPrimitiveFlag(typ, r.Fixed64)
	case *flags.RepeatedFlags_Sfixed32:
		m.MustType(typ, pgs.SFixed32, pgs.UnknownWKT)
		m.CheckPrimitiveFlag(typ, r.Sfixed32)
	case *flags.RepeatedFlags_Sfixed64:
		m.MustType(typ, pgs.SFixed64, pgs.UnknownWKT)
		m.CheckPrimitiveFlag(typ, r.Sfixed64)
	case *flags.RepeatedFlags_Bool:
		m.MustType(typ, pgs.BoolT, pgs.BoolValueWKT)
		m.CheckPrimitiveFlag(typ, r.Bool)
	case *flags.RepeatedFlags_String_:
		m.MustType(typ, pgs.StringT, pgs.StringValueWKT)
		m.CheckPrimitiveFlag(typ, r.String_)
	case *flags.RepeatedFlags_Bytes:
		m.MustType(typ, pgs.BytesT, pgs.BytesValueWKT)
		m.CheckBytes(typ, r.Bytes)
	case *flags.RepeatedFlags_Enum:
		m.MustType(typ, pgs.EnumT, pgs.UnknownWKT)
		m.CheckEnum(typ, r.Enum)
		m.CheckPrimitiveFlag(typ, r.Enum)
	case *flags.RepeatedFlags_Duration:
		m.CheckDuration(typ, r.Duration)
		m.CheckPrimitiveFlag(typ, r.Duration)
	case *flags.RepeatedFlags_Timestamp:
		m.CheckTimestamp(typ, r.Timestamp)
	default:
		m.Failf("unknown repeated flag type (%T)", repeated.Type)
	}
}

func (m *Module) mustFieldType(ft FieldType) pgs.FieldType {
	typ, ok := ft.(pgs.FieldType)
	if !ok {
		m.Failf("unexpected field type (%T)", ft)
	}
	return typ
}
