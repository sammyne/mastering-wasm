package validator

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/sammyne/mastering-wasm/wavm/types"
)

type codeValidator struct {
	OperandStack OperandStack
	ControlStack ControlStack

	moduleValidator *moduleValidator
	Idx             int
	InstructionPath map[int]string

	localLen int
}

func (cv *codeValidator) checkAlign(bitWidth int, args interface{}) error {
	align := args.(types.MemoryArg).Align
	if a, b := 1<<align, bitWidth/8; a > b {
		return fmt.Errorf("alignment(%d) must be smaller than natural alignment(%d)", a, b)
	}

	return nil
}

func (cv *codeValidator) f32Load(args interface{}, bitWidth int) error {
	return cv.load(types.ValueTypeF32, bitWidth, args)
}

func (cv *codeValidator) f32Store(args interface{}, bits int) error {
	return cv.store(types.ValueTypeF32, bits, args)
}

func (cv *codeValidator) f64Load(args interface{}, bits int) error {
	return cv.load(types.ValueTypeF64, bits, args)
}

func (cv *codeValidator) f64Store(args interface{}, bits int) error {
	return cv.store(types.ValueTypeF64, bits, args)
}

func (cv *codeValidator) getControlFrame(idx int) (ControlFrame, error) {
	if idx >= len(cv.ControlStack) {
		return ControlFrame{}, fmt.Errorf("idx bound is %d: %w", len(cv.ControlStack), ErrIndexOutOfBound)
	}

	return cv.ControlStack[idx], nil
}

func (cv *codeValidator) hasMemory() bool {
	return cv.moduleValidator.getMemoryLen() > 0
}

func (cv *codeValidator) i32Load(args interface{}, bitWidth int) error {
	return cv.load(types.ValueTypeI32, bitWidth, args)
}

func (cv *codeValidator) i32Store(args interface{}, bits int) error {
	return cv.store(types.ValueTypeI32, bits, args)
}

func (cv *codeValidator) i64Load(args interface{}, bits int) error {
	return cv.load(types.ValueTypeI64, bits, args)
}

func (cv *codeValidator) i64Store(args interface{}, bits int) error {
	return cv.store(types.ValueTypeI64, bits, args)
}

func (cv *codeValidator) load(vt types.ValueType, bits int, args interface{}) error {
	if cv.moduleValidator.getMemoryLen() == 0 {
		return errors.New("no usable memory")
	}
	if err := cv.checkAlign(bits, args); err != nil {
		return fmt.Errorf("bad alignment: %w", err)
	}
	if err := cv.popI32(); err != nil {
		return fmt.Errorf("pop i32: %w", err)
	}
	cv.pushOperand(vt)

	return nil
}

func (cv *codeValidator) popControlFrame() (ControlFrame, error) {
	f, err := cv.getControlFrame(0)
	if err != nil {
		return ControlFrame{}, fmt.Errorf("pop from stack: %w", err)
	}

	if err := cv.popOperands(f.EndTypes); err != nil {
		return ControlFrame{}, fmt.Errorf("pop results: %w", err)
	}

	if len(cv.OperandStack) != f.Height {
		return ControlFrame{}, fmt.Errorf("wrong height(%d!=%d) of operand stack", f.Height, len(cv.OperandStack))
	}

	cv.ControlStack = cv.ControlStack[:len(cv.ControlStack)-1]
	return f, nil
}

func (cv *codeValidator) popF32() error {
	_, err := cv.popTypeSpecificOperand(types.ValueTypeF32)
	return err
}

func (cv *codeValidator) popF64() error {
	_, err := cv.popTypeSpecificOperand(types.ValueTypeF64)
	return err
}

func (cv *codeValidator) popI32() error {
	_, err := cv.popTypeSpecificOperand(types.ValueTypeI32)
	return err
}

func (cv *codeValidator) popI64() error {
	_, err := cv.popTypeSpecificOperand(types.ValueTypeI64)
	return err
}

func (cv *codeValidator) popOperand() (types.ValueType, error) {
	if f, err := cv.getControlFrame(0); err != nil {
		return types.ValueTypeUnknown, fmt.Errorf("get the latest control frame: %w", err)
	} else if f.Height == len(cv.OperandStack) {
		if f.Unreachable {
			return types.ValueTypeUnknown, nil
		}

		return types.ValueTypeUnknown, errors.New("type mismatch")
	}

	r := cv.OperandStack[len(cv.OperandStack)-1]
	cv.OperandStack = cv.OperandStack[:len(cv.OperandStack)-1]
	return r, nil
}

func (cv *codeValidator) popOperands(expected []types.ValueType) error {
	for i := len(expected) - 1; i >= 0; i-- {
		if _, err := cv.popTypeSpecificOperand(expected[i]); err != nil {
			return fmt.Errorf("pop %d-th operand: %w", i, err)
		}
	}

	return nil
}

