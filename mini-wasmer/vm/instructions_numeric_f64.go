package vm

import (
	"math"
)

func F64Abs(vm *VM, _ interface{}) error {
	v, ok := vm.PopFloat64()
	if !ok {
		return ErrOperandPop
	}

	vm.PushFloat64(math.Abs(v))
	return nil
}

func F64Add(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowFloat64()
	if err != nil {
		return err
	}

	vm.PushFloat64(v1 + v2)
	return nil
}

func F64Ceil(vm *VM, _ interface{}) error {
	v, ok := vm.PopFloat64()
	if !ok {
		return ErrOperandPop
	}

	vm.PushFloat64(math.Ceil(v))
	return nil
}

func F64Const(vm *VM, v interface{}) error {
	vm.PushFloat64(v.(float64))
	return nil
}

func F64ConvertI32S(vm *VM, _ interface{}) error {
	v, ok := vm.PopInt32()
	if !ok {
		return ErrOperandPop
	}

	vm.PushFloat64(float64(v))
	return nil
}

func F64ConvertI32U(vm *VM, _ interface{}) error {
	v, ok := vm.PopUint32()
	if !ok {
		return ErrOperandPop
	}

	vm.PushFloat64(float64(v))
	return nil
}

func F64ConvertI64S(vm *VM, _ interface{}) error {
	v, ok := vm.PopInt64()
	if !ok {
		return ErrOperandPop
	}

	vm.PushFloat64(float64(v))
	return nil
}

func F64ConvertI64U(vm *VM, _ interface{}) error {
	v, ok := vm.PopUint64()
	if !ok {
		return ErrOperandPop
	}

	vm.PushFloat64(float64(v))
	return nil
}

func F64CopySign(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowFloat64()
	if err != nil {
		return err
	}

	vm.PushFloat64(math.Copysign(v1, v2))
	return nil
}

func F64Div(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowFloat64()
	if err != nil {
		return err
	}

	vm.PushFloat64(v1 / v2)
	return nil
}

func F64Eq(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowFloat64()
	if err != nil {
		return err
	}

	vm.PushBool(v1 == v2)
	return nil
}

func F64Floor(vm *VM, _ interface{}) error {
	v, ok := vm.PopFloat64()
	if !ok {
		return ErrOperandPop
	}

	vm.PushFloat64(math.Floor(v))
	return nil
}

func F64Ge(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowFloat64()
	if err != nil {
		return err
	}

	vm.PushBool(v1 >= v2)
	return nil
}

func F64Gt(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowFloat64()
	if err != nil {
		return err
	}

	vm.PushBool(v1 > v2)
	return nil
}

func F64Ne(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowFloat64()
	if err != nil {
		return err
	}

	vm.PushBool(v1 != v2)
	return nil
}

func F64Le(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowFloat64()
	if err != nil {
		return err
	}

	vm.PushBool(v1 <= v2)
	return nil
}

func F64Lt(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowFloat64()
	if err != nil {
		return err
	}

	vm.PushBool(v1 < v2)
	return nil
}

func F64Max(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowFloat64()
	if err != nil {
		return err
	}

	if v1 > v2 {
		vm.PushFloat64(v1)
	} else {
		vm.PushFloat64(v2)
	}

	return nil
}

func F64Min(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowFloat64()
	if err != nil {
		return err
	}

	if v1 < v2 {
		vm.PushFloat64(v1)
	} else {
		vm.PushFloat64(v2)
	}

	return nil
}

func F64Mul(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowFloat64()
	if err != nil {
		return err
	}

	vm.PushFloat64(v1 * v2)
	return nil
}

func F64Nearest(vm *VM, _ interface{}) error {
	v, ok := vm.PopFloat64()
	if !ok {
		return ErrOperandPop
	}

	vm.PushFloat64(math.RoundToEven(v))
	return nil
}

func F64Neg(vm *VM, _ interface{}) error {
	v, ok := vm.PopFloat64()
	if !ok {
		return ErrOperandPop
	}

	vm.PushFloat64(-v)
	return nil
}

func F64PromoteF32(vm *VM, _ interface{}) error {
	v, ok := vm.PopFloat32()
	if !ok {
		return ErrOperandPop
	}

	vm.PushFloat64(float64(v))
	return nil
}

func F64Sqrt(vm *VM, _ interface{}) error {
	v, ok := vm.PopFloat64()
	if !ok {
		return ErrOperandPop
	}

	vm.PushFloat64(math.Sqrt(v))
	return nil
}

func F64Sub(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowFloat64()
	if err != nil {
		return err
	}

	vm.PushFloat64(v1 - v2)
	return nil
}

func F64Trunc(vm *VM, _ interface{}) error {
	v, ok := vm.PopFloat64()
	if !ok {
		return ErrOperandPop
	}

	vm.PushFloat64(math.Trunc(v))
	return nil
}
