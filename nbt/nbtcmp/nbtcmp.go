package nbtcmp

import (
	"github.com/google/go-cmp/cmp"
	"github.com/toros100/easynbt/nbt"
)

// NBTOptionTypeCmpOption returns a cmp.Option that can be used by github.com/google/go-cmp/cmp (in tests!)
// to compare nbt.Option[T] types and makes two values of type nbt.Option[T] satisfy cmp.Equal,
// if both Ok fields are false, ignoring the value fields.
func NBTOptionTypeCmpOption() cmp.Option {
	// great function name

	return cmp.FilterPath(func(p cmp.Path) bool {
		if len(p) < 2 {
			// can not possibly be inside nbt.Option[T] struct
			return false
		}

		parent := p[len(p)-2]

		if !nbt.IsOptionType(parent.Type()) {
			return false
		}

		v1, v2 := parent.Values()

		// field with name "Ok" is guaranteed to be valid
		// (unless nbt.Option[T] is redefined, so if the outer function panics, that must be why)
		if v1.FieldByName("Ok").IsZero() && v2.FieldByName("Ok").IsZero() {
			// because the Ok field is bool, this condition is equivalent to both bools being false
			// and in that case, we ignore this path
			// technically we have not checked whether we are on the Ok or Value or any other field inside
			// the nbt.Option[T] struct, but this works out either way.
			return true
		}

		return false

	}, cmp.Ignore())
}
