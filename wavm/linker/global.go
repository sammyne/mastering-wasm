package linker

import "github.com/sammyne/mastering-wasm/wavm/types"

type Global interface {
	Get() (types.WasmVal, error)
	GetAsUint64() uint64
	Set(v types.WasmVal) error
	SetAsUint64(v uint64) error
	Type() types.GlobalType
}
