package vm

import "fmt"

func GlobalGet(vm *VM, args interface{}) error {
	idx, ok := args.(uint32)
	if !ok {
		return fmt.Errorf("index must be uint32: %w", ErrBadArgs)
	}

	if ell := uint32(len(vm.globals)); idx >= ell {
		return fmt.Errorf("index out of bound(%d): %w", ell, ErrBadArgs)
	}

	vm.OperandStack.PushUint64(vm.globals[idx].ToUint64())
	return nil
}

func GlobalSet(vm *VM, args interface{}) error {
	idx, ok := args.(uint32)
	if !ok {
		return fmt.Errorf("index must be uint32: %w", ErrBadArgs)
	}

	if ell := uint32(len(vm.globals)); idx >= ell {
		return fmt.Errorf("index out of bound(%d): %w", ell, ErrBadArgs)
	}

	val, ok := vm.OperandStack.PopUint64()
	if !ok {
		return ErrOperandPop
	}

	return vm.globals[idx].FromUint64(val)
}

func LocalGet(vm *VM, args interface{}) error {
	idx, ok := args.(uint32)
	if !ok {
		return ErrBadArgs
	}

	val, ok := vm.OperandStack.Get(vm.local0Idx + idx)
	if !ok {
		return ErrOperandPop
	}

	vm.OperandStack.PushUint64(val)
	return nil
}

func LocalSet(vm *VM, args interface{}) error {
	idx, ok := args.(uint32)
	if !ok {
		return fmt.Errorf("wrong type of index: %w", ErrBadArgs)
	}

	val, ok := vm.PopUint64()
	if !ok {
		return ErrOperandPop
	}

	if ok := vm.OperandStack.Set(vm.local0Idx+idx, val); !ok {
		return fmt.Errorf("bad idx to set: %w", ErrBadArgs)
	}
	return nil
}

func LocalTee(vm *VM, args interface{}) error {
	idx, ok := args.(uint32)
	if !ok {
		return fmt.Errorf("wrong type of index: %w", ErrBadArgs)
	}

	val, ok := vm.PopUint64()
	if !ok {
		return ErrOperandPop
	}

	vm.OperandStack.PushUint64(val)

	if ok := vm.OperandStack.Set(vm.local0Idx+idx, val); !ok {
		return fmt.Errorf("bad idx to set: %w", ErrBadArgs)
	}
	return nil
}
