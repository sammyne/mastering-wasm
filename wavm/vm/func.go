package vm

import (
	"github.com/sammyne/mastering-wasm/wavm/linker"
	"github.com/sammyne/mastering-wasm/wavm/types"
)

type Func struct {
	type_      types.FuncType
	code       types.Code
	externalFn linker.Function // efn is an external function
	ctx        *VM
}

func (f Func) Call(args ...types.WasmVal) ([]types.WasmVal, error) {
	if f.externalFn != nil {
		return f.externalFn.Call(args...)
	}

	return f.call(args)
}

func (f Func) Type() types.FuncType {
	return f.type_
}
