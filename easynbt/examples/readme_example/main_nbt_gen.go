package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"slices"
	"unicode/utf8"

	"github.com/toros100/easynbt/nbt"
)

var (
	_ nbt.Unmarshaler = (*Data)(nil)
)

func (d *Data) UnmarshalPayload(data []byte) (int, error) {
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

		case "hello":
			if foundFields0[0] {
				return 0, fmt.Errorf("on field Hello (nbt: hello): %w", nbt.ErrDuplicateValue)
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

			nbt.StringFromBytes(&d.Hello, data[off:off+strLen])
			off += strLen

			foundFields0[0] = true

		case "pos":
			if foundFields0[1] {
				return 0, fmt.Errorf("on field Position (nbt: pos): %w", nbt.ErrDuplicateValue)
			}

			if tag0 != nbt.TagCompound {
				return 0, nbt.ErrUnexpectedTag
			}

			foundFields1 := [2]bool{}

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

					if tag1 != nbt.TagInt {
						return 0, nbt.ErrUnexpectedTag
					}

					if off+4 > len(data) {
						return 0, nbt.ErrUnexpectedEOF
					}
					nbt.IntPayloadFromBytes(&d.Position.X, data[off:])
					off += 4

					foundFields1[0] = true

				case "Y":
					if foundFields1[1] {
						return 0, fmt.Errorf("on field Y: %w", nbt.ErrDuplicateValue)
					}

					if tag1 != nbt.TagInt {
						return 0, nbt.ErrUnexpectedTag
					}

					if off+4 > len(data) {
						return 0, nbt.ErrUnexpectedEOF
					}
					nbt.IntPayloadFromBytes(&d.Position.Y, data[off:])
					off += 4

					foundFields1[1] = true

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
					return 0, fmt.Errorf("on field %s: %w", [2]string{"X", "Y"}[i], nbt.ErrMissingValue)
				}
			}
			foundFields0[1] = true

		case "Numbers":
			if foundFields0[2] {
				return 0, fmt.Errorf("on field Numbers: %w", nbt.ErrDuplicateValue)
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

			if d.Numbers.Value == nil {
				d.Numbers.Value = make([]int8, int(listLen1))
			} else {
				d.Numbers.Value = slices.Grow(d.Numbers.Value, max(0, int(listLen1)))[:int(listLen1)]
			}
			list1 := d.Numbers.Value

			for i1 := range list1 {
				_ = i1

				if off+1 > len(data) {
					return 0, nbt.ErrUnexpectedEOF
				}
				nbt.BytePayloadFromBytes(&list1[i1], data[off:])
				off += 1

			}

			d.Numbers.Ok = true
			foundFields0[2] = true

		default:
			k, err := nbt.SkipPayload(tag0, data[off:])
			if err != nil {
				return 0, err
			}
			off += k
		}
	}

	required0 := [3]bool{true, true, false}
	for i := range foundFields0 {
		if !foundFields0[i] && required0[i] {
			return 0, fmt.Errorf("on field %s: %w", [3]string{"Hello (nbt: hello)", "Position (nbt: pos)", "Numbers"}[i], nbt.ErrMissingValue)
		}
	}

	return off, nil
}

func (d *Data) TagType() byte {
	return nbt.TagCompound
}
