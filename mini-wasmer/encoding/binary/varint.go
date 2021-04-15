package binary

import (
	"encoding/binary"
	"io"
)

type bits uint8

const (
	Bits32 bits = (1 << 5) << iota
	Bits64
)

func ReadUvarint(r io.ByteReader, bits bits) (uint64, error) {
	maxLen := binary.MaxVarintLen32
	if bits == Bits64 {
		maxLen = binary.MaxVarintLen64
	}

	var out uint64
	var s uint
	ubits := uint(bits)
	for i := 0; i < maxLen; i++ {
		b, err := r.ReadByte()
		if err != nil {
			return 0, err
		}

		if b >= 0x80 {
			out, s = out|uint64(b&0x7f)<<s, s+7
			continue
		}

		if i == maxLen-1 && (b>>(ubits-s) > 0) {
			return 0, ErrOverflow
		}

		return out | uint64(b&0x7f)<<s, nil
	}

	return 0, ErrOverflow
}

func ReadVarint(r io.ByteReader, bits bits) (int64, error) {
	maxLen := binary.MaxVarintLen32
	if bits == Bits64 {
		maxLen = binary.MaxVarintLen64
	}

	var out int64
	var s uint
	ubits := uint(bits)
	for i := 0; i < maxLen; i++ {
		b, err := r.ReadByte()
		if err != nil {
			return 0, err
		}

		if b >= 0x80 {
			out, s = out|int64(b&0x7f)<<s, s+7
			continue
		}

		if i == maxLen-1 && (b>>(ubits-s) > 0) {
			return 0, ErrOverflow
		}

		out, s = out|int64(b&0x7f)<<s, s+7
		out |= (int64(b) >> 6) * (-1 << s)
		return out, nil
	}

	return 0, ErrOverflow
}
