package linker

import "github.com/sammyne/mastering-wasm/wavm/types"

type Memory interface {
	Grow(v uint32) uint32
	Read(offset uint64, buf []byte) error
	Size() uint32
	Type() types.Memory
	Write(offset uint64, buf []byte) error
}
