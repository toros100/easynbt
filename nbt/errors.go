package nbt

import "errors"

var ErrMissingValue = errors.New("no value found")
var ErrDuplicateValue = errors.New("duplicate value found")
var ErrUnexpectedEOF = errors.New("unexpected EOF")
var ErrUnexpectedTag = errors.New("unexpected tag type byte")
var ErrInvalidUTF8 = errors.New("string bytes are not valid UTF-8")
var ErrInvalidLength = errors.New("invalid list or array length (negative)")
