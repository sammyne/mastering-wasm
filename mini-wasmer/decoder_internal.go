package wasmer

import (
	"fmt"
	"math"

	"github.com/sammyne/mastering-wasm/mini-wasmer/types"
)

func (d *Decoder) decodeCode() (*types.Code, error) {
	data, err := d.DecodeBytes()
	if err != nil {
		return nil, fmt.Errorf("decode data: %w", err)
	}

	dd := NewDecoder(data)

	locals, err := dd.decodeLocalsVec()
	if err != nil {
		return nil, fmt.Errorf("decode locals vec: %w", err)
	}

	out := &types.Code{Locals: locals}
	if n := out.LocalCount(); n >= math.MaxUint32 {
		return nil, fmt.Errorf("too many locals: %d", n)
	}

	return out, nil
}

func (d *Decoder) decodeCodes() ([]types.Code, error) {
	n, err := d.DecodeUvarint32()
	if err != nil {
		return nil, fmt.Errorf("decode #(code): %w", err)
	}

	out := make([]types.Code, n)
	for i := range out {
		c, err := d.decodeCode()
		if err != nil {
			return nil, fmt.Errorf("%d-th code: %w", i, err)
		}
		out[i] = *c
	}

	return out, nil
}

func (d *Decoder) decodeCustom() (*types.Custom, error) {
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
	n, err := d.DecodeUvarint32()
	if err != nil {
		return nil, fmt.Errorf("decode #(data): %w", err)
	}

	out := make([]types.Data, n)
	for i := range out {
		v, err := d.decodeDatum()
		if err != nil {
			return nil, fmt.Errorf("%d-th datum: %w", i, err)
		}
		out[i] = *v
	}

	return out, nil
}

func (d *Decoder) decodeDatum() (*types.Data, error) {
	memoryIdx, err := d.DecodeUvarint32()
	if err != nil {
		return nil, fmt.Errorf("decode memory idx: %w", err)
	}

	offset, err := d.decodeExpr()
	if err != nil {
		return nil, fmt.Errorf("decode expr: %w", err)
	}

	init, err := d.DecodeBytes()
	if err != nil {
		return nil, fmt.Errorf("decode init: %w", err)
	}

	out := &types.Data{MemoryIdx: memoryIdx, Offset: offset, Init: init}
	return out, nil
}

func (d *Decoder) decodeElement() (*types.Element, error) {
	tableIdx, err := d.DecodeUint32()
	if err != nil {
		return nil, fmt.Errorf("decode table index: %w", err)
	}

	offset, err := d.decodeExpr()
	if err != nil {
		return nil, fmt.Errorf("decode expr: %w", err)
	}

	init, err := d.decodeIndices()
	if err != nil {
		return nil, fmt.Errorf("decode indices: %w", err)
	}

	out := &types.Element{TableIdx: tableIdx, Offset: offset, Init: init}
	return out, nil
}

func (d *Decoder) decodeElements() ([]types.Element, error) {
	n, err := d.DecodeUvarint32()
	if err != nil {
		return nil, fmt.Errorf("decode #(element): %w", err)
	}

	out := make([]types.Element, n)
	for i := range out {
		v, err := d.decodeElement()
		if err != nil {
			return nil, fmt.Errorf("%d-th element: %w", i, err)
		}
		out[i] = *v
	}

	return out, nil
}

func (d *Decoder) decodeExport() (*types.Export, error) {
	name, err := d.DecodeName()
	if err != nil {
		return nil, fmt.Errorf("decode name: %w", err)
	}

	description, err := d.decodeExportDescription()
	if err != nil {
		return nil, fmt.Errorf("decode export description: %w", err)
	}

	out := &types.Export{Name: name, Description: *description}
	return out, nil
}

func (d *Decoder) decodeExportDescription() (*types.ExportDescription, error) {
	tag, err := d.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("decode tag: %w", err)
	}

	switch tag {
	case types.PortTagFunc, types.PortTagTable, types.PortTagMemory, types.PortTagGlobal:
	default:
		return nil, fmt.Errorf("invalid tag: %v", tag)
	}

	idx, err := d.DecodeUint32()
	if err != nil {
		return nil, fmt.Errorf("decode uint32: %w", err)
	}

	out := &types.ExportDescription{Tag: tag, Idx: idx}
	return out, nil
}

