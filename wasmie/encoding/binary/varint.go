package binary

import (
	"io"
	"math/bits"
	"unsafe"
)

type bitsLen uint8

const (
	BitsLen32 bitsLen = (1 << 5) << iota
	BitsLen64
)

func ReadUvarint(r io.ByteReader, bitsLen bitsLen) (uint64, error) {
	out, _, err := readUvarint(r, bitsLen)
	return out, err
}

func ReadVarint(r io.ByteReader, bitsLen bitsLen) (int64, error) {
	out, last, err := readUvarint(r, bitsLen)
	if err != nil {
		return 0, err
	}

	out |= (uint64(last) >> 6) * ^(1<<bits.Len64(out) - 1)

	return *(*int64)(unsafe.Pointer(&out)), nil
}
