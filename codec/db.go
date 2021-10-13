package codec

import (
	"reflect"
)

var (
	_interfDB map[string][]reflect.Type = make(map[string][]reflect.Type) // interf name -> [impl type, ...]
	_implDB   map[string]int            = make(map[string]int)            // struct name -> id
)

// returns dummy bool, so we can call function outside main()
func RegisterInterface(interf interface{}, impls ...interface{}) bool {
	interfName := reflect.TypeOf(interf).Elem().Name()

	lst := make([]reflect.Type, len(impls))

	for id, impl := range impls {
		t := reflect.TypeOf(impl).Elem()

		name := t.Name()

		if name == "" {
			panic("empty impl name during registration")
		}

		if _, ok := _implDB[name]; ok {
			panic("impl " + name + " already registered")
		}

		_implDB[name] = id

		lst[id] = t

		_interfDB[interfName] = lst
	}

	return true
}
