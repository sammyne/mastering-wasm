package types

type (
	FuncIdx   = uint32
	GlobalIdx = uint32
	LabelIdx  = uint32
	LocalIdx  = uint32
	MemoryIdx = uint32
	TableIdx  = uint32
	TypeIdx   = uint32
)

type Code struct {
	Locals []Locals
	Expr   Expr
}

type Custom struct {
	Name  string
	Bytes []byte
}

type Data struct {
	MemoryIdx MemoryIdx
	Offset    Expr
	Init      []byte
}

type Element struct {
	Table  TableIdx
	Offset Expr
	Init   []FuncIdx
}

type Export struct {
	Name        string
	Description ExportDescription
}

type ExportDescription struct {
	Tag PortTag
	Idx uint32
}

type Expr interface{}

type FuncType struct {
	Tag         byte
	ParamTypes  []ValueType
	ResultTypes []ValueType
}

type Global struct {
	Type GlobalType
	Init Expr
}

type GlobalType struct {
	ValueType ValueType
	Mutable   byte
}

type Import struct {
	Module      string
	Name        string
	Description ImportDescription
}

type ImportDescription struct {
	Tag    PortTag
	Func   TypeIdx
	Table  Table
	Memory Memory
	Global GlobalType
}

type Limits struct {
	Tag byte
	Min uint32
	Max uint32
}

type Locals struct {
	N    uint32
	Type ValueType
}

type Memory = Limits

type Table struct {
	ElementType byte
	Limits      Limits
}
