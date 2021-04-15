package wasmer

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

type ValueType = byte

const (
	ValueTypeI32 ValueType = 0x7F - iota
	ValueTypeI64
	ValueTypeF32
	ValueTypeF64
)
