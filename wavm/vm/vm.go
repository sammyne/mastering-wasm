package vm

import (
	"fmt"

	"github.com/sammyne/mastering-wasm/wavm"
	"github.com/sammyne/mastering-wasm/wavm/linker"
	"github.com/sammyne/mastering-wasm/wavm/types"
	"github.com/sammyne/mastering-wasm/wavm/validator"
)

type VM struct {
	ControlStack
	OperandStack

	globals   []linker.Global
	local0Idx uint32 // operand stack index for first operand
	memory    linker.Memory
	module    *wavm.Module
	funcs     []Func
	table     linker.Table
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

func (vm *VM) GetGlobalVal(name string) (types.WasmVal, error) {
	v, err := vm.GetMember(name)
	if err != nil {
		return nil, fmt.Errorf("global not found: %w", err)
	}

	g, ok := v.(linker.Global)
	if !ok {
		return nil, fmt.Errorf("member isn't global: %s", name)
	}

	return g.Get()
}

func (vm *VM) GetMember(name string) (interface{}, error) {
	for _, export := range vm.module.Exports {
		if export.Name != name {
			continue
		}

		idx := export.Description.Idx
		switch export.Description.Tag {
		case types.PortTagFunc:
			return vm.funcs[idx], nil
		case types.PortTagTable:
			return vm.table, nil
		case types.PortTagMemory:
			return vm.memory, nil
		case types.PortTagGlobal:
			return vm.globals[idx], nil
		default:
			return nil, ErrUnimplemented
		}
	}

	return nil, fmt.Errorf("no found: %w", ErrBadArgs)
}

func (vm *VM) InvokeFunc(name string, args ...types.WasmVal) ([]types.WasmVal, error) {
	fn, err := vm.GetMember(name)
	if err != nil {
		return nil, fmt.Errorf("func not found: %w", err)
	}

	f, ok := fn.(linker.Function)
	if !ok {
		return nil, fmt.Errorf("module member isn't func: %s", name)
	}

	return f.Call(args...)
}

func (vm *VM) SetGlobalVal(name string, val types.WasmVal) error {
	v, err := vm.GetMember(name)
	if err != nil {
		return fmt.Errorf("global not found: %w", err)
	}

	g, ok := v.(linker.Global)
	if !ok {
		return fmt.Errorf("global not found: %s", name)
	}

	return g.Set(val)
}

func NewVM(m *wavm.Module, externals map[string]linker.Module) (linker.Module, error) {
	if err := validator.Validate(*m); err != nil {
		return nil, fmt.Errorf("invalid main module: %w", err)
	}

	vm := &VM{module: m}

	if err := vm.linkImports(externals); err != nil {
		return nil, fmt.Errorf("link imports: %w", err)
	}
	if err := vm.initMemory(); err != nil {
		return nil, fmt.Errorf("init memory: %w", err)
	}
	if err := vm.initGlobals(); err != nil {
		return nil, fmt.Errorf("init globals: %w", err)
	}
	if err := vm.initFuncs(); err != nil {
		return nil, fmt.Errorf("init funcs: %w", err)
	}
	if err := vm.initTable(); err != nil {
		return nil, fmt.Errorf("init table: %w", err)
	}

	if err := vm.execStartFunc(); err != nil {
		return nil, fmt.Errorf("exec start func: %w", err)
	}

	return vm, nil
}

func Run(m *wavm.Module) error {
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
