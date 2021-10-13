package codec

import "reflect"

// the reflect.Indirect only works for pointers, but we want the same behavior for interfaces
func indirect(v reflect.Value) reflect.Value {
	switch v.Kind() {
	case reflect.Ptr:
		return v.Elem()
	case reflect.Interface:
		return v.Elem()
	default:
		return v
	}
}

func indirectType(t reflect.Type) reflect.Type {
	switch t.Kind() {
	case reflect.Ptr:
		return t.Elem()
	case reflect.Interface:
		return t.Elem()
	default:
		return t
	}
}
