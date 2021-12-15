package validator

import (
	"errors"
	"fmt"

	"github.com/sammyne/mastering-wasm/wavm/types"
)

func (v *moduleValidator) getFuncType(idx int) (types.FuncType, bool) {
	switch {
	case idx < len(v.importedFuncs):
		fnIdx := v.importedFuncs[idx].Description.Func
		return v.module.Types[fnIdx], true
	case idx < v.getFuncLen():
		fnIdx := v.module.Functions[idx-len(v.importedFuncs)]
		return v.module.Types[fnIdx], true
	default:
	}

	return types.FuncType{}, false
}

func (v *moduleValidator) getFuncLen() int {
	return len(v.importedFuncs) + len(v.module.Functions)
}

func (v *moduleValidator) getGlobalLen() int {
	return len(v.importedGlobals) + len(v.module.Globals)
}

func (v *moduleValidator) getMemoryLen() int {
	ell := len(v.module.Memories)
	if v.importedMemory != nil {
		ell++
	}

	return ell
}

func (v *moduleValidator) getTableLen() int {
	ell := len(v.module.Tables)
	if v.importedTable != nil {
		ell++
	}

	return ell
}

func (v *moduleValidator) validateCode(idx int, code types.Code, funcType types.FuncType) error {
	cv := &codeValidator{
		moduleValidator: v,
		Idx:             idx,
		InstructionPath: make(map[int]string),
	}

	return cv.validateCode(code, funcType)
}

func (v *moduleValidator) validateCodes() error {
	if len(v.module.Codes) != len(v.module.Functions) {
		return fmt.Errorf("#(code)=%d != #(func)=%d", len(v.module.Codes), len(v.module.Functions))
	}

	for i, c := range v.module.Codes {
		fnIdx := v.module.Functions[i]
		fnType := v.module.Types[fnIdx]
		if err := v.validateCode(i, c, fnType); err != nil {
			return fmt.Errorf("validate %d-th code: %w", i, err)
		}
	}

	return nil
}

func (v *moduleValidator) validateConstExpr(exprs []types.Instruction,
	expectedType types.ValueType) error {
	if len(exprs) == 0 {
		return nil
	}

	if len(exprs) > 1 {
		for i, instr := range exprs {
			switch instr.Opcode {
			case types.OpcodeI32Const, types.OpcodeI64Const, types.OpcodeF32Const, types.OpcodeF64Const,
				types.OpcodeGlobalGet:
			default:
				return fmt.Errorf("%d-th instruction has non-constant opcode(%d)", i, instr.Opcode)
			}
		}
		return errors.New("type mismatch") // TODO
	}

	var actualType byte = 0
	switch exprs[0].Opcode {
	case types.OpcodeI32Const:
		actualType = types.ValueTypeI32
	case types.OpcodeI64Const:
		actualType = types.ValueTypeI64
	case types.OpcodeF32Const:
		actualType = types.ValueTypeF32
	case types.OpcodeF64Const:
		actualType = types.ValueTypeF64
	case types.OpcodeGlobalGet:
		gIdx := exprs[0].Args.(uint32)
		if int(gIdx) >= len(v.globalTypes) {
			return fmt.Errorf("unknown global: %d", gIdx)
		}
		actualType = v.globalTypes[gIdx].ValueType
	default:
		return errors.New("constant expression required")
	}
	if actualType != expectedType {
		return errors.New("type mismatch") // TODO
	}

	return nil
}

func (v *moduleValidator) validateData() error {
	for i, data := range v.module.Data {
		if int(data.MemoryIdx) >= v.getMemoryLen() {
			return fmt.Errorf("data[%d]: unknown memory: %d", i, data.MemoryIdx)
		}
		if err := v.validateConstExpr(data.Offset, types.ValueTypeI32); err != nil {
			return fmt.Errorf("data[%d] has invalid const expr: %w", i, err)
		}
	}

	return nil
}

func (v *moduleValidator) validateElements() error {
	for i, elem := range v.module.Elements {
		if int(elem.TableIdx) >= v.getTableLen() {
			return fmt.Errorf("elem[%d] has unknown table %d", i, elem.TableIdx)
		}
		if err := v.validateConstExpr(elem.Offset, types.ValueTypeI32); err != nil {
			return fmt.Errorf("elem[%d] has invalid const expr: %s", i, err)
		}
		for j, funcIdx := range elem.Init {
			if int(funcIdx) >= v.getFuncLen() {
				return fmt.Errorf("elem[%d][%d] has unknown init function: %d", i, j, funcIdx)
			}
		}
	}

	return nil
}