func (cv *codeValidator) popTowThenPushOneType(t types.ValueType) error {
	if _, err := cv.popTypeSpecificOperand(t); err != nil {
		return fmt.Errorf("pop 1st operand: %w", err)
	}
	if _, err := cv.popTypeSpecificOperand(t); err != nil {
		return fmt.Errorf("pop 2nd operand: %w", err)
	}

	cv.pushOperand(t)

	return nil
}

func (cv *codeValidator) popThenPush(popType, pushType types.ValueType) error {
	if _, err := cv.popTypeSpecificOperand(popType); err != nil {
		return fmt.Errorf("pop operand: %w", err)
	}
	cv.pushOperand(pushType)
	return nil
}

func (cv *codeValidator) popTypeSpecificOperand(expect types.ValueType) (types.ValueType, error) {
	got, err := cv.popOperand()
	if err != nil {
		return types.ValueTypeUnknown, fmt.Errorf("pop generic operand: %w", err)
	}

	switch {
	case got == types.ValueTypeUnknown:
		return expect, nil
	case expect == types.ValueTypeUnknown:
		return got, nil
	case expect != got:
		return types.ValueTypeUnknown, fmt.Errorf("got type=%d: %w", got, ErrTypeMismatch)
	default:
	}

	return got, nil
}

func (cv *codeValidator) pushControlFrame(opcode byte, in, out []types.ValueType) {
	f := ControlFrame{
		Opcode:      opcode,
		SatrtTypes:  in,
		EndTypes:    out,
		Height:      len(cv.OperandStack),
		Unreachable: false,
	}

	cv.ControlStack = append(cv.ControlStack, f)
	cv.pushOperands(in)
}

func (cv *codeValidator) pushOperand(v types.ValueType) {
	cv.OperandStack = append(cv.OperandStack, v)
}

func (cv *codeValidator) pushOperands(values []types.ValueType) {
	for _, v := range values {
		cv.pushOperand(v)
	}
}

func (cv *codeValidator) store(vt types.ValueType, bits int, args interface{}) error {
	if cv.moduleValidator.getMemoryLen() == 0 {
		return errors.New("no usable memory")
	}
	if err := cv.checkAlign(bits, args); err != nil {
		return fmt.Errorf("bad alignment: %w", err)
	}

	if _, err := cv.popTypeSpecificOperand(vt); err != nil {
		return fmt.Errorf("pop operand: %w", err)
	}
	if err := cv.popI32(); err != nil {
		return fmt.Errorf("pop i32: %w", err)
	}

	return nil
}

func (cv *codeValidator) unreachable() error {
	f, err := cv.getControlFrame(0)
	if err != nil {
		return fmt.Errorf("get control frame: %w", err)
	}

	cv.OperandStack = cv.OperandStack[:f.Height]
	cv.ControlStack[len(cv.ControlStack)-1].Unreachable = true

	return nil
}

func (cv *codeValidator) validateBlockOrLoop(instr types.Instruction) error {
	block := instr.Args.(types.Block)
	bt, err := cv.moduleValidator.module.GetBlockType(block.BlockType)
	if err != nil {
		return fmt.Errorf("get block type: %w", err)
	}
	if err := cv.popOperands(bt.ParamTypes); err != nil {
		return fmt.Errorf("pop operands for block/loop: %w", err)
	}
	cv.pushControlFrame(instr.Opcode, bt.ParamTypes, bt.ResultTypes)
	if err := cv.validateExprs(block.Instructions); err != nil {
		return fmt.Errorf("invalid exprs for block/loop: %w", err)
	}
	cf, err := cv.popControlFrame()
	if err != nil {
		return fmt.Errorf("pop control frame: %w", err)
	}
	cv.pushOperands(cf.EndTypes)

	return nil
}

func (cv *codeValidator) validateBreak(instr types.Instruction) error {
	n := int(instr.Args.(uint32))
	if len(cv.ControlStack) < n {
		return fmt.Errorf("unknown label: %v", n)
	}

	cf, err := cv.getControlFrame(n)
	if err != nil {
		return fmt.Errorf("get control frame: %w", err)
	}

	if err := cv.popOperands(cf.LabelTypes()); err != nil {
		return fmt.Errorf("pop operands: %w", err)
	}

	if err := cv.unreachable(); err != nil {
		return fmt.Errorf("mark frame unreachable: %w", err)
	}

	return nil
}

func (cv *codeValidator) validateBreakIf(instr types.Instruction) error {
	n := int(instr.Args.(uint32))
	if len(cv.ControlStack) < n {
		return fmt.Errorf("label(%d) is larger than max=%d", n, len(cv.ControlStack))
	}

	if _, err := cv.popTypeSpecificOperand(types.ValueTypeI32); err != nil {
		return fmt.Errorf("pop i32: %w", err)
	}

	f, err := cv.getControlFrame(n)
	if err != nil {
		return fmt.Errorf("get control frame(%d): %w", n, err)
	}

	if err := cv.popOperands(f.LabelTypes()); err != nil {
		return fmt.Errorf("pop operands: %w", err)
	}
	cv.pushOperands(f.LabelTypes())

	return nil
}

