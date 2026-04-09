package pointer

import (
	"encoding/binary"
	"fmt"

	"github.com/toros100/easynbt/nbt"
)

var (
	_ nbt.Unmarshaller = (*List)(nil)
)

func (l *List) UnmarshalPayload(data []byte) (int, error) {
	off := 0

	foundFields0 := [2]bool{}

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

		case "Head":
			if foundFields0[0] {
				return 0, fmt.Errorf("on field Head: %w", nbt.ErrDuplicateValue)
			}

			if tag0 != nbt.TagByte {
				return 0, nbt.ErrUnexpectedTag
			}

			if off+1 > len(data) {
				return 0, nbt.ErrUnexpectedEOF
			}
			nbt.BytePayloadFromBytes(&l.Head, data[off:])
			off += 1

			foundFields0[0] = true

		case "Tail":
			if foundFields0[1] {
				return 0, fmt.Errorf("on field Tail: %w", nbt.ErrDuplicateValue)
			}

			if tag0 != (*List)(nil).TagType() {
				return 0, nbt.ErrUnexpectedTag
			}

			l.Tail.Value = new(List)

			k, err := l.Tail.Value.UnmarshalPayload(data[off:])
			if err != nil {
				return 0, err
			}
			off += k

			l.Tail.Ok = true
			foundFields0[1] = true

		default:
			k, err := nbt.SkipPayload(tag0, data[off:])
			if err != nil {
				return 0, err
			}
			off += k
		}
	}

	required0 := [2]bool{true, false}
	for i := range foundFields0 {
		if !foundFields0[i] && required0[i] {
			return 0, fmt.Errorf("on field %s: %w", [2]string{"Head", "Tail"}[i], nbt.ErrMissingValue)
		}
	}

	return off, nil
}

func (l *List) TagType() byte {
	return nbt.TagCompound
}
