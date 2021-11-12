package wasmer

import (
	"fmt"
	"os"

	"github.com/sammyne/mastering-wasm/wavm/types"
)

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
	Start     *types.FuncIdx
	Elements  []types.Element
	Codes     []types.Code
	Data      []types.Data
}

func DecodeModuleFromFile(filename string) (*Module, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	return NewDecoder(data).DecodeModule()
}
