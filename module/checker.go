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

	for _, f := range msg.Fields() {
		m.Push(f.Name().String())

		var field flags.FieldFlags
		_, err = f.Extension(flags.E_Value, &field)
		m.CheckErr(err, "unable to read flags from field")
		m.CheckFieldRules(f, &field)
		m.Pop()
	}
}

func (m *Module) CheckFieldRules(f pgs.Field, field *flags.FieldFlags) {
	if field == nil {
		return
	}
	typ := f.Type()

	switch r := field.Type.(type) {
	case *flags.FieldFlags_Float:
		m.checkCommon(typ, r.Float, pgs.FloatT, pgs.FloatValueWKT, false)
	case *flags.FieldFlags_Double:
		m.checkCommon(typ, r.Double, pgs.DoubleT, pgs.DoubleValueWKT, false)
	case *flags.FieldFlags_Int32:
		m.checkCommon(typ, r.Int32, pgs.Int32T, pgs.Int32ValueWKT, false)
	case *flags.FieldFlags_Int64:
		m.checkCommon(typ, r.Int64, pgs.Int64T, pgs.Int64ValueWKT, false)
	case *flags.FieldFlags_Uint32:
		m.checkCommon(typ, r.Uint32, pgs.UInt32T, pgs.UInt32ValueWKT, false)
	case *flags.FieldFlags_Uint64:
		m.checkCommon(typ, r.Uint64, pgs.UInt64T, pgs.UInt64ValueWKT, false)
	case *flags.FieldFlags_Sint32:
		m.checkCommon(typ, r.Sint32, pgs.SInt32, pgs.UnknownWKT, false)
	case *flags.FieldFlags_Sint64:
		m.checkCommon(typ, r.Sint64, pgs.SInt64, pgs.UnknownWKT, false)
	case *flags.FieldFlags_Fixed32:
		m.checkCommon(typ, r.Fixed32, pgs.Fixed32T, pgs.UnknownWKT, false)
	case *flags.FieldFlags_Fixed64:
		m.checkCommon(typ, r.Fixed64, pgs.Fixed64T, pgs.UnknownWKT, false)
	case *flags.FieldFlags_Sfixed32:
		m.checkCommon(typ, r.Sfixed32, pgs.SFixed32, pgs.UnknownWKT, false)
	case *flags.FieldFlags_Sfixed64:
		m.checkCommon(typ, r.Sfixed64, pgs.SFixed64, pgs.UnknownWKT, false)
	case *flags.FieldFlags_Bool:
		m.checkCommon(typ, r.Bool, pgs.BoolT, pgs.BoolValueWKT, false)
	case *flags.FieldFlags_String_:
		m.checkCommon(typ, r.String_, pgs.StringT, pgs.StringValueWKT, false)
	case *flags.FieldFlags_Bytes:
		m.checkCommon(typ, r.Bytes, pgs.BytesT, pgs.BytesValueWKT, false)
	case *flags.FieldFlags_Enum:
		m.checkEnum(typ, r.Enum, pgs.EnumT, pgs.UnknownWKT)
	case *flags.FieldFlags_Duration:
		m.checkCommon(typ, r.Duration, pgs.MessageT, pgs.DurationWKT, false)
	case *flags.FieldFlags_Timestamp:
		m.checkTimestamp(typ, r.Timestamp)
	case *flags.FieldFlags_Repeated:
		el, ok := typ.(Element)
		if !ok || el.Element() == nil {
			m.Failf("field '%s' is not a repeated field (actual type: %v), "+
				"but repeated flag configuration was specified", f.Name(), typ.ProtoType())
			return
		}
		m.CheckRepeatedFlag(el.Element(), r.Repeated)
	case *flags.FieldFlags_Message:
		current := m.ctx.ImportPath(f).String()
		if i := m.ctx.ImportPath(typ.Embed()).String(); i != current {
			m.imports[i] = struct{}{}
		}
		m.checkMessage(typ, r.Message)
	case *flags.FieldFlags_Map:
		m.checkMap(typ, r.Map)
	case nil: // noop
	default:
		m.Failf("unknown rule type (%T)", field.Type)
	}
}