func (cv *codeValidator) validateBreakTable(instr types.Instruction) error {
	args := instr.Args.(types.BreakTable)
	m := int(args.Default)
	if len(cv.ControlStack) <= m {
		return fmt.Errorf("default label(%d) larger than max(%d)", m, len(cv.ControlStack))
	}

	defaultFrame, err := cv.getControlFrame(m)
	if err != nil {
		return fmt.Errorf("get default control frame labeled by %d: %w", m, err)
	}

	for _, n := range args.Labels {
		if len(cv.ControlStack) < int(n) {
			return fmt.Errorf("label(%d) larger than max(%d)", n, len(cv.ControlStack))
		}

		f1, err := cv.getControlFrame(int(n))
		if err != nil {
			return fmt.Errorf("get control frame labeled by %d: %w", n, err)
		}

		if !bytes.Equal(f1.LabelTypes(), defaultFrame.LabelTypes()) {
			return fmt.Errorf("inconsistent label types: %v != %v", f1.LabelTypes(), defaultFrame.LabelTypes())
		}
	}

	if _, err := cv.popTypeSpecificOperand(types.ValueTypeI32); err != nil {
		return fmt.Errorf("pop target label: %w", err)
	}

	if err := cv.popOperands(defaultFrame.LabelTypes()); err != nil {
		return fmt.Errorf("pop operands: %w", err)
	}

	return cv.unreachable()
}

func (cv *codeValidator) validateCall(instr types.Instruction) error {
	fIdx := instr.Args.(uint32)
	ft, ok := cv.moduleValidator.getFuncType(int(fIdx))
	if !ok {
		return fmt.Errorf("bad func type idx = %d", fIdx)
	}

	if err := cv.popOperands(ft.ParamTypes); err != nil {
		return fmt.Errorf("pop operands: %w", err)
	}

	cv.pushOperands(ft.ResultTypes)

	return nil
}

func (cv *codeValidator) validateCallIndirect(instr types.Instruction) error {
	if cv.moduleValidator.getTableLen() == 0 {
		return errors.New("no table")
	}

	ftIdx := instr.Args.(uint32)
	if int(ftIdx) >= len(cv.moduleValidator.module.Types) {
		return fmt.Errorf("type idx %d larger than max(%d)", ftIdx, len(cv.moduleValidator.module.Types))
	}
	ft := cv.moduleValidator.module.Types[ftIdx]

	if err := cv.popI32(); err != nil {
		return fmt.Errorf("pop elem: %w", err)
	}

	if err := cv.popOperands(ft.ParamTypes); err != nil {
		return fmt.Errorf("pop operands: %w", err)
	}

	cv.pushOperands(ft.ResultTypes)

	return nil
}

func (cv *codeValidator) validateCode(code types.Code, funcType types.FuncType) error {
	cv.popOperands(funcType.ParamTypes)
	cv.localLen = len(funcType.ParamTypes)

	for _, v := range code.Locals {
		for i := 0; i < int(v.N); i++ {
			cv.pushOperand(v.Type)
			cv.localLen++
		}
	}

	cv.pushControlFrame(types.OpcodeBlock, nil, funcType.ResultTypes)
	if err := cv.validateExprs(code.Expr); err != nil {
		return fmt.Errorf("validate exprs: %w", err)
	}

	cf, err := cv.popControlFrame()
	if err != nil {
		return fmt.Errorf("pop control frame: %w", err)
	}

	cv.pushOperands(cf.EndTypes)

	return nil
}

func (cv *codeValidator) validateExprs(exprs []types.Instruction) error {
	depth := len(cv.InstructionPath)

	for i, v := range exprs {
		cv.InstructionPath[depth] = v.GetOpname()
		if err := cv.validateInstr(v); err != nil {
			return fmt.Errorf("valida %d-th instruction: %w", i, err)
		}
	}

	return nil
}

func (cv *codeValidator) validateF32Cmp() error {
	if err := cv.popF32(); err != nil {
		return fmt.Errorf("pop 2nd operand: %w", err)
	}
	if err := cv.popF32(); err != nil {
		return fmt.Errorf("pop 1st operand: %w", err)
	}
	cv.pushOperand(types.ValueTypeI32)
	return nil
}

func (cv *codeValidator) validateF64Cmp() error {
	if err := cv.popF64(); err != nil {
		return fmt.Errorf("pop 2nd operand: %w", err)
	}
	if err := cv.popF64(); err != nil {
		return fmt.Errorf("pop 1st operand: %w", err)
	}
	cv.pushOperand(types.ValueTypeI32)
	return nil
}

