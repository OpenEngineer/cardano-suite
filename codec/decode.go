package codec

import (
	"errors"
	"reflect"

	cbor "github.com/fxamacker/cbor/v2"
)

// Deserialize a Cardano-specific CBOR object into a Golang value
func FromCBOR(b []byte, dst interface{}) error {
	var src interface{} = nil

	if err := cbor.Unmarshal(b, &src); err != nil {
		return err
	}

	return FromUntyped(src, dst)
}

func FromUntyped(src interface{}, dst interface{}) error {
	return fromUntyped(reflect.ValueOf(src), reflect.ValueOf(dst))
}

func fromUntyped(src reflect.Value, dst reflect.Value) error {
	src = indirect(src)
	dst = indirect(dst)

	srcT := src.Type()
	dstT := dst.Type()

	if dstT.Name() == "" {
		panic("empty dstT")
	}

	if !dst.CanSet() {
		return errors.New(dstT.Name() + " unsettable")
	}

	switch dstT.Kind() {
	case reflect.Interface:
		switch srcT.Kind() {
		case reflect.Slice:
			interfName := dstT.Name()

			implTypes, ok := _interfDB[interfName]
			if !ok {
				return errors.New("interface " + interfName + " not registered")
			}

			if src.Len() < 1 {
				return errors.New("expected slice larger than 0")
			}

			id_ := src.Index(0).Elem()

			var id int
			switch id_.Type().Kind() {
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				id = int(id_.Uint())
			default:
				return errors.New("type index isn't an uint " + id_.Type().Name())
			}

			if id >= len(implTypes) {
				return errors.New("type id out of range")
			}

			implType := implTypes[id]

			return structFromCBOR(1, src, implType, dst)
		default:
			return errors.New("expected slice")
		}
	case reflect.Struct:
		switch srcT.Kind() {
		case reflect.Slice:
			return structFromCBOR(0, src, dstT, dst)
		default:
			return errors.New("expected slice")
		}
	case reflect.Array:
		switch srcT.Kind() {
		case reflect.Array, reflect.Slice:
			n := dst.Len()

			if n != src.Len() {
				return errors.New("wrong array length")
			}

			arrayT := reflect.ArrayOf(n, dstT.Elem())

			tmp := reflect.New(arrayT)

			for i := 0; i < n; i++ {
				item := indirect(tmp).Index(i).Addr()

				if err := fromUntyped(src.Index(i), item); err != nil {
					return err
				}
			}
		default:
			return errors.New("expected array or slice")
		}
	case reflect.Slice:
		switch srcT.Kind() {
		case reflect.Slice:
			n := src.Len()

			sliceT := reflect.SliceOf(dstT.Elem())

			tmp := reflect.MakeSlice(sliceT, 0, n)

			for i := 0; i < n; i++ {
				item := reflect.New(dstT.Elem())

				if err := fromUntyped(src.Index(i), item); err != nil {
					return err
				}

				tmp = reflect.Append(tmp, indirect(item))
			}

			dst.Set(tmp)
		default:
			return errors.New("expected slice")
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch srcT.Kind() {
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			dst.Set(src)
		default:
			return errors.New("not an unsigned int")
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch srcT.Kind() {
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			dst.Set(reflect.ValueOf(int(src.Uint())))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			dst.Set(src)
		default:
			return errors.New("not an int")
		}
	case reflect.Bool:
		switch srcT.Kind() {
		case reflect.Bool:
			dst.Set(src)
		default:
			return errors.New("not a bool")
		}
	case reflect.String:
		switch srcT.Kind() {
		case reflect.String:
			dst.Set(src)
		default:
			return errors.New("not a string")
		}
	default:
		return errors.New("Internal error: \"" + dstT.Name() + "\" unhandled")
	}

	return nil
}

// dstType is a concrete type
// src must be a slice (otherwise this function panics)
func structFromCBOR(offset int, src reflect.Value, dstT reflect.Type, dst reflect.Value) error {
	tmp := reflect.New(dstT)

	n := dstT.NumField()

Outer:
	for i := 0; i < n; i++ {
		dstF := indirect(tmp).Field(i)
		tag := dstT.Field(i).Tag

		dstFT := indirectType(dstT.Field(i).Type)

		setDstF := func(srcF reflect.Value) error {
			// XXX: for some reason we must use reflect.ValueOf(dstF.Addr().Interface()) to get the underlying struct in case dstF is actually an interface
			dstFTmp := reflect.New(dstFT)

			if err := fromUntyped(srcF, reflect.ValueOf(dstFTmp.Interface())); err != nil {
				return err
			}

			dstF_ := indirect(dstF.Addr())

			if dstF_.Kind() == reflect.Ptr || dstF_.Kind() == reflect.Interface {
				dstF_.Set(dstFTmp)
			} else {
				dstF_.Set(indirect(dstFTmp))
			}

			return nil
		}

		switch {
		case tag == "wrapped":
			if i >= src.Len()-offset {
				return errors.New("not enough fields")
			}

			srcF := src.Index(i + offset)

			x, ok := srcF.Interface().(cbor.Tag)
			if !ok {
				return errors.New("not a cbor.Tag")
			}

			if x.Number != 24 {
				return errors.New("not a cbor.Tag{24})")
			}

			data, ok := x.Content.([]byte)
			if !ok {
				return errors.New("cbor.Tag content not a byte list")
			}

			if err := cbor.Unmarshal(data, srcF.Addr().Interface()); err != nil {
				return err
			}

			if err := setDstF(srcF); err != nil {
				return err
			}
		case tag == "emptyiftrue":
			if i != 0 {
				return errors.New("emptyiftrue can only be used for first tag")
			}

			if src.Len() == 0 {
				srcF := reflect.ValueOf(true)

				if err := setDstF(srcF); err != nil {
					return err
				}

				break Outer
			} else {
				srcF := reflect.ValueOf(false)
				offset--

				if err := setDstF(srcF); err != nil {
					return err
				}
			}

		case tag == "":
			if i >= src.Len()-offset {
				return errors.New("not enough fields")
			}

			srcF := src.Index(i + offset)

			if err := setDstF(srcF); err != nil {
				return err
			}
		default:
			return errors.New("unrecognized struct tag \"" + string(tag) + "\"")
		}
	}

	dst.Set(indirect(tmp))

	return nil
}
