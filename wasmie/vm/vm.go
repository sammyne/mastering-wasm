package vm

import (
	"fmt"

	wasmer "github.com/sammyne/mastering-wasm/wasmie"
	"github.com/sammyne/mastering-wasm/wasmie/types"
)

type VM struct {
	ControlStack
	OperandStack

	globals   []GlobalVar
	local0Idx uint32 // operand stack index for first operand
	memory    *Memory
	module    *wasmer.Module
}

func (vm *VM) ExecuteCode(idx int) error {
	code := vm.module.Codes[idx]

	for i, v := range code.Expr {
		//opname, _ := types.GetOpname(v.Opcode)
		//fmt.Println(i, opname)
		if err := vm.ExecuteInstruction(v); err != nil {
			opname, _ := types.GetOpname(v.Opcode)
			return fmt.Errorf("exec %d-th instruction(%s): %w", i, opname, err)
		}
	}

	return nil
}

func (vm *VM) ExecuteInstruction(i types.Instruction) error {
	run := instructionTable[i.Opcode]
	if run == nil {
		return ErrUnimplemented
	}

	return run(vm, i.Args)
}

func Run(m *wasmer.Module) error {
	if m.Start == nil {
		return ErrNoStartFunc
	}

	vm := &VM{module: m}

	if err := vm.initMemory(); err != nil {
		return fmt.Errorf("init memory for VM: %w", err)
	}
	if err := vm.initGlobals(); err != nil {
		return fmt.Errorf("init globals for VM: %w", err)
	}

	if err := Call(vm, *m.Start); err != nil {
		return fmt.Errorf("set up control stack: %w", err)
	}

	return vm.loop()
}
