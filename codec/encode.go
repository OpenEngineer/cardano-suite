package codec

import (
	"reflect"

	cbor "github.com/fxamacker/cbor/v2"
)

// Serialize a Golang value into a Cardano-specific CBOR encoded object
func ToCBOR(x interface{}) []byte {
	untyped := ToUntyped(x)

	b, err := cbor.Marshal(untyped)
	if err != nil {
		panic(err)
	}

	return b
}

func ToUntyped(x interface{}) interface{} {
	return toUntyped(indirect(reflect.ValueOf(x)))
}

func toUntyped(x reflect.Value) interface{} {
	x = indirect(x)
	t := x.Type()

	switch t.Kind() {
	case reflect.Interface, reflect.Ptr:
		panic("should've been indirected")

		/*structName := t.Name()

		id, ok := _implDB[structName]
		if !ok {
			panic("\"" + structName + "\" not registered")
		}

		return structToUntyped(id, x)*/
	case reflect.Struct:
		structName := t.Name()

		if id, ok := _implDB[structName]; ok {
			return structToUntyped(id, x)
		} else {
			return structToUntyped(-1, x)
		}
	case reflect.Slice:
		n := x.Len()

		y := make([]interface{}, n)

		for i := 0; i < n; i++ {
			y[i] = toUntyped(x.Index(i))
		}

		return y
	default:
		return x.Interface()
	}
}

func structToUntyped(id int, x reflect.Value) interface{} {
	y := make([]interface{}, 0)

	if id >= 0 {
		y = append(y, id)
	}

	// the struct fields are marshalled consequtively into the array

	n := x.NumField()
	t := x.Type()

	for i := 0; i < n; i++ {
		// XXX: a trick to get the underlying struct in case the field is actually an interface
		f := reflect.ValueOf(x.Field(i).Interface())

		tag := t.Field(i).Tag

		switch {
		case tag == "wrapped":
			b, err := cbor.Marshal(f.Interface())
			if err != nil {
				panic(err)
			}

			y = append(y, cbor.Tag{24, b})
		case tag == "emptyiftrue":
			if i != 0 {
				panic("\"emptyiftrue\" tag can only be used for first field of struct")
			}

			if f.Bool() {
				return y
			}
		case tag != "":
			panic("unrecognised tag " + tag)
		default:
			item := toUntyped(f)

			y = append(y, item)
		}
	}

	return y
}
