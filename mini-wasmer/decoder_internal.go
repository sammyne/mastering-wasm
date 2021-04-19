package wasmer

import (
	"fmt"
	"math"

	"github.com/sammyne/mastering-wasm/mini-wasmer/types"
)

func (d *Decoder) decodeCode() (*types.Code, error) {
	panic("todo")
}

func (d *Decoder) decodeCodes() ([]types.Code, error) {
	panic("todo")
}

func (d *Decoder) decodeCustomSection() (*types.Custom, error) {
	buf, err := d.DecodeBytes()
	if err != nil {
		return nil, fmt.Errorf("decode bytes: %w", err)
	}

	dd := NewDecoder(buf)
	name, err := dd.DecodeName()
	if err != nil {
		return nil, fmt.Errorf("decode name: %w", err)
	}

	data := make([]byte, dd.Len())
	_, _ = dd.Read(data)

	out := &types.Custom{Name: name, Bytes: data}
	return out, nil
}

func (d *Decoder) decodeData() ([]types.Data, error) {
	panic("todo")
}

func (d *Decoder) decodeDatum() (*types.Data, error) {
	panic("todo")
}

func (d *Decoder) decodeElement() (*types.Element, error) {
	panic("todo")
}

func (d *Decoder) decodeElements() ([]types.Element, error) {
	panic("todo")
}

func (d *Decoder) decodeExport() (*types.Export, error) {
	panic("todo")
}

func (d *Decoder) decodeExportDescription() (*types.ExportDescription, error) {
	panic("todo")
}

func (d *Decoder) decodeExports() ([]types.Export, error) {
	panic("todo")
}

func (d *Decoder) decodeExpr() (*types.Expr, error) {
	panic("todo")
}

func (d *Decoder) decodeFuncType() (*types.FuncType, error) {
	tag, err := d.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("read tag: %w", err)
	}

	paramTypes, err := d.decodeValueTypes()
	if err != nil {
		return nil, fmt.Errorf("decode parameter types: %w", err)
	}

	resultTypes, err := d.decodeValueTypes()
	if err != nil {
		return nil, fmt.Errorf("decode result types: %w", err)
	}

	out := &types.FuncType{Tag: tag, ParamTypes: paramTypes, ResultTypes: resultTypes}
	return out, nil
}

func (d *Decoder) decodeGlobal() (*types.Global, error) {
	panic("todo")
}

func (d *Decoder) decodeGlobalType() (*types.GlobalType, error) {
	panic("todo")
}

func (d *Decoder) decodeGlobals() ([]types.Global, error) {
	panic("todo")
}

func (d *Decoder) decodeImport() (*types.Import, error) {
	panic("todo")
}

func (d *Decoder) decodeImportDescription() (*types.ImportDescription, error) {
	panic("todo")
}

func (d *Decoder) decodeImports() ([]types.Import, error) {
	panic("todo")
}

func (d *Decoder) decodeIndices() ([]uint32, error) {
	panic("todo")
}

func (d *Decoder) decodeLimits() (*types.Limits, error) {
	panic("todo")
}

func (d *Decoder) decodeLocals() (*types.Locals, error) {
	panic("todo")
}

func (d *Decoder) decodeLocalsVec() ([]types.Locals, error) {
	panic("todo")
}

func (d *Decoder) decodeMemories() ([]types.Memory, error) {
	panic("todo")
}

func (d *Decoder) decodeNonCustomSectionIntoModule(ID byte, m *Module) error {
	var err error

	switch ID {
	case types.SectionIDType:
		m.Types, err = d.decodeTypes()
	case types.SectionIDImport:
		m.Imports, err = d.decodeImports()
	case types.SectionIDFunc:
		m.Functions, err = d.decodeIndices()
	case types.SectionIDTable:
		m.Tables, err = d.decodeTables()
	case types.SectionIDMemory:
		m.Memories, err = d.decodeMemories()
	case types.SectionIDGlobal:
		m.Globals, err = d.decodeGlobals()
	case types.SectionIDExport:
		m.Exports, err = d.decodeExports()
	case types.SectionIDStart:
		m.Start, err = d.decodeStart()
		if err != nil {
			m.Start = math.MaxUint32
		}
	case types.SectionIDElement:
		m.Elements, err = d.decodeElements()
	case types.SectionIDCode:
		m.Codes, err = d.decodeCodes()
	case types.SectionIDData:
		m.Data, err = d.decodeData()
	default:
		err = fmt.Errorf("invalid section ID(%v)", ID)
	}

	return err
}

func (d *Decoder) decodeStart() (uint32, error) {
	panic("todo")
}

func (d *Decoder) decodeTables() ([]types.Table, error) {
	panic("todo")
}

func (r *Decoder) decodeTypes() ([]types.FuncType, error) {
	n, err := r.DecodeVarint32()
	if err != nil {
		return nil, fmt.Errorf("read #(types): %w", err)
	}

	out := make([]types.FuncType, n)
	for i := range out {
		v, err := r.decodeFuncType()
		if err != nil {
			return nil, fmt.Errorf("decode %d-th func type: %w", i, err)
		}
		out[i] = *v
	}

	return out, nil
}

func (d *Decoder) decodeValueType() (*types.ValueType, error) {
	panic("todo")
}

func (d *Decoder) decodeValueTypes() ([]types.ValueType, error) {
	panic("todo")
}
