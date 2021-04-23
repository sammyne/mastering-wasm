package vm

import "fmt"

func Call(vm *VM, arg interface{}) error {
	idx, ok := arg.(uint32)
	if !ok {
		return ErrBadArgs
	}

	if idx >= uint32(len(vm.module.Imports)) {
		return fmt.Errorf("import idx out of bound: %w", ErrBadArgs)
	}

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

func assertEq(a, b interface{}) {
	if a != b {
		panic(fmt.Errorf("%v != %v", a, b))
	}
}
