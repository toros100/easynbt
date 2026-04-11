package option

import (
	"encoding/binary"
	"fmt"
	"math"
	"slices"
	"unicode/utf8"

	"github.com/toros100/easynbt/nbt"
)

var (
	_ nbt.Unmarshaler = (*AllTypesStruct)(nil)
	_ nbt.Unmarshaler = (*Named)(nil)
	_ nbt.Unmarshaler = (*SomeStruct)(nil)
)

func (a *AllTypesStruct) UnmarshalPayload(data []byte) (int, error) {
	off := 0

	foundFields0 := [11]bool{}

	for {
		if off >= len(data) {
			return 0, nbt.ErrUnexpectedEOF
		}

		tag0 := data[off]
		off += 1

		if tag0 == nbt.TagEnd {
			break
		}

		if off+2 > len(data) {
			return 0, nbt.ErrUnexpectedEOF
		}

		strLen := int(binary.BigEndian.Uint16(data[off:]))
		off += 2

		if off+strLen > len(data) {
			return 0, nbt.ErrUnexpectedEOF
		}

		strData := data[off : off+strLen]
		off += strLen

		switch string(strData) {

		case "Byte":
			if foundFields0[0] {
				return 0, fmt.Errorf("on field Byte: %w", nbt.ErrDuplicateValue)
			}

			if tag0 != nbt.TagByte {
				return 0, nbt.ErrUnexpectedTag
			}

			if off+1 > len(data) {
				return 0, nbt.ErrUnexpectedEOF
			}
			nbt.BytePayloadFromBytes(&a.Byte.Value, data[off:])
			off += 1

			a.Byte.Ok = true
			foundFields0[0] = true

		case "Short":
			if foundFields0[1] {
				return 0, fmt.Errorf("on field Short: %w", nbt.ErrDuplicateValue)
			}

			if tag0 != nbt.TagShort {
				return 0, nbt.ErrUnexpectedTag
			}

			if off+2 > len(data) {
				return 0, nbt.ErrUnexpectedEOF
			}
			nbt.ShortPayloadFromBytes(&a.Short.Value, data[off:])
			off += 2

			a.Short.Ok = true
			foundFields0[1] = true

		case "Int":
			if foundFields0[2] {
				return 0, fmt.Errorf("on field Int: %w", nbt.ErrDuplicateValue)
			}

			if tag0 != nbt.TagInt {
				return 0, nbt.ErrUnexpectedTag
			}

			if off+4 > len(data) {
				return 0, nbt.ErrUnexpectedEOF
			}
			nbt.IntPayloadFromBytes(&a.Int.Value, data[off:])
			off += 4

			a.Int.Ok = true
			foundFields0[2] = true

		case "Long":
			if foundFields0[3] {
				return 0, fmt.Errorf("on field Long: %w", nbt.ErrDuplicateValue)
			}

			if tag0 != nbt.TagLong {
				return 0, nbt.ErrUnexpectedTag
			}

			if off+8 > len(data) {
				return 0, nbt.ErrUnexpectedEOF
			}
			nbt.LongPayloadFromBytes(&a.Long.Value, data[off:])
			off += 8

			a.Long.Ok = true
			foundFields0[3] = true

		case "Float":
			if foundFields0[4] {
				return 0, fmt.Errorf("on field Float: %w", nbt.ErrDuplicateValue)
			}

			if tag0 != nbt.TagFloat {
				return 0, nbt.ErrUnexpectedTag
			}

			if off+4 > len(data) {
				return 0, nbt.ErrUnexpectedEOF
			}
			nbt.FloatPayloadFromBytes(&a.Float.Value, data[off:])
			off += 4

			a.Float.Ok = true
			foundFields0[4] = true

		case "Double":
			if foundFields0[5] {
				return 0, fmt.Errorf("on field Double: %w", nbt.ErrDuplicateValue)
			}

			if tag0 != nbt.TagDouble {
				return 0, nbt.ErrUnexpectedTag
			}

			if off+8 > len(data) {
				return 0, nbt.ErrUnexpectedEOF
			}
			nbt.DoublePayloadFromBytes(&a.Double.Value, data[off:])
			off += 8

			a.Double.Ok = true
			foundFields0[5] = true

		case "String":
			if foundFields0[6] {
				return 0, fmt.Errorf("on field String: %w", nbt.ErrDuplicateValue)
			}

			if tag0 != nbt.TagString {
				return 0, nbt.ErrUnexpectedTag
			}

			if off+2 > len(data) {
				return 0, nbt.ErrUnexpectedEOF
			}

			strLen := int(binary.BigEndian.Uint16(data[off:]))
			off += 2

			if off+strLen > len(data) {
				return 0, nbt.ErrUnexpectedEOF
			}

			if !utf8.Valid(data[off : off+strLen]) {
				return 0, nbt.ErrInvalidUTF8
			}

			nbt.StringFromBytes(&a.String.Value, data[off:off+strLen])
			off += strLen

			a.String.Ok = true
			foundFields0[6] = true

		case "Compound":
			if foundFields0[7] {
				return 0, fmt.Errorf("on field Compound: %w", nbt.ErrDuplicateValue)
			}

			if tag0 != nbt.TagCompound {
				return 0, nbt.ErrUnexpectedTag
			}

			foundFields1 := [1]bool{}

			for {
				if off >= len(data) {
					return 0, nbt.ErrUnexpectedEOF
				}

				tag1 := data[off]
				off += 1

				if tag1 == nbt.TagEnd {
					break
				}

				if off+2 > len(data) {
					return 0, nbt.ErrUnexpectedEOF
				}

				strLen := int(binary.BigEndian.Uint16(data[off:]))
				off += 2

				if off+strLen > len(data) {
					return 0, nbt.ErrUnexpectedEOF
				}

				strData := data[off : off+strLen]
				off += strLen

				switch string(strData) {

				case "X":
					if foundFields1[0] {
						return 0, fmt.Errorf("on field X: %w", nbt.ErrDuplicateValue)
					}

					if tag1 != nbt.TagByte {
						return 0, nbt.ErrUnexpectedTag
					}

					if off+1 > len(data) {
						return 0, nbt.ErrUnexpectedEOF
					}
					nbt.BytePayloadFromBytes(&a.Compound.Value.X, data[off:])
					off += 1

					foundFields1[0] = true

				default:
					k, err := nbt.SkipPayload(tag1, data[off:])
					if err != nil {
						return 0, err
					}
					off += k
				}
			}

			for i := range foundFields1 {
				if !foundFields1[i] {
					return 0, fmt.Errorf("on field %s: %w", [1]string{"X"}[i], nbt.ErrMissingValue)
				}
			}

			a.Compound.Ok = true
			foundFields0[7] = true

		case "ListPrimitive":
			if foundFields0[8] {
				return 0, fmt.Errorf("on field ListPrimitive: %w", nbt.ErrDuplicateValue)
			}

			if tag0 != nbt.TagList {
				return 0, nbt.ErrUnexpectedTag
			}

			if off+5 > len(data) {
				return 0, nbt.ErrUnexpectedEOF
			}

			elemTag1 := data[off]
			off += 1

			listLen1 := binary.BigEndian.Uint32(data[off:])
			off += 4

			if listLen1 > math.MaxInt32 {
				return 0, nbt.ErrInvalidLength
			}

			if elemTag1 != nbt.TagByte && !(listLen1 == 0 && elemTag1 == nbt.TagEnd) {
				return 0, nbt.ErrUnexpectedTag
			}

			if a.ListPrimitive.Value == nil {
				a.ListPrimitive.Value = make([]int8, int(listLen1))
			} else {
				a.ListPrimitive.Value = slices.Grow(a.ListPrimitive.Value, max(0, int(listLen1)))[:int(listLen1)]
			}
			list1 := a.ListPrimitive.Value

			for i1 := range list1 {
				_ = i1

				if off+1 > len(data) {
					return 0, nbt.ErrUnexpectedEOF
				}
				nbt.BytePayloadFromBytes(&list1[i1], data[off:])
				off += 1

			}

			a.ListPrimitive.Ok = true
			foundFields0[8] = true

		case "ListNamed":
			if foundFields0[9] {
				return 0, fmt.Errorf("on field ListNamed: %w", nbt.ErrDuplicateValue)
			}

			if tag0 != nbt.TagList {
				return 0, nbt.ErrUnexpectedTag
			}

			if off+5 > len(data) {
				return 0, nbt.ErrUnexpectedEOF
			}

			elemTag1 := data[off]
			off += 1

			listLen1 := binary.BigEndian.Uint32(data[off:])
			off += 4

			if listLen1 > math.MaxInt32 {
				return 0, nbt.ErrInvalidLength
			}

			if elemTag1 != (*Named)(nil).TagType() && !(listLen1 == 0 && elemTag1 == nbt.TagEnd) {
				return 0, nbt.ErrUnexpectedTag
			}

			if a.ListNamed.Value == nil {
				a.ListNamed.Value = make([]Named, int(listLen1))
			} else {
				a.ListNamed.Value = slices.Grow(a.ListNamed.Value, max(0, int(listLen1)))[:int(listLen1)]
			}
			list1 := a.ListNamed.Value

			for i1 := range list1 {
				_ = i1

				k, err := list1[i1].UnmarshalPayload(data[off:])
				if err != nil {
					return 0, err
				}
				off += k
			}

			a.ListNamed.Ok = true
			foundFields0[9] = true

		case "Named":
			if foundFields0[10] {
				return 0, fmt.Errorf("on field Named: %w", nbt.ErrDuplicateValue)
			}

			if tag0 != (*Named)(nil).TagType() {
				return 0, nbt.ErrUnexpectedTag
			}

			k, err := a.Named.Value.UnmarshalPayload(data[off:])
			if err != nil {
				return 0, err
			}
			off += k

			a.Named.Ok = true
			foundFields0[10] = true

		default:
			k, err := nbt.SkipPayload(tag0, data[off:])
			if err != nil {
				return 0, err
			}
			off += k
		}
	}

	return off, nil
}

