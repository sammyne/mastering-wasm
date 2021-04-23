package vm

import "errors"

var (
	ErrBadArgs       = errors.New("bad args")
	ErrBadSubOpcode  = errors.New("bad sub-opcode saturated trunc")
	ErrNoStartFunc   = errors.New("missing start func")
	ErrOperandPop    = errors.New("pop operands")
	ErrUnimplemented = errors.New("not implemented")
)