func (cv *codeValidator) validateGlobalGet(instr types.Instruction) error {
	n := int(instr.Args.(uint32))
	if n >= len(cv.moduleValidator.globalTypes) {
		return fmt.Errorf("local idx(%d) >= max(%d)", n, len(cv.moduleValidator.globalTypes))
	}
	cv.pushOperand(cv.moduleValidator.globalTypes[n].ValueType)
	return nil
}

func (cv *codeValidator) validateGlobalSet(instr types.Instruction) error {
	n := int(instr.Args.(uint32))
	if n >= len(cv.moduleValidator.globalTypes) {
		return fmt.Errorf("local idx(%d) >= max(%d)", n, len(cv.moduleValidator.globalTypes))
	}
	gt := cv.moduleValidator.globalTypes[n]
	if gt.Mutable != 1 {
		return errors.New("immutable")
	}
	if _, err := cv.popTypeSpecificOperand(gt.ValueType); err != nil {
		return fmt.Errorf("pop operand: %w", err)
	}
	return nil
}

func (cv *codeValidator) validateIf(instr types.Instruction) error {
	blockIf := instr.Args.(types.BlockIf)
	bt, err := cv.moduleValidator.module.GetBlockType(blockIf.BlockType)
	if err != nil {
		return fmt.Errorf("get block type: %w", err)
	}

	if _, err := cv.popTypeSpecificOperand(types.ValueTypeI32); err != nil {
		return fmt.Errorf("pop i32: %w", err)
	}

	if err := cv.popOperands(bt.ParamTypes); err != nil {
		return fmt.Errorf("pop param types: %w", err)
	}

	cv.pushControlFrame(types.OpcodeIf, bt.ParamTypes, bt.ResultTypes)
	if err := cv.validateExprs(blockIf.Instructions1); err != nil {
		return fmt.Errorf("validate exprs: %w", err)
	}

	cf, err := cv.popControlFrame()
	if err != nil {
		return fmt.Errorf("pop control frame for if: %w", err)
	} else if cf.Opcode != types.OpcodeIf {
		return fmt.Errorf("opcode isn't if: %v", cf.Opcode)
	}

	cv.pushControlFrame(types.OpcodeElse, cf.SatrtTypes, cf.EndTypes)
	if err := cv.validateExprs(blockIf.Instructions2); err != nil {
		return fmt.Errorf("bad else block: %w", err)
	}

	blockElse, err := cv.popControlFrame()
	if err != nil {
		return fmt.Errorf("pop control frame for else: %w", err)
	}

	cv.pushOperands(blockElse.EndTypes)
	return nil
}

