package option

import "github.com/toros100/easynbt/nbt"

//go:generate go run ../../. -types=SomeStruct,AllTypesStruct,Named

type SomeStruct struct {
	OptionalIntViaType nbt.Option[int32]
	OptionalIntViaTag  int32 `nbtoptional:""`
	SomeString         string
}

type AllTypesStruct struct {
	Byte          nbt.Option[int8]
	Short         nbt.Option[int16]
	Int           nbt.Option[int32]
	Long          nbt.Option[int64]
	Float         nbt.Option[float32]
	Double        nbt.Option[float64]
	String        nbt.Option[string]
	Compound      nbt.Option[struct{ X int8 }]
	ListPrimitive nbt.Option[[]int8]
	ListNamed     nbt.Option[[]Named]
	Named         nbt.Option[Named]
	// TODO: array types, pointers
}

type Named int16
