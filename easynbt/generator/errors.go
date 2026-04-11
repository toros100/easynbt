package generator

import "errors"

var ErrUnexpectedType = errors.New("unexpected type")
var ErrFieldNames = errors.New("invalid field names")
var ErrTypeNotFound = errors.New("type not found")
var ErrMethodCollision = errors.New("nbt.Unmarshaler method collision")
var ErrInvalidInput = errors.New("invalid input")