func (cv *codeValidator) validateInstr(instr types.Instruction) error {
	switch instr.Opcode {
	case types.OpcodeUnreachable:
		cv.unreachable()
	case types.OpcodeNop:
	case types.OpcodeBlock, types.OpcodeLoop:
		if err := cv.validateBlockOrLoop(instr); err != nil {
			return fmt.Errorf("bad block or loop: %w", err)
		}
	case types.OpcodeIf:
		if err := cv.validateIf(instr); err != nil {
			return fmt.Errorf("bad if block: %w", err)
		}
	case types.OpcodeBr:
		if err := cv.validateBreak(instr); err != nil {
			return fmt.Errorf("bad br: %w", err)
		}
	case types.OpcodeBrIf:
		if err := cv.validateBreakIf(instr); err != nil {
			return fmt.Errorf("bad br if: %w", err)
		}
	case types.OpcodeBrTable:
		if err := cv.validateBreakTable(instr); err != nil {
			return fmt.Errorf("validate br table: %w", err)
		}
	case types.OpcodeReturn:
		if err := cv.validateReturn(instr); err != nil {
			return fmt.Errorf("validate return: %w", err)
		}
	case types.OpcodeCall:
		if err := cv.validateCall(instr); err != nil {
			return fmt.Errorf("validate call: %w", err)
		}
	case types.OpcodeCallIndirect:
		if err := cv.validateCallIndirect(instr); err != nil {
			return fmt.Errorf("bad call indirect: %w", err)
		}
	case types.OpcodeDrop:
		if _, err := cv.popOperand(); err != nil {
			return fmt.Errorf("no operand to drop: %w", err)
		}
	case types.OpcodeSelect:
		if err := cv.validateSelect(instr); err != nil {
			return fmt.Errorf("bad select: %w", err)
		}
	case types.OpcodeLocalGet:
		if err := cv.validateLocalGet(instr); err != nil {
			return fmt.Errorf("bad local.get: %w", err)
		}
	case types.OpcodeLocalSet:
		if err := cv.validateLocalSet(instr); err != nil {
			return fmt.Errorf("bad local.set: %w", err)
		}
	case types.OpcodeLocalTee:
		if err := cv.validateLocalTee(instr); err != nil {
			return fmt.Errorf("bad local.tee: %w", err)
		}
	case types.OpcodeGlobalGet:
		if err := cv.validateGlobalGet(instr); err != nil {
			return fmt.Errorf("bad global.get: %w", err)
		}
	case types.OpcodeGlobalSet:
		if err := cv.validateGlobalSet(instr); err != nil {
			return fmt.Errorf("bad global.set: %w", err)
		}
	case types.OpcodeI32Load:
		if err := cv.i32Load(instr.Args, 32); err != nil {
			return fmt.Errorf("bad i32.load: %w", err)
		}
	case types.OpcodeF32Load:
		if err := cv.f32Load(instr.Args, 32); err != nil {
			return fmt.Errorf("bad f32.load: %w", err)
		}
	case types.OpcodeI64Load:
		if err := cv.i64Load(instr.Args, 64); err != nil {
			return fmt.Errorf("bad i64.load: %w", err)
		}
	case types.OpcodeF64Load:
		if err := cv.f64Load(instr.Args, 64); err != nil {
			return fmt.Errorf("bad f64.load: %w", err)
		}
	case types.OpcodeI32Load8S, types.OpcodeI32Load8U:
		if err := cv.i32Load(instr.Args, 8); err != nil {
			return fmt.Errorf("bad i32.load_8s/i32.load_8u: %w", err)
		}
	case types.OpcodeI32Load16S, types.OpcodeI32Load16U:
		if err := cv.i32Load(instr.Args, 16); err != nil {
			return fmt.Errorf("bad i32.load_16s/i32.load_16u: %w", err)
		}
	case types.OpcodeI64Load8S, types.OpcodeI64Load8U:
		if err := cv.i64Load(instr.Args, 8); err != nil {
			return fmt.Errorf("bad i64.load_8s/i64.load_8u: %w", err)
		}
	case types.OpcodeI64Load16S, types.OpcodeI64Load16U:
		if err := cv.i64Load(instr.Args, 16); err != nil {
			return fmt.Errorf("bad i64.load_16s/i64.load_16u: %w", err)
		}
	case types.OpcodeI64Load32S, types.OpcodeI64Load32U:
		if err := cv.i64Load(instr.Args, 32); err != nil {
			return fmt.Errorf("bad i64.load_32s/i64.load_32u: %w", err)
		}
	case types.OpcodeI32Store:
		if err := cv.i32Store(instr.Args, 32); err != nil {
			return fmt.Errorf("bad i32.store: %w", err)
		}
	case types.OpcodeI64Store:
		if err := cv.i64Store(instr.Args, 64); err != nil {
			return fmt.Errorf("bad i64.store: %w", err)
		}
	case types.OpcodeF32Store:
		if err := cv.f32Store(instr.Args, 32); err != nil {
			return fmt.Errorf("bad f32.store: %w", err)
		}
	case types.OpcodeF64Store:
		if err := cv.f64Store(instr.Args, 64); err != nil {
			return fmt.Errorf("bad f64.store: %w", err)
		}
	case types.OpcodeI32Store8:
		if err := cv.i32Store(instr.Args, 8); err != nil {
			return fmt.Errorf("bad i32.store8: %w", err)
		}
	case types.OpcodeI32Store16:
		if err := cv.i32Store(instr.Args, 16); err != nil {
			return fmt.Errorf("bad i32.store16: %w", err)
		}
	case types.OpcodeI64Store8:
		if err := cv.i64Store(instr.Args, 8); err != nil {
			return fmt.Errorf("bad i64.store8: %w", err)
		}
	case types.OpcodeI64Store16:
		if err := cv.i64Store(instr.Args, 16); err != nil {
			return fmt.Errorf("bad i64.store16: %w", err)
		}
	case types.OpcodeI64Store32:
		if err := cv.i64Store(instr.Args, 32); err != nil {
			return fmt.Errorf("bad i64.store32: %w", err)
		}
	case types.OpcodeMemorySize:
		if yes := cv.hasMemory(); !yes {
			return errors.New("no usable memory")
		}
		cv.pushOperand(types.ValueTypeI32)
	case types.OpcodeMemoryGrow:
		if err := cv.validateMemoryGrow(instr); err != nil {
			return fmt.Errorf("bad memory.grow: %w", err)
		}
	case types.OpcodeI32Const:
		cv.pushOperand(types.ValueTypeI32)
	case types.OpcodeI64Const:
		cv.pushOperand(types.ValueTypeI64)
	case types.OpcodeF32Const:
		cv.pushOperand(types.ValueTypeF32)
	case types.OpcodeF64Const:
		cv.pushOperand(types.ValueTypeF64)
	case types.OpcodeI32Eqz:
		if err := cv.validateI32Eqz(); err != nil {
			return fmt.Errorf("bad i32.eqz: %w", err)
		}
	case types.OpcodeI32Eq, types.OpcodeI32Ne, types.OpcodeI32LtS, types.OpcodeI32LtU,
		types.OpcodeI32GtS, types.OpcodeI32GtU, types.OpcodeI32LeS, types.OpcodeI32LeU,
		types.OpcodeI32GeS, types.OpcodeI32GeU:
		if err := cv.validateI32Cmp(); err != nil {
			return fmt.Errorf("bad i32.eq/.../i32.ge_u: %w", err)
		}
	case types.OpcodeI64Eqz:
		if err := cv.validateI64Eqz(); err != nil {
			return fmt.Errorf("bad i64.eqz: %w", err)
		}
	case types.OpcodeI64Eq, types.OpcodeI64Ne, types.OpcodeI64LtS, types.OpcodeI64LtU,
		types.OpcodeI64GtS, types.OpcodeI64GtU, types.OpcodeI64LeS, types.OpcodeI64LeU,
		types.OpcodeI64GeS, types.OpcodeI64GeU:
		if err := cv.validateI64Cmp(); err != nil {
			return fmt.Errorf("bad i64.eq/.../i64.ge_u: %w", err)
		}
	case types.OpcodeF32Eq, types.OpcodeF32Ne, types.OpcodeF32Lt, types.OpcodeF32Gt,
		types.OpcodeF32Le, types.OpcodeF32Ge:
		if err := cv.validateF32Cmp(); err != nil {
			return fmt.Errorf("bad f32.eq/.../f32.ge: %w", err)
		}
	case types.OpcodeF64Eq, types.OpcodeF64Ne, types.OpcodeF64Lt, types.OpcodeF64Gt,
		types.OpcodeF64Le, types.OpcodeF64Ge:
		if err := cv.validateF64Cmp(); err != nil {
			return fmt.Errorf("bad f64.eq/.../f64.ge: %w", err)
		}
	case types.OpcodeI32Clz, types.OpcodeI32Ctz, types.OpcodeI32PopCnt:
		if err := cv.popThenPush(types.ValueTypeI32, types.ValueTypeI32); err != nil {
			return fmt.Errorf("bad i32.{clz,ctz,pop_cnt}: %w", err)
		}
	case types.OpcodeI32Add, types.OpcodeI32Sub, types.OpcodeI32Mul, types.OpcodeI32DivS,
		types.OpcodeI32DivU, types.OpcodeI32RemS, types.OpcodeI32RemU, types.OpcodeI32And,
		types.OpcodeI32Or, types.OpcodeI32Xor, types.OpcodeI32Shl, types.OpcodeI32ShrS,
		types.OpcodeI32ShrU, types.OpcodeI32Rotl, types.OpcodeI32Rotr:
		if err := cv.popTowThenPushOneType(types.ValueTypeI32); err != nil {
			return fmt.Errorf("bad binary op for i32: %w", err)
		}
	case types.OpcodeI64Clz, types.OpcodeI64Ctz, types.OpcodeI64PopCnt:
		if err := cv.popThenPush(types.ValueTypeI64, types.ValueTypeI64); err != nil {
			return fmt.Errorf("bad i64.{clz,ctz,pop_cnt}: %w", err)
		}
	case types.OpcodeI64Add, types.OpcodeI64Sub, types.OpcodeI64Mul, types.OpcodeI64DivS,
		types.OpcodeI64DivU, types.OpcodeI64RemS, types.OpcodeI64RemU, types.OpcodeI64And,
		types.OpcodeI64Or, types.OpcodeI64Xor, types.OpcodeI64Shl, types.OpcodeI64ShrS,
		types.OpcodeI64ShrU, types.OpcodeI64Rotl, types.OpcodeI64Rotr:
		if err := cv.popTowThenPushOneType(types.ValueTypeI64); err != nil {
			return fmt.Errorf("bad binary op for i64: %w", err)
		}
	case types.OpcodeF32Abs, types.OpcodeF32Neg, types.OpcodeF32Ceil, types.OpcodeF32Floor,
		types.OpcodeF32Trunc, types.OpcodeF32Nearest, types.OpcodeF32Sqrt:
		if err := cv.popThenPush(types.ValueTypeF32, types.ValueTypeF32); err != nil {
			return fmt.Errorf("bad unary op for f32: %w", err)
		}
	case types.OpcodeF32Add, types.OpcodeF32Sub, types.OpcodeF32Mul, types.OpcodeF32Div,
		types.OpcodeF32Min, types.OpcodeF32Max, types.OpcodeF32CopySign:
		if err := cv.popTowThenPushOneType(types.ValueTypeF32); err != nil {
			return fmt.Errorf("bad binary op for f32: %w", err)
		}
	case types.OpcodeF64Abs, types.OpcodeF64Neg, types.OpcodeF64Ceil, types.OpcodeF64Floor,
		types.OpcodeF64Trunc, types.OpcodeF64Nearest, types.OpcodeF64Sqrt:
		if err := cv.popThenPush(types.ValueTypeF64, types.ValueTypeF64); err != nil {
			return fmt.Errorf("bad unary op for f64: %w", err)
		}
	case types.OpcodeF64Add, types.OpcodeF64Sub, types.OpcodeF64Mul, types.OpcodeF64Div,
		types.OpcodeF64Min, types.OpcodeF64Max, types.OpcodeF64CopySign:
		if err := cv.popTowThenPushOneType(types.ValueTypeF64); err != nil {
			return fmt.Errorf("bad binary op for f64: %w", err)
		}
	case types.OpcodeI32WrapI64:
		if err := cv.popThenPush(types.ValueTypeI64, types.ValueTypeI32); err != nil {
			return fmt.Errorf("bad i32.wrap_i64: %w", err)
		}
	case types.OpcodeI32TruncF32S, types.OpcodeI32TruncF32U:
		if err := cv.popThenPush(types.ValueTypeF32, types.ValueTypeI32); err != nil {
			return fmt.Errorf("bad i32.trunc_f32{s,u}: %w", err)
		}
	case types.OpcodeI32TruncF64S, types.OpcodeI32TruncF64U:
		if err := cv.popThenPush(types.ValueTypeF64, types.ValueTypeI32); err != nil {
			return fmt.Errorf("bad i32.trunc_f64{s,u}: %w", err)
		}
	case types.OpcodeI64ExtendI32S, types.OpcodeI64ExtendI32U:
		if err := cv.popThenPush(types.ValueTypeI32, types.ValueTypeI64); err != nil {
			return fmt.Errorf("bad i64.extend_i32{s,u}: %w", err)
		}
	case types.OpcodeI64TruncF32S, types.OpcodeI64TruncF32U:
		if err := cv.popThenPush(types.ValueTypeF32, types.ValueTypeI64); err != nil {
			return fmt.Errorf("bad i64.trunc_f32{s,u}: %w", err)
		}
	case types.OpcodeI64TruncF64S, types.OpcodeI64TruncF64U:
		if err := cv.popThenPush(types.ValueTypeF64, types.ValueTypeI64); err != nil {
			return fmt.Errorf("bad i64.trunc_f64{s,u}: %w", err)
		}
	case types.OpcodeF32ConvertI32S, types.OpcodeF32ConvertI32U:
		if err := cv.popThenPush(types.ValueTypeI32, types.ValueTypeF32); err != nil {
			return fmt.Errorf("bad f32.convert_i32{s,u}: %w", err)
		}
	case types.OpcodeF32ConvertI64S, types.OpcodeF32ConvertI64U:
		if err := cv.popThenPush(types.ValueTypeI64, types.ValueTypeF32); err != nil {
			return fmt.Errorf("bad f32.convert_i64{s,u}: %w", err)
		}
	case types.OpcodeF32DemoteF64:
		if err := cv.popThenPush(types.ValueTypeF64, types.ValueTypeF32); err != nil {
			return fmt.Errorf("bad f32.demote_f64: %w", err)
		}
	case types.OpcodeF64ConvertI32S, types.OpcodeF64ConvertI32U:
		if err := cv.popThenPush(types.ValueTypeI32, types.ValueTypeF64); err != nil {
			return fmt.Errorf("bad f64.convert_i32{s,u}: %w", err)
		}
	case types.OpcodeF64ConvertI64S, types.OpcodeF64ConvertI64U:
		if err := cv.popThenPush(types.ValueTypeI64, types.ValueTypeF64); err != nil {
			return fmt.Errorf("bad f64.convert_i64{s,u}: %w", err)
		}
	case types.OpcodeF64PromoteF32:
		if err := cv.popThenPush(types.ValueTypeF32, types.ValueTypeF64); err != nil {
			return fmt.Errorf("bad f64.promote_f32: %w", err)
		}
	case types.OpcodeI32ReinterpretF32:
		if err := cv.popThenPush(types.ValueTypeF32, types.ValueTypeI32); err != nil {
			return fmt.Errorf("bad i32.reinterpret_f32: %w", err)
		}
	case types.OpcodeI64ReinterpretF64:
		if err := cv.popThenPush(types.ValueTypeF64, types.ValueTypeI64); err != nil {
			return fmt.Errorf("bad i64.reinterpret_f64: %w", err)
		}
	case types.OpcodeF32ReinterpretI32:
		if err := cv.popThenPush(types.ValueTypeI32, types.ValueTypeF32); err != nil {
			return fmt.Errorf("bad f32.reinterpret_i32: %w", err)
		}
	case types.OpcodeF64ReinterpretI64:
		if err := cv.popThenPush(types.ValueTypeI64, types.ValueTypeF64); err != nil {
			return fmt.Errorf("bad f64.reinterpret_i64: %w", err)
		}
	case types.OpcodeI32Extend8S, types.OpcodeI32Extend16S:
		if err := cv.popThenPush(types.ValueTypeI32, types.ValueTypeI32); err != nil {
			return fmt.Errorf("bad i32.extend_{8,16}s: %w", err)
		}
	case types.OpcodeI64Extend8S, types.OpcodeI64Extend16S, types.OpcodeI64Extend32S:
		if err := cv.popThenPush(types.ValueTypeI64, types.ValueTypeI64); err != nil {
			return fmt.Errorf("bad i64.extend_{8,16,32}s: %w", err)
		}
	case types.OpcodeTruncSat:
		var err error
		subOpcode := instr.Args.(byte)
		switch subOpcode {
		case 0, 1:
			err = cv.popThenPush(types.ValueTypeF32, types.ValueTypeI32)
		case 2, 3:
			err = cv.popThenPush(types.ValueTypeF64, types.ValueTypeI32)
		case 4, 5:
			err = cv.popThenPush(types.ValueTypeF32, types.ValueTypeI64)
		case 6, 7:
			err = cv.popThenPush(types.ValueTypeF64, types.ValueTypeI64)
		default:
			err = errors.New("unknown subopcode")
		}
		if err != nil {
			return fmt.Errorf("bad trunc_sat with sub-opcode=%d: %w", subOpcode, err)
		}
	default:
		return fmt.Errorf("unknown opcode: 0x%x", instr.Opcode)
	}

	return nil
}

