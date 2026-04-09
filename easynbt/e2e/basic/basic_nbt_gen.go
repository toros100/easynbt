package basic

import (
	"encoding/binary"
	"fmt"
	"math"
	"slices"
	"unicode/utf8"

	"github.com/toros100/easynbt/nbt"
)

var (
	_ nbt.Unmarshaller = (*BytePayload)(nil)
	_ nbt.Unmarshaller = (*CompoundPayload)(nil)
	_ nbt.Unmarshaller = (*DoublePayload)(nil)
	_ nbt.Unmarshaller = (*FloatPayload)(nil)
	_ nbt.Unmarshaller = (*IntPayload)(nil)
	_ nbt.Unmarshaller = (*ListOfIntPayload)(nil)
	_ nbt.Unmarshaller = (*LongPayload)(nil)
	_ nbt.Unmarshaller = (*ShortPayload)(nil)
	_ nbt.Unmarshaller = (*StringPayload)(nil)
)

func (b *BytePayload) UnmarshalPayload(data []byte) (int, error) {
	off := 0

	if off+1 > len(data) {
		return 0, nbt.ErrUnexpectedEOF
	}
	nbt.BytePayloadFromBytes(b, data[off:])
	off += 1

	return off, nil
}

func (b *BytePayload) TagType() byte {
	return nbt.TagByte
}

func (c *CompoundPayload) UnmarshalPayload(data []byte) (int, error) {
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

		case "IntPayloadField":
			if foundFields0[0] {
				return 0, fmt.Errorf("on field IntPayloadField: %w", nbt.ErrDuplicateValue)
			}

			if tag0 != (*IntPayload)(nil).TagType() {
				return 0, nbt.ErrUnexpectedTag
			}

			k, err := c.IntPayloadField.UnmarshalPayload(data[off:])
			if err != nil {
				return 0, err
			}
			off += k
			foundFields0[0] = true

		case "RawIntField":
			if foundFields0[1] {
				return 0, fmt.Errorf("on field RawIntField: %w", nbt.ErrDuplicateValue)
			}

			if tag0 != nbt.TagInt {
				return 0, nbt.ErrUnexpectedTag
			}

			if off+4 > len(data) {
				return 0, nbt.ErrUnexpectedEOF
			}
			nbt.IntPayloadFromBytes(&c.RawIntField, data[off:])
			off += 4

			foundFields0[1] = true

		case "InnerCompoundField":
			if foundFields0[2] {
				return 0, fmt.Errorf("on field InnerCompoundField: %w", nbt.ErrDuplicateValue)
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

				case "ListOfByteField":
					if foundFields1[0] {
						return 0, fmt.Errorf("on field ListOfByteField: %w", nbt.ErrDuplicateValue)
					}

					if tag1 != nbt.TagList {
						return 0, nbt.ErrUnexpectedTag
					}

					if off+5 > len(data) {
						return 0, nbt.ErrUnexpectedEOF
					}

					elemTag2 := data[off]
					off += 1

					listLen2 := binary.BigEndian.Uint32(data[off:])
					off += 4

					if listLen2 > math.MaxInt32 {
						return 0, nbt.ErrInvalidLength
					}

					if elemTag2 != nbt.TagByte && !(listLen2 == 0 && elemTag2 == nbt.TagEnd) {
						return 0, nbt.ErrUnexpectedTag
					}
					if listLen2 == 0 {
						continue
					}

					if c.InnerCompoundField.ListOfByteField == nil {
						c.InnerCompoundField.ListOfByteField = make([]int8, int(listLen2))
					} else {
						c.InnerCompoundField.ListOfByteField = slices.Grow(c.InnerCompoundField.ListOfByteField, int(listLen2))[:int(listLen2)]
					}
					list2 := c.InnerCompoundField.ListOfByteField

					for i2 := range list2 {

						if off+1 > len(data) {
							return 0, nbt.ErrUnexpectedEOF
						}
						nbt.BytePayloadFromBytes(&list2[i2], data[off:])
						off += 1

					}
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
					return 0, fmt.Errorf("on field %s: %w", [1]string{"ListOfByteField"}[i], nbt.ErrMissingValue)
				}
			}
			foundFields0[2] = true

		default:
			k, err := nbt.SkipPayload(tag0, data[off:])
			if err != nil {
				return 0, err
			}
			off += k
		}
	}

	for i := range foundFields0 {
		if !foundFields0[i] {
			return 0, fmt.Errorf("on field %s: %w", [3]string{"IntPayloadField", "RawIntField", "InnerCompoundField"}[i], nbt.ErrMissingValue)
		}
	}

	return off, nil
}

