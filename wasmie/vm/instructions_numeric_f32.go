package vm

import (
	"math"
)

func F32Abs(vm *VM, _ interface{}) error {
	v, ok := vm.PopFloat32()
	if !ok {
		return ErrOperandPop
	}

	vm.PushFloat32(float32(math.Abs(float64(v))))
	return nil
}

func F32Add(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowFloat32()
	if err != nil {
		return err
	}

	vm.PushFloat32(v1 + v2)
	return nil
}

func F32Ceil(vm *VM, _ interface{}) error {
	v, ok := vm.PopFloat32()
	if !ok {
		return ErrOperandPop
	}

	vm.PushFloat32(float32(math.Ceil(float64(v))))
	return nil
}

func F32Const(vm *VM, v interface{}) error {
	vm.PushFloat32(v.(float32))
	return nil
}

func F32ConvertI32S(vm *VM, _ interface{}) error {
	v, ok := vm.PopInt32()
	if !ok {
		return ErrOperandPop
	}

	vm.PushFloat32(float32(v))
	return nil
}

func F32ConvertI32U(vm *VM, _ interface{}) error {
	v, ok := vm.PopUint32()
	if !ok {
		return ErrOperandPop
	}

	vm.PushFloat32(float32(v))
	return nil
}

func F32ConvertI64S(vm *VM, _ interface{}) error {
	v, ok := vm.PopInt64()
	if !ok {
		return ErrOperandPop
	}

	vm.PushFloat32(float32(v))
	return nil
}

func F32ConvertI64U(vm *VM, _ interface{}) error {
	v, ok := vm.PopUint64()
	if !ok {
		return ErrOperandPop
	}

	vm.PushFloat32(float32(v))
	return nil
}

func F32CopySign(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowFloat32()
	if err != nil {
		return err
	}

	vm.PushFloat32(float32(math.Copysign(float64(v1), float64(v2))))
	return nil
}

func F32DemoteF64(vm *VM, _ interface{}) error {
	v, ok := vm.PopFloat64()
	if !ok {
		return ErrOperandPop
	}

	vm.PushFloat32(float32(v))
	return nil
}

func F32Div(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowFloat32()
	if err != nil {
		return err
	}

	vm.PushFloat32(v1 / v2)
	return nil
}

func F32Eq(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowFloat32()
	if err != nil {
		return err
	}

	vm.PushBool(v1 == v2)
	return nil
}

func F32Floor(vm *VM, _ interface{}) error {
	v, ok := vm.PopFloat32()
	if !ok {
		return ErrOperandPop
	}

	vm.PushFloat32(float32(math.Floor(float64(v))))
	return nil
}

func F32Ge(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowFloat32()
	if err != nil {
		return err
	}

	vm.PushBool(v1 >= v2)
	return nil
}

func F32Gt(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowFloat32()
	if err != nil {
		return err
	}

	vm.PushBool(v1 > v2)
	return nil
}

func F32Load(vm *VM, arg interface{}) error {
	return I32Load(vm, arg)
}

func F32Ne(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowFloat32()
	if err != nil {
		return err
	}

	vm.PushBool(v1 != v2)
	return nil
}

func F32Le(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowFloat32()
	if err != nil {
		return err
	}

	vm.PushBool(v1 <= v2)
	return nil
}

func F32Lt(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowFloat32()
	if err != nil {
		return err
	}

	vm.PushBool(v1 < v2)
	return nil
}

func F32Max(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowFloat32()
	if err != nil {
		return err
	}

	if v1 > v2 {
		vm.PushFloat32(v1)
	} else {
		vm.PushFloat32(v2)
	}

	return nil
}

func F32Min(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowFloat32()
	if err != nil {
		return err
	}

	if v1 < v2 {
		vm.PushFloat32(v1)
	} else {
		vm.PushFloat32(v2)
	}

	return nil
}

func F32Mul(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowFloat32()
	if err != nil {
		return err
	}

	vm.PushFloat32(v1 * v2)
	return nil
}

func F32Nearest(vm *VM, _ interface{}) error {
	v, ok := vm.PopFloat32()
	if !ok {
		return ErrOperandPop
	}

	vm.PushFloat32(float32(math.RoundToEven(float64(v))))
	return nil
}

func F32Neg(vm *VM, _ interface{}) error {
	v, ok := vm.PopFloat32()
	if !ok {
		return ErrOperandPop
	}

	vm.PushFloat32(-v)
	return nil
}

func F32Sqrt(vm *VM, _ interface{}) error {
	v, ok := vm.PopFloat32()
	if !ok {
		return ErrOperandPop
	}

	vm.PushFloat32(float32(math.Sqrt(float64(v))))
	return nil
}

func F32Store(vm *VM, arg interface{}) error {
	return I32Store(vm, arg)
}

func F32Sub(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowFloat32()
	if err != nil {
		return err
	}

	vm.PushFloat32(v1 - v2)
	return nil
}

func F32Trunc(vm *VM, _ interface{}) error {
	v, ok := vm.PopFloat32()
	if !ok {
		return ErrOperandPop
	}

	vm.PushFloat32(float32(math.Trunc(float64(v))))
	return nil
}
