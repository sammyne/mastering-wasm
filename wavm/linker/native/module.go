package native

import (
	"errors"
	"fmt"

	"github.com/sammyne/mastering-wasm/wavm/linker"
	"github.com/sammyne/mastering-wasm/wavm/types"
)

type Module struct {
	exported map[string]types.WasmVal
}

func (m Module) GetGlobalVal(name string) (types.WasmVal, error) {
	v, err := m.GetMember(name)
	if err != nil {
		return nil, fmt.Errorf("get member: %w", err)
	}

	vv, ok := v.(linker.Global)
	if !ok {
		return nil, errors.New("not global")
	}

	return vv.Get()
}

func (m Module) GetMember(name string) (types.WasmVal, error) {
	out, ok := m.exported[name]
	if !ok {
		return nil, errors.New("not found")
	}

	return out, nil
}

func (m Module) InvokeFunc(name string, args ...types.WasmVal) ([]types.WasmVal, error) {
	v, ok := m.exported[name]
	if !ok {
		return nil, errors.New("member not found")
	}

	fn, ok := v.(linker.Function)
	if !ok {
		return nil, errors.New("not func")
	}

	return fn.Call(args...)
}

func (m Module) Register(name string, x types.WasmVal) {
	m.exported[name] = x
}

func (m Module) RegisterFunc(nameAndSig string, f GoFunc) {
	name, sig := parseNameAndSig(nameAndSig)
	m.exported[name] = Function{type_: sig, f: f}
}

func (m Module) SetGlobalVal(name string, vv types.WasmVal) error {
	raw, err := m.GetMember(name)
	if err != nil {
		return fmt.Errorf("get member: %w", err)
	}

	v, ok := raw.(linker.Global)
	if !ok {
		return errors.New("non-global")
	}

	return v.Set(vv)
}

func NewModule() Module {
	return Module{exported: make(map[string]types.WasmVal)}
}
