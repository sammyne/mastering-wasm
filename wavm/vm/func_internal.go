package vm

import (
	"fmt"

	"github.com/sammyne/mastering-wasm/wavm/linker"
	"github.com/sammyne/mastering-wasm/wavm/types"
)

func (f Func) call(args []types.WasmVal) ([]types.WasmVal, error) {
	if err := f.ctx.pushArgs(f.type_.ParamTypes, args); err != nil {
		return nil, fmt.Errorf("push args: %w", err)
	}

	callFunc(f.ctx, f)

	if f.externalFn == nil {
		if err := f.ctx.loop(); err != nil {
			return nil, fmt.Errorf("loop VM: %w", err)
		}
	}

	return f.ctx.popResults(f.type_.ResultTypes)
}

func newExternalFunc(t types.FuncType, f linker.Function, ctx *VM) Func {
	return Func{type_: t, externalFn: f, ctx: ctx}
}

func newInternalFunc(t types.FuncType, code types.Code, ctx *VM) Func {
	return Func{type_: t, code: code, ctx: ctx}
}
