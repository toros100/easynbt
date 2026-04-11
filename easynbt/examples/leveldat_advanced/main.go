package main

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"log"
	"os"
	"time"

	"github.com/k0kubun/pp/v3"
	"github.com/toros100/easynbt/nbt"
)

//go:generate easynbt -types=Root,Version

// cf. https://minecraft.wiki/w/Java_Edition_level_format#level.dat_format
type Root struct {
	Data struct {
		LevelName          string
		LastPlayed         Timestamp
		DifficultySettings DifficultySettings `nbt:"difficulty_settings"`
		AllowCommands      Bool               `nbt:"allowCommands"`
		Version            Version
	}
}

// a named type can be defined to "extract" a nested struct type
// caution: in this case, the type does need to implement nbt.Unmarshaler,
// which we accomplished by adding it to the target types of easynbt
type Version struct {
	Id       int32
	Name     string
	Series   string
	Snapshot Bool
}

// similarly, alias types may be used to improve readability of nested structs,
// but without any additional requirements: it will behave exactly as if the aliased was inlined
type DifficultySettings = struct {
	Difficulty string `nbt:"difficulty"`
	Locked     Bool   `nbt:"locked"`
}

// a timestamp type with custom nbt.Unmarshaler implementation that represents a long payload as a string
type Timestamp string

func (t *Timestamp) TagType() byte {
	// needs to return nbt.TagLong, because the corresponding nbt data is still a long payload
	// note that the underlying type of Timestamp is "string", which would correspond with
	// nbt.TagString, if Timestamp were a code generation target. but when writing a custom Unmarshaler,
	// the underlying type does not need to "match" the value returned by TagType
	return nbt.TagLong
}

func (t *Timestamp) UnmarshalPayload(data []byte) (int, error) {
	if len(data) < 8 {
		return 0, nbt.ErrUnexpectedEOF
	}
	val := int64(binary.BigEndian.Uint64(data))
	*t = Timestamp(time.UnixMilli(val).Local().Format(time.RFC1123))
	return 8, nil
}

// alternatively:
// type Timestamp2 int64
// use easynbt to generate nbt.Unmarshaler methods for this type
// implement a .Format() or .String() method, e.g:
// func (t *Timestamp2) Format(layout string) string {
// 		return time.UnixMilli(int64(*t)).Format(layout)
// }

// the NBT format does not have a boolean type, but sometimes byte payloads
// are interpreted as such
type Bool bool

func (b *Bool) TagType() byte {
	return nbt.TagByte
}

func (b *Bool) UnmarshalPayload(data []byte) (int, error) {
	if len(data) == 0 {
		return 0, nbt.ErrUnexpectedEOF
	}
	if data[0] == 0 {
		*b = false
	} else {
		*b = true
	}
	return 1, nil
}

func main() {
	data, err := os.ReadFile("../testdata/level.dat")
	if err != nil {
		log.Fatal(err)
	}

	r, err := gzip.NewReader(bytes.NewReader(data))

	if err != nil {
		log.Fatal(err)
	}

	buf := bytes.NewBuffer(nil)

	_, err = buf.ReadFrom(r)
	r.Close()

	decompressed := buf.Bytes()

	root := new(Root)

	_, _, err = nbt.Unmarshal(root, decompressed)

	if err != nil {
		log.Fatal(err)
	}

	pp.Println(root)
}
