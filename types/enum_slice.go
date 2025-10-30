package types

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var _ pflag.Value = (*EnumSliceValue)(nil)

type EnumSliceValue struct {
	value        reflect.Value // reflect.Value of the slice pointer
	changed      bool
	allowedTypes []string // List of valid enum value names
	typ          protoreflect.EnumType
	descriptors  protoreflect.EnumValueDescriptors // Enum value descriptors for validation
}

func (s *EnumSliceValue) parseEnum(val string) (protoreflect.Enum, error) {
	val = strings.TrimSpace(val)
	if enumVal := s.descriptors.ByName(protoreflect.Name(val)); enumVal != nil {
		return s.typ.New(enumVal.Number()), nil
	}
	if numVal, err := strconv.Atoi(val); err == nil {
		if enumVal := s.descriptors.ByNumber(protoreflect.EnumNumber(numVal)); enumVal != nil {
			return s.typ.New(enumVal.Number()), nil
		}
	}
	return nil, fmt.Errorf("invalid enum value:%s, allowed values are: %s", val, strings.Join(s.allowedTypes, ","))
}

func (s *EnumSliceValue) Set(val string) error {
	// Get the slice we're working with
	slice := s.value.Elem()
	val = strings.TrimSpace(val)

	if !s.changed {
		newSlice := reflect.MakeSlice(slice.Type(), 0, 1)
		slice.Set(newSlice)
		// Try to find the enum by name
		newEnum, err := s.parseEnum(val)
		if err != nil {
			return err
		}
		slice.Set(reflect.Append(slice, reflect.ValueOf(newEnum)))
		s.changed = true
	} else {
		// Try to find the enum by name
		newEnum, err := s.parseEnum(val)
		if err != nil {
			return err
		}
		slice.Set(reflect.Append(slice, reflect.ValueOf(newEnum)))
	}
	return nil
}

func (s *EnumSliceValue) Append(val string) error {
	// Get the slice we're working with
	slice := s.value.Elem()
	newEnum, err := s.parseEnum(val)
	if err != nil {
		return err
	}
	slice.Set(reflect.Append(slice, reflect.ValueOf(newEnum)))
	return nil
}

func (s *EnumSliceValue) Replace(val []string) error {
	// Get the slice we're working with
	slice := s.value.Elem()
	newSlice := reflect.MakeSlice(slice.Type(), 0, len(val))
	for _, d := range val {
		// Try to find the enum by name
		newEnum, err := s.parseEnum(d)
		if err != nil {
			return err
		}
		newSlice.Set(reflect.Append(newSlice, reflect.ValueOf(newEnum)))
	}
	slice.Set(newSlice)
	return nil
}

func (s *EnumSliceValue) GetSlice() []string {
	// Get the slice we're working with
	slice := s.value.Elem()

	if !slice.IsValid() || slice.Len() == 0 {
		return []string{}
	}

	out := make([]string, slice.Len())
	for i := 0; i < slice.Len(); i++ {
		e := slice.Index(i).Interface().(protoreflect.Enum)
		ev := s.descriptors.ByNumber(e.Number())
		if ev != nil {
			out[i] = string(ev.Name())
		} else {
			out[i] = fmt.Sprintf("%d", e.Number())
		}
	}
	return out
}

// Type returns a string that uniquely represents this flag's type.
func (s *EnumSliceValue) Type() string {
	return "enumSliceValue"
}

// String defines a "native" format for this enum slice flag value.
func (s *EnumSliceValue) String() string {
	if !s.value.IsValid() || s.value.Elem().Len() == 0 {
		return "[]"
	}

	slice := s.value.Elem()
	enumStrSlice := make([]string, slice.Len())
	for i := 0; i < slice.Len(); i++ {
		e := slice.Index(i).Interface().(protoreflect.Enum)
		ev := s.descriptors.ByNumber(e.Number())
		if ev != nil {
			enumStrSlice[i] = string(ev.Name())
		} else {
			enumStrSlice[i] = fmt.Sprintf("%d", e.Number())
		}
	}

	return "[" + strings.Join(enumStrSlice, ",") + "]"
}

func EnumSlice(arrPtr any) *EnumSliceValue {
	if arrPtr == nil {
		panic("nil array pointer")
	}

	v := reflect.ValueOf(arrPtr)

	if v.Kind() != reflect.Pointer {
		panic("argument must be a pointer to enum slice")
	}

	elem := v.Elem()
	if !elem.CanSet() {
		panic("cannot set enum slice value (ensure pointer is to a settable field)")
	}

	// Get the slice type and element type
	sliceType := elem.Type()
	if sliceType.Kind() != reflect.Slice {
		panic("argument must be a pointer to a slice")
	}

	itemType := sliceType.Elem()

	// Create a sample enum value to get type information
	sampleVal := reflect.New(itemType).Interface()
	enumVal, ok := sampleVal.(protoreflect.Enum)
	if !ok {
		panic("slice element type must implement protoreflect.Enum interface")
	}

	values := enumVal.Descriptor().Values()
	allowed := make([]string, 0, values.Len())
	for i := 0; i < values.Len(); i++ {
		allowed = append(allowed, string(values.Get(i).Name()))
	}

	return &EnumSliceValue{
		value:        v,
		changed:      false,
		typ:          enumVal.Type(),
		descriptors:  values,
		allowedTypes: allowed,
	}
}
