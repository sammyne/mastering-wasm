package wasmer

import (
	"errors"
	"fmt"
	"math"

	"github.com/sammyne/mastering-wasm/wavm/tools"
	"github.com/sammyne/mastering-wasm/wavm/types"
)

func (d *Decoder) decodeArgs(opcode byte) (interface{}, error) {
	var out interface{}
	var err error

	switch opcode {
	case types.OpcodeBlock, types.OpcodeLoop:
		out, err = d.decodeBlock()
	case types.OpcodeIf:
		out, err = d.decodeBlockIf()
	case types.OpcodeBr, types.OpcodeBrIf:
		out, err = d.DecodeUvarint32()
	case types.OpcodeBrTable:
		out, err = d.decodeBreakTable()
	case types.OpcodeCall:
		out, err = d.DecodeUvarint32()
	case types.OpcodeCallIndirect:
		out, err = d.decodeCallIndirectArgs()
	case types.OpcodeLocalGet, types.OpcodeLocalSet, types.OpcodeLocalTee:
		out, err = d.DecodeUvarint32()
	case types.OpcodeGlobalGet, types.OpcodeGlobalSet:
		out, err = d.DecodeUvarint32()
	case types.OpcodeMemorySize, types.OpcodeMemoryGrow:
		out, err = 0, d.decodeZero()
	case types.OpcodeI32Const:
		out, err = d.DecodeVarint32()
	case types.OpcodeI64Const:
		out, err = d.DecodeVarint64()
	case types.OpcodeF32Const:
		out, err = d.DecodeFloat32()
	case types.OpcodeF64Const:
		out, err = d.DecodeFloat64()
	case types.OpcodeTruncSat:
		out, err = d.ReadByte()
	default:
		if opcode >= types.OpcodeI32Load && opcode <= types.OpcodeI64Store32 {
			out, err = d.decodeMemoryArg()
		}
	}

	return out, err
}

func (d *Decoder) decodeBlock() (*types.Block, error) {
	blockType, err := d.decodeBlockType()
	if err != nil {
		return nil, fmt.Errorf("decode block type: %w", err)
	}

	instructions, endOpcode, err := d.decodeInstructions()
	if err != nil {
		return nil, fmt.Errorf("decode instructions: %w", err)
	} else if endOpcode != types.OpcodeEnd {
		return nil, fmt.Errorf("invalid end opcode: %v", endOpcode)
	}

	out := &types.Block{BlockType: blockType, Instructions: instructions}
	return out, nil
}

func (d *Decoder) decodeBlockIf() (*types.BlockIf, error) {
	blockType, err := d.decodeBlockType()
	if err != nil {
		return nil, fmt.Errorf("decode block type: %w", err)
	}

	instructions1, endOpcode, err := d.decodeInstructions()
	if err != nil {
		return nil, fmt.Errorf("decode instructions: %w", err)
	}

	var instructions2 []types.Instruction
	if endOpcode == types.OpcodeElse {
		if instructions2, endOpcode, err = d.decodeInstructions(); err != nil {
			return nil, fmt.Errorf("decode else block: %w", err)
		} else if endOpcode != types.OpcodeEnd {
			return nil, fmt.Errorf("invalid end opcode: %v", endOpcode)
		}
	}

	out := &types.BlockIf{
		BlockType:     blockType,
		Instructions1: instructions1,
		Instructions2: instructions2,
	}

	return out, nil
}

func (d *Decoder) decodeBlockType() (int32, error) {
	t, err := d.DecodeVarint32()
	if err != nil {
		return 0, fmt.Errorf("decode varint32: %w", err)
	}

	if t < 0 {
		switch t {
		case types.BlockTypeI32, types.BlockTypeI64, types.BlockTypeF32, types.BlockTypeF64,
			types.BlockTypeEmpty:
		default:
			return 0, fmt.Errorf("bad block type: %d", t)
		}
	}

	return t, nil
}

