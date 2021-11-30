package linker

import "github.com/sammyne/mastering-wasm/wavm/types"

type Table interface {
	Grow(n uint32)
	Size() uint32
	Type() types.Table
	GetElem(idx uint32) (Function, error)
	SetElem(idx uint32, elem Function) error
}
