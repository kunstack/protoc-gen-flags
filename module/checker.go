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

		var fieldDefaults flags.FieldFlags
		_, err = f.Extension(flags.E_Field, &fieldDefaults)
		m.CheckErr(err, "unable to read flags from field")

		m.CheckFieldRules(f.Type(), &fieldDefaults)
		m.Pop()
	}
}

func (m *Module) CheckFieldRules(typ FieldType, fieldDefaults *flags.FieldFlags) {
	if fieldDefaults == nil {
		return
	}

	switch r := fieldDefaults.Type.(type) {
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
		m.CheckPrimitiveFlag(typ, r.Bytes)
	case *flags.FieldFlags_Enum:
		m.MustType(typ, pgs.EnumT, pgs.UnknownWKT)
		m.CheckEnum(typ, r.Enum)
		m.CheckPrimitiveFlag(typ, r.Enum)
	case *flags.FieldFlags_Duration:
		m.CheckDuration(typ, r.Duration)
		m.CheckPrimitiveFlag(typ, r.Duration)
	case *flags.FieldFlags_Timestamp:
		m.CheckTimestamp(typ, r.Timestamp)
		m.CheckPrimitiveFlag(typ, r.Timestamp)
	//case *flags.FieldFlags_Repeated:
	//	m.MustType(typ, pgs.MessageT, pgs.UnknownWKT)
	case nil: // noop
	default:
		m.Failf("unknown rule type (%T)", fieldDefaults.Type)
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

func (m *Module) CheckEnum(ft FieldType, r *flags.PrimitiveFlag) {
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
	if flag.GetMessage().Disabled {
		return
	}
	current := m.ctx.ImportPath(f.Message()).String()
	if i := m.ctx.ImportPath(f.Type().Embed()).String(); i != current {
		m.imports[i] = struct{}{}
	}
}

func (m *Module) CheckDuration(ft FieldType, r *flags.PrimitiveFlag) {
	if embed := ft.Embed(); embed == nil || embed.WellKnownType() != pgs.DurationWKT {
		m.Failf("unexpected field type (%T) for Duration, expected google.protobuf.Duration ", ft)
	}
}

func (m *Module) CheckPrimitiveFlag(ft FieldType, r *flags.PrimitiveFlag) {
	// Check if deprecated flag has proper deprecation usage message
	if r.Deprecated && r.DeprecatedUsage == "" {
		m.Failf("deprecated flag must provide deprecated_usage message")
	}
}

func (m *Module) CheckTimestamp(ft FieldType, r *flags.PrimitiveFlag) {
	if embed := ft.Embed(); embed == nil || embed.WellKnownType() != pgs.TimestampWKT {
		m.Failf("unexpected field type (%T) for Timestamp, expected google.protobuf.Timestamp ", ft)
	}
}

func (m *Module) mustFieldType(ft FieldType) pgs.FieldType {
	typ, ok := ft.(pgs.FieldType)
	if !ok {
		m.Failf("unexpected field type (%T)", ft)
	}
	return typ
}
