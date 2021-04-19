package wasmer

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"

	localBinaryPkg "github.com/sammyne/mastering-wasm/mini-wasmer/encoding/binary"
	"github.com/sammyne/mastering-wasm/mini-wasmer/types"
)

type Decoder struct {
	*bytes.Reader
}

func (r *Decoder) DecodeBytes() ([]byte, error) {
	n, err := r.DecodeVarint32()
	if err != nil {
		return nil, fmt.Errorf("read bytes length: %w", err)
	}

	out := make([]byte, n)
	if _, err := io.ReadFull(r.Reader, out); err != nil {
		return nil, fmt.Errorf("read full bytes: %w", err)
	}

	return out, nil
}

func (r *Decoder) DecodeFloat32() (float32, error) {
	var out uint32
	if err := binary.Read(r.Reader, binary.LittleEndian, &out); err != nil {
		return 0, err
	}

	return math.Float32frombits(out), nil
}

func (r *Decoder) DecodeFloat64() (float64, error) {
	var out uint64
	if err := binary.Read(r.Reader, binary.LittleEndian, &out); err != nil {
		return 0, err
	}

	return math.Float64frombits(out), nil
}

func (d *Decoder) DecodeModule() (*Module, error) {
	magic, err := d.DecodeUint32()
	if err != nil {
		return nil, fmt.Errorf("decode magic: %w", err)
	}

	version, err := d.DecodeUint32()
	if err != nil {
		return nil, fmt.Errorf("decode version: %w", err)
	}

	// sections
	out := &Module{Magic: magic, Version: version}

	var prevSectionID byte
	for d.Len() > 0 {
		sectionID, err := d.ReadByte()
		if err != nil {
			return nil, fmt.Errorf("read section ID: %w", err)
		}

		if sectionID == types.SectionIDCustom {
			s, err := d.decodeCustomSection()
			if err != nil {
				return nil, fmt.Errorf("decode custom section: %w", err)
			}
			out.Customs = append(out.Customs, *s)
			continue
		}

		if sectionID <= prevSectionID {
			return nil, fmt.Errorf("malformed section ID: %v", sectionID)
		}
		prevSectionID = sectionID

		if err := d.decodeNonCustomSectionIntoModule(sectionID, out); err != nil {
			return nil, fmt.Errorf("bad section(%d): %w", sectionID, err)
		}
	}

	return out, nil
}

func (r *Decoder) DecodeName() (string, error) {
	buf, err := r.DecodeBytes()
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func (r *Decoder) DecodeUint32() (uint32, error) {
	var out uint32
	if err := binary.Read(r.Reader, binary.LittleEndian, &out); err != nil {
		return 0, err
	}

	return out, nil
}

func (r *Decoder) DecodeUvarint32() (uint32, error) {
	out, err := localBinaryPkg.ReadUvarint(r.Reader, localBinaryPkg.BitsLen32)
	return uint32(out), err
}

func (r *Decoder) DecodeVarint32() (int32, error) {
	out, err := localBinaryPkg.ReadVarint(r.Reader, localBinaryPkg.BitsLen32)
	return int32(out), err
}

func (r *Decoder) DecodeVarint64() (int64, error) {
	out, err := localBinaryPkg.ReadVarint(r.Reader, localBinaryPkg.BitsLen64)
	return int64(out), err
}

func NewDecoder(buf []byte) *Decoder {
	return &Decoder{Reader: bytes.NewReader(buf)}
}
