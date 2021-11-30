package native

import "github.com/sammyne/mastering-wasm/wavm/types"

type GoFunc = func(args []types.WasmVal) ([]types.WasmVal, error)

type Function struct {
	type_ types.FuncType
	f     GoFunc
}

func (f Function) Call(args ...types.WasmVal) ([]types.WasmVal, error) {
	return f.f(args)
}

func (f Function) Type() types.FuncType {
	return f.type_
}
