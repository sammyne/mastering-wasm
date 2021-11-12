package vm

import "github.com/sammyne/mastering-wasm/wavm/types"

type Memory struct {
	Type types.Memory
	Data []byte
}

func (m *Memory) Grow(n uint32) uint32 {
	oldSize := m.Size()
	old := m.Data
	m.Data = append(m.Data, make([]byte, n*types.PageSize)...)
	copy(m.Data, old)

	return uint32(oldSize)
}

func (m *Memory) Read(offset uint64, buf []byte) {
	copy(buf, m.Data[offset:])
}

func (m *Memory) Size() uint32 {
	return uint32(len(m.Data) / types.PageSize)
}

func (m *Memory) Write(offset uint64, data []byte) {
	copy(m.Data[offset:], data)
}

func NewMemory(t types.Memory) *Memory {
	return &Memory{Type: t, Data: make([]byte, t.Min*types.PageSize)}
}
