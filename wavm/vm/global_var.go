package vm

import (
	"fmt"

	"github.com/sammyne/mastering-wasm/wavm/types"
)

type GlobalVar struct {
	Type_ types.GlobalType
	Value uint64
}

func (v *GlobalVar) FromUint64(val uint64) error {
	if v.Type_.Mutable == 0 {
		return ErrVarImmutable
	}

	v.Value = val
	return nil
}

func (v *GlobalVar) Get() (types.WasmVal, error) {
	return wrapUint64(v.Type_.ValueType, v.Value)
}

func (v *GlobalVar) GetAsUint64() uint64 {
	return v.Value
}

func (v *GlobalVar) SetAsUint64(val uint64) error {
	if v.Type_.Mutable != 1 {
		return ErrVarImmutable
	}
	v.Value = val
	return nil
}

func (v *GlobalVar) Set(value types.WasmVal) error {
	vv, err := unwrapUint64(v.Type_.ValueType, value)
	if err != nil {
		return fmt.Errorf("unwrap as uint64: %w", err)
	}

	v.Value = vv
	return nil
}

func (v *GlobalVar) Type() types.GlobalType {
	return v.Type_
}

func (v *GlobalVar) ToUint64() uint64 {
	return v.Value
}

func NewGlobalVar(t types.GlobalType, v uint64) *GlobalVar {
	return &GlobalVar{Type_: t, Value: v}
}
