package nbt

import (
	"bytes"
	"encoding/binary"
	"iter"
	"math"
	"slices"
)

var (
	_ Unmarshaller = (*ByteArray)(nil)
	_ Unmarshaller = (*IntArray)(nil)
	_ Unmarshaller = (*LongArray)(nil)
)

type ByteArray struct {
	backing []byte
}

func (a *ByteArray) TagType() byte {
	return TagByteArray
}

func (a *ByteArray) UnmarshalPayload(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, ErrUnexpectedEOF
	}

	l := binary.BigEndian.Uint32(data)
	if l > math.MaxInt32 {
		return 0, ErrInvalidLength
	}

	if len(data)-4 < int(l) {
		return 0, ErrUnexpectedEOF
	}
	a.backing = data[4 : 4+l]

	return 4 + int(l), nil
}

func (a *ByteArray) Len() int {
	return len(a.backing)
}

func (a *ByteArray) At(idx int) int8 {
	if idx >= len(a.backing) {
		panic("index out of range")
	}
	return int8(a.backing[idx])
}

func (a *ByteArray) ToSlice() []int8 {
	sl := make([]int8, len(a.backing))
	for i := range a.backing {
		sl[i] = int8(a.backing[i])
	}
	return sl
}

func (a *ByteArray) Iter() iter.Seq2[int, int8] {
	return func(yield func(int, int8) bool) {
		for i := range a.Len() {
			if ok := yield(i, a.At(i)); !ok {
				return
			}
		}
	}
}

type IntArray struct {
	backing []byte
}

func (a *IntArray) TagType() byte {
	return TagIntArray
}

func (a *IntArray) UnmarshalPayload(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, ErrUnexpectedEOF
	}

	l := binary.BigEndian.Uint32(data)
	if l > math.MaxInt32 {
		return 0, ErrInvalidLength
	}

	if len(data)-4 < int(l)*4 {
		return 0, ErrUnexpectedEOF
	}
	a.backing = data[4 : 4+l*4]
	return 4 + int(l)*4, nil
}

func (a *IntArray) Len() int {
	return len(a.backing) / 4
}

func (a *IntArray) At(idx int) int32 {
	i := idx * 4
	if i+4 > len(a.backing) {
		panic("index out of range")
	}
	return int32(binary.BigEndian.Uint32(a.backing[i:]))
}

func (a *IntArray) ToSlice() []int32 {
	sl := make([]int32, len(a.backing)/4)
	for i := range sl {
		sl[i] = int32(binary.BigEndian.Uint32(a.backing[i*4:]))
	}
	return sl
}

func (a *IntArray) Iter() iter.Seq2[int, int32] {
	return func(yield func(int, int32) bool) {
		for i := range a.Len() {
			if ok := yield(i, a.At(i)); !ok {
				return
			}
		}
	}
}

type LongArray struct {
	backing []byte
}

func (a *LongArray) TagType() byte {
	return TagLongArray
}

func (a *LongArray) UnmarshalPayload(data []byte) (int, error) {
	if len(data) < 4 {
		return 0, ErrUnexpectedEOF
	}

	l := binary.BigEndian.Uint32(data)
	if l > math.MaxInt32 {
		return 0, ErrInvalidLength
	}

	if len(data)-4 < int(l)*8 {
		return 0, ErrUnexpectedEOF
	}
	a.backing = data[4 : 4+l*8]
	return 4 + int(l)*8, nil
}

func (a *LongArray) Len() int {
	return len(a.backing) / 8
}

func (a *LongArray) At(idx int) int64 {
	i := idx * 8
	if i+8 > len(a.backing) {
		panic("index out of range")
	}
	return int64(binary.BigEndian.Uint64(a.backing[i:]))
}

func (a *LongArray) ToSlice() []int64 {
	sl := make([]int64, len(a.backing)/8)
	for i := range sl {
		sl[i] = int64(binary.BigEndian.Uint64(a.backing[i*8:]))
	}
	return sl
}

func (a *LongArray) Iter() iter.Seq2[int, int64] {
	return func(yield func(int, int64) bool) {
		for i := range a.Len() {
			if ok := yield(i, a.At(i)); !ok {
				return
			}
		}
	}
}

// following methods are mainly for being able to construct arrays for comparison in testing
// i dont expect them to be useful anywhere else

func (a *ByteArray) FromSlice(sl []int8) *ByteArray {
	if a == nil {
		panic("receiver nil")
	}
	buf := make([]byte, len(sl))
	for i := range buf {
		buf[i] = byte(sl[i])
	}
	a.backing = buf
	return a
}

func (a *IntArray) FromSlice(sl []int32) *IntArray {
	if a == nil {
		panic("receiver nil")
	}
	buf := bytes.NewBuffer(make([]byte, 0, len(sl)*4))
	_ = binary.Write(buf, binary.BigEndian, sl)
	a.backing = buf.Bytes()
	return a
}

func (a *LongArray) FromSlice(sl []int64) *LongArray {
	if a == nil {
		panic("receiver nil")
	}
	buf := bytes.NewBuffer(make([]byte, 0, len(sl)*8))
	_ = binary.Write(buf, binary.BigEndian, sl)
	a.backing = buf.Bytes()
	return a
}

func (a *ByteArray) Equal(other *ByteArray) bool {
	if (a == nil) != (other == nil) {
		return false
	}
	if a == nil {
		// both nil
		return true
	}
	return slices.Equal(a.backing, other.backing)
}

func (a *IntArray) Equal(other *IntArray) bool {
	if (a == nil) != (other == nil) {
		return false
	}
	if a == nil {
		// both nil
		return true
	}
	return slices.Equal(a.backing, other.backing)
}

func (a *LongArray) Equal(other *LongArray) bool {
	if (a == nil) != (other == nil) {
		return false
	}
	if a == nil {
		// both nil
		return true
	}
	return slices.Equal(a.backing, other.backing)
}
