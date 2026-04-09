package main

import (
	"bytes"
	"compress/gzip"
	"log"
	"os"

	"github.com/k0kubun/pp/v3"
	"github.com/toros100/easynbt/nbt"
)

//go:generate easynbt -types=Root

// cf. https://minecraft.wiki/w/Java_Edition_level_format#level.dat_format
type Root struct {
	Data struct {
		LevelName          string
		LastPlayed         int64
		DifficultySettings struct {
			Difficulty string `nbt:"difficulty"`
			Locked     int8   `nbt:"locked"`
		} `nbt:"difficulty_settings"`
		AllowCommands int8 `nbt:"allowCommands"`
		Version       struct {
			Id       int32
			Name     string
			Series   string
			Snapshot int8
		}
	}
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
