package validator

import (
	"fmt"

	"github.com/sammyne/mastering-wasm/wavm"
	"github.com/sammyne/mastering-wasm/wavm/types"
)

type moduleValidator struct {
	module wavm.Module

	importedFuncs   []types.Import
	importedTable   *types.Import
	importedMemory  *types.Import
	importedGlobals []types.Import
	globalTypes     []types.GlobalType
}

func (v *moduleValidator) Validate() error {
	if err := v.validateImports(); err != nil {
		return fmt.Errorf("bad imports: %w", err)
	}
	if err := v.validateFunctions(); err != nil {
		return fmt.Errorf("bad functions: %w", err)
	}
	if err := v.validateTable(); err != nil {
		return fmt.Errorf("bad table: %w", err)
	}
	if err := v.validateMemory(); err != nil {
		return fmt.Errorf("bad memory: %w", err)
	}
	if err := v.validateGlobals(); err != nil {
		return fmt.Errorf("bad global: %w", err)
	}
	if err := v.validateExports(); err != nil {
		return fmt.Errorf("bad exports: %w", err)
	}
	if err := v.validateStart(); err != nil {
		return fmt.Errorf("bad start: %w", err)
	}
	if err := v.validateElements(); err != nil {
		return fmt.Errorf("bad elements: %w", err)
	}
	if err := v.validateCodes(); err != nil {
		return fmt.Errorf("bad codes: %w", err)
	}
	if err := v.validateData(); err != nil {
		return fmt.Errorf("bad data: %w", err)
	}

	return nil
}

func Validate(m wavm.Module) error {
	v := moduleValidator{module: m}
	return v.Validate()
}