func (m *Module) mustType(typ FieldType, pt pgs.ProtoType, wrapper pgs.WellKnownType) {
	if emb := typ.Embed(); emb != nil && emb.IsWellKnown() && emb.WellKnownType() == wrapper {
		if wrapper != pgs.DurationWKT && wrapper != pgs.TimestampWKT {
			m.mustType(emb.Fields()[0].Type(), pt, pgs.UnknownWKT)
			return
		}
	}

	m.Assert(typ.ProtoType() == pt,
		" expected flags for ",
		typ.ProtoType().Proto(),
		" but got ",
		pt.Proto(),
	)
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

func (m *Module) CheckRepeatedFlag(typ FieldType, repeated *flags.RepeatedFlags) {
	if repeated == nil {
		return
	}

	if typ, ok := typ.(Repeatable); ok {
		m.Assert(typ.IsRepeated(), "repeated flag should be used for repeated fields")
	}

	switch r := repeated.Type.(type) {
	case *flags.RepeatedFlags_Float:
		m.checkCommon(typ, r.Float, pgs.FloatT, pgs.FloatValueWKT, true)
	case *flags.RepeatedFlags_Double:
		m.checkCommon(typ, r.Double, pgs.DoubleT, pgs.DoubleValueWKT, true)
	case *flags.RepeatedFlags_Int32:
		m.checkCommon(typ, r.Int32, pgs.Int32T, pgs.Int32ValueWKT, true)
	case *flags.RepeatedFlags_Int64:
		m.checkCommon(typ, r.Int64, pgs.Int64T, pgs.Int64ValueWKT, true)
	case *flags.RepeatedFlags_Uint32:
		m.checkCommon(typ, r.Uint32, pgs.UInt32T, pgs.UInt32ValueWKT, true)
	case *flags.RepeatedFlags_Uint64:
		m.checkCommon(typ, r.Uint64, pgs.UInt64T, pgs.UInt32ValueWKT, true)
	case *flags.RepeatedFlags_Sint32:
		m.checkCommon(typ, r.Sint32, pgs.SInt32, pgs.UnknownWKT, true)
	case *flags.RepeatedFlags_Sint64:
		m.checkCommon(typ, r.Sint64, pgs.SInt64, pgs.UnknownWKT, true)
	case *flags.RepeatedFlags_Fixed32:
		m.checkCommon(typ, r.Fixed32, pgs.Fixed32T, pgs.UnknownWKT, true)
	case *flags.RepeatedFlags_Fixed64:
		m.checkCommon(typ, r.Fixed64, pgs.Fixed64T, pgs.UnknownWKT, true)
	case *flags.RepeatedFlags_Sfixed32:
		m.checkCommon(typ, r.Sfixed32, pgs.SFixed32, pgs.UnknownWKT, true)
	case *flags.RepeatedFlags_Sfixed64:
		m.checkCommon(typ, r.Sfixed64, pgs.SFixed64, pgs.UnknownWKT, true)
	case *flags.RepeatedFlags_Bool:
		m.checkCommon(typ, r.Bool, pgs.BoolT, pgs.BoolValueWKT, true)
	case *flags.RepeatedFlags_String_:
		m.checkCommon(typ, r.String_, pgs.StringT, pgs.StringValueWKT, true)
	case *flags.RepeatedFlags_Bytes:
		m.checkCommon(typ, r.Bytes, pgs.BytesT, pgs.BytesValueWKT, true)
	case *flags.RepeatedFlags_Enum:
		m.checkEnumSlice(typ, r.Enum, pgs.EnumT, pgs.UnknownWKT)
	case *flags.RepeatedFlags_Duration:
		m.checkCommon(typ, r.Duration, pgs.MessageT, pgs.DurationWKT, true)
	case *flags.RepeatedFlags_Timestamp:
		m.checkTimestampSlice(typ, r.Timestamp)
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