func (cv *codeValidator) validateI32Cmp() error {
	if err := cv.popI32(); err != nil {
		return fmt.Errorf("pop 2nd operand: %w", err)
	}
	if err := cv.popI32(); err != nil {
		return fmt.Errorf("pop 1st operand: %w", err)
	}
	cv.pushOperand(types.ValueTypeI32)
	return nil
}

func (cv *codeValidator) validateI32Eqz() error {
	if err := cv.popI32(); err != nil {
		return fmt.Errorf("pop i32: %w", err)
	}
	cv.pushOperand(types.ValueTypeI32)
	return nil
}

func (cv *codeValidator) validateI64Cmp() error {
	if err := cv.popI64(); err != nil {
		return fmt.Errorf("pop 2nd operand: %w", err)
	}
	if err := cv.popI64(); err != nil {
		return fmt.Errorf("pop 1st operand: %w", err)
	}
	cv.pushOperand(types.ValueTypeI32)
	return nil
}

func (cv *codeValidator) validateI64Eqz() error {
	if err := cv.popI64(); err != nil {
		return fmt.Errorf("pop i64: %w", err)
	}
	cv.pushOperand(types.ValueTypeI32)
	return nil
}

func (cv *codeValidator) validateMemoryGrow(instr types.Instruction) error {
	if yes := cv.hasMemory(); !yes {
		return errors.New("no usable memory")
	}
	if err := cv.popI32(); err != nil {
		return fmt.Errorf("pop i32: %w", err)
	}
	cv.pushOperand(types.ValueTypeI32)
	return nil
}

