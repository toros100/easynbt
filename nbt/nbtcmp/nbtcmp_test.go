package nbtcmp

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/toros100/easynbt/nbt"
)

func TestCmpOption(t *testing.T) {

	opt := NBTOptionTypeCmpOption()

	o1 := nbt.Option[int8]{
		Value: 7,
		Ok:    false,
	}

	o2 := nbt.Option[int8]{
		Value: 0,
		Ok:    false,
	}

	if cmp.Equal(o1, o2) {
		t.Fatalf("expected structs to not be equal without the option")
	}

	if !cmp.Equal(o1, o2, opt) {
		t.Fatalf("expected structs to be equal with the option")
	}
}
