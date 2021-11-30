package linker

import "github.com/sammyne/mastering-wasm/wavm/types"

type Module interface {
	GetGlobalVal(name string) (types.WasmVal, error)
	GetMember(name string) (interface{}, error)
	InvokeFunc(name string, args ...types.WasmVal) ([]types.WasmVal, error)
	SetGlobalVal(name string, val types.WasmVal) error
}
