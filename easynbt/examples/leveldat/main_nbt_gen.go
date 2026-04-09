package main

import (
	"encoding/binary"
	"fmt"
	"unicode/utf8"

	"github.com/toros100/easynbt/nbt"
)

var (
	_ nbt.Unmarshaller = (*Root)(nil)
)

func (r *Root) UnmarshalPayload(data []byte) (int, error) {
	off := 0

	foundFields0 := [1]bool{}

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

		case "Data":
			if foundFields0[0] {
				return 0, fmt.Errorf("on field Data: %w", nbt.ErrDuplicateValue)
			}

			if tag0 != nbt.TagCompound {
				return 0, nbt.ErrUnexpectedTag
			}

			foundFields1 := [5]bool{}

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

				case "LevelName":
					if foundFields1[0] {
						return 0, fmt.Errorf("on field LevelName: %w", nbt.ErrDuplicateValue)
					}

					if tag1 != nbt.TagString {
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

					nbt.StringFromBytes(&r.Data.LevelName, data[off:off+strLen])
					off += strLen

					foundFields1[0] = true

				case "LastPlayed":
					if foundFields1[1] {
						return 0, fmt.Errorf("on field LastPlayed: %w", nbt.ErrDuplicateValue)
					}

					if tag1 != nbt.TagLong {
						return 0, nbt.ErrUnexpectedTag
					}

					if off+8 > len(data) {
						return 0, nbt.ErrUnexpectedEOF
					}
					nbt.LongPayloadFromBytes(&r.Data.LastPlayed, data[off:])
					off += 8

					foundFields1[1] = true

				case "difficulty_settings":
					if foundFields1[2] {
						return 0, fmt.Errorf("on field DifficultySettings (nbt: difficulty_settings): %w", nbt.ErrDuplicateValue)
					}

					if tag1 != nbt.TagCompound {
						return 0, nbt.ErrUnexpectedTag
					}

					foundFields2 := [2]bool{}

					for {
						if off >= len(data) {
							return 0, nbt.ErrUnexpectedEOF
						}

						tag2 := data[off]
						off += 1

						if tag2 == nbt.TagEnd {
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

						case "difficulty":
							if foundFields2[0] {
								return 0, fmt.Errorf("on field Difficulty (nbt: difficulty): %w", nbt.ErrDuplicateValue)
							}

							if tag2 != nbt.TagString {
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

							nbt.StringFromBytes(&r.Data.DifficultySettings.Difficulty, data[off:off+strLen])
							off += strLen

							foundFields2[0] = true

						case "locked":
							if foundFields2[1] {
								return 0, fmt.Errorf("on field Locked (nbt: locked): %w", nbt.ErrDuplicateValue)
							}

							if tag2 != nbt.TagByte {
								return 0, nbt.ErrUnexpectedTag
							}

							if off+1 > len(data) {
								return 0, nbt.ErrUnexpectedEOF
							}
							nbt.BytePayloadFromBytes(&r.Data.DifficultySettings.Locked, data[off:])
							off += 1

							foundFields2[1] = true

						default:
							k, err := nbt.SkipPayload(tag2, data[off:])
							if err != nil {
								return 0, err
							}
							off += k
						}
					}

					for i := range foundFields2 {
						if !foundFields2[i] {
							return 0, fmt.Errorf("on field %s: %w", [2]string{"Difficulty (nbt: difficulty)", "Locked (nbt: locked)"}[i], nbt.ErrMissingValue)
						}
					}
					foundFields1[2] = true

				case "allowCommands":
					if foundFields1[3] {
						return 0, fmt.Errorf("on field AllowCommands (nbt: allowCommands): %w", nbt.ErrDuplicateValue)
					}

					if tag1 != nbt.TagByte {
						return 0, nbt.ErrUnexpectedTag
					}

					if off+1 > len(data) {
						return 0, nbt.ErrUnexpectedEOF
					}
					nbt.BytePayloadFromBytes(&r.Data.AllowCommands, data[off:])
					off += 1

					foundFields1[3] = true

				case "Version":
					if foundFields1[4] {
						return 0, fmt.Errorf("on field Version: %w", nbt.ErrDuplicateValue)
					}

					if tag1 != nbt.TagCompound {
						return 0, nbt.ErrUnexpectedTag
					}

					foundFields2 := [4]bool{}

					for {
						if off >= len(data) {
							return 0, nbt.ErrUnexpectedEOF
						}

						tag2 := data[off]
						off += 1

						if tag2 == nbt.TagEnd {
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

						case "Id":
							if foundFields2[0] {
								return 0, fmt.Errorf("on field Id: %w", nbt.ErrDuplicateValue)
							}

							if tag2 != nbt.TagInt {
								return 0, nbt.ErrUnexpectedTag
							}

							if off+4 > len(data) {
								return 0, nbt.ErrUnexpectedEOF
							}
							nbt.IntPayloadFromBytes(&r.Data.Version.Id, data[off:])
							off += 4

							foundFields2[0] = true

						case "Name":
							if foundFields2[1] {
								return 0, fmt.Errorf("on field Name: %w", nbt.ErrDuplicateValue)
							}

							if tag2 != nbt.TagString {
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

							nbt.StringFromBytes(&r.Data.Version.Name, data[off:off+strLen])
							off += strLen

							foundFields2[1] = true

						case "Series":
							if foundFields2[2] {
								return 0, fmt.Errorf("on field Series: %w", nbt.ErrDuplicateValue)
							}

							if tag2 != nbt.TagString {
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

							nbt.StringFromBytes(&r.Data.Version.Series, data[off:off+strLen])
							off += strLen

							foundFields2[2] = true

						case "Snapshot":
							if foundFields2[3] {
								return 0, fmt.Errorf("on field Snapshot: %w", nbt.ErrDuplicateValue)
							}

							if tag2 != nbt.TagByte {
								return 0, nbt.ErrUnexpectedTag
							}

							if off+1 > len(data) {
								return 0, nbt.ErrUnexpectedEOF
							}
							nbt.BytePayloadFromBytes(&r.Data.Version.Snapshot, data[off:])
							off += 1

							foundFields2[3] = true

						default:
							k, err := nbt.SkipPayload(tag2, data[off:])
							if err != nil {
								return 0, err
							}
							off += k
						}
					}

					for i := range foundFields2 {
						if !foundFields2[i] {
							return 0, fmt.Errorf("on field %s: %w", [4]string{"Id", "Name", "Series", "Snapshot"}[i], nbt.ErrMissingValue)
						}
					}
					foundFields1[4] = true

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
					return 0, fmt.Errorf("on field %s: %w", [5]string{"LevelName", "LastPlayed", "DifficultySettings (nbt: difficulty_settings)", "AllowCommands (nbt: allowCommands)", "Version"}[i], nbt.ErrMissingValue)
				}
			}
			foundFields0[0] = true

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
			return 0, fmt.Errorf("on field %s: %w", [1]string{"Data"}[i], nbt.ErrMissingValue)
		}
	}

	return off, nil
}

func (r *Root) TagType() byte {
	return nbt.TagCompound
}