func (v *moduleValidator) validateExports() error {
	exported := map[string]bool{}
	for i, w := range v.module.Exports {
		if exported[w.Name] {
			return fmt.Errorf("duplicate export name: %s", w.Name)
		}
		exported[w.Name] = true

		switch w.Description.Tag {
		case types.PortTagFunc:
			if int(w.Description.Idx) >= v.getFuncLen() {
				return fmt.Errorf("export[%d] refs unknown function %d", i, w.Description.Idx)
			}
		case types.PortTagTable:
			if int(w.Description.Idx) >= v.getTableLen() {
				return fmt.Errorf("export[%d] refs unknown table %d", i, w.Description.Idx)
			}
		case types.PortTagMemory:
			if int(w.Description.Idx) >= v.getMemoryLen() {
				return fmt.Errorf("export[%d] refs unknown memory %d", i, w.Description.Idx)
			}
		case types.PortTagGlobal:
			if int(w.Description.Idx) >= v.getGlobalLen() {
				return fmt.Errorf("export[%d] refs unknown global %d", i, w.Description.Idx)
			}
		}
	}

	return nil
}

func (v *moduleValidator) validateFunctions() error {
	typesLen := uint32(len(v.module.Types))
	for i, fnIdx := range v.module.Functions {
		if fnIdx >= typesLen {
			return fmt.Errorf("%d-th func idx(%d) out of bound(%d)", i, fnIdx, typesLen)
		}
	}

	return nil
}

func (v *moduleValidator) validateGlobals() error {
	importedGlobalsLen := len(v.importedGlobals)
	for i, g := range v.module.Globals {
		if err := v.validateConstExpr(g.Init, g.Type.ValueType); err != nil {
			return fmt.Errorf("global[%d]: %w", i+importedGlobalsLen, err)
		}
		v.globalTypes = append(v.globalTypes, g.Type)
	}

	return nil
}

func (v *moduleValidator) validateImports() error {
	for i, vv := range v.module.Imports {
		switch vv.Description.Tag {
		case types.PortTagFunc:
			v.importedFuncs = append(v.importedFuncs, vv)
			if int(vv.Description.Func) >= len(v.module.Types) {
				return fmt.Errorf("import[%d]: unknown type: %d", i, vv.Description.Func)
			}
		case types.PortTagTable:
			if v.importedTable != nil {
				return errors.New("multiple tables")
			}
			if err := validateTableLimits(vv.Description.Table.Limits); err != nil {
				return fmt.Errorf("bad table limit for import[%d]: %w", i, err)
			}
			v.importedTable = &v.module.Imports[i]
		case types.PortTagMemory:
			if v.importedMemory != nil {
				return fmt.Errorf("multiple memories")
			}
			if err := validateMemoryLimits(vv.Description.Memory); err != nil {
				return fmt.Errorf("bad memory limits for import[%d]: %s", i, err)
			}
			v.importedMemory = &v.module.Imports[i]
		case types.PortTagGlobal:
			v.importedGlobals = append(v.importedGlobals, vv)
			v.globalTypes = append(v.globalTypes, vv.Description.Global)
		}
	}

	return nil
}

func (v *moduleValidator) validateMemory() error {
	memLen := len(v.module.Memories)
	if v.importedMemory != nil {
		memLen++
	}
	if memLen == 0 {
		return nil
	} else if memLen > 1 {
		return errors.New("multiple memory sections")
	}

	return validateMemoryLimits(v.module.Memories[0])
}

func (v *moduleValidator) validateStart() error {
	if v.module.Start == nil {
		return nil
	}

	idx := *v.module.Start
	t, ok := v.getFuncType(int(idx))
	if !ok {
		return errors.New("get func type")
	} else if len(t.ParamTypes) != 0 || len(t.ResultTypes) != 0 {
		return errors.New("must have no params and return no results")
	}

	return nil
}

func (v *moduleValidator) validateTable() error {
	tableLen := len(v.module.Tables)
	if v.importedTable != nil {
		tableLen++
	}
	if tableLen == 0 {
		return nil
	} else if tableLen > 1 {
		return errors.New("multiple table sections")
	}

	if err := validateTableLimits(v.module.Tables[0].Limits); err != nil {
		return fmt.Errorf("bad table limits: %w", err)
	}

	return nil
}

func validateMemoryLimits(limits types.Limits) error {
	const maxLen = 1 << 16

	switch {
	case limits.Min > maxLen:
		return fmt.Errorf("min(%d) oversizes", limits.Min)
	case limits.Tag == 0:
	case limits.Min > limits.Max:
		return fmt.Errorf("wrong limit range(%d, %d)", limits.Min, limits.Max)
	default:
	}

	return nil
}

func validateTableLimits(limits types.Limits) error {
	if limits.Tag == 0 {
		return nil
	}

	if limits.Min > limits.Max {
		return fmt.Errorf("wrong limit range(%d, %d)", limits.Min, limits.Max)
	}

	return nil
}
