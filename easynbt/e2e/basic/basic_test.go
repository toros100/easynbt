package basic

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/toros100/easynbt/nbt"
)

// TODO: tests for remaining types

func expectTagType(t *testing.T, u nbt.Unmarshaller, tagType byte) {
	t.Helper()

	if tt := u.TagType(); tt != tagType {
		t.Fatalf("expected tag type %v on %T but found %v", tagType, u, tt)
	}
}

func TestTagType(t *testing.T) {
	expectTagType(t, (*BytePayload)(nil), nbt.TagByte)
	expectTagType(t, (*ShortPayload)(nil), nbt.TagShort)
	expectTagType(t, (*IntPayload)(nil), nbt.TagInt)
	expectTagType(t, (*LongPayload)(nil), nbt.TagLong)
	expectTagType(t, (*FloatPayload)(nil), nbt.TagFloat)
	expectTagType(t, (*DoublePayload)(nil), nbt.TagDouble)
	expectTagType(t, (*ListOfIntPayload)(nil), nbt.TagList)
	expectTagType(t, (*CompoundPayload)(nil), nbt.TagCompound)
	expectTagType(t, (*StringPayload)(nil), nbt.TagString)
	expectTagType(t, (*nbt.ByteArray)(nil), nbt.TagByteArray)
	expectTagType(t, (*nbt.IntArray)(nil), nbt.TagIntArray)
	expectTagType(t, (*nbt.LongArray)(nil), nbt.TagLongArray)
}

func TestUnmarshalByte(t *testing.T) {
	data := []byte{3}

	p := new(BytePayload)
	k, err := p.UnmarshalPayload(data)
	expectErr(t, err, nil)
	expectEq(t, k, len(data), "bytes read")
	expectEq(t, *p, 3, "value")

	_, err = p.UnmarshalPayload([]byte{})
	expectErr(t, err, nbt.ErrUnexpectedEOF)
}

func TestUnmarshalShort(t *testing.T) {
	data := []byte{127, 255}

	p := new(ShortPayload)
	k, err := p.UnmarshalPayload(data)
	expectErr(t, err, nil)
	expectEq(t, k, len(data), "bytes read")
	expectEq(t, *p, 32767, "value")

	_, err = p.UnmarshalPayload(data[1:])
	expectErr(t, err, nbt.ErrUnexpectedEOF)
}

func TestUnmarshalString(t *testing.T) {
	data := []byte{0, 12, 'h', 'e', 'l', 'l', 'o', ',', ' ', 'w', 'o', 'r', 'l', 'd'}

	p := new(StringPayload)
	k, err := p.UnmarshalPayload(data)
	expectErr(t, err, nil)
	expectEq(t, k, 14, "bytes read")
	expectEq(t, *p, "hello, world", "value")
}

func TestUnmarshalList(t *testing.T) {

	data := []byte{nbt.TagInt, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 1, 0, 255, 255, 255, 255}

	p := new(ListOfIntPayload)

	k, err := p.UnmarshalPayload(data)
	expectErr(t, err, nil)
	expectEq(t, k, len(data), "bytes read")
	expectEq(t, *p, []int32{0, 256, -1}, "v")
}

func TestUnmarshalByteArray(t *testing.T) {
	data := []byte{0, 0, 0, 8, 0, 1, 2, 3, 4, 5, 6, 7}

	p := new(nbt.ByteArray)

	k, err := p.UnmarshalPayload(data)

	expectErr(t, err, nil)
	expectEq(t, k, len(data), "bytes read")

	expected := new(nbt.ByteArray).FromSlice([]int8{0, 1, 2, 3, 4, 5, 6, 7})

	expectEq(t, p, expected, "array data")
}

func TestUnmarshalIntArray(t *testing.T) {
	data := []byte{0, 0, 0, 2, 0, 1, 2, 3, 4, 5, 6, 7}

	p := new(nbt.IntArray)

	k, err := p.UnmarshalPayload(data)

	expectErr(t, err, nil)
	expectEq(t, k, len(data), "bytes read")

	expected := new(nbt.IntArray).FromSlice([]int32{66051, 67438087})

	expectEq(t, p, expected, "array data")
}

func TestUnmarshalLongArray(t *testing.T) {
	data := []byte{0, 0, 0, 2, 0, 1, 2, 3, 4, 5, 6, 7, 0, 0, 0, 0, 0, 0, 0, 1}

	p := new(nbt.LongArray)

	k, err := p.UnmarshalPayload(data)

	expectErr(t, err, nil)
	expectEq(t, k, len(data), "bytes read")

	expected := new(nbt.LongArray).FromSlice([]int64{283686952306183, 1})

	expectEq(t, p, expected, "array data")
}

func expectEq[T any](t *testing.T, x, y T, message string) {
	t.Helper()
	if !cmp.Equal(x, y) {
		t.Fatalf("%s: %v != %v", message, x, y)
	}
}

func expectErr(t *testing.T, actual, expected error) {
	t.Helper()

	if !errors.Is(actual, expected) {
		t.Fatalf("expected error %v but found %v", expected, actual)
	}
}
