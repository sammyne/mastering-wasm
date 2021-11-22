package vm

import (
	"fmt"

	wasmer "github.com/sammyne/mastering-wasm/wavm"
	"github.com/sammyne/mastering-wasm/wavm/types"
)

type VM struct {
	ControlStack
	OperandStack

	globals   []GlobalVar
	local0Idx uint32 // operand stack index for first operand
	memory    *Memory
	module    *wasmer.Module
	funcs     []Func
	table     *Table
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
	var startFuncIdx uint32
	if m.Start != nil {
		startFuncIdx = *m.Start
	} else {
		var err error
		if startFuncIdx, err = getMainFuncIdx(m.Exports); err != nil {
			return ErrNoStartFunc
		}
	}

	vm := &VM{module: m}

	if err := vm.initMemory(); err != nil {
		return fmt.Errorf("init memory: %w", err)
	}
	if err := vm.initGlobals(); err != nil {
		return fmt.Errorf("init globals: %w", err)
	}
	if err := vm.initFuncs(); err != nil {
		return fmt.Errorf("init funcs: %w", err)
	}
	if err := vm.initTable(); err != nil {
		return fmt.Errorf("init table: %w", err)
	}

	if err := Call(vm, startFuncIdx); err != nil {
		return fmt.Errorf("set up control stack: %w", err)
	}

	return vm.loop()
}
