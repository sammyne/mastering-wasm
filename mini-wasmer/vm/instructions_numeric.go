package vm

import "math"

func Nop(vm *VM, _ interface{}) error {
	return nil
}

func TruncSat(vm *VM, subOpcode interface{}) error {
	if vm.OperandStack.Len() == 0 {
		return ErrOperandPop
	}

	// @TODO: add the must API for OperandStack to tidy up
	switch subOpcode.(byte) {
	case 0x00:
		v, _ := vm.PopFloat32()
		vm.PushInt32(int32(truncSatS(float64(v), 32)))
	case 0x01:
		v, _ := vm.PopFloat32()
		vm.PushUint32(uint32(truncSatU(float64(v), 32)))
	case 0x02:
		v, _ := vm.PopFloat64()
		vm.PushInt32(int32(truncSatS(v, 32)))
	case 0x03:
		v, _ := vm.PopFloat64()
		vm.PushUint32(uint32(truncSatU(v, 32)))
	case 0x04:
		v, _ := vm.PopFloat32()
		vm.PushInt64(int64(truncSatS(float64(v), 64)))
	case 0x05:
		v, _ := vm.PopFloat32()
		vm.PushUint64(uint64(truncSatU(float64(v), 64)))
	case 0x06:
		v, _ := vm.PopFloat64()
		vm.PushInt64(int64(truncSatS(v, 64)))
	case 0x07:
		v, _ := vm.PopFloat64()
		vm.PushUint64(uint64(truncSatU(v, 64)))
	default:
		return ErrBadSubOpcode
	}

	return nil
}

func truncSatS(z float64, n int) int64 {
	if math.IsNaN(z) {
		return 0
	}

	min := -(int64(1) << (n - 1))
	max := (int64(1) << (n - 1)) - 1

	switch {
	case math.IsInf(z, -1):
		return min
	case math.IsInf(z, 1):
		return max
	default:
	}

	x := math.Trunc(z)

	switch {
	case x < float64(min):
		return min
	case x >= float64(max):
		return max
	default:
	}

	return int64(x)
}

func truncSatU(z float64, n int) uint64 {
	if math.IsNaN(z) || math.IsInf(z, -1) {
		return 0
	}

	max := (uint64(1) << n) - 1

	if math.IsInf(z, 1) {
		return max
	}

	x := math.Trunc(z)
	switch {
	case x < 0:
		return 0
	case x >= float64(max):
		return max
	default:
	}

	return uint64(x)
}
