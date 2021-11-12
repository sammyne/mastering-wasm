package vm

import "errors"

var (
	ErrBadArgs          = errors.New("bad args")
	ErrBadSubOpcode     = errors.New("bad sub-opcode saturated trunc")
	ErrMissingCallFrame = errors.New("miss call frame")
	ErrNoStartFunc      = errors.New("missing start func")
	ErrOperandPop       = errors.New("pop operands")
	ErrUnimplemented    = errors.New("not implemented")
	ErrVarImmutable     = errors.New("immutable variables")
)