func (d *Decoder) decodeExports() ([]types.Export, error) {
	n, err := d.DecodeUvarint32()
	if err != nil {
		return nil, fmt.Errorf("decode #(export): %w", err)
	}

	out := make([]types.Export, n)
	for i := range out {
		v, err := d.decodeExport()
		if err != nil {
			return nil, fmt.Errorf("%d-th export: %w", i, err)
		}
		out[i] = *v
	}

	return out, nil
}

func (d *Decoder) decodeExpr() (*types.Expr, error) {
	var b byte
	var err error
	for b != 0x0B {
		b, err = d.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("read byte: %w", err)
		}
	}

	return nil, nil
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

func (d *Decoder) decodeGlobalType() (*types.GlobalType, error) {
	valueType, err := d.decodeValueType()
	if err != nil {
		return nil, fmt.Errorf("decode value type: %w", err)
	}

	mutable, err := d.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("read mut: %w", err)
	}
	switch mutable {
	case types.MutConst, types.MutVar:
	default:
		return nil, fmt.Errorf("bad mutability: %d", mutable)
	}

	out := &types.GlobalType{ValueType: valueType, Mutable: mutable}
	return out, nil
}

func (d *Decoder) decodeGlobals() ([]types.Global, error) {
	n, err := d.DecodeUvarint32()
	if err != nil {
		return nil, fmt.Errorf("decode #(global): %w", err)
	}

	out := make([]types.Global, n)
	for i := range out {
		t, err := d.decodeGlobalType()
		if err != nil {
			return nil, fmt.Errorf("decode %d-th global's type: %w", i, err)
		}

		init, err := d.decodeExpr()
		if err != nil {
			return nil, fmt.Errorf("decode %d-th global's init: %w", i, err)
		}

		out[i] = types.Global{Type: *t, Init: init}
	}

	return out, nil
}

func (d *Decoder) decodeImport() (*types.Import, error) {
	module, err := d.DecodeName()
	if err != nil {
		return nil, fmt.Errorf("decode module name: %w", err)
	}

	name, err := d.DecodeName()
	if err != nil {
		return nil, fmt.Errorf("decode name: %w", err)
	}

	description, err := d.decodeImportDescription()
	if err != nil {
		return nil, fmt.Errorf("decode description: %w", err)
	}

	out := &types.Import{Module: module, Name: name, Description: *description}
	return out, nil
}

func (d *Decoder) decodeImportDescription() (*types.ImportDescription, error) {
	tag, err := d.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("decode tag: %w", err)
	}

	out := &types.ImportDescription{Tag: tag}
	switch tag {
	case types.PortTagFunc:
		out.Func, err = d.DecodeUvarint32()
	case types.PortTagTable:
		err = d.decodeTable(&out.Table)
	case types.PortTagMemory:
		var v *types.Memory
		v, err = d.decodeLimits()
		if err == nil {
			out.Memory = *v
		}
	case types.PortTagGlobal:
		var v *types.GlobalType
		v, err = d.decodeGlobalType()
		if err == nil {
			out.Global = *v
		}
	default:
		return nil, fmt.Errorf("bad tag: %d", tag)
	}

	return out, err
}

func (d *Decoder) decodeImports() ([]types.Import, error) {
	n, err := d.DecodeUvarint32()
	if err != nil {
		return nil, fmt.Errorf("decode #(import): %w", err)
	}

	out := make([]types.Import, n)
	for i := range out {
		v, err := d.decodeImport()
		if err != nil {
			return nil, fmt.Errorf("%d-th import: %w", i, err)
		}
		out[i] = *v
	}

	return out, nil
}

func (d *Decoder) decodeIndices() ([]uint32, error) {
	n, err := d.DecodeUvarint32()
	if err != nil {
		return nil, fmt.Errorf("decode #(indices): %w", err)
	}

	out := make([]uint32, n)
	for i := range out {
		out[i], err = d.DecodeUvarint32()
		if err != nil {
			return nil, fmt.Errorf("%d-th index: %w", i, err)
		}
	}

	return out, nil
}

