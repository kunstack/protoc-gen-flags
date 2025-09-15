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
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/kunstack/protoc-gen-flags/flags"
	pgs "github.com/lyft/protoc-gen-star"
	"github.com/prometheus/common/model"
)

func (m *Module) genFieldFlags(f pgs.Field, genOneOfField ...bool) (string, bool) {
	m.Push(f.Name().String())
	defer m.Pop()
	var fieldFlags flags.FieldFlags
	ok, err := f.Extension(flags.E_Field, &fieldFlags)
	if err != nil || !ok {
		return "", false
	}
	wk := pgs.UnknownWKT
	if emb := f.Type().Embed(); emb != nil {
		wk = emb.WellKnownType()
	}
	//if !isOk(genOneOfField) && f.InRealOneOf() {
	//	if m.isOneOfDone(f.OneOf()) {
	//		return "", false
	//	}
	//	m.setOneOfDone(f.OneOf())
	//	var out string
	//	var oneOfDefault string
	//	if _, err := f.OneOf().Extension(flags.E_Oneof, &oneOfDefault); err != nil {
	//		m.Fail(err)
	//	}
	//	var defaultField pgs.Field
	//	for _, f := range f.OneOf().Fields() {
	//		if f.Name().String() == oneOfDefault {
	//			defaultField = f
	//		}
	//	}
	//	if defaultField != nil {
	//		out += fmt.Sprint(`
	//			if x.`, m.ctx.Name(f.OneOf()), ` == nil {
	//				x.`, m.ctx.Name(f.OneOf()), ` = &`, m.ctx.OneofOption(defaultField), `{}
	//			}`)
	//	}
	//	out += fmt.Sprint(`
	//		switch x := x.`, m.ctx.Name(f.OneOf()), `.(type) {`)
	//	for _, f := range f.OneOf().Fields() {
	//		def, ok := m.genFieldFlags(f, true)
	//		if !ok {
	//			continue
	//		}
	//		out += fmt.Sprint(`
	//			case *`, m.ctx.OneofOption(f), `: `, def)
	//	}
	//	out += `}`
	//	return out, true
	//}
	name := m.ctx.Name(f)
	switch r := fieldFlags.Type.(type) {
	case *flags.FieldFlags_Float:
		return m.simpleFlags(f, 0, fieldFlags.GetFloat(), wk), true
	case *flags.FieldFlags_Double:
		return m.simpleFlags(f, 0, fieldFlags.GetDouble(), wk), true
	case *flags.FieldFlags_Int32:
		return m.simpleFlags(f, 0, fieldFlags.GetInt32(), wk), true
	case *flags.FieldFlags_Int64:
		return m.simpleFlags(f, 0, fieldFlags.GetInt64(), wk), true
	case *flags.FieldFlags_Uint32:
		return m.simpleFlags(f, 0, fieldFlags.GetUint32(), wk), true
	case *flags.FieldFlags_Uint64:
		return m.simpleFlags(f, 0, fieldFlags.GetUint64(), wk), true
	case *flags.FieldFlags_Sint32:
		return m.simpleFlags(f, 0, fieldFlags.GetSint32(), wk), true
	case *flags.FieldFlags_Sint64:
		return m.simpleFlags(f, 0, fieldFlags.GetSint64(), wk), true
	case *flags.FieldFlags_Fixed32:
		return m.simpleFlags(f, 0, fieldFlags.GetFixed32(), wk), true
	case *flags.FieldFlags_Fixed64:
		return m.simpleFlags(f, 0, fieldFlags.GetFixed32(), wk), true
	case *flags.FieldFlags_Sfixed32:
		return m.simpleFlags(f, 0, fieldFlags.GetSfixed32(), wk), true
	case *flags.FieldFlags_Sfixed64:
		return m.simpleFlags(f, 0, fieldFlags.GetSfixed64(), wk), true
	case *flags.FieldFlags_Bool:
		return m.simpleFlags(f, false, fieldFlags.GetBool(), wk), true
	case *flags.FieldFlags_String_:
		return m.simpleFlags(f, `""`, fmt.Sprint(`"`, fieldFlags.GetString_(), `"`), wk), true
	case *flags.FieldFlags_Bytes:
		if wk == pgs.UnknownWKT {
			return fmt.Sprint(`
				if len(x.`, name, `) == 0 {
					x.`, name, ` = []byte("`, string(fieldFlags.GetBytes()), `")
				}`), true
		}
		return fmt.Sprint(`
				if x.`, name, ` == nil {
					x.`, name, ` = &wrapperspb.BytesValue{Value: []byte("`, string(fieldFlags.GetBytes()), `")}
				}`), true
	case *flags.FieldFlags_Enum:
		return m.simpleFlags(f, 0, fieldFlags.GetEnum(), wk), true
	case *flags.FieldFlags_Duration:
		d, err := model.ParseDuration(fieldFlags.GetDuration())
		if err != nil {
			m.Failf("invalid duration: %s %v", fieldFlags.GetDuration(), err)
		}
		return m.simpleFlags(f, `nil`, fmt.Sprint(`durationpb.New(`, int64(d), `)`), pgs.UnknownWKT), true
	case *flags.FieldFlags_Timestamp:
		v := strings.TrimSpace(fieldFlags.GetTimestamp())
		if strings.ToLower(v) == "now" {
			return m.simpleFlags(f, `nil`, `timestamppb.Now()`, pgs.UnknownWKT), true
		}
		t, err := parseTime(v)
		if err != nil {
			m.Failf("invalid timestamp: %s %v", fieldFlags.GetTimestamp(), err)
		}
		v = fmt.Sprint(`&timestamppb.Timestamp{Seconds: `, t.Unix(), `, Nanos: `, t.Nanosecond(), `}
			`)
		return m.simpleFlags(f, `nil`, v, pgs.UnknownWKT), true
	case *flags.FieldFlags_Message:
		if fieldFlags.GetMessage() != nil && !fieldFlags.GetMessage().AddFlags {
			return fmt.Sprint("\n// ", name, ": flags disabled by [(flags.value).message = {add_flags: false}]"), true
		}
		prefix := fieldFlags.GetMessage().Prefix
		if prefix == "" {
			// use field name instead, convert to kebab-case
			prefix = strings.ToLower(f.Name().String())
		}
		var decl string
		return decl + fmt.Sprint(`
			if v, ok := interface{}(x.`, name, `).(flags.Interface); ok && x.`, name, ` != nil {
				v.AddFlags(fs,`, prefix, ` )
			}`), true
	case nil: // noop
	default:
		_ = r
		m.Failf("unknown rule type (%T)", fieldFlags.Type)
	}
	return fmt.Sprint("\n// ", f.Name()), true
}

