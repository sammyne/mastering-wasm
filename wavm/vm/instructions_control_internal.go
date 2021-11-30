package vm

import (
	"fmt"

	"github.com/sammyne/mastering-wasm/wavm/linker"
	"github.com/sammyne/mastering-wasm/wavm/tools"
	"github.com/sammyne/mastering-wasm/wavm/types"
)

func callFunc(vm *VM, f Func) error {
	if f.externalFn != nil {
		return callExternalFunc(vm, f.externalFn)
	}

	return callInternalFunc(vm, f)
}

func callExternalFunc(vm *VM, f linker.Function) error {
	args, err := vm.popArgs(f.Type())
	if err != nil {
		return fmt.Errorf("pop args: %w", err)
	}

	results, err := f.Call(args...)
	if err != nil {
		return fmt.Errorf("call go func: %w", err)
	}

	if err := vm.pushResults(f.Type(), results); err != nil {
		return fmt.Errorf("push results: %w", err)
	}

	return nil
}

func callInternalFunc(vm *VM, f Func) error {
	vm.enterBlock(types.OpcodeCall, f.type_, f.code.Expr)

	for i := tools.CountLocals(f.code.Locals); i > 0; i-- {
		vm.PushUint64(0)
	}

	return nil
}