func (d *Decoder) decodeLimits() (*types.Limits, error) {
	tag, err := d.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("decode tag: %w", err)
	}

	min, err := d.DecodeUvarint32()
	if err != nil {
		return nil, fmt.Errorf("decode min: %w", err)
	}

	var max uint32
	if tag == 1 {
		if max, err = d.DecodeUvarint32(); err != nil {
			return nil, fmt.Errorf("decode max: %w", err)
		}
	}

	out := &types.Limits{Tag: tag, Min: min, Max: max}
	return out, nil
}

func (d *Decoder) decodeLocals() (*types.Locals, error) {
	n, err := d.DecodeUvarint32()
	if err != nil {
		return nil, fmt.Errorf("decode N: %w", err)
	}

	_type, err := d.decodeValueType()
	if err != nil {
		return nil, fmt.Errorf("decode value type: %w", err)
	}

	out := &types.Locals{N: n, Type: _type}
	return out, nil
}

func (d *Decoder) decodeLocalsVec() ([]types.Locals, error) {
	n, err := d.DecodeUvarint32()
	if err != nil {
		return nil, fmt.Errorf("decode #(locals): %w", err)
	}

	out := make([]types.Locals, n)
	for i := range out {
		v, err := d.decodeLocals()
		if err != nil {
			return nil, fmt.Errorf("decode %d-th locals: %w", i, err)
		}
		out[i] = *v
	}

	return out, nil
}

func (d *Decoder) decodeMemories() ([]types.Memory, error) {
	n, err := d.DecodeUvarint32()
	if err != nil {
		return nil, fmt.Errorf("decode #(memory): %w", err)
	}

	out := make([]types.Memory, n)
	for i := range out {
		v, err := d.decodeLimits()
		if err != nil {
			return nil, fmt.Errorf("decode %d-th limits: %w", i, err)
		}
		out[i] = *v
	}

	return out, nil
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
	funcIdx, err := d.DecodeUvarint32()
	if err != nil {
		funcIdx = math.MaxUint32
	}

	return funcIdx, err
}

func (d *Decoder) decodeTable(t *types.Table) error {
	elemType, err := d.ReadByte()
	if err != nil {
		return fmt.Errorf("decode element type: %w", err)
	} else if elemType != types.FuncRef {
		return fmt.Errorf("invalid element type: expect %d, got %d", types.FuncRef, elemType)
	}

	limits, err := d.decodeLimits()
	if err != nil {
		return fmt.Errorf("decode limits: %w", err)
	}

	t.ElementType, t.Limits = elemType, *limits
	return nil
}

func (d *Decoder) decodeTables() ([]types.Table, error) {
	n, err := d.DecodeUvarint32()
	if err != nil {
		return nil, fmt.Errorf("decode #(table): %w", err)
	}

	out := make([]types.Table, n)
	for i := range out {
		if err := d.decodeTable(&out[i]); err != nil {
			return nil, fmt.Errorf("decode %d-th table: %w", i, err)
		}
	}

	return out, nil
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

func (d *Decoder) decodeValueType() (types.ValueType, error) {
	t, err := d.ReadByte()
	if err != nil {
		return types.ValueTypeUnknown, fmt.Errorf("decode type: %w", err)
	}

	switch t {
	case types.ValueTypeI32, types.ValueTypeI64, types.ValueTypeF32, types.ValueTypeF64:
	default:
		return types.ValueTypeUnknown, fmt.Errorf("invalid type: %d", t)
	}

	return t, nil
}

func (d *Decoder) decodeValueTypes() ([]types.ValueType, error) {
	n, err := d.DecodeUvarint32()
	if err != nil {
		return nil, fmt.Errorf("decode #(table): %w", err)
	}

	out := make([]types.ValueType, n)
	for i := range out {
		if out[i], err = d.decodeValueType(); err != nil {
			return nil, fmt.Errorf("decode %d-th table: %w", i, err)
		}
	}

	return out, nil
}
