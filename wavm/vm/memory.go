package vm

import "github.com/sammyne/mastering-wasm/wavm/types"

type Memory struct {
	Type_ types.Memory
	Data  []byte
}

func (m *Memory) Grow(n uint32) uint32 {
	oldSize := m.Size()
	old := m.Data
	m.Data = append(m.Data, make([]byte, n*types.PageSize)...)
	copy(m.Data, old)

	return uint32(oldSize)
}

func (m *Memory) Read(offset uint64, buf []byte) error {
	copy(buf, m.Data[offset:])
	return nil // TODO: bound check
}

func (m *Memory) Size() uint32 {
	return uint32(len(m.Data) / types.PageSize)
}

func (m *Memory) Type() types.Memory {
	return m.Type_
}

func (m *Memory) Write(offset uint64, data []byte) error {
	copy(m.Data[offset:], data)
	return nil // TODO: bound check
}

func NewMemory(t types.Memory) *Memory {
	return &Memory{Type_: t, Data: make([]byte, t.Min*types.PageSize)}
}
