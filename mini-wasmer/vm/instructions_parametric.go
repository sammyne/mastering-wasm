package vm

import "fmt"

func Drop(vm *VM, _ interface{}) error {
	if _, ok := vm.PopUint64(); !ok {
		return ErrOperandPop
	}

	return nil
}

func Select(vm *VM, _ interface{}) error {
	v1, ok := vm.PopBool()
	if !ok {
		return fmt.Errorf("pop 1st-operand: %w", ErrOperandPop)
	}

	v2, ok := vm.PopUint64()
	if !ok {
		return fmt.Errorf("pop 2nd-operand: %w", ErrOperandPop)
	}

	v3, ok := vm.PopUint64()
	if !ok {
		return fmt.Errorf("pop 3rd-operand: %w", ErrOperandPop)
	}

	if v1 {
		vm.PushUint64(v3)
	} else {
		vm.PushUint64(v2)
	}

	return nil
}
