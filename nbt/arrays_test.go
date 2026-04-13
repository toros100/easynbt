package nbt

import (
	"bytes"
	"encoding/binary"
	"errors"
	"math"
	"slices"
	"testing"
)

func TestTagType(t *testing.T) {
	if (*Array[int8])(nil).TagType() != TagByteArray {
		t.Fail()
	}
	if (*Array[int32])(nil).TagType() != TagIntArray {
		t.Fail()
	}
	if (*Array[int64])(nil).TagType() != TagLongArray {
		t.Fail()
	}
}

func TestByteArray(t *testing.T) {
	testHelper(t, []int8{})
	testHelper(t, []int8{-3, 4, 5, 6, 7})
}
func TestIntArray(t *testing.T) {
	testHelper(t, []int32{})
	testHelper(t, []int32{10230, 24334, -2341324, 23})
}

func TestLongArray(t *testing.T) {
	testHelper(t, []int64{})
	testHelper(t, []int64{430334, 444444, 0, -348983428})
}

func testHelper[T int8 | int32 | int64](t *testing.T, arr []T) {
	t.Helper()

	data := writeArray(arr)

	a := new(Array[T])

	k, err := a.UnmarshalPayload(data)

	if err != nil {
		t.Fatalf("expected err nil but found %v", err)
	}

	if k != len(data) {
		t.Fatal("did not consume the expected number of bytes")
	}

	if !slices.Equal(*a, arr) {
		t.Fatal("unmarshaled array values do not match expected values")
	}

	if len(data) == 0 {
		return
	}

	dataShort := data[:len(data)-1]

	a2 := new(Array[T])

	_, err = a2.UnmarshalPayload(dataShort)

	if !errors.Is(err, ErrUnexpectedEOF) {
		t.Fatalf("expected err %v but found %v", ErrUnexpectedEOF, err)
	}

}

func writeArray[T int8 | int32 | int64](arr []T) []byte {
	if len(arr) > math.MaxInt32 {
		panic("bad test: arr too long")
	}

	l := uint32(len(arr))
	buf := bytes.NewBuffer(nil)
	_ = binary.Write(buf, binary.BigEndian, l)
	_ = binary.Write(buf, binary.BigEndian, arr)
	return buf.Bytes()
}
