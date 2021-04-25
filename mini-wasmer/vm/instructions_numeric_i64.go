package vm

import (
	"fmt"
	"math"
	"math/bits"
)

func I64Add(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowUint64()
	if err != nil {
		return err
	}

	vm.PushUint64(v1 + v2)
	return nil
}

func I64And(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowUint64()
	if err != nil {
		return err
	}

	vm.PushUint64(v1 & v2)
	return nil
}

func I64Clz(vm *VM, _ interface{}) error {
	v, ok := vm.PopUint64()
	if !ok {
		return ErrOperandPop
	}

	vm.PushUint64(uint64(bits.LeadingZeros64(v)))
	return nil
}

func I64Ctz(vm *VM, _ interface{}) error {
	v, ok := vm.PopUint64()
	if !ok {
		return ErrOperandPop
	}

	vm.PushUint64(uint64(bits.TrailingZeros64(v)))
	return nil
}

func I64Const(vm *VM, v interface{}) error {
	vm.PushInt64(v.(int64))
	return nil
}

func I64DivS(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowInt64()
	if err != nil {
		return err
	}

	vm.PushInt64(v1 / v2)
	return nil
}

func I64DivU(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowUint64()
	if err != nil {
		return err
	}

	vm.PushUint64(v1 / v2)
	return nil
}

func I64Eq(vm *VM, _ interface{}) error {
	v2, ok := vm.PopUint64()
	if !ok {
		return fmt.Errorf("pop 1st operand: %w", ErrOperandPop)
	}

	v1, ok := vm.PopUint64()
	if !ok {
		return fmt.Errorf("pop 2nd operand: %w", ErrOperandPop)
	}

	vm.PushBool(v1 == v2)
	return nil
}

func I64Eqz(vm *VM, _ interface{}) error {
	v, ok := vm.PopUint64()
	if !ok {
		return ErrOperandPop
	}

	vm.PushBool(v == 0)
	return nil
}

func I64ExtendI32S(vm *VM, _ interface{}) error {
	v, ok := vm.PopInt32()
	if !ok {
		return ErrOperandPop
	}

	vm.PushInt64(int64(v))
	return nil
}

func I64ExtendI32U(vm *VM, _ interface{}) error {
	v, ok := vm.PopUint32()
	if !ok {
		return ErrOperandPop
	}

	vm.PushUint64(uint64(v))
	return nil
}

func I64Extend16S(vm *VM, _ interface{}) error {
	v, ok := vm.PopInt64()
	if !ok {
		return ErrOperandPop
	}

	vm.PushInt64(int64(int16(v)))
	return nil
}

func I64Extend32S(vm *VM, _ interface{}) error {
	v, ok := vm.PopInt64()
	if !ok {
		return ErrOperandPop
	}

	vm.PushInt64(int64(int32(v)))
	return nil
}

func I64Extend8S(vm *VM, _ interface{}) error {
	v, ok := vm.PopInt64()
	if !ok {
		return ErrOperandPop
	}

	vm.PushInt64(int64(int8(v)))
	return nil
}

func I64GeS(vm *VM, _ interface{}) error {
	v2, ok := vm.PopInt64()
	if !ok {
		return fmt.Errorf("pop 1st operand: %w", ErrOperandPop)
	}

	v1, ok := vm.PopInt64()
	if !ok {
		return fmt.Errorf("pop 2nd operand: %w", ErrOperandPop)
	}

	vm.PushBool(v1 > v2)
	return nil
}

func I64GeU(vm *VM, _ interface{}) error {
	v2, ok := vm.PopUint64()
	if !ok {
		return fmt.Errorf("pop 1st operand: %w", ErrOperandPop)
	}

	v1, ok := vm.PopUint64()
	if !ok {
		return fmt.Errorf("pop 2nd operand: %w", ErrOperandPop)
	}

	vm.PushBool(v1 > v2)
	return nil
}

func I64GtS(vm *VM, _ interface{}) error {
	v2, ok := vm.PopInt64()
	if !ok {
		return fmt.Errorf("pop 1st operand: %w", ErrOperandPop)
	}

	v1, ok := vm.PopInt64()
	if !ok {
		return fmt.Errorf("pop 2nd operand: %w", ErrOperandPop)
	}

	vm.PushBool(v1 > v2)
	return nil
}