func (c *CompoundPayload) TagType() byte {
	return nbt.TagCompound
}

func (d *DoublePayload) UnmarshalPayload(data []byte) (int, error) {
	off := 0

	if off+8 > len(data) {
		return 0, nbt.ErrUnexpectedEOF
	}
	nbt.DoublePayloadFromBytes(d, data[off:])
	off += 8

	return off, nil
}

func (d *DoublePayload) TagType() byte {
	return nbt.TagDouble
}

func (f *FloatPayload) UnmarshalPayload(data []byte) (int, error) {
	off := 0

	if off+4 > len(data) {
		return 0, nbt.ErrUnexpectedEOF
	}
	nbt.FloatPayloadFromBytes(f, data[off:])
	off += 4

	return off, nil
}

func (f *FloatPayload) TagType() byte {
	return nbt.TagFloat
}

func (i *IntPayload) UnmarshalPayload(data []byte) (int, error) {
	off := 0

	if off+4 > len(data) {
		return 0, nbt.ErrUnexpectedEOF
	}
	nbt.IntPayloadFromBytes(i, data[off:])
	off += 4

	return off, nil
}

func (i *IntPayload) TagType() byte {
	return nbt.TagInt
}

func (l *ListOfIntPayload) UnmarshalPayload(data []byte) (int, error) {
	off := 0

	if off+5 > len(data) {
		return 0, nbt.ErrUnexpectedEOF
	}

	elemTag0 := data[off]
	off += 1

	listLen0 := binary.BigEndian.Uint32(data[off:])
	off += 4

	if listLen0 > math.MaxInt32 {
		return 0, nbt.ErrInvalidLength
	}

	if elemTag0 != nbt.TagInt && !(listLen0 == 0 && elemTag0 == nbt.TagEnd) {
		return 0, nbt.ErrUnexpectedTag
	}
	if listLen0 == 0 {
		return off, nil
	}

	if *l == nil {
		*l = make([]int32, int(listLen0))
	} else {
		*l = slices.Grow(*l, int(listLen0))[:int(listLen0)]
	}
	list0 := *l

	for i0 := range list0 {

		if off+4 > len(data) {
			return 0, nbt.ErrUnexpectedEOF
		}
		nbt.IntPayloadFromBytes(&list0[i0], data[off:])
		off += 4

	}

	return off, nil
}

func (l *ListOfIntPayload) TagType() byte {
	return nbt.TagList
}

func (l *LongPayload) UnmarshalPayload(data []byte) (int, error) {
	off := 0

	if off+8 > len(data) {
		return 0, nbt.ErrUnexpectedEOF
	}
	nbt.LongPayloadFromBytes(l, data[off:])
	off += 8

	return off, nil
}

func (l *LongPayload) TagType() byte {
	return nbt.TagLong
}

func (s *ShortPayload) UnmarshalPayload(data []byte) (int, error) {
	off := 0

	if off+2 > len(data) {
		return 0, nbt.ErrUnexpectedEOF
	}
	nbt.ShortPayloadFromBytes(s, data[off:])
	off += 2

	return off, nil
}

func (s *ShortPayload) TagType() byte {
	return nbt.TagShort
}

func (s *StringPayload) UnmarshalPayload(data []byte) (int, error) {
	off := 0

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

	nbt.StringFromBytes(s, data[off:off+strLen])
	off += strLen

	return off, nil
}

func (s *StringPayload) TagType() byte {
	return nbt.TagString
}
