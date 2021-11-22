package vm

import (
	"fmt"

	"github.com/sammyne/mastering-wasm/wavm/types"
)

func assertEq(a, b types.WasmVal) {
	if a != b {
		panic(fmt.Errorf("%v != %v", a, b))
	}
}

func assertTrue(args []types.WasmVal) ([]types.WasmVal, error) {
	assertEq(args[0].(int32), int32(1))
	return nil, nil
}

func assertFalse(args []types.WasmVal) ([]types.WasmVal, error) {
	assertEq(args[0].(int32), int32(0))
	return nil, nil
}

func assertEqI32(args []types.WasmVal) ([]types.WasmVal, error) {
	assertEq(args[0].(int32), args[1].(int32))
	return nil, nil
}

func assertEqI64(args []types.WasmVal) ([]types.WasmVal, error) {
	assertEq(args[0].(int64), args[1].(int64))
	return nil, nil
}

func assertEqF32(args []types.WasmVal) ([]types.WasmVal, error) {
	assertEq(args[0].(float32), args[1].(float32))
	return nil, nil
}

func assertEqF64(args []types.WasmVal) ([]types.WasmVal, error) {
	assertEq(args[0].(float64), args[1].(float64))
	return nil, nil
}

func printChar(args []types.WasmVal) ([]types.WasmVal, error) {
	fmt.Printf("%c", args[0].(int32))
	return nil, nil
}
