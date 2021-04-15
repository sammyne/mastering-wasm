package wasmer

type Module struct {
	Magic     uint32
	Version   uint32
	Customs   []Custom
	Types     []FuncType
	Imports   []Import
	Functions []TypeIdx
	Tables    []TableType
	Memories  []MemoryType
	Globals   []Global
	Exports   []Export
	Start     *FuncIdx
	Elements  []Element
	Codes     []Code
	Data      []Data
}
