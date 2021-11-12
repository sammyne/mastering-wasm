package vm

func MemoryGrow(vm *VM, _ interface{}) error {
	n, ok := vm.PopUint32()
	if !ok {
		return ErrOperandPop
	}

	vm.PushUint32(vm.memory.Grow(n))
	return nil
}

func MemorySize(vm *VM, _ interface{}) error {
	vm.PushUint32(vm.memory.Size())
	return nil
}