func (d *Decoder) decodeBreakTable() (*types.BreakTable, error) {
	labels, err := d.decodeIndices()
	if err != nil {
		return nil, fmt.Errorf("decode indices: %w", err)
	}

	defaultIdx, err := d.DecodeUvarint32()
	if err != nil {
		return nil, fmt.Errorf("decode default: %w", err)
	}

	out := &types.BreakTable{Labels: labels, Default: defaultIdx}
	return out, nil
}

func (d *Decoder) decodeCallIndirectArgs() (uint32, error) {
	typeIdx, err := d.DecodeUvarint32()
	if err != nil {
		return 0, fmt.Errorf("decode type idx: %w", err)
	}

	if err := d.decodeZero(); err != nil {
		return 0, fmt.Errorf("decode zero: %w", err)
	}

	return typeIdx, nil
}

func (d *Decoder) decodeCode(out *types.Code) error {
	n, err := d.DecodeUvarint32()
	if err != nil {
		return fmt.Errorf("decode byte count: %w", err)
	}
	remainingLen := d.Len()

	locals, err := d.decodeLocalsVec()
	if err != nil {
		return fmt.Errorf("decode locals: %w", err)
	}
	if nLocal := tools.CountLocals(locals); nLocal >= math.MaxUint32 {
		return fmt.Errorf("too many locals: %d", n)
	}

	var expr types.Expr
	if err := d.decodeExpr(&expr); err != nil {
		return fmt.Errorf("decode expr: %w", err)
	}

	if d.Len()+int(n) != remainingLen {
		return errors.New("invalid code length")
	}

	out.Locals, out.Expr = locals, expr
	return nil
}

func (d *Decoder) decodeCodes() ([]types.Code, error) {
	n, err := d.DecodeUvarint32()
	if err != nil {
		return nil, fmt.Errorf("decode #(code): %w", err)
	}

	out := make([]types.Code, n)
	for i := range out {
		if err := d.decodeCode(&out[i]); err != nil {
			return nil, fmt.Errorf("%d-th code: %w", i, err)
		}
	}

	return out, nil
}

func (d *Decoder) decodeCustom(out *types.Custom) error {
	buf, err := d.DecodeBytes()
	if err != nil {
		return fmt.Errorf("decode bytes: %w", err)
	}

	dd := NewDecoder(buf)
	name, err := dd.DecodeName()
	if err != nil {
		return fmt.Errorf("decode name: %w", err)
	}

	data := make([]byte, dd.Len())
	_, _ = dd.Read(data)

	out.Name, out.Bytes = name, data
	return nil
}

func (d *Decoder) decodeData() ([]types.Data, error) {
	n, err := d.DecodeUvarint32()
	if err != nil {
		return nil, fmt.Errorf("decode #(data): %w", err)
	}

	out := make([]types.Data, n)
	for i := range out {
		if err := d.decodeDatum(&out[i]); err != nil {
			return nil, fmt.Errorf("%d-th datum: %w", i, err)
		}
	}

	return out, nil
}

func (d *Decoder) decodeDatum(out *types.Data) error {
	memoryIdx, err := d.DecodeUvarint32()
	if err != nil {
		return fmt.Errorf("decode memory idx: %w", err)
	}

	var offset types.Expr
	if err := d.decodeExpr(&offset); err != nil {
		return fmt.Errorf("decode expr: %w", err)
	}

	init, err := d.DecodeBytes()
	if err != nil {
		return fmt.Errorf("decode init: %w", err)
	}

	out.MemoryIdx, out.Offset, out.Init = memoryIdx, offset, init
	return nil
}

func (d *Decoder) decodeElement(out *types.Element) error {
	tableIdx, err := d.DecodeUvarint32()
	if err != nil {
		return fmt.Errorf("decode table index: %w", err)
	}

	var offset types.Expr
	if err := d.decodeExpr(&offset); err != nil {
		return fmt.Errorf("decode expr: %w", err)
	}

	init, err := d.decodeIndices()
	if err != nil {
		return fmt.Errorf("decode indices: %w", err)
	}

	out.TableIdx, out.Offset, out.Init = tableIdx, offset, init
	return nil
}

