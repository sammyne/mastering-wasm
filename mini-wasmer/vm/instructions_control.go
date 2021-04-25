package vm

import (
	"fmt"

	"github.com/sammyne/mastering-wasm/mini-wasmer/tools"
	"github.com/sammyne/mastering-wasm/mini-wasmer/types"
)

func BreakIf(vm *VM, _ interface{}) error {
	yes, ok := vm.PopBool()
	if !ok {
		return ErrOperandPop
	}

	if !yes {
		return nil
	}

	if err := vm.exitBlock(); err != nil {
		return fmt.Errorf("exit block: %w", err)
	}

	return nil
}

func Call(vm *VM, arg interface{}) error {
	idx, ok := arg.(uint32)
	if !ok {
		return ErrBadArgs
	}

	var err error
	if ell := uint32(len(vm.module.Imports)); idx < ell {
		err = callAssertFunc(vm, idx)
	} else {
		err = callInternalFunc(vm, idx-ell)
	}
	if err != nil {
		return fmt.Errorf("run func(%d): %w", idx, err)
	}

	return nil
}

func assertEq(a, b interface{}) {
	if a != b {
		panic(fmt.Errorf("%v != %v", a, b))
	}
}

func callAssertFunc(vm *VM, idx uint32) error {
	var err error
	switch vm.module.Imports[idx].Name {
	case "assert_true":
		assertEq(vm.mustPopBool(), true)
	case "assert_false":
		assertEq(vm.mustPopBool(), false)
	case "assert_eq_i32":
		assertEq(vm.mustPopUint32(), vm.mustPopUint32())
	case "assert_eq_i64":
		assertEq(vm.mustPopUint64(), vm.mustPopUint64())
	case "assert_eq_f32":
		assertEq(vm.mustPopFloat32(), vm.mustPopFloat32())
	case "assert_eq_f64":
		assertEq(vm.mustPopFloat64(), vm.mustPopFloat64())
	default:
		err = ErrBadArgs
	}

	return err
}

func callInternalFunc(vm *VM, idx uint32) error {
	if ell := uint32(len(vm.module.Functions)); idx >= ell {
		return fmt.Errorf("func idx out of bound(%d): %w", ell, ErrBadArgs)
	}
	if ell := uint32(len(vm.module.Codes)); idx >= ell {
		return fmt.Errorf("code idx out of bound(%d): %w", ell, ErrBadArgs)
	}

	typeIdx := vm.module.Functions[idx]
	if ell := uint32(len(vm.module.Types)); typeIdx >= ell {
		return fmt.Errorf("func type idx out of bound(%d): %w", ell, ErrBadArgs)
	}

	code := vm.module.Codes[idx]

	vm.enterBlock(types.OpcodeCall, vm.module.Types[typeIdx], code.Expr)

	for i := tools.CountLocals(code.Locals); i > 0; i-- {
		vm.PushUint64(0)
	}

	return nil
}
