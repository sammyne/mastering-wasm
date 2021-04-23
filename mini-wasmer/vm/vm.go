package vm

import (
	"errors"
	"fmt"

	wasmer "github.com/sammyne/mastering-wasm/mini-wasmer"
	"github.com/sammyne/mastering-wasm/mini-wasmer/types"
)

type VM struct {
	OperandStack

	module *wasmer.Module
}

func (vm *VM) ExecuteCode(idx int) error {
	code := vm.module.Codes[idx]

	for i, v := range code.Expr {
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
		return errors.New("unimplemented instruction")
	}

	return run(vm, i.Args)
}

func Run(m *wasmer.Module) error {
	if m.Start == nil {
		return ErrNoStartFunc
	}

	vm := &VM{module: m}

	return vm.ExecuteCode(int(*m.Start) - len(m.Imports))
}
