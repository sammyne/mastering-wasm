package wasmer

import "github.com/sammyne/mastering-wasm/mini-wasmer/types"

type Module struct {
	Magic     uint32
	Version   uint32
	Customs   []types.Custom
	Types     []types.FuncType
	Imports   []types.Import
	Functions []types.TypeIdx
	Tables    []types.Table
	Memories  []types.Memory
	Globals   []types.Global
	Exports   []types.Export
	Start     types.FuncIdx // math.MaxUint32 means no start
	Elements  []types.Element
	Codes     []types.Code
	Data      []types.Data
}
