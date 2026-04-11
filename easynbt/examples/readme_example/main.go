package main

import (
	"fmt"
	"log"

	"github.com/k0kubun/pp/v3"
	"github.com/toros100/easynbt/nbt"
)

//go:generate easynbt -types=Data

type Data struct {
	Hello    string `nbt:"hello"`
	Position struct {
		X int32
		Y int32
	} `nbt:"pos"`
	Numbers nbt.Option[[]int8]
}

// Data describes a compound tag with the following child tags:
// 1. a string tag with the name "hello" (the struct tag nbt overrides the field name)
// 2. a compound tag with the name "pos", that itself has two int child tags with names "X" and "Y"
// 3. a list tag, which has the element tag type byte (i.e. a sequence of signed 8 bit integers)
// Because the latter uses nbt.Option as a wrapper type around []int8, this child tag is implicitly optional
// (which could also be achieved by using the struct tag nbtoptional)

func main() {

	data1 := new(Data)

	// We use the helper function nbt.Unmarshal to unmarshal the "full" compound tag, including its name.
	// This internally uses the generated method data1.UnmarshalPayload(...), which only unmarshals
	// the payload portion, without the byte for the tag type and the name.
	// Recall that a compound tags payload is a collection of fully formed tags, in this case the tags
	// described by the fields of the underlying struct type of Data.
	name1, _, err := nbt.Unmarshal(data1, bytes1)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("unmarshaled a tag with name %q:\n", name1)
	pp.Println(data1)

	data2 := new(Data)

	name2, _, err := nbt.Unmarshal(data2, bytes2)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("unmarshaled a tag with name %q:\n", name2)
	pp.Println(data2)

	// for data2, we will see that data1.Numbers.Ok == false,
	// which indicates that the optional child tag Numbers was not found.
	// In this particular case, it would have been possible to detect this by seeing
	// that the value is nil, because that can never be the result of unmarshaling
	// a list payload (even a list tag with length 0 would result in an actual empty slice, not nil).
	// For other types, such as int32, (int tag), we would not be able to distinguish
	// between a 0 that we actually unmarshaled from bytes, or that just being the
	// zero value of the type int32 on a field that was not found.
}

var bytes1 = []byte{
	nbt.TagCompound,
	0, 4,
	'r', 'o', 'o', 't',
	nbt.TagString,
	0, 5, 'h', 'e', 'l', 'l', 'o', 0, 5, 'w', 'o', 'r', 'l', 'd',
	nbt.TagCompound,
	0, 3, 'p', 'o', 's', nbt.TagInt, 0, 1, 'X', 0, 0, 40, 120, nbt.TagInt, 0, 1, 'Y', 255, 0, 0, 0, nbt.TagEnd,
	nbt.TagList, 0, 7, 'N', 'u', 'm', 'b', 'e', 'r', 's', nbt.TagByte, 0, 0, 0, 5, 0, 1, 2, 3, 4,
	nbt.TagEnd,
}

var bytes2 = []byte{
	nbt.TagCompound,
	0, 0,
	nbt.TagString,
	0, 5, 'h', 'e', 'l', 'l', 'o', 0, 5, 'a', 'g', 'a', 'i', 'n',
	nbt.TagCompound,
	0, 3, 'p', 'o', 's', nbt.TagInt, 0, 1, 'X', 0, 0, 55, 0, nbt.TagInt, 0, 1, 'Y', 5, 240, 0, 0, nbt.TagEnd,
	nbt.TagEnd,
}
