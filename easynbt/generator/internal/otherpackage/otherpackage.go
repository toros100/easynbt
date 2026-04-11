package otherpackage

type SomeType struct { // underlying acceptable for codegen
	X int32
}

type IsMarshaler struct { // does implement unmarshaler, but its underlying type is not acceptable
	BadField ***uintptr
}

func (i *IsMarshaler) TagType() byte {
	return 0
}
func (i *IsMarshaler) UnmarshalPayload(data []byte) (int, error) {
	return 0, nil
}
