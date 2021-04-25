package vm

import "math"

type OperandStack struct {
	slots []uint64
}

func (s *OperandStack) Get(idx uint32) (uint64, bool) {
	if idx >= uint32(len(s.slots)) {
		return 0, false
	}

	return s.slots[idx], true
}

func (s *OperandStack) Len() int {
	return len(s.slots)
}

func (s *OperandStack) PopBool() (bool, bool) {
	v, ok := s.PopUint64()
	return v != 0, ok
}

func (s *OperandStack) PopFloat32() (float32, bool) {
	v, ok := s.PopUint32()
	return math.Float32frombits(v), ok
}

func (s *OperandStack) PopFloat64() (float64, bool) {
	v, ok := s.PopUint64()
	return math.Float64frombits(v), ok
}

func (s *OperandStack) PopInt32() (int32, bool) {
	v, ok := s.PopUint32()
	return int32(v), ok
}

func (s *OperandStack) PopInt64() (int64, bool) {
	v, ok := s.PopUint64()
	return int64(v), ok
}

func (s *OperandStack) PopUint32() (uint32, bool) {
	v, ok := s.PopUint64()
	return uint32(v), ok
}

func (s *OperandStack) PopUint64() (uint64, bool) {
	if len(s.slots) == 0 {
		return 0, false
	}

	val := s.slots[len(s.slots)-1]
	s.slots = s.slots[:len(s.slots)-1]

	return val, true
}

func (s *OperandStack) PopUint64s(n int) ([]uint64, bool) {
	if len(s.slots) < n {
		return nil, false
	}

	dividerIdx := len(s.slots) - n
	out := s.slots[dividerIdx:]
	s.slots = s.slots[:dividerIdx]

	return out, true
}

func (s *OperandStack) PushBool(val bool) {
	if val {
		s.PushUint64(1)
	} else {
		s.PushUint64(0)
	}
}

func (s *OperandStack) PushFloat32(val float32) {
	s.PushUint32(math.Float32bits(val))
}

func (s *OperandStack) PushFloat64(val float64) {
	s.PushUint64(math.Float64bits(val))
}

func (s *OperandStack) PushInt32(val int32) {
	s.PushUint32(uint32(val))
}

func (s *OperandStack) PushInt64(val int64) {
	s.PushUint64(uint64(val))
}

func (s *OperandStack) PushUint32(val uint32) {
	s.PushUint64(uint64(val))
}

func (s *OperandStack) PushUint64(val uint64) {
	s.slots = append(s.slots, val)
}

func (s *OperandStack) PushUint64s(values ...uint64) {
	s.slots = append(s.slots, values...)
}

func (s *OperandStack) Set(idx uint32, v uint64) bool {
	if idx > uint32(s.Len()) {
		return false
	}

	s.slots[idx] = v
	return true
}
