package nbt

import (
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	"math"
)

const (
	// i don't like using iota if the values are actually meaningful
	// (although i really doubt they would ever change the way iota behaves)
	TagEnd       byte = 0
	TagByte      byte = 1
	TagShort     byte = 2
	TagInt       byte = 3
	TagLong      byte = 4
	TagFloat     byte = 5
	TagDouble    byte = 6
	TagByteArray byte = 7
	TagString    byte = 8
	TagList      byte = 9
	TagCompound  byte = 10
	TagIntArray  byte = 11
	TagLongArray byte = 12
)

func isValidTag(tag byte) bool {
	return tag <= 12
}

type Unmarshaller interface {
	// TagType is expected to return the same byte value in {0, ..., 12} across calls
	// on the same implementing type, otherwise the behaviour of the generated code is undefined.
	// Further, TagType must be callable safely on a nil receiver.
	TagType() byte

	// UnmarshalPayload is called to unmarshal a NBT payload from the provided bytes.
	// It will only be called if the associated NBT payload type (byte) matches the value
	// reported by the TagType method. If err == nil, then n must be the number of bytes consumed.
	// Note that acceptable values of n are uniquely determined by the NBT payload type
	// and data, based on the semantics of the Minecraft (Java Edition) NBT format.
	UnmarshalPayload(data []byte) (n int, err error)
}

// Unmarshal is usually used on a type representing the outermost/root compound tag.
// If err == nil, then name is the name of the root tag (often empty) and n the number of bytes
// consumed.
func Unmarshal[T Unmarshaller](t T, data []byte) (name string, n int, err error) {
	if len(data) == 0 {
		return "", 0, ErrUnexpectedEOF
	}

	typ := data[0]
	off := 1

	if !isValidTag(typ) {
		return "", 0, ErrUnexpectedTag
	}

	if typ != t.TagType() {
		return "", 0, ErrUnexpectedTag
	}

	if off+2 > len(data) {
		return "", 0, ErrUnexpectedEOF
	}

	nameLen := int(binary.BigEndian.Uint16(data[off:]))
	off += 2

	if off+nameLen > len(data) {
		return "", 0, ErrUnexpectedEOF
	}

	nm := string(data[off : off+nameLen])

	off += nameLen

	k, err := (t).UnmarshalPayload(data[off:])

	if err != nil {
		return "", 0, err
	}

	return nm, off + k, nil
}

func skipString(data []byte) (n int, err error) {
	if len(data) < 2 {
		return 0, ErrUnexpectedEOF
	}
	skip := int(binary.BigEndian.Uint16(data)) + 2
	if skip > len(data) {
		return 0, ErrUnexpectedEOF
	}
	return skip, nil
}

// SkipPayload is used by generated code to skip an unneeded NBT payload.
// More precisely, payloads of tags with names that did not match any
// expected name in a compound tag that is being unmarshalled.
func SkipPayload(tag byte, data []byte) (n int, err error) {
	if !isValidTag(tag) || tag == TagEnd {
		return 0, ErrUnexpectedTag
	}

	switch tag {

	case TagByte, TagShort, TagInt, TagLong, TagFloat, TagDouble:
		skip := fixedSizePayloadSize(tag)
		if skip > len(data) {
			return 0, ErrUnexpectedEOF
		}
		return skip, nil

	case TagString:
		return skipString(data)
	case TagByteArray:
		return skipFixedSizePayloadSequence(TagByte, data)
	case TagIntArray:
		return skipFixedSizePayloadSequence(TagInt, data)
	case TagLongArray:
		return skipFixedSizePayloadSequence(TagLong, data)
	case TagList:
		// first byte: element type byte
		// next 4 bytes: length
		if len(data) < 5 {
			return 0, ErrUnexpectedEOF
		}
		elemTag := data[0]

		if !isValidTag(elemTag) {
			return 0, fmt.Errorf("invalid list element tag: %v", elemTag)
		}

		switch elemTag {
		case TagEnd, TagByte, TagShort, TagInt, TagLong, TagFloat, TagDouble:
			k, err := skipFixedSizePayloadSequence(elemTag, data[1:])
			return k + 1, err

		default:
			listLen := binary.BigEndian.Uint32(data[1:])
			if listLen > math.MaxInt32 {
				return 0, ErrInvalidLength
			}
			skip := 5

			for range listLen {
				k, err := SkipPayload(elemTag, data[skip:])
				if err != nil {
					return 0, err
				}
				skip += k
			}
			return skip, nil
		}

	case TagCompound:
		skip := 0
		for {
			if skip >= len(data) {
				return 0, ErrUnexpectedEOF
			}

			t := data[skip]
			skip += 1

			if t == TagEnd {
				break
			}

			if !isValidTag(t) {
				return 0, ErrUnexpectedTag
			}

			k, err := skipString(data[skip:])

			if err != nil {
				return 0, err
			}

			skip += k

			l, err := SkipPayload(t, data[skip:])
			if err != nil {
				return 0, err
			}
			skip += l
		}
		return skip, nil
	default:
		// should be unreachable
		panic("programmer error")
	}

}

func fixedSizePayloadSize(tag byte) int {
	switch tag {
	case TagByte:
		return 1
	case TagShort:
		return 2
	case TagInt:
		return 4
	case TagLong:
		return 8
	case TagFloat:
		return 4
	case TagDouble:
		return 8
	default:
		log.Panicf("programmer error: %v is not associated with a fixed size payload", tag)
		return -1 // unreachable
	}
}

// primary use: skipping nbt array payloads
// can be reused for lists if element tag has fixed size payload (anything except list, compound, string),
// because lists and arrays are both length-prefixed with a 32 bit big endian signed int
func skipFixedSizePayloadSequence(elemTag byte, data []byte) (int, error) {
	if len(data) < 4 {
		return 0, ErrUnexpectedEOF
	}

	seqLen := binary.BigEndian.Uint32(data)

	if seqLen > math.MaxInt32 {
		return 0, ErrInvalidLength
	}

	if elemTag == TagEnd {
		if seqLen != 0 {
			return 0, errors.New("list with element tag 0 (end tag) must have length 0")
		}
		return 4, nil
	}

	elemSize := fixedSizePayloadSize(elemTag)
	n := int(seqLen)*elemSize + 4

	if n >= len(data) {
		return 0, ErrUnexpectedEOF
	}

	return n, nil
}

// helper methods to be used in generated code
// deliberately minimal to enable inlining
// (external error handling, length checks required)

func BytePayloadFromBytes[T ~int8](v *T, data []byte) {
	*v = T(data[0])
}

func ShortPayloadFromBytes[T ~int16](v *T, data []byte) {
	*v = T(binary.BigEndian.Uint16(data))
}

func IntPayloadFromBytes[T ~int32](v *T, data []byte) {
	*v = T(binary.BigEndian.Uint32(data))
}

func LongPayloadFromBytes[T ~int64](v *T, data []byte) {
	*v = T(binary.BigEndian.Uint64(data))
}

func FloatPayloadFromBytes[T ~float32](v *T, data []byte) {
	*v = T(math.Float32frombits(binary.BigEndian.Uint32(data)))
}

func DoublePayloadFromBytes[T ~float64](v *T, data []byte) {
	*v = T(math.Float64frombits(binary.BigEndian.Uint64(data)))
}

// caution: this does not read a full string payload from bytes, as the
// two bytes for the strings length are read before calling this method
func StringFromBytes[T ~string](v *T, strData []byte) {
	*v = T(strData)
}
