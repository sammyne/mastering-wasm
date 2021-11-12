package types

// Block may be block or loop
type Block struct {
	BlockType    BlockType
	Instructions []Instruction
}

type BlockIf struct {
	BlockType     BlockType
	Instructions1 []Instruction
	Instructions2 []Instruction
}

type BreakTable struct {
	Labels  []LabelIdx
	Default LabelIdx
}

type Instruction struct {
	Opcode byte
	Args   interface{}
}

type MemoryArg struct {
	Align  uint32
	Offset uint32
}

func (i Instruction) GetOpname() string {
	return opnames[i.Opcode]
}
