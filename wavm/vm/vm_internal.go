package vm

import (
	"errors"
	"fmt"

	"github.com/sammyne/mastering-wasm/wavm/linker"
	"github.com/sammyne/mastering-wasm/wavm/types"
)

func (vm *VM) clearBlock(f ControlFrame) error {
	results, ok := vm.PopUint64s(len(f.BlockType.ResultTypes))
	if !ok {
		return fmt.Errorf("pop results: %w", ErrOperandPop)
	}

	vm.PopUint64s(vm.OperandStack.Len() - f.BP)
	vm.PushUint64s(results...)

	if f.Opcode != types.OpcodeCall || vm.ControlStack.Len() == 0 {
		return nil
	}

	lastCallFrame, _, ok := vm.TopCallFrame()
	if !ok {
		return ErrMissingCallFrame
	}

	vm.local0Idx = uint32(lastCallFrame.BP)
	return nil
}

func (vm *VM) enterBlock(opcode byte, type_ types.FuncType, expr []types.Instruction) {
	bp := vm.OperandStack.Len() - len(type_.ParamTypes)
	frame := NewControlFrame(opcode, type_, expr, bp)
	vm.ControlStack.Push(frame)

	if opcode == types.OpcodeCall {
		vm.local0Idx = uint32(bp)
	}
}

func (vm *VM) exitBlock() error {
	frame, ok := vm.ControlStack.Pop()
	if !ok {
		return nil
	}

	return vm.clearBlock(frame)
}

func (vm *VM) execStartFunc() error {
	idx := vm.module.Start
	if idx == nil {
		return nil
	}

	_, err := vm.funcs[*idx].call(nil)
	return err
}

func (vm *VM) initFuncs() error {
	for i, v := range vm.module.Functions {
		t := vm.module.Types[v]
		code := vm.module.Codes[i]
		vm.funcs = append(vm.funcs, newInternalFunc(t, code, vm))
	}

	return nil
}

func (vm *VM) initGlobals() error {
	for i, v := range vm.module.Globals {
		for j, w := range v.Init {
			if err := vm.ExecuteInstruction(w); err != nil {
				return fmt.Errorf("exec %d-th instruction for %d-th global: %w", j, i, err)
			}

			vv, ok := vm.OperandStack.PopUint64()
			if !ok {
				return fmt.Errorf("expect global at stack top: %w", ErrOperandPop)
			}

			vm.globals = append(vm.globals, NewGlobalVar(v.Type, vv))
		}
	}

	return nil
}

func (vm *VM) initMemory() error {
	if len(vm.module.Memories) <= 0 {
		return nil
	}

	vm.memory = NewMemory(vm.module.Memories[0])

	for i, v := range vm.module.Data {
		for j, vv := range v.Offset {
			if err := vm.ExecuteInstruction(vv); err != nil {
				return fmt.Errorf("eval %d-th offset for data[%d]: %w", j, i, err)
			}

			offset, ok := vm.PopUint64()
			if !ok {
				return fmt.Errorf("pop %d-th offset for data[%d]: %w", j, i, ErrOperandPop)
			}

			vm.memory.Write(offset, v.Init)
		}
	}

	return nil
}

func (vm *VM) initTable() error {
	if len(vm.module.Tables) > 0 {
		vm.table = newTable(vm.module.Tables[0])
	}

	for i, v := range vm.module.Elements {
		for k, vv := range v.Offset {
			if err := vm.ExecuteInstruction(vv); err != nil {
				return fmt.Errorf("eval %d-th offset for %d-th elem: %w", k, i, err)
			}
		}

		offset, ok := vm.PopUint32()
		if !ok {
			return fmt.Errorf("eval offset for %d-th elem: %w", i, ErrOperandPop)
		}

		for j, funcIdx := range v.Init {
			vm.table.SetElem(offset+uint32(j), vm.funcs[funcIdx])
		}
	}

	return nil
}

func (vm *VM) linkImport(m linker.Module, i types.Import) error {
	exported, err := m.GetMember(i.Name)
	if err != nil {
		return fmt.Errorf("get module member: %w", err)
	}

	switch x := exported.(type) {
	case linker.Function:
		vm.funcs = append(vm.funcs, newExternalFunc(x.Type(), x, vm))
	case linker.Global:
		vm.globals = append(vm.globals, x)
	case linker.Memory:
		vm.memory = x
	case linker.Table:
		vm.table = x
	default:
		return fmt.Errorf("unknown member type: %T", x)
	}

	return nil
}

func (vm *VM) linkImports(imports map[string]linker.Module) error {
	for _, v := range vm.module.Imports {
		m, ok := imports[v.Module]
		if !ok {
			return fmt.Errorf("module(%s) not found", v.Module)
		}
		if err := vm.linkImport(m, v); err != nil {
			return fmt.Errorf("link module(%s): %w", v.Module, err)
		}
	}

	return nil
}

func (vm *VM) loop() error {
	depth := vm.ControlStack.Len()
	for vm.ControlStack.Len() >= depth {
		f, ok := vm.ControlStack.Top()
		if !ok {
			return errors.New("miss control frame")
		}

		if f.PC == len(f.Expr) {
			if err := vm.exitBlock(); err != nil {
				return fmt.Errorf("exit block: %w", err)
			}
			continue
		}

		instruction := f.Expr[f.PC]
		f.PC++
		if err := vm.ExecuteInstruction(instruction); err != nil {
			return fmt.Errorf("exec instruction of PC(%d): %w", f.PC-1, err)
		}
	}

	return nil
}

