package vm

import (
	"fmt"

	"github.com/sammyne/mastering-wasm/wavm/tools"
	"github.com/sammyne/mastering-wasm/wavm/types"
)

func callFunc(vm *VM, f Func) error {
	if f.goFunc != nil {
		return callExternalFunc(vm, f.type_, f.goFunc)
	}

	return callInternalFunc(vm, f)
}

func callExternalFunc(vm *VM, t types.FuncType, f types.GoFunc) error {
	args, err := vm.popArgs(t)
	if err != nil {
		return fmt.Errorf("pop args: %w", err)
	}

	results, err := f(args)
	if err != nil {
		return fmt.Errorf("call go func: %w", err)
	}

	if err := vm.pushResults(t, results); err != nil {
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
