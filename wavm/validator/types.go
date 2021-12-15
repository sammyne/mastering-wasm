package validator

import "github.com/sammyne/mastering-wasm/wavm/types"

type ControlFrame struct {
	Opcode      byte
	SatrtTypes  []types.ValueType
	EndTypes    []types.ValueType
	Height      int
	Unreachable bool
}

type ControlStack = []ControlFrame

type OperandStack = []types.ValueType

func (f *ControlFrame) LabelTypes() []types.ValueType {
	if f.Opcode == types.OpcodeLoop {
		return f.SatrtTypes
	}

	return f.EndTypes
}