func (vm *VM) popArgs(t types.FuncType) ([]types.WasmVal, error) {
	out := make([]types.WasmVal, len(t.ParamTypes))
	for i := len(t.ParamTypes) - 1; i >= 0; i-- {
		v, ok := vm.PopUint64()
		if !ok {
			return nil, fmt.Errorf("miss arg[%d]: %w", i, ErrOperandPop)
		}
		var err error
		if out[i], err = wrapUint64(t.ParamTypes[i], v); err != nil {
			return nil, fmt.Errorf("parse args[%d]: %w", i, err)
		}
	}

	return out, nil
}

func (vm *VM) popResults(resultTypes []types.ValueType) ([]types.WasmVal, error) {
	out := make([]types.WasmVal, len(resultTypes))
	for i := len(out) - 1; i >= 0; i-- {
		v, ok := vm.PopUint64()
		if !ok {
			return nil, fmt.Errorf("pop %d-th result: %w", i, ErrOperandPop)
		}
		vv, err := wrapUint64(resultTypes[i], v)
		if err != nil {
			return nil, fmt.Errorf("wrap result as uint64: %w", err)
		}
		out[i] = vv
	}

	return out, nil
}

func (vm *VM) popTowFloat32() (float32, float32, error) {
	v2, ok := vm.PopFloat32()
	if !ok {
		return 0, 0, fmt.Errorf("pop 1st operand: %w", ErrOperandPop)
	}

	v1, ok := vm.PopFloat32()
	if !ok {
		return 0, 0, fmt.Errorf("pop 2nd operand: %w", ErrOperandPop)
	}

	return v1, v2, nil
}

func (vm *VM) popTowFloat64() (float64, float64, error) {
	v2, ok := vm.PopFloat64()
	if !ok {
		return 0, 0, fmt.Errorf("pop 1st operand: %w", ErrOperandPop)
	}

	v1, ok := vm.PopFloat64()
	if !ok {
		return 0, 0, fmt.Errorf("pop 2nd operand: %w", ErrOperandPop)
	}

	return v1, v2, nil
}

func (vm *VM) popTowInt32() (int32, int32, error) {
	v2, ok := vm.PopInt32()
	if !ok {
		return 0, 0, fmt.Errorf("pop 1st operand: %w", ErrOperandPop)
	}

	v1, ok := vm.PopInt32()
	if !ok {
		return 0, 0, fmt.Errorf("pop 2nd operand: %w", ErrOperandPop)
	}

	return v1, v2, nil
}

func (vm *VM) popTowInt64() (int64, int64, error) {
	v2, ok := vm.PopInt64()
	if !ok {
		return 0, 0, fmt.Errorf("pop 1st operand: %w", ErrOperandPop)
	}

	v1, ok := vm.PopInt64()
	if !ok {
		return 0, 0, fmt.Errorf("pop 2nd operand: %w", ErrOperandPop)
	}

	return v1, v2, nil
}

func (vm *VM) popTowUint32() (uint32, uint32, error) {
	v2, ok := vm.PopUint32()
	if !ok {
		return 0, 0, fmt.Errorf("pop 1st operand: %w", ErrOperandPop)
	}

	v1, ok := vm.PopUint32()
	if !ok {
		return 0, 0, fmt.Errorf("pop 2nd operand: %w", ErrOperandPop)
	}

	return v1, v2, nil
}

func (vm *VM) popTowUint64() (uint64, uint64, error) {
	v2, ok := vm.PopUint64()
	if !ok {
		return 0, 0, fmt.Errorf("pop 1st operand: %w", ErrOperandPop)
	}

	v1, ok := vm.PopUint64()
	if !ok {
		return 0, 0, fmt.Errorf("pop 2nd operand: %w", ErrOperandPop)
	}

	return v1, v2, nil
}

func (vm *VM) pushArgs(paramTypes []types.ValueType, args []types.WasmVal) error {
	if len(paramTypes) != len(args) {
		return fmt.Errorf("len(paramTypes)=%d != len(args)=%d", len(paramTypes), len(args))
	}

	for i, v := range paramTypes {
		w, err := unwrapUint64(v, args[i])
		if err != nil {
			return fmt.Errorf("unwrap %d-th arg: %w", i, err)
		}
		vm.PushUint64(w)
	}

	return nil
}

func (vm *VM) pushResults(t types.FuncType, results []types.WasmVal) error {
	for i, v := range results {
		vv, err := unwrapUint64(t.ResultTypes[i], v)
		if err != nil {
			return fmt.Errorf("push results[%d]: %w", i, err)
		}
		vm.PushUint64(vv)
	}

	return nil
}

func (vm *VM) resetBlock(f *ControlFrame) error {
	results, ok := vm.OperandStack.PopUint64s(len(f.BlockType.ParamTypes))
	if !ok {
		return fmt.Errorf("pop result: %w", ErrOperandPop)
	}

	if _, ok := vm.OperandStack.PopUint64s(vm.OperandStack.Len() - f.BP); !ok {
		return fmt.Errorf("clean up call control stack: %w", ErrOperandPop)
	}

	vm.OperandStack.PushUint64s(results...)

	return nil
}
