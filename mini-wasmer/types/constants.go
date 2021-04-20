package types

const FuncRef = 0x70

const (
	Magic   = 0x6D736100 // `\0asm`
	Version = 0x00000001
)

const (
	MutConst byte = 0
	MutVar   byte = 1
)

type PortTag = byte

const (
	PortTagFunc = iota
	PortTagTable
	PortTagMemory
	PortTagGlobal
)

const (
	SectionIDCustom = iota
	SectionIDType
	SectionIDImport
	SectionIDFunc
	SectionIDTable
	SectionIDMemory
	SectionIDGlobal
	SectionIDExport
	SectionIDStart
	SectionIDElement
	SectionIDCode
	SectionIDData
)

type ValueType = byte

const (
	ValueTypeUnknown ValueType = 0
	ValueTypeI32     ValueType = 0x80 - iota
	ValueTypeI64
	ValueTypeF32
	ValueTypeF64
)

func StringifyValueType(t ValueType) string {
	switch t {
	case ValueTypeI32:
		return "i32"
	case ValueTypeI64:
		return "i64"
	case ValueTypeF32:
		return "f32"
	case ValueTypeF64:
		return "f64"
	default:
	}

	return "unknown"
}