func (d *Decoder) decodeElements() ([]types.Element, error) {
	n, err := d.DecodeUvarint32()
	if err != nil {
		return nil, fmt.Errorf("decode #(element): %w", err)
	}

	out := make([]types.Element, n)
	for i := range out {
		if err := d.decodeElement(&out[i]); err != nil {
			return nil, fmt.Errorf("%d-th element: %w", i, err)
		}
	}

	return out, nil
}

func (d *Decoder) decodeExport(out *types.Export) error {
	name, err := d.DecodeName()
	if err != nil {
		return fmt.Errorf("decode name: %w", err)
	}

	description, err := d.decodeExportDescription()
	if err != nil {
		return fmt.Errorf("decode export description: %w", err)
	}

	out.Name, out.Description = name, *description
	return nil
}

func (d *Decoder) decodeExportDescription() (*types.ExportDescription, error) {
	tag, err := d.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("decode tag: %w", err)
	}

	switch tag {
	case types.PortTagFunc, types.PortTagTable, types.PortTagMemory, types.PortTagGlobal:
	default:
		return nil, fmt.Errorf("invalid tag: %02x", tag)
	}

	idx, err := d.DecodeUvarint32()
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
		if err := d.decodeExport(&out[i]); err != nil {
			return nil, fmt.Errorf("%d-th export: %w", i, err)
		}
	}

	return out, nil
}

