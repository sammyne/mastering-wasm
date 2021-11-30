package vm

import (
	"github.com/sammyne/mastering-wasm/wavm/linker"
	"github.com/sammyne/mastering-wasm/wavm/types"
)

type Table struct {
	type_ types.Table
	elems []linker.Function
}

func (t *Table) GetElem(i uint32) (linker.Function, error) {
	if i >= uint32(len(t.elems)) {
		return nil, ErrIndexOutOfBound
	}
	return t.elems[i], nil
}

func (t *Table) Grow(n uint32) {
	t.elems = append(t.elems, make([]linker.Function, n)...)
}

func (t *Table) SetElem(i uint32, e linker.Function) error {
	if i >= uint32(len(t.elems)) {
		return ErrIndexOutOfBound
	}

	t.elems[i] = e
	return nil
}

func (t *Table) Size() uint32 {
	return uint32(len(t.elems))
}

func (t *Table) Type() types.Table {
	return t.type_
}

func newTable(t types.Table) *Table {
	out := &Table{
		type_: t,
		elems: make([]linker.Function, t.Limits.Min),
	}
	return out
}
