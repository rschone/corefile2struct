package config

import (
	"reflect"
)

func isPointerToStruct(ps interface{}) bool {
	t := reflect.TypeOf(ps)
	if t.Kind() != reflect.Pointer {
		return false
	}

	s := t.Elem()
	return s != nil && s.Kind() == reflect.Struct
}
