package basic

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/toros100/easynbt/nbt"
)

// TODO: tests for remaining types

func expectTagType(t *testing.T, u nbt.Unmarshaler, tagType byte) {
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
