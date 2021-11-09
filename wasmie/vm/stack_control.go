package vm

import "github.com/sammyne/mastering-wasm/wasmie/types"

type ControlFrame struct {
	Opcode    byte
	BlockType types.FuncType
	Expr      []types.Instruction
	BP        int
	PC        int
}

type ControlStack struct {
	frames []ControlFrame
}

func (s *ControlStack) Len() int {
	return len(s.frames)
}

func (s *ControlStack) Pop() (ControlFrame, bool) {
	ell := s.Len()
	if ell == 0 {
		return ControlFrame{}, false
	}

	out := s.frames[ell-1]
	s.frames = s.frames[:ell-1]

	return out, true
}

func (s *ControlStack) Push(f ControlFrame) {
	s.frames = append(s.frames, f)
}

func (s *ControlStack) Top() (*ControlFrame, bool) {
	ell := s.Len()
	if ell == 0 {
		return nil, false
	}

	return &s.frames[ell-1], true
}

// TopCallFrame return the top-most call frame the its depth counting from top.
func (s *ControlStack) TopCallFrame() (ControlFrame, int, bool) {
	for i := s.Len() - 1; i >= 0; i-- {
		if v := s.frames[i]; v.Opcode == types.OpcodeCall {
			return v, s.Len() - i - 1, true
		}
	}

	return ControlFrame{}, -1, false
}

func NewControlFrame(
	opcode byte, blockType types.FuncType, expr []types.Instruction, BP int) ControlFrame {
	out := ControlFrame{Opcode: opcode, BlockType: blockType, Expr: expr, BP: BP, PC: 0}
	return out
}
