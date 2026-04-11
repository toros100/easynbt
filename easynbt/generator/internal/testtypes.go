package testtypes

// TODO: this is a bad setup

import (
	"github.com/toros100/easynbt/easynbt/generator/internal/otherpackage"
	"github.com/toros100/easynbt/nbt"
)

const SomeConst = "hello"

var SomeVar = 1

type SneakyNonLocalAlias = otherpackage.SomeType

type UnderlyingOfExternalOk otherpackage.SomeType
type UnderlyingOfExternalMarshaler otherpackage.IsMarshaler

type UnderlyingOfOption nbt.Option[string]

type Generic[T int32 | string] struct {
	Field T
}

type Interface interface {
	cool()
}

type StructEmbeddedField struct {
	FieldType
}

type BadPrimitive uint32

type StructEmbeddedFieldInner struct {
	Inner struct {
		string
	}
}

type NonUTF8Field struct {
	Field float32 `nbt:"\xff"`
}

type StructEmpty struct {
	Field string `nbtignore:""`
	_     int32  // also ignored
}

type IgnoredBadTypeField struct {
	Field uint `nbtignore:""`
	Other string
}

type BadTypeField struct {
	Field uint
	Other string
}

type NameCollision struct {
	Field int32
	Other string `nbt:"Field"`
}

type SavedNameCollision struct {
	Field int32  `nbtignore:""`
	Other string `nbt:"Field"`
}

type EmptyFieldName struct {
	Field float64 `nbt:""`
}

type StructBadPointer struct {
	Field *int32
}

type StructFieldNotUnmarshaler struct {
	Field FieldType
}

type FieldType string

type AlreadyHasTagType string

func (a *AlreadyHasTagType) TagType() byte {
	return 10
}

type AlreadyHasUnmarshalPayload int32

func (a AlreadyHasUnmarshalPayload) UnmarshalPayload(args, do, not, matter string) {}

type HasTagTypeButOnlyInTargetFile int8

type IsMarshaler struct{}

func (i *IsMarshaler) TagType() byte {
	return 0
}
func (i *IsMarshaler) UnmarshalPayload(data []byte) (int, error) {
	return 0, nil
}

type TargetType int8
type PrimitiveAlias = int8
type TargetAlias = TargetType
type MarshalerAlias = IsMarshaler
type OptionAlias = nbt.Option[int8]
type ListAlias = []int8
type CompoundAlias = struct {
	X int32
	Y string
}
type OptionInnerAlias = nbt.Option[CompoundAlias]

type EverythingStruct struct {
	_                    BadPrimitive
	B                    BadPrimitive `nbtignore:""`
	OkButIgnored         int8         `nbtignore:""`
	NonUTF8ButIgnored    int8         `nbt:"\xff" nbtignore:""`
	PrimitiveAlias       PrimitiveAlias
	TargetAlias          TargetAlias
	MarshalerAlias       MarshalerAlias
	OptionAlias          OptionAlias
	ListAlias            ListAlias
	CompoundAlias        CompoundAlias
	Byte                 int8
	Short                int16
	Int                  int32
	Long                 int64
	Float                float32
	Double               float64
	String               string
	Compound             struct{ X int8 }
	Named                IsMarshaler
	PointerNamed         *IsMarshaler
	NamedExternal        otherpackage.IsMarshaler
	PointerNamedExternal *otherpackage.IsMarshaler
	ListPrimitive        []int8
	ListNamed            []IsMarshaler
	ListPointer          []*IsMarshaler
	ListList             [][]int8
	ListNamedTarget      []TargetType
	ListPointerTarget    []*TargetType
	ListCompound         []struct{ S string }
	OptionByte           nbt.Option[int8]
	OptionShort          nbt.Option[int16]
	OptionInt            nbt.Option[int32]
	OptionLong           nbt.Option[int64]
	OptionFloat          nbt.Option[float32]
	OptionDouble         nbt.Option[float64]
	OptionString         nbt.Option[string]
	OptionList           nbt.Option[[]int8]
	OptionNamed          nbt.Option[IsMarshaler]
	OptionPointer        nbt.Option[*IsMarshaler]
	OptionTarget         nbt.Option[TargetType]
	OptionPointerTarget  nbt.Option[*TargetType]
	OptionCompound       nbt.Option[struct{ X int8 }]
}
