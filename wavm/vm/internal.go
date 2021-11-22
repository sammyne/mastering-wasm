package vm

import (
	"fmt"
	"math"

	"github.com/sammyne/mastering-wasm/wavm/types"
)

func getMainFuncIdx(exports []types.Export) (uint32, error) {
	for _, v := range exports {
		if v.Description.Tag == types.PortTagFunc && v.Name == "main" {
			return v.Description.Idx, nil
		}
	}

	return 0, fmt.Errorf("'main' is not found")
}

func unwrapUint64(t types.ValueType, v types.WasmVal) (uint64, error) {
	switch t {
	case types.ValueTypeI32:
		vv, ok := v.(int32)
		if !ok {
			return 0, fmt.Errorf("value %v isn't of type i32: %w", v, ErrBadValue)
		}
		return uint64(vv), nil
	case types.ValueTypeI64:
		vv, ok := v.(int64)
		if !ok {
			return 0, fmt.Errorf("value %v isn't of type i64: %w", v, ErrBadValue)
		}
		return uint64(vv), nil
	case types.ValueTypeF32:
		vv, ok := v.(float32)
		if !ok {
			return 0, fmt.Errorf("value %v isn't of type f32: %w", v, ErrBadValue)
		}
		return uint64(math.Float32bits(vv)), nil
	case types.ValueTypeF64:
		vv, ok := v.(float64)
		if !ok {
			return 0, fmt.Errorf("value %v isn't of type f64: %w", v, ErrBadValue)
		}
		return uint64(math.Float64bits(vv)), nil
	default:
	}

	return 0, ErrBadValueType
}

func wrapUint64(t types.ValueType, v uint64) (types.WasmVal, error) {
	switch t {
	case types.ValueTypeI32:
		return int32(v), nil
	case types.ValueTypeI64:
		return int64(v), nil
	case types.ValueTypeF32:
		return math.Float32frombits(uint32(v)), nil
	case types.ValueTypeF64:
		return math.Float64frombits(uint64(v)), nil
	default:
	}

	return nil, ErrBadValueType
}
