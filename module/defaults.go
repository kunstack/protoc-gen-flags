package module

import (
	"fmt"

	"github.com/kunstack/protoc-gen-flags/flags"
	pgs "github.com/lyft/protoc-gen-star"
)

func (m *Module) genFieldDefaults(f pgs.Field) string {
	m.Push(f.Name().String())
	defer m.Pop()
	var field flags.FieldFlags
	ok, err := f.Extension(flags.E_Value, &field)
	if err != nil || !ok {
		return ""
	}

	wk := pgs.UnknownWKT
	if emb := f.Type().Embed(); emb != nil {
		wk = emb.WellKnownType()
	}

	// 检查是否是oneof字段，发出警告并跳过处理
	if f.InRealOneOf() {
		m.Failf("oneof field '%s' is not supported for flag generation. Please use regular fields instead.", f.Name())
		return ""
	}

	name := m.ctx.Name(f)
	switch r := field.Type.(type) {
	case *flags.FieldFlags_Float:
		return m.genCommonDefaults(f, name, 0, r.Float.Default, wk)
	case *flags.FieldFlags_Double:
		return m.genCommonDefaults(f, name, 0.0, r.Double.Default, wk)
	case *flags.FieldFlags_Int32:
		return m.genCommonDefaults(f, name, 0, r.Int32.Default, wk)
	case *flags.FieldFlags_Int64:
		return m.genCommonDefaults(f, name, 0, r.Int64.Default, wk)
	case *flags.FieldFlags_Uint32:
		return m.genCommonDefaults(f, name, 0, r.Uint32.Default, wk)
	case *flags.FieldFlags_Uint64:
		return m.genCommonDefaults(f, name, 0, r.Uint64.Default, wk)
	case *flags.FieldFlags_Sint32:
		return m.genCommonDefaults(f, name, 0, r.Sint32.Default, wk)
	case *flags.FieldFlags_Sint64:
		return m.genCommonDefaults(f, name, 0, r.Sint64.Default, wk)
	case *flags.FieldFlags_Fixed32:
		return m.genCommonDefaults(f, name, 0, r.Fixed32.Default, wk)
	case *flags.FieldFlags_Fixed64:
		return m.genCommonDefaults(f, name, 0, r.Fixed64.Default, wk)
	case *flags.FieldFlags_Sfixed32:
		return m.genCommonDefaults(f, name, 0, r.Sfixed32.Default, wk)
	case *flags.FieldFlags_Sfixed64:
		return m.genCommonDefaults(f, name, 0, r.Sfixed64.Default, wk)
	case *flags.FieldFlags_Bool:
		return m.genCommonDefaults(f, name, false, r.Bool.Default, wk)
	case *flags.FieldFlags_String_:
		return ""
		//return m.genCommonDefaults(f, name, `""`, r.String_.Default, wk)
	case *flags.FieldFlags_Bytes:
		return ""
	case *flags.FieldFlags_Enum:
		return ""
	case *flags.FieldFlags_Duration:
		return ""
	case *flags.FieldFlags_Timestamp:
		return ""
	case *flags.FieldFlags_Message:
		return m.genMessageDefaults(f, name, r.Message)
	case *flags.FieldFlags_Map:
		return ""
	case *flags.FieldFlags_Repeated:
		return ""
	case nil: // noop
	default:
		m.Failf("unknown rule type (%T)", field.Type)
	}
	return fmt.Sprint("\n// ", f.Name())
}