func (a *AllTypesStruct) TagType() byte {
	return nbt.TagCompound
}

func (n *Named) UnmarshalPayload(data []byte) (int, error) {
	off := 0

	if off+2 > len(data) {
		return 0, nbt.ErrUnexpectedEOF
	}
	nbt.ShortPayloadFromBytes(n, data[off:])
	off += 2

	return off, nil
}

func (n *Named) TagType() byte {
	return nbt.TagShort
}

func (s *SomeStruct) UnmarshalPayload(data []byte) (int, error) {
	off := 0

	foundFields0 := [3]bool{}

	for {
		if off >= len(data) {
			return 0, nbt.ErrUnexpectedEOF
		}

		tag0 := data[off]
		off += 1

		if tag0 == nbt.TagEnd {
			break
		}

		if off+2 > len(data) {
			return 0, nbt.ErrUnexpectedEOF
		}

		strLen := int(binary.BigEndian.Uint16(data[off:]))
		off += 2

		if off+strLen > len(data) {
			return 0, nbt.ErrUnexpectedEOF
		}

		strData := data[off : off+strLen]
		off += strLen

		switch string(strData) {

		case "OptionalIntViaType":
			if foundFields0[0] {
				return 0, fmt.Errorf("on field OptionalIntViaType: %w", nbt.ErrDuplicateValue)
			}

			if tag0 != nbt.TagInt {
				return 0, nbt.ErrUnexpectedTag
			}

			if off+4 > len(data) {
				return 0, nbt.ErrUnexpectedEOF
			}
			nbt.IntPayloadFromBytes(&s.OptionalIntViaType.Value, data[off:])
			off += 4

			s.OptionalIntViaType.Ok = true
			foundFields0[0] = true

		case "OptionalIntViaTag":
			if foundFields0[1] {
				return 0, fmt.Errorf("on field OptionalIntViaTag: %w", nbt.ErrDuplicateValue)
			}

			if tag0 != nbt.TagInt {
				return 0, nbt.ErrUnexpectedTag
			}

			if off+4 > len(data) {
				return 0, nbt.ErrUnexpectedEOF
			}
			nbt.IntPayloadFromBytes(&s.OptionalIntViaTag, data[off:])
			off += 4

			foundFields0[1] = true

		case "SomeString":
			if foundFields0[2] {
				return 0, fmt.Errorf("on field SomeString: %w", nbt.ErrDuplicateValue)
			}

			if tag0 != nbt.TagString {
				return 0, nbt.ErrUnexpectedTag
			}

			if off+2 > len(data) {
				return 0, nbt.ErrUnexpectedEOF
			}

			strLen := int(binary.BigEndian.Uint16(data[off:]))
			off += 2

			if off+strLen > len(data) {
				return 0, nbt.ErrUnexpectedEOF
			}

			if !utf8.Valid(data[off : off+strLen]) {
				return 0, nbt.ErrInvalidUTF8
			}

			nbt.StringFromBytes(&s.SomeString, data[off:off+strLen])
			off += strLen

			foundFields0[2] = true

		default:
			k, err := nbt.SkipPayload(tag0, data[off:])
			if err != nil {
				return 0, err
			}
			off += k
		}
	}

	required0 := [3]bool{false, false, true}
	for i := range foundFields0 {
		if !foundFields0[i] && required0[i] {
			return 0, fmt.Errorf("on field %s: %w", [3]string{"OptionalIntViaType", "OptionalIntViaTag", "SomeString"}[i], nbt.ErrMissingValue)
		}
	}

	return off, nil
}

func (s *SomeStruct) TagType() byte {
	return nbt.TagCompound
}
