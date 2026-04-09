package option

import (
	"encoding/binary"
	"errors"
	"github.com/google/go-cmp/cmp"
	"github.com/toros100/easynbt/nbt"
	"github.com/toros100/easynbt/nbt/nbtcmp"
	"math"
	"testing"
)

func TestOptionalPresent(t *testing.T) {

	data := []byte{
		nbt.TagCompound,
		0, 4,
		'r', 'o', 'o', 't',
		nbt.TagInt,
		0, 18,
		'O', 'p', 't', 'i', 'o', 'n', 'a', 'l', 'I', 'n', 't', 'V', 'i', 'a', 'T', 'y', 'p', 'e',
		0, 0, 0, 123,
		nbt.TagInt,
		0, 17,
		'O', 'p', 't', 'i', 'o', 'n', 'a', 'l', 'I', 'n', 't', 'V', 'i', 'a', 'T', 'a', 'g',
		0, 0, 0, 99,
		nbt.TagString,
		0, 10,
		'S', 'o', 'm', 'e', 'S', 't', 'r', 'i', 'n', 'g',
		0, 5,
		'H', 'E', 'L', 'L', 'O',
		nbt.TagEnd,
	}

	expected := SomeStruct{
		OptionalIntViaType: nbt.Some[int32](123),
		OptionalIntViaTag:  99,
		SomeString:         "HELLO",
	}

	s := SomeStruct{}

	name, k, err := nbt.Unmarshal(&s, data)

	if err != nil {
		t.Fatalf("expected err to be nil, but found %v", err)
	}

	if name != "root" {
		t.Fatalf("expected name to be %q, but found %q", "root", name)
	}

	if k != len(data) {
		t.FailNow()
	}

	if s != expected {
		t.Fatalf("expected %+v == %+v", s, expected)
	}

}

func TestOptionalAbsent(t *testing.T) {

	data := []byte{
		nbt.TagCompound,
		0, 4,
		'r', 'o', 'o', 't',
		nbt.TagString,
		0, 10,
		'S', 'o', 'm', 'e', 'S', 't', 'r', 'i', 'n', 'g',
		0, 5,
		'H', 'E', 'L', 'L', 'O',
		nbt.TagEnd,
	}

	expected := SomeStruct{
		SomeString:         "HELLO",
		OptionalIntViaType: nbt.None[int32](), // = nbt.Option[int32]{Value: 0, Ok: false}
		OptionalIntViaTag:  0,
	}

	s := SomeStruct{}

	name, k, err := nbt.Unmarshal(&s, data)

	if err != nil {
		t.Fatalf("expected err to be nil, but found %v", err)
	}

	if name != "root" {
		t.Fatalf("expected name to be %q, but found %q", "root", name)
	}

	if k != len(data) {
		t.FailNow()
	}

	if s != expected {
		t.Fatalf("expected %+v == %+v", s, expected)
	}
}

func TestRequiredAbsent(t *testing.T) {

	data := []byte{
		nbt.TagCompound,
		0, 4,
		'r', 'o', 'o', 't',
		nbt.TagInt,
		0, 18,
		'O', 'p', 't', 'i', 'o', 'n', 'a', 'l', 'I', 'n', 't', 'V', 'i', 'a', 'T', 'y', 'p', 'e',
		0, 0, 0, 123,
		nbt.TagInt,
		0, 17,
		'O', 'p', 't', 'i', 'o', 'n', 'a', 'l', 'I', 'n', 't', 'V', 'i', 'a', 'T', 'a', 'g',
		0, 0, 0, 99,
		nbt.TagEnd,
	}

	s := SomeStruct{}

	_, _, err := nbt.Unmarshal(&s, data)

	if !errors.Is(err, nbt.ErrMissingValue) {
		t.Fatalf("expected err to be %v, but found %v", nbt.ErrMissingValue, err)
	}

}

func TestAllOptionTypes(t *testing.T) {

	a := AllTypesStruct{}

	expected := AllTypesStruct{
		Byte:   nbt.Some[int8](1),
		Short:  nbt.Some[int16](2),
		Int:    nbt.Some[int32](3),
		Long:   nbt.Some[int64](4),
		Float:  nbt.Some[float32](5.0),
		Double: nbt.Some[float64](6.0),
	}

	floatBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(floatBytes, math.Float32bits(5.0))
	doubleBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(doubleBytes, math.Float64bits(6.0))

	_ = expected

	data := []byte{
		nbt.TagCompound,
		0, 0,
		nbt.TagByte, 0, 4, 'B', 'y', 't', 'e', 1,
		nbt.TagShort, 0, 5, 'S', 'h', 'o', 'r', 't', 0, 2,
		nbt.TagInt, 0, 3, 'I', 'n', 't', 0, 0, 0, 3,
		nbt.TagLong, 0, 4, 'L', 'o', 'n', 'g', 0, 0, 0, 0, 0, 0, 0, 4,
		nbt.TagFloat, 0, 5, 'F', 'l', 'o', 'a', 't', floatBytes[0], floatBytes[1], floatBytes[2], floatBytes[3],
		nbt.TagDouble, 0, 6, 'D', 'o', 'u', 'b', 'l', 'e', doubleBytes[0], doubleBytes[1], doubleBytes[2], doubleBytes[3], doubleBytes[4], doubleBytes[5], doubleBytes[6], doubleBytes[7],
		nbt.TagEnd,
	}

	name, k, err := nbt.Unmarshal(&a, data)

	if err != nil {
		t.Fatalf("expected err to be nil, but found %v", err)
	}

	if name != "" {
		t.Fail()
	}
	if k != len(data) {
		t.Fail()
	}

	opt := nbtcmp.NBTOptionTypeCmpOption()
	_ = opt
	if !cmp.Equal(a, expected) {
		t.Fail()
	}

}
