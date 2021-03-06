package vm

import (
	"fmt"

	"github.com/sammyne/mastering-wasm/wavm/tools"
	"github.com/sammyne/mastering-wasm/wavm/types"
)

func Block(vm *VM, arg interface{}) error {
	b, ok := arg.(*types.Block)
	if !ok {
		return fmt.Errorf("expect *types.Block: %w", ErrBadArgs)
	}

	blockType := tools.ParseBlockSig(b.BlockType, vm.module.Types)
	vm.enterBlock(types.OpcodeBlock, blockType, b.Instructions)

	return nil
}

func BlockIf(vm *VM, arg interface{}) error {
	a, ok := arg.(*types.BlockIf)
	if !ok {
		return fmt.Errorf("expect *types.BlockIf: %w", ErrBadArgs)
	}

	blockType := tools.ParseBlockSig(a.BlockType, vm.module.Types)

	yes, ok := vm.PopBool()
	if !ok {
		return ErrOperandPop
	}

	if yes {
		vm.enterBlock(types.OpcodeIf, blockType, a.Instructions1)
	} else {
		vm.enterBlock(types.OpcodeIf, blockType, a.Instructions2)
	}

	return nil
}

func Break(vm *VM, arg interface{}) error {
	labelIdx, ok := arg.(uint32)
	if !ok {
		return fmt.Errorf("expect index as uint32: %w", ErrBadArgs)
	}

	for i := 0; i < int(labelIdx); i++ {
		vm.ControlStack.Pop()
	}

	f, ok := vm.ControlStack.Top()
	if !ok {
		return fmt.Errorf("miss frame: %w", ErrOperandPop)
	}

	if f.Opcode != types.OpcodeLoop {
		if err := vm.exitBlock(); err != nil {
			return fmt.Errorf("run non-loop block: %w", err)
		}
		return nil
	}

	if err := vm.resetBlock(f); err != nil {
		return fmt.Errorf("reset block: %w", err)
	}
	f.PC = 0

	return nil
}

func BreakIf(vm *VM, arg interface{}) error {
	yes, ok := vm.PopBool()
	if !ok {
		return fmt.Errorf("pop condition: %w", ErrOperandPop)
	}

	if !yes {
		return nil
	}

	return Break(vm, arg)
}

func BreakTable(vm *VM, arg interface{}) error {
	a, ok := arg.(*types.BreakTable)
	if !ok {
		return fmt.Errorf("expect *types.BreakTable: %w", ErrBadArgs)
	}

	n, ok := vm.PopUint32()
	if !ok {
		return fmt.Errorf("pop label: %w", ErrOperandPop)
	}

	if n < uint32(len(a.Labels)) {
		return Break(vm, a.Labels[n])
	}

	return Break(vm, a.Default)
}

func Call(vm *VM, arg interface{}) error {
	idx, ok := arg.(uint32)
	if !ok || idx >= uint32(len(vm.funcs)) {
		return ErrBadArgs
	}

	return callFunc(vm, vm.funcs[idx])
}

func CallIndirect(vm *VM, arg interface{}) error {
	elemIdx, ok := vm.PopUint32()
	if !ok {
		return fmt.Errorf("no elem idx: %w", ErrOperandPop)
	}

	f, err := vm.table.GetElem(elemIdx)
	if err != nil {
		return fmt.Errorf("miss elem at %d: %w", elemIdx, err)
	}

	typeIdx, ok := arg.(uint32)
	if !ok {
		return fmt.Errorf("arg should be u32: %w", ErrBadArgs)
	}
	funcType := vm.module.Types[typeIdx] // bound check

	if expect, got := funcType.String(), f.Type().String(); expect != got {
		return fmt.Errorf("mismatch func type: expect %s, got %s, %w", expect, got, ErrBadArgs)
	}

	if v, ok := f.(Func); ok {
		return callFunc(vm, v)
	}

	return callExternalFunc(vm, f)
}

func Loop(vm *VM, arg interface{}) error {
	a, ok := arg.(*types.Block)
	if !ok {
		return fmt.Errorf("expect *types.Block: %w", ErrBadArgs)
	}

	blockType := tools.ParseBlockSig(a.BlockType, vm.module.Types)
	vm.enterBlock(types.OpcodeLoop, blockType, a.Instructions)

	return nil
}

func Nop(vm *VM, _ interface{}) error {
	return nil
}

func Return(vm *VM, _ interface{}) error {
	_, labelIdx, ok := vm.TopCallFrame()
	if !ok {
		return fmt.Errorf("missing top call frame: %w", ErrOperandPop)
	}

	return Break(vm, uint32(labelIdx))
}

func Unreachable(vm *VM, _ interface{}) error {
	panic("unreachble")
}
