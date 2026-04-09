package generator

import (
	"errors"
	"testing"
)

// TODO: i do not like this test setup
// there has to be a better way than "hiding" stuff in a testdata directory 
// (which i am doing to make the test package(s) not importable)
// (also, this takes annoyingly long to run)

func TestNonUTF8(t *testing.T) {
	helper(t, ErrFieldNames, "NonUTF8Field")
}

func TestSneaky(t *testing.T) {
	helper(t, ErrUnexpectedType, "SneakyNonLocalAlias")
	helper(t, nil, "UnderlyingOfExternalOk")
	helper(t, ErrUnexpectedType, "UnderlyingOfExternalMarshaller")
	helper(t, ErrUnexpectedType, "UnderlyingOfOption")
}

func TestUnsupportedSpecial(t *testing.T) {
	helper(t, ErrUnexpectedType, "Generic")
	helper(t, ErrUnexpectedType, "Interface")
	helper(t, ErrInvalidInput, "SomeConst")
	helper(t, ErrInvalidInput, "SomeVar")

}

func TestNameCollision(t *testing.T) {
	helper(t, ErrFieldNames, "NameCollision")
	helper(t, nil, "SavedNameCollision")
}

func TestEmptyFieldName(t *testing.T) {
	helper(t, nil, "EmptyFieldName")
}

func TestTypeNotFound(t *testing.T) {
	helper(t, ErrTypeNotFound, "DoesNotExist9328459328950")
}

func TestMethodCollision(t *testing.T) {
	helper(t, ErrMethodCollision, "AlreadyHasTagType")
	helper(t, ErrMethodCollision, "AlreadyHasUnmarshalPayload")
	helper(t, nil, "HasTagTypeButOnlyInTargetFile")
}

func TestAllowedTypes(t *testing.T) {
	helper(t, ErrUnexpectedType, "StructBadPointer")
	helper(t, ErrUnexpectedType, "BadTypeField")
	helper(t, ErrUnexpectedType, "BadPrimitive")
	helper(t, ErrUnexpectedType, "StructEmbeddedField")
	helper(t, ErrUnexpectedType, "StructEmbeddedFieldInner")
	helper(t, nil, "StructEmpty")
	helper(t, ErrUnexpectedType, "StructFieldNotUnmarshaller")
	helper(t, nil, "StructFieldNotUnmarshaller", "FieldType")
	helper(t, nil, "IgnoredBadTypeField")
	helper(t, nil, "EverythingStruct", "TargetType")
}

func helper(t *testing.T, expectedErr error, types ...string) {
	t.Helper()

	g := New()
	opts := &Options{
		DryRun:  true,
		OutFile: "./internal/testtypes_nbt_gen.go",
	}

	err := g.Generate(opts, "./internal/", types)

	if !errors.Is(err, expectedErr) {
		t.Fatalf("expected error %v but found %v", expectedErr, err)
	}
}
