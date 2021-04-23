package vm

import (
	"math"
	"math/bits"
)

func I32Add(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowUint32()
	if err != nil {
		return err
	}

	vm.PushUint32(v1 + v2)
	return nil
}

func I32And(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowUint32()
	if err != nil {
		return err
	}

	vm.PushUint32(v1 & v2)
	return nil
}

func I32Clz(vm *VM, _ interface{}) error {
	v, ok := vm.PopUint32()
	if !ok {
		return ErrOperandPop
	}

	vm.PushUint32(uint32(bits.LeadingZeros32(v)))
	return nil
}

func I32Ctz(vm *VM, _ interface{}) error {
	v, ok := vm.PopUint32()
	if !ok {
		return ErrOperandPop
	}

	vm.PushUint32(uint32(bits.TrailingZeros32(v)))
	return nil
}

func I32Const(vm *VM, v interface{}) error {
	vm.PushInt32(v.(int32))
	return nil
}

func I32DivS(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowInt32()
	if err != nil {
		return err
	}

	vm.PushInt32(v1 / v2)
	return nil
}

func I32DivU(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowUint32()
	if err != nil {
		return err
	}

	vm.PushUint32(v1 / v2)
	return nil
}

func I32Eq(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowUint32()
	if err != nil {
		return err
	}

	vm.PushBool(v1 == v2)
	return nil
}

func I32Eqz(vm *VM, _ interface{}) error {
	v, ok := vm.PopUint32()
	if !ok {
		return ErrOperandPop
	}

	vm.PushBool(v == 0)
	return nil
}

func I32Extend16S(vm *VM, _ interface{}) error {
	v, ok := vm.PopInt32()
	if !ok {
		return ErrOperandPop
	}

	vm.PushInt32(int32(int16(v)))
	return nil
}

func I32Extend8S(vm *VM, _ interface{}) error {
	v, ok := vm.PopInt32()
	if !ok {
		return ErrOperandPop
	}

	vm.PushInt32(int32(int8(v)))
	return nil
}

func I32GeS(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowInt32()
	if err != nil {
		return err
	}

	vm.PushBool(v1 > v2)
	return nil
}

func I32GeU(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowUint32()
	if err != nil {
		return err
	}

	vm.PushBool(v1 > v2)
	return nil
}

func I32GtS(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowInt32()
	if err != nil {
		return err
	}

	vm.PushBool(v1 > v2)
	return nil
}

func I32GtU(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowUint32()
	if err != nil {
		return err
	}

	vm.PushBool(v1 > v2)
	return nil
}

func I32LeS(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowInt32()
	if err != nil {
		return err
	}

	vm.PushBool(v1 <= v2)
	return nil
}

func I32LeU(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowUint32()
	if err != nil {
		return err
	}

	vm.PushBool(v1 <= v2)
	return nil
}

func I32LtS(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowInt32()
	if err != nil {
		return err
	}

	vm.PushBool(v1 < v2)
	return nil
}

func I32LtU(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowUint32()
	if err != nil {
		return err
	}

	vm.PushBool(v1 < v2)
	return nil
}

func I32Mul(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowUint32()
	if err != nil {
		return err
	}

	vm.PushUint32(v1 * v2)
	return nil
}

func I32Ne(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowUint32()
	if err != nil {
		return err
	}

	vm.PushBool(v1 != v2)
	return nil
}

func I32Or(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowUint32()
	if err != nil {
		return err
	}

	vm.PushUint32(v1 | v2)
	return nil
}

func I32PopCnt(vm *VM, _ interface{}) error {
	v, ok := vm.PopUint32()
	if !ok {
		return ErrOperandPop
	}

	vm.PushUint32(uint32(bits.OnesCount32(v)))
	return nil
}

func I32RemS(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowInt32()
	if err != nil {
		return err
	}

	vm.PushInt32(v1 % v2)
	return nil
}

func I32RemU(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowUint32()
	if err != nil {
		return err
	}

	vm.PushUint32(v1 % v2)
	return nil
}

func I32Rotl(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowUint32()
	if err != nil {
		return err
	}

	vm.PushUint32(bits.RotateLeft32(v1, int(v2)))
	return nil
}

func I32Rotr(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowUint32()
	if err != nil {
		return err
	}

	vm.PushUint32(bits.RotateLeft32(v1, -int(v2)))
	return nil
}

func I32Shl(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowUint32()
	if err != nil {
		return err
	}

	vm.PushUint32(v1 << (v2 % 32))
	return nil
}

func I32ShrS(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowInt32()
	if err != nil {
		return err
	}

	vm.PushInt32(v1 >> (v2 % 32))
	return nil
}

func I32ShrU(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowUint32()
	if err != nil {
		return err
	}

	vm.PushUint32(v1 >> (v2 % 32))
	return nil
}

func I32Sub(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowUint32()
	if err != nil {
		return err
	}

	vm.PushUint32(v1 - v2)
	return nil
}

func I32TruncF32S(vm *VM, _ interface{}) error {
	v, ok := vm.PopFloat32()
	if !ok {
		return ErrOperandPop
	}

	vm.PushInt32(int32(math.Trunc(float64(v))))
	return nil
}

func I32TruncF32U(vm *VM, _ interface{}) error {
	v, ok := vm.PopFloat32()
	if !ok {
		return ErrOperandPop
	}

	vm.PushUint32(uint32(math.Trunc(float64(v))))
	return nil
}

func I32TruncF64S(vm *VM, _ interface{}) error {
	v, ok := vm.PopFloat64()
	if !ok {
		return ErrOperandPop
	}

	vm.PushInt32(int32(math.Trunc(v)))
	return nil
}

func I32TruncF64U(vm *VM, _ interface{}) error {
	v, ok := vm.PopFloat64()
	if !ok {
		return ErrOperandPop
	}

	vm.PushUint32(uint32(math.Trunc(v)))
	return nil
}

func I32WrapI64(vm *VM, _ interface{}) error {
	v, ok := vm.PopUint32()
	if !ok {
		return ErrOperandPop
	}

	vm.PushUint32(v)
	return nil
}

func I32Xor(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowUint32()
	if err != nil {
		return err
	}

	vm.PushUint32(v1 ^ v2)
	return nil
}