func I64GtU(vm *VM, _ interface{}) error {
	v2, ok := vm.PopUint64()
	if !ok {
		return fmt.Errorf("pop 1st operand: %w", ErrOperandPop)
	}

	v1, ok := vm.PopUint64()
	if !ok {
		return fmt.Errorf("pop 2nd operand: %w", ErrOperandPop)
	}

	vm.PushBool(v1 > v2)
	return nil
}

func I64LeS(vm *VM, _ interface{}) error {
	v2, ok := vm.PopInt64()
	if !ok {
		return fmt.Errorf("pop 1st operand: %w", ErrOperandPop)
	}

	v1, ok := vm.PopInt64()
	if !ok {
		return fmt.Errorf("pop 2nd operand: %w", ErrOperandPop)
	}

	vm.PushBool(v1 <= v2)
	return nil
}

func I64LeU(vm *VM, _ interface{}) error {
	v2, ok := vm.PopUint64()
	if !ok {
		return fmt.Errorf("pop 1st operand: %w", ErrOperandPop)
	}

	v1, ok := vm.PopUint64()
	if !ok {
		return fmt.Errorf("pop 2nd operand: %w", ErrOperandPop)
	}

	vm.PushBool(v1 <= v2)
	return nil
}

func I64Load(vm *VM, arg interface{}) error {
	v, err := readUint64(vm, arg)
	if err != nil {
		return fmt.Errorf("read uint64: %w", err)
	}

	vm.PushUint64(v)
	return nil
}

func I64Load16S(vm *VM, arg interface{}) error {
	v, err := readUint16(vm, arg)
	if err != nil {
		return fmt.Errorf("read uint16: %w", err)
	}

	vm.PushInt64(int64(int16(v)))
	return nil
}

func I64Load16U(vm *VM, arg interface{}) error {
	v, err := readUint16(vm, arg)
	if err != nil {
		return fmt.Errorf("read uint16: %w", err)
	}

	vm.PushUint64(uint64(v))
	return nil
}

func I64Load32S(vm *VM, arg interface{}) error {
	v, err := readUint32(vm, arg)
	if err != nil {
		return fmt.Errorf("read uint32: %w", err)
	}

	vm.PushInt64(int64(int32(v)))
	return nil
}

func I64Load32U(vm *VM, arg interface{}) error {
	v, err := readUint32(vm, arg)
	if err != nil {
		return fmt.Errorf("read uint32: %w", err)
	}

	vm.PushUint64(uint64(v))
	return nil
}

func I64Load8S(vm *VM, arg interface{}) error {
	v, err := readUint8(vm, arg)
	if err != nil {
		return fmt.Errorf("read uint8: %w", err)
	}

	vm.PushInt64(int64(int8(v)))
	return nil
}

func I64Load8U(vm *VM, arg interface{}) error {
	v, err := readUint8(vm, arg)
	if err != nil {
		return fmt.Errorf("read uint8: %w", err)
	}

	vm.PushUint64(uint64(v))
	return nil
}

func I64LtS(vm *VM, _ interface{}) error {
	v2, ok := vm.PopInt64()
	if !ok {
		return fmt.Errorf("pop 1st operand: %w", ErrOperandPop)
	}

	v1, ok := vm.PopInt64()
	if !ok {
		return fmt.Errorf("pop 2nd operand: %w", ErrOperandPop)
	}

	vm.PushBool(v1 < v2)
	return nil
}

func I64LtU(vm *VM, _ interface{}) error {
	v2, ok := vm.PopUint64()
	if !ok {
		return fmt.Errorf("pop 1st operand: %w", ErrOperandPop)
	}

	v1, ok := vm.PopUint64()
	if !ok {
		return fmt.Errorf("pop 2nd operand: %w", ErrOperandPop)
	}

	vm.PushBool(v1 < v2)
	return nil
}

func I64Mul(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowUint64()
	if err != nil {
		return err
	}

	vm.PushUint64(v1 * v2)
	return nil
}

func I64Ne(vm *VM, _ interface{}) error {
	v2, ok := vm.PopUint64()
	if !ok {
		return fmt.Errorf("pop 1st operand: %w", ErrOperandPop)
	}

	v1, ok := vm.PopUint64()
	if !ok {
		return fmt.Errorf("pop 2nd operand: %w", ErrOperandPop)
	}

	vm.PushBool(v1 != v2)
	return nil
}

