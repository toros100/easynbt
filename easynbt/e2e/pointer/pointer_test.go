package pointer

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/toros100/easynbt/nbt"
)

func TestPointerList(t *testing.T) {
	testListHelper(t, []int8{1})
	testListHelper(t, []int8{1, 4, 5, 54, 34, 4, 34, 5, 45, 55, 43, 5})
	testListHelper(t, []int8{1, 23, 4, 54, 54, 4, 2, 43, 43, 5, 34, 43,
		43, 3, 32, 5, 23, 3, 32, 32, 76, 85, 64, 54, 25, 43})
}

func testListHelper(t *testing.T, bs []int8) {
	t.Helper()

	buf := bytes.NewBuffer(nil)
	writeList(buf, bs)

	data := buf.Bytes()

	g := new(List)

	_, k, err := nbt.Unmarshal(g, data)
	if err != nil {
		t.Fatal(err)
	}

	if k != len(data) {
		t.Fatal()
	}

	s := getListStruct(bs)

	if !cmp.Equal(s, g) {
		t.Fatal("unmarshalled not equal to expected")
	}
}

// len(bs) > 0 required
func writeList(buf *bytes.Buffer, bs []int8) {
	if len(bs) == 0 {
		panic("programmer error")
	}
	buf.Write([]byte{nbt.TagCompound, 0, 0})
	writeListInner(buf, bs)
	buf.WriteByte(nbt.TagEnd)
}

func writeListInner(buf *bytes.Buffer, bs []int8) {
	switch len(bs) {
	case 1:
		buf.Write([]byte{nbt.TagByte, 0, 4, 'H', 'e', 'a', 'd', byte(bs[0])})
	default:
		buf.Write([]byte{nbt.TagByte, 0, 4, 'H', 'e', 'a', 'd', byte(bs[0]),
			nbt.TagCompound, 0, 4, 'T', 'a', 'i', 'l'})
		writeListInner(buf, bs[1:])
		buf.WriteByte(nbt.TagEnd)
	}
}

// len(bs) > 0 required
func getListStruct(bs []int8) *List {
	switch len(bs) {
	case 0:
		panic("programmer error")
	case 1:
		return &List{
			Head: int8(bs[0]),
			Tail: nbt.None[*List](),
		}
	default:
		return &List{
			Head: int8(bs[0]),
			Tail: nbt.Some(getListStruct(bs[1:])),
		}

	}
}
