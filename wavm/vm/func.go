package vm

import "github.com/sammyne/mastering-wasm/wavm/types"

type Func struct {
	type_  types.FuncType
	code   types.Code
	goFunc types.GoFunc
}

func newExternalFunc(t types.FuncType, f types.GoFunc) Func {
	return Func{type_: t, goFunc: f}
}

func newInternalFunc(t types.FuncType, code types.Code) Func {
	return Func{type_: t, code: code}
}
