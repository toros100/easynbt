package otherpackage

type SomeType struct { // underlying acceptable for codegen
	X int32
}

type IsMarshaller struct { // does implement unmarshaller, but its underlying type is not acceptable
	BadField ***uintptr
}

func (i *IsMarshaller) TagType() byte {
	return 0
}
func (i *IsMarshaller) UnmarshalPayload(data []byte) (int, error) {
	return 0, nil
}
