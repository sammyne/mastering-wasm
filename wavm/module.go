package wavm

import (
	"errors"
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

func (m *Module) GetBlockType(t types.BlockType) (types.FuncType, error) {
	switch t {
	case types.BlockTypeI32:
		return types.FuncType{ResultTypes: []types.ValueType{types.ValueTypeI32}}, nil
	case types.BlockTypeI64:
		return types.FuncType{ResultTypes: []types.ValueType{types.ValueTypeI64}}, nil
	case types.BlockTypeF32:
		return types.FuncType{ResultTypes: []types.ValueType{types.ValueTypeF32}}, nil
	case types.BlockTypeF64:
		return types.FuncType{ResultTypes: []types.ValueType{types.ValueTypeF64}}, nil
	case types.BlockTypeEmpty:
		return types.FuncType{}, nil
	default:
	}

	if int(t) >= len(m.Types) {
		return types.FuncType{}, errors.New("index out of bound")
	}

	return types.FuncType{}, errors.New("no such block")
}

func DecodeModuleFromFile(filename string) (*Module, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	return NewDecoder(data).DecodeModule()
}
