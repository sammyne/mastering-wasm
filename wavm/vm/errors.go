package vm

import "errors"

var (
	ErrBadArgs          = errors.New("bad args")
	ErrBadSubOpcode     = errors.New("bad sub-opcode saturated trunc")
	ErrBadValue         = errors.New("bad value")
	ErrBadValueType     = errors.New("bad value type")
	ErrIndexOutOfBound  = errors.New("index out of bound")
	ErrMissingCallFrame = errors.New("miss call frame")
	ErrNoStartFunc      = errors.New("missing start func")
	ErrOperandPop       = errors.New("pop operands")
	ErrUnimplemented    = errors.New("not implemented")
	ErrVarImmutable     = errors.New("immutable variables")
)
