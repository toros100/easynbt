package nbt

import (
	"fmt"
	"reflect"
)

type optionIdentity struct{}

type Option[T any] struct {
	// The blank field with an arbitrary (but unique) unexported type serves as a marker and
	// allows the Option type to be "recognized", even when inspecting the instantiated
	// underlying struct type (without type name information)
	_     optionIdentity
	Value T
	Ok    bool
}

func (o Option[T]) String() string {
	if o.Ok {
		return fmt.Sprintf("Some[%s](%v)", reflect.TypeFor[T]().String(), o.Value)
	} else {
		return fmt.Sprintf("None[%s]()", reflect.TypeFor[T]().String())
	}
}

func Some[T any](val T) Option[T] {
	return Option[T]{
		Value: val,
		Ok:    true,
	}
}

func None[T any]() Option[T] {
	return Option[T]{}
}

func IsOptionType(typ reflect.Type) bool {
	if typ.Kind() != reflect.Struct {
		return false
	}

	if typ.NumField() != 3 {
		return false
	}

	return typ.Field(0).Type == reflect.TypeFor[optionIdentity]()
}
