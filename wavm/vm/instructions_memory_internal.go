package vm

import (
	"encoding/binary"
	"fmt"

	"github.com/sammyne/mastering-wasm/wavm/types"
)

var byteOrder = binary.LittleEndian

func getOffset(vm *VM, arg interface{}) (uint64, error) {
	offset1 := arg.(types.MemoryArg).Offset

	offset2, ok := vm.PopUint32()
	if !ok {
		return 0, fmt.Errorf("pop operand offset: %w", ErrOperandPop)
	}

	return uint64(offset1) + uint64(offset2), nil
}

func readUint16(vm *VM, arg interface{}) (uint16, error) {
	offset, err := getOffset(vm, arg)
	if err != nil {
		return 0, fmt.Errorf("get offset: %d", offset)
	}

	var buf [2]byte
	vm.memory.Read(offset, buf[:])

	return byteOrder.Uint16(buf[:]), nil
}

func readUint32(vm *VM, arg interface{}) (uint32, error) {
	offset, err := getOffset(vm, arg)
	if err != nil {
		return 0, fmt.Errorf("get offset: %d", offset)
	}

	var buf [4]byte
	vm.memory.Read(offset, buf[:])

	return byteOrder.Uint32(buf[:]), nil
}

func readUint64(vm *VM, arg interface{}) (uint64, error) {
	offset, err := getOffset(vm, arg)
	if err != nil {
		return 0, fmt.Errorf("get offset: %d", offset)
	}

	var buf [8]byte
	vm.memory.Read(offset, buf[:])

	return byteOrder.Uint64(buf[:]), nil
}

func readUint8(vm *VM, arg interface{}) (byte, error) {
	offset, err := getOffset(vm, arg)
	if err != nil {
		return 0, fmt.Errorf("get offset: %d", offset)
	}

	var buf [1]byte
	vm.memory.Read(offset, buf[:])
	return buf[0], nil
}

func writeUint16(vm *VM, arg interface{}, v uint16) error {
	offset, err := getOffset(vm, arg)
	if err != nil {
		return fmt.Errorf("get offset: %w", err)
	}

	var buf [2]byte
	byteOrder.PutUint16(buf[:], v)

	vm.memory.Write(offset, buf[:])
	return nil
}

func writeUint32(vm *VM, arg interface{}, v uint32) error {
	offset, err := getOffset(vm, arg)
	if err != nil {
		return fmt.Errorf("get offset: %w", err)
	}

	var buf [4]byte
	byteOrder.PutUint32(buf[:], v)

	vm.memory.Write(offset, buf[:])
	return nil
}

func writeUint64(vm *VM, arg interface{}, v uint64) error {
	offset, err := getOffset(vm, arg)
	if err != nil {
		return fmt.Errorf("get offset: %w", err)
	}

	var buf [8]byte
	byteOrder.PutUint64(buf[:], v)

	vm.memory.Write(offset, buf[:])
	return nil
}

func writeUint8(vm *VM, arg interface{}, v byte) error {
	offset, err := getOffset(vm, arg)
	if err != nil {
		return fmt.Errorf("get offset: %w", err)
	}

	buf := [...]byte{v}
	vm.memory.Write(offset, buf[:])
	return nil
}
