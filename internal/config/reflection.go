package config

import (
	"errors"
	"fmt"
	"net"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	cfTag      = "cf"
	defaultTag = "default"
)

var (
	fieldNotFound = errors.New("field not found")
	zeroValue     = reflect.Value{}
)

func isPointerToStruct(ps interface{}) bool {
	t := reflect.TypeOf(ps)
	if t.Kind() != reflect.Pointer {
		return false
	}

	s := t.Elem()
	return s != nil && s.Kind() == reflect.Struct
}

func checkIfAssignable(target reflect.Value, data interface{}) error {
	if !target.CanSet() {
		return errors.New("target is not settable")
	}

	dataType := reflect.TypeOf(data)
	valueType := target.Type()
	if !dataType.AssignableTo(valueType) {
		return errors.New("target is not assignable due the type mishmash")
	}

	return nil
}

func assignTo(target reflect.Value, data interface{}) error {
	err := checkIfAssignable(target, data)
	if err == nil {
		dataVal := reflect.ValueOf(data)
		target.Set(dataVal)
	}
	return err
}

func assignToField(structVal reflect.Value, fieldName string, data interface{}) error {
	field := structVal.FieldByName(fieldName)
	if field != zeroValue {
		return assignTo(field, data)
	}
	return fieldNotFound
}

func findFieldByTag(structVal reflect.Value, name string) reflect.Value {
	t := structVal.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if tag, ok := field.Tag.Lookup(cfTag); ok {
			if ok && tag == name {
				return structVal.Field(i)
			}
		}
	}
	return zeroValue
}

func assignFromString(value reflect.Value, input string) error {
	switch value.Kind() {
	case reflect.String:
		value.SetString(input)
	case reflect.Int:
		fallthrough
	case reflect.Int64:
		if value.Type().Name() == "Duration" {
			durationValue, err := time.ParseDuration(input)
			if err != nil {
				return err
			}
			value.SetInt(int64(durationValue))
		} else {
			intValue, err := strconv.Atoi(input)
			if err != nil {
				return err
			}
			value.SetInt(int64(intValue))
		}
	case reflect.Float64:
		floatValue, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return err
		}
		value.SetFloat(floatValue)
	case reflect.Bool:
		boolValue, err := strconv.ParseBool(input)
		if err != nil {
			return err
		}
		value.SetBool(boolValue)
	case reflect.Slice:
		if value.Type().Elem().Kind() == reflect.String {
			stringSlice := strings.Split(input, ",")
			value.Set(reflect.ValueOf(stringSlice))
		} else if value.Type().Elem().Kind() == reflect.Int {
			intSlice := []int{}
			for _, str := range strings.Split(input, ",") {
				intValue, err := strconv.Atoi(str)
				if err != nil {
					return err
				}
				intSlice = append(intSlice, intValue)
			}
			value.Set(reflect.ValueOf(intSlice))
		} else if value.Type().Elem().Kind() == reflect.Uint8 { // net.IP
			// TODO: is there any other way how to determine net.IP type?
			ip := net.ParseIP(input)
			value.Set(reflect.ValueOf(ip))
		}
	default:
		return fmt.Errorf("unsupported type: %v", value.Type())
	}
	return nil
}
