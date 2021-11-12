package tools

import "github.com/sammyne/mastering-wasm/wavm/types"

func CountLocals(localsVec []types.Locals) uint64 {
	var n uint64
	for _, v := range localsVec {
		n += uint64(v.N)
	}

	return n
}

func ParseBlockSig(t types.BlockType, moduleTypes []types.FuncType) types.FuncType {
	switch t {
	case types.BlockTypeI32:
		return types.FuncType{ResultTypes: []types.ValueType{types.ValueTypeI32}}
	case types.BlockTypeI64:
		return types.FuncType{ResultTypes: []types.ValueType{types.ValueTypeI64}}
	case types.BlockTypeF32:
		return types.FuncType{ResultTypes: []types.ValueType{types.ValueTypeF32}}
	case types.BlockTypeF64:
		return types.FuncType{ResultTypes: []types.ValueType{types.ValueTypeF64}}
	case types.BlockTypeEmpty:
		return types.FuncType{}
	default:
	}

	return moduleTypes[t]
}
