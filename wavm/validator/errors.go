package validator

import "errors"

var (
	ErrIndexOutOfBound = errors.New("index out of bound")
	ErrTypeMismatch    = errors.New("type mismatch")
)
