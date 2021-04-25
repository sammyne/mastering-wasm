package vm

import "github.com/sammyne/mastering-wasm/mini-wasmer/types"

type GlobalVar struct {
	Type  types.GlobalType
	Value uint64
}

func (v *GlobalVar) FromUint64(val uint64) error {
	if v.Type.Mutable == 0 {
		return ErrVarImmutable
	}

	v.Value = val
	return nil
}

func (v *GlobalVar) ToUint64() uint64 {
	return v.Value
}

func NewGlobalVar(t types.GlobalType, v uint64) GlobalVar {
	return GlobalVar{Type: t, Value: v}
}