func (cv *codeValidator) validateLocalGet(instr types.Instruction) error {
	n := int(instr.Args.(uint32))
	if n >= cv.localLen {
		return fmt.Errorf("bad label(%d)>=%d", n, cv.localLen)
	}
	cv.pushOperand(cv.OperandStack[n])
	return nil
}

func (cv *codeValidator) validateLocalSet(instr types.Instruction) error {
	n := int(instr.Args.(uint32))
	if n >= cv.localLen {
		return fmt.Errorf("bad label(%d)>=%d", n, cv.localLen)
	}
	if _, err := cv.popTypeSpecificOperand(cv.OperandStack[n]); err != nil {
		return fmt.Errorf("pop operand: %w", err)
	}
	return nil
}

func (cv *codeValidator) validateLocalTee(instr types.Instruction) error {
	n := int(instr.Args.(uint32))
	if n >= cv.localLen {
		return fmt.Errorf("bad label(%d)>=%d", n, cv.localLen)
	}
	if _, err := cv.popTypeSpecificOperand(cv.OperandStack[n]); err != nil {
		return fmt.Errorf("pop stack top: %w", err)
	}
	cv.pushOperand(cv.OperandStack[n])
	return nil
}

func (cv *codeValidator) validateReturn(instr types.Instruction) error {
	n := len(cv.ControlStack) - 1

	f, err := cv.getControlFrame(n)
	if err != nil {
		return fmt.Errorf("get control frame: %w", err)
	}

	if err := cv.popOperands(f.LabelTypes()); err != nil {
		return fmt.Errorf("pop operands: %w", err)
	}

	return cv.unreachable()
}

func (cv *codeValidator) validateSelect(instr types.Instruction) error {
	if err := cv.popI32(); err != nil {
		return fmt.Errorf("pop condition: %w", err)
	}

	t1, err := cv.popOperand()
	if err != nil {
		return fmt.Errorf("pop 2nd operand: %w", err)
	}
	t2, err := cv.popTypeSpecificOperand(t1)
	if err != nil {
		return fmt.Errorf("pop 1st operand: %w", err)
	}
	cv.pushOperand(t2)

	return nil
}
