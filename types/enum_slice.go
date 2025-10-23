package types

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
	"google.golang.org/protobuf/reflect/protoreflect"

	"github.com/kunstack/protoc-gen-flags/utils"
)

var _ pflag.Value = (*EnumSliceValue)(nil)

type EnumSliceValue struct {
	value        reflect.Value // reflect.Value of the slice pointer
	itemType     reflect.Type
	changed      bool
	allowedTypes []string // List of valid enum value names
	typ          protoreflect.EnumType
	descriptors  protoreflect.EnumValueDescriptors // Enum value descriptors for validation
}

// Set converts, and assigns, the comma-separated enum argument string representation as the []protoreflect.Enum value of this flag.
// If Set is called on a flag that already has a []protoreflect.Enum assigned, the newly converted values will be appended.
func (s *EnumSliceValue) Set(val string) error {
	if s == nil {
		return fmt.Errorf("enum slice value is nil or invalid")
	}

	if !s.value.IsValid() {
		return errors.New("enum slice value is not initialized")
	}

	// remove all quote characters
	rmQuote := strings.NewReplacer(`"`, "", `'`, "", "`", "")

	// Parse the comma-separated values
	vals, err := utils.ReadAsCSV(rmQuote.Replace(val))
	if err != nil {
		return err
	}

	// Get the slice we're working with
	slice := s.value.Elem()

	// Clear the slice if this is the first call (not an append)
	if !s.changed {
		newSlice := reflect.MakeSlice(slice.Type(), 0, len(vals))
		slice.Set(newSlice)
		s.changed = true
	}

	// Parse each enum value
	for _, strVal := range vals {
		// Try to find the enum by name
		if enumVal := s.descriptors.ByName(protoreflect.Name(strVal)); enumVal != nil {
			newEnum := s.typ.New(enumVal.Number())
			slice.Set(reflect.Append(slice, reflect.ValueOf(newEnum)))
			continue
		}

		// Try to parse as a number
		if numVal, err := strconv.Atoi(strVal); err == nil {
			if enumVal := s.descriptors.ByNumber(protoreflect.EnumNumber(numVal)); enumVal != nil {
				newEnum := s.typ.New(enumVal.Number())
				slice.Set(reflect.Append(slice, reflect.ValueOf(newEnum)))
				continue
			}
			return fmt.Errorf("invalid enum number %q, allowed values are: %s", strVal, strings.Join(s.allowedTypes, ", "))
		}
		return fmt.Errorf("invalid enum value %q, allowed values are: %s", strVal, strings.Join(s.allowedTypes, ", "))
	}

	return nil
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
	out, _ := utils.WriteAsCSV(enumStrSlice)

	return "[" + out + "]"
}

func EnumSlice(arrPtr any) *EnumSliceValue {
	if arrPtr == nil {
		panic("nil array pointer")
	}

	v := reflect.ValueOf(arrPtr)

	if v.Kind() != reflect.Ptr {
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
	if enumVal, ok := sampleVal.(protoreflect.Enum); ok {
		values := enumVal.Descriptor().Values()
		allowed := make([]string, 0, values.Len())
		for i := 0; i < values.Len(); i++ {
			allowed = append(allowed, string(values.Get(i).Name()))
		}

		return &EnumSliceValue{
			value:        v, // Store the reflect.Value of the slice pointer
			itemType:     itemType,
			changed:      false,
			typ:          enumVal.Type(),
			descriptors:  values,
			allowedTypes: allowed,
		}
	}

	panic("slice element type must implement protoreflect.Enum interface")
}