func (m *Module) simpleFlags(f pgs.Field, zero, value interface{}, wk pgs.WellKnownType) string {
	name := m.ctx.Name(f).String()
	if wk != "" && wk != pgs.UnknownWKT {
		return fmt.Sprint(`
			if x.`, name, ` == nil {
				x.`, name, ` = &wrapperspb.`, wk, `{Value: `, value, `}
			}`)
	}
	if f.HasOptionalKeyword() {
		zero = "nil"
		return fmt.Sprint(`
		if x.`, name, ` == `, zero, ` {
			v := `, m.ctx.Type(f).Value(), `(`, value, `)
			x.`, name, ` = &v 
		}`)
	}
	return fmt.Sprint(`
		if x.`, name, ` == `, zero, ` {
			x.`, name, ` = `, value, `
		}`)
}

func parseTime(s string) (time.Time, error) {
	for _, format := range []string{
		time.RFC822,
		time.RFC822Z,
		time.RFC850,
		time.RFC1123,
		time.RFC1123Z,
		time.RFC3339,
	} {
		t, err := time.Parse(format, s)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, errors.New("cannot parse timestamp, timestamp supported format: RFC822 / RFC822Z / RFC850 / RFC1123 / RFC1123Z / RFC3339")
}

func isOk(b []bool) bool {
	return len(b) > 0 && b[0]
}
