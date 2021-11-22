package vm

import "github.com/sammyne/mastering-wasm/wavm/types"

type Table struct {
	type_ types.Table
	elems []Func
}

func (t *Table) GetElem(i uint32) (Func, error) {
	if i >= uint32(len(t.elems)) {
		return Func{}, ErrIndexOutOfBound
	}
	return t.elems[i], nil
}

func (t *Table) SetElem(i uint32, e Func) error {
	if i >= uint32(len(t.elems)) {
		return ErrIndexOutOfBound
	}

	t.elems[i] = e
	return nil
}

func newTable(t types.Table) *Table {
	out := &Table{
		type_: t,
		elems: make([]Func, t.Limits.Min),
	}
	return out
}
