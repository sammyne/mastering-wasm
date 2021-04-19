package binary

import (
	"encoding/binary"
	"io"
)

func readUvarint(r io.ByteReader, bitsLen bitsLen) (uint64, byte, error) {
	maxLen := binary.MaxVarintLen32
	if bitsLen == BitsLen64 {
		maxLen = binary.MaxVarintLen64
	}

	var out uint64
	var s uint
	ubits := uint(bitsLen)
	for i := 0; i < maxLen; i++ {
		b, err := r.ReadByte()
		if err != nil {
			return 0, 0, err
		}

		if b >= 0x80 {
			out, s = out|uint64(b&0x7f)<<s, s+7
			continue
		}

		if i == maxLen-1 && (b>>(ubits-s) > 0) {
			return 0, 0, ErrOverflow
		}

		return out | uint64(b&0x7f)<<s, b, nil
	}

	return 0, 0, ErrOverflow
}