func (d *Decoder) decodeExpr(out *types.Expr) error {
	instructions, opcode, err := d.decodeInstructions()
	if err != nil {
		return fmt.Errorf("decode instructions: %w", err)
	}

	if opcode != types.OpcodeEnd {
		return fmt.Errorf("invalid opcode: expect %v, got %v", types.OpcodeEnd, opcode)
	}

	*out = instructions
	return nil
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

func (d *Decoder) decodeGlobalType(out *types.GlobalType) error {
	valueType, err := d.decodeValueType()
	if err != nil {
		return fmt.Errorf("decode value type: %w", err)
	}

	mutable, err := d.ReadByte()
	if err != nil {
		return fmt.Errorf("read mut: %w", err)
	}
	switch mutable {
	case types.MutConst, types.MutVar:
	default:
		return fmt.Errorf("bad mutability: %d", mutable)
	}

	out.ValueType, out.Mutable = valueType, mutable
	return nil
}

func (d *Decoder) decodeGlobals() ([]types.Global, error) {
	n, err := d.DecodeUvarint32()
	if err != nil {
		return nil, fmt.Errorf("decode #(global): %w", err)
	}

	out := make([]types.Global, n)
	for i := range out {
		if err := d.decodeGlobalType(&out[i].Type); err != nil {
			return nil, fmt.Errorf("decode %d-th global's type: %w", i, err)
		}

		if err := d.decodeExpr(&out[i].Init); err != nil {
			return nil, fmt.Errorf("decode %d-th global's init: %w", i, err)
		}
	}

	return out, nil
}

func (d *Decoder) decodeImport(out *types.Import) error {
	module, err := d.DecodeName()
	if err != nil {
		return fmt.Errorf("decode module name: %w", err)
	}

	name, err := d.DecodeName()
	if err != nil {
		return fmt.Errorf("decode name: %w", err)
	}

	description, err := d.decodeImportDescription()
	if err != nil {
		return fmt.Errorf("decode description: %w", err)
	}

	out.Module, out.Name, out.Description = module, name, *description
	return nil
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
		err = d.decodeLimits(&out.Memory)
	case types.PortTagGlobal:
		err = d.decodeGlobalType(&out.Global)
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
		if err := d.decodeImport(&out[i]); err != nil {
			return nil, fmt.Errorf("%d-th import: %w", i, err)
		}
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

func (d *Decoder) decodeInstruction(out *types.Instruction) error {
	opcode, err := d.ReadByte()
	if err != nil {
		return fmt.Errorf("decode opcode: %w", err)
	} else if _, ok := types.GetOpname(opcode); !ok {
		return fmt.Errorf("unknown opcode: %02x", opcode)
	}

	args, err := d.decodeArgs(opcode)
	if err != nil {
		return fmt.Errorf("decode args: %w", err)
	}

	out.Opcode, out.Args = opcode, args
	return nil
}

func (d *Decoder) decodeInstructions() ([]types.Instruction, byte, error) {
	var out []types.Instruction

	var t types.Instruction
	for {
		if err := d.decodeInstruction(&t); err != nil {
			return nil, 0, fmt.Errorf("decode instruction: %w", err)
		}

		if t.Opcode == types.OpcodeElse || t.Opcode == types.OpcodeEnd {
			return out, t.Opcode, nil
		}

		out = append(out, t)
	}
}

func (d *Decoder) decodeLimits(out *types.Limits) error {
	tag, err := d.ReadByte()
	if err != nil {
		return fmt.Errorf("decode tag: %w", err)
	}

	min, err := d.DecodeUvarint32()
	if err != nil {
		return fmt.Errorf("decode min: %w", err)
	}

	var max uint32
	if tag == 1 {
		if max, err = d.DecodeUvarint32(); err != nil {
			return fmt.Errorf("decode max: %w", err)
		}
	}

	out.Tag, out.Min, out.Max = tag, min, max
	return nil
}

func (d *Decoder) decodeLocals(out *types.Locals) error {
	n, err := d.DecodeUvarint32()
	if err != nil {
		return fmt.Errorf("decode N: %w", err)
	}

	_type, err := d.decodeValueType()
	if err != nil {
		return fmt.Errorf("decode value type: %w", err)
	}

	out.N, out.Type = n, _type
	return nil
}

func (d *Decoder) decodeLocalsVec() ([]types.Locals, error) {
	n, err := d.DecodeUvarint32()
	if err != nil {
		return nil, fmt.Errorf("decode #(locals): %w", err)
	}

	out := make([]types.Locals, n)
	for i := range out {
		if err := d.decodeLocals(&out[i]); err != nil {
			return nil, fmt.Errorf("decode %d-th locals: %w", i, err)
		}
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
		if err := d.decodeLimits(&out[i]); err != nil {
			return nil, fmt.Errorf("decode %d-th limits: %w", i, err)
		}
	}

	return out, nil
}

func (d *Decoder) decodeMemoryArg() (types.MemoryArg, error) {
	align, err := d.DecodeUvarint32()
	if err != nil {
		return types.MemoryArg{}, fmt.Errorf("decode align: %w", err)
	}

	offset, err := d.DecodeUvarint32()
	if err != nil {
		return types.MemoryArg{}, fmt.Errorf("decode offset: %w", err)
	}

	out := types.MemoryArg{Align: align, Offset: offset}
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
		var idx uint32
		idx, err = d.decodeStart()
		if err == nil {
			m.Start = &idx
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

	var limits types.Limits
	if err := d.decodeLimits(&limits); err != nil {
		return fmt.Errorf("decode limits: %w", err)
	}

	t.ElementType, t.Limits = elemType, limits
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
		return types.ValueTypeUnknown, fmt.Errorf("invalid type: %02x-%02x", t, types.ValueTypeI32)
	}

	return t, nil
}

func (d *Decoder) decodeValueTypes() ([]types.ValueType, error) {
	n, err := d.DecodeUvarint32()
	if err != nil {
		return nil, fmt.Errorf("decode #(value types): %w", err)
	}

	out := make([]types.ValueType, n)
	for i := range out {
		if out[i], err = d.decodeValueType(); err != nil {
			return nil, fmt.Errorf("decode %d-th value types: %w", i, err)
		}
	}

	return out, nil
}

func (d *Decoder) decodeZero() error {
	v, err := d.ReadByte()
	if err != nil {
		return fmt.Errorf("read byte: %w", err)
	} else if v != 0 {
		return errors.New("non-zero")
	}

	return nil
}
