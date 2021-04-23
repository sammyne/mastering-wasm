package vm

import "math"

func (s *OperandStack) mustPopBool() bool {
	v, ok := s.PopUint64()
	if !ok {
		panic(ErrOperandPop)
	}

	return v != 0
}

func (s *OperandStack) mustPopFloat32() float32 {
	v, ok := s.PopUint32()
	if !ok {
		panic(ErrOperandPop)
	}

	return math.Float32frombits(v)
}

func (s *OperandStack) mustPopFloat64() float64 {
	v, ok := s.PopUint64()
	if !ok {
		panic(ErrOperandPop)
	}

	return math.Float64frombits(v)
}

func (s *OperandStack) mustPopInt32() int32 {
	v, ok := s.PopUint32()
	if !ok {
		panic(ErrOperandPop)
	}

	return int32(v)
}

func (s *OperandStack) mustPopInt64() int64 {
	v, ok := s.PopUint64()
	if !ok {
		panic(ErrOperandPop)
	}
	return int64(v)
}

func (s *OperandStack) mustPopUint32() uint32 {
	v, ok := s.PopUint64()
	if !ok {
		panic(ErrOperandPop)
	}
	return uint32(v)
}

func (s *OperandStack) mustPopUint64() uint64 {
	if len(s.slots) == 0 {
		panic(ErrOperandPop)
	}

	val := s.slots[len(s.slots)-1]
	s.slots = s.slots[:len(s.slots)-1]

	return val
}
