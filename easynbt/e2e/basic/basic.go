package basic

//go:generate go run ../../. -types=BytePayload,ShortPayload,IntPayload,LongPayload,FloatPayload,DoublePayload,StringPayload,ListOfIntPayload,CompoundPayload

type BytePayload int8
type ShortPayload int16
type IntPayload int32
type LongPayload int64
type FloatPayload float32
type DoublePayload float64
type StringPayload string

type ListOfIntPayload []int32

type CompoundPayload struct {
	IntPayloadField    IntPayload
	RawIntField        int32
	InnerCompoundField struct {
		ListOfByteField []int8
	}
}
