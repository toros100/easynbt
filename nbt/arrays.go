package nbt

import (
	"bytes"
	"encoding/binary"
	"math"
)

type Array[T int8 | int32 | int64] []T

type ByteArray = Array[int8]
type IntArray = Array[int32]
type LongArray = Array[int64]

var (
	_ Unmarshaler = (*Array[int8])(nil)
	_ Unmarshaler = (*Array[int32])(nil)
	_ Unmarshaler = (*Array[int64])(nil)
)

func (a *Array[T]) TagType() byte {
	var t T
	switch any(t).(type) {
	case int8:
		return TagByteArray
	case int32:
		return TagIntArray
	case int64:
		return TagLongArray
	default:
		panic("unreachable")
	}
}

func (a *Array[T]) UnmarshalPayload(data []byte) (int, error) {

	if len(data) < 4 {
		return 0, ErrUnexpectedEOF
	}

	l := binary.BigEndian.Uint32(data)
	if l > math.MaxInt32 {
		return 0, ErrInvalidLength
	}
	byteSize := binary.Size(*new(T))

	if len(data)-4 < int(l)*byteSize {
		return 0, ErrUnexpectedEOF
	}

	r := bytes.NewReader(data[4:])

	*a = make(Array[T], int(l))

	// err is always nil here
	_ = binary.Read(r, binary.BigEndian, a)

	return 4 + int(l)*byteSize, nil
}
