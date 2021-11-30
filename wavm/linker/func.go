package linker

import "github.com/sammyne/mastering-wasm/wavm/types"

type Function interface {
	Type() types.FuncType
	Call(args ...types.WasmVal) ([]types.WasmVal, error)
}
