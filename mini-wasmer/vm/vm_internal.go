package vm

import (
	"errors"
	"fmt"

	"github.com/sammyne/mastering-wasm/mini-wasmer/types"
)

func (vm *VM) clearBlock(f ControlFrame) error {
	results, ok := vm.PopUint64s(len(f.BlockType.ResultTypes))
	if !ok {
		return fmt.Errorf("pop results: %w", ErrOperandPop)
	}

	vm.PopUint64s(vm.OperandStack.Len() - f.BP)
	vm.PushUint64s(results...)

	if f.Opcode != types.OpcodeCall || vm.ControlStack.Len() == 0 {
		return nil
	}

	lastCallFrame, _, ok := vm.TopCallFrame()
	if !ok {
		return ErrMissingCallFrame
	}

	vm.local0Idx = uint32(lastCallFrame.BP)
	return nil
}

func (vm *VM) enterBlock(opcode byte, type_ types.FuncType, expr []types.Instruction) {
	bp := vm.OperandStack.Len() - len(type_.ParamTypes)
	frame := NewControlFrame(opcode, type_, expr, bp)
	vm.ControlStack.Push(frame)

	if opcode == types.OpcodeCall {
		vm.local0Idx = uint32(bp)
	}
}

func (vm *VM) exitBlock() error {
	frame, ok := vm.ControlStack.Pop()
	if !ok {
		return nil
	}

	return vm.clearBlock(frame)
}

func (vm *VM) initGlobals() error {
	for i, v := range vm.module.Globals {
		for j, w := range v.Init {
			if err := vm.ExecuteInstruction(w); err != nil {
				return fmt.Errorf("exec %d-th instruction for %d-th global: %w", j, i, err)
			}

			vv, ok := vm.OperandStack.PopUint64()
			if !ok {
				return fmt.Errorf("expect global at stack top: %w", ErrOperandPop)
			}

			vm.globals = append(vm.globals, NewGlobalVar(v.Type, vv))
		}
	}

	return nil
}

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

func (vm *VM) loop() error {
	depth := vm.ControlStack.Len()
	for vm.ControlStack.Len() >= depth {
		f, ok := vm.ControlStack.Top()
		if !ok {
			return errors.New("miss control frame")
		}

		if f.PC == len(f.Expr) {
			if err := vm.exitBlock(); err != nil {
				return fmt.Errorf("exit block: %w", err)
			}
			continue
		}

		instruction := f.Expr[f.PC]
		f.PC++
		if err := vm.ExecuteInstruction(instruction); err != nil {
			return fmt.Errorf("exec instruction of PC(%d): %w", f.PC-1, err)
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

func (vm *VM) resetBlock(f *ControlFrame) error {
	results, ok := vm.OperandStack.PopUint64s(len(f.BlockType.ParamTypes))
	if !ok {
		return fmt.Errorf("pop result: %w", ErrOperandPop)
	}

	if _, ok := vm.OperandStack.PopUint64s(vm.OperandStack.Len() - f.BP); !ok {
		return fmt.Errorf("clean up call control stack: %w", ErrOperandPop)
	}

	vm.OperandStack.PushUint64s(results...)

	return nil
}
