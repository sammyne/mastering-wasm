package vm

import "fmt"

func (vm *VM) initMemory() error {
	if len(vm.module.Memories) <= 0 {
		return nil
	}

	vm.memory = NewMemory(vm.module.Memories[0])

	for i, v := range vm.module.Data {
		for j, vv := range v.Offset {
			if err := vm.ExecuteInstruction(vv); err != nil {
				return fmt.Errorf("eval %d-th offset for data[%d]: %w", j, i, err)
			}

			offset, ok := vm.PopUint64()
			if !ok {
				return fmt.Errorf("pop %d-th offset for data[%d]: %w", j, i, ErrOperandPop)
			}

			vm.memory.Write(offset, v.Init)
		}
	}

	return nil
}

func (vm *VM) popTowFloat32() (float32, float32, error) {
	v2, ok := vm.PopFloat32()
	if !ok {
		return 0, 0, fmt.Errorf("pop 1st operand: %w", ErrOperandPop)
	}

	v1, ok := vm.PopFloat32()
	if !ok {
		return 0, 0, fmt.Errorf("pop 2nd operand: %w", ErrOperandPop)
	}

	return v1, v2, nil
}

func (vm *VM) popTowFloat64() (float64, float64, error) {
	v2, ok := vm.PopFloat64()
	if !ok {
		return 0, 0, fmt.Errorf("pop 1st operand: %w", ErrOperandPop)
	}

	v1, ok := vm.PopFloat64()
	if !ok {
		return 0, 0, fmt.Errorf("pop 2nd operand: %w", ErrOperandPop)
	}

	return v1, v2, nil
}

func (vm *VM) popTowInt32() (int32, int32, error) {
	v2, ok := vm.PopInt32()
	if !ok {
		return 0, 0, fmt.Errorf("pop 1st operand: %w", ErrOperandPop)
	}

	v1, ok := vm.PopInt32()
	if !ok {
		return 0, 0, fmt.Errorf("pop 2nd operand: %w", ErrOperandPop)
	}

	return v1, v2, nil
}

func (vm *VM) popTowInt64() (int64, int64, error) {
	v2, ok := vm.PopInt64()
	if !ok {
		return 0, 0, fmt.Errorf("pop 1st operand: %w", ErrOperandPop)
	}

	v1, ok := vm.PopInt64()
	if !ok {
		return 0, 0, fmt.Errorf("pop 2nd operand: %w", ErrOperandPop)
	}

	return v1, v2, nil
}

func (vm *VM) popTowUint32() (uint32, uint32, error) {
	v2, ok := vm.PopUint32()
	if !ok {
		return 0, 0, fmt.Errorf("pop 1st operand: %w", ErrOperandPop)
	}

	v1, ok := vm.PopUint32()
	if !ok {
		return 0, 0, fmt.Errorf("pop 2nd operand: %w", ErrOperandPop)
	}

	return v1, v2, nil
}

func (vm *VM) popTowUint64() (uint64, uint64, error) {
	v2, ok := vm.PopUint64()
	if !ok {
		return 0, 0, fmt.Errorf("pop 1st operand: %w", ErrOperandPop)
	}

	v1, ok := vm.PopUint64()
	if !ok {
		return 0, 0, fmt.Errorf("pop 2nd operand: %w", ErrOperandPop)
	}

	return v1, v2, nil
}