func I64Or(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowUint64()
	if err != nil {
		return err
	}

	vm.PushUint64(v1 | v2)
	return nil
}

func I64PopCnt(vm *VM, _ interface{}) error {
	v, ok := vm.PopUint64()
	if !ok {
		return ErrOperandPop
	}

	vm.PushUint64(uint64(bits.OnesCount64(v)))
	return nil
}

func I64RemS(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowInt64()
	if err != nil {
		return err
	}

	vm.PushInt64(v1 % v2)
	return nil
}

func I64RemU(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowUint64()
	if err != nil {
		return err
	}

	vm.PushUint64(v1 % v2)
	return nil
}

func I64Rotl(vm *VM, _ interface{}) error {
	v2, ok := vm.PopUint64()
	if !ok {
		return fmt.Errorf("pop 1st operand: %w", ErrOperandPop)
	}

	v1, ok := vm.PopUint64()
	if !ok {
		return fmt.Errorf("pop 2nd operand: %w", ErrOperandPop)
	}

	vm.PushUint64(bits.RotateLeft64(v1, int(v2)))
	return nil
}

func I64Rotr(vm *VM, _ interface{}) error {
	v2, ok := vm.PopUint64()
	if !ok {
		return fmt.Errorf("pop 1st operand: %w", ErrOperandPop)
	}

	v1, ok := vm.PopUint64()
	if !ok {
		return fmt.Errorf("pop 2nd operand: %w", ErrOperandPop)
	}

	vm.PushUint64(bits.RotateLeft64(v1, -int(v2)))
	return nil
}

func I64Shl(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowUint64()
	if err != nil {
		return err
	}

	vm.PushUint64(v1 << (v2 % 64))
	return nil
}

func I64ShrS(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowInt64()
	if err != nil {
		return err
	}

	vm.PushInt64(v1 >> (v2 % 64))
	return nil
}

func I64ShrU(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowUint64()
	if err != nil {
		return err
	}

	vm.PushUint64(v1 >> (v2 % 64))
	return nil
}

func I64Store(vm *VM, arg interface{}) error {
	v, ok := vm.PopUint64()
	if !ok {
		return fmt.Errorf("pop uint64: %w", ErrOperandPop)
	}

	return writeUint64(vm, arg, v)
}

func I64Store16(vm *VM, arg interface{}) error {
	v, ok := vm.PopUint64()
	if !ok {
		return fmt.Errorf("pop uint64: %w", ErrOperandPop)
	}

	return writeUint16(vm, arg, uint16(v))
}

func I64Store32(vm *VM, arg interface{}) error {
	v, ok := vm.PopUint64()
	if !ok {
		return fmt.Errorf("pop uint64: %w", ErrOperandPop)
	}

	return writeUint32(vm, arg, uint32(v))
}

func I64Store8(vm *VM, arg interface{}) error {
	v, ok := vm.PopUint64()
	if !ok {
		return fmt.Errorf("pop uint64: %w", ErrOperandPop)
	}

	return writeUint8(vm, arg, byte(v))
}

func I64Sub(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowUint64()
	if err != nil {
		return err
	}

	vm.PushUint64(v1 - v2)
	return nil
}

func I64TruncF32S(vm *VM, _ interface{}) error {
	v, ok := vm.PopFloat32()
	if !ok {
		return ErrOperandPop
	}

	vm.PushInt64(int64(math.Trunc(float64(v))))
	return nil
}

func I64TruncF32U(vm *VM, _ interface{}) error {
	v, ok := vm.PopFloat32()
	if !ok {
		return ErrOperandPop
	}

	vm.PushUint64(uint64(math.Trunc(float64(v))))
	return nil
}

func I64TruncF64S(vm *VM, _ interface{}) error {
	v, ok := vm.PopFloat64()
	if !ok {
		return ErrOperandPop
	}

	vm.PushInt64(int64(math.Trunc(v)))
	return nil
}

func I64TruncF64U(vm *VM, _ interface{}) error {
	v, ok := vm.PopFloat64()
	if !ok {
		return ErrOperandPop
	}

	vm.PushUint64(uint64(math.Trunc(v)))
	return nil
}

func I64Xor(vm *VM, _ interface{}) error {
	v1, v2, err := vm.popTowUint64()
	if err != nil {
		return err
	}

	vm.PushUint64(v1 ^ v2)
	return nil
}
