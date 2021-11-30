package wavm

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"

	localBinaryPkg "github.com/sammyne/mastering-wasm/wavm/encoding/binary"
	"github.com/sammyne/mastering-wasm/wavm/types"
)

type Decoder struct {
	*bytes.Reader
}

func (d *Decoder) DecodeBytes() ([]byte, error) {
	n, err := d.DecodeUvarint32()
	if err != nil {
		return nil, fmt.Errorf("read bytes length: %w", err)
	}

	out := make([]byte, n)
	if _, err := io.ReadFull(d.Reader, out); err != nil {
		return nil, fmt.Errorf("read full bytes: %w", err)
	}

	return out, nil
}

func (d *Decoder) DecodeFloat32() (float32, error) {
	var out uint32
	if err := binary.Read(d.Reader, binary.LittleEndian, &out); err != nil {
		return 0, err
	}

	return math.Float32frombits(out), nil
}

func (d *Decoder) DecodeFloat64() (float64, error) {
	var out uint64
	if err := binary.Read(d.Reader, binary.LittleEndian, &out); err != nil {
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
			var c types.Custom
			if err := d.decodeCustom(&c); err != nil {
				return nil, fmt.Errorf("decode custom section: %w", err)
			}
			out.Customs = append(out.Customs, c)
			continue
		}

		if sectionID <= prevSectionID {
			return nil, fmt.Errorf("malformed section ID: %v", sectionID)
		}
		prevSectionID = sectionID

		sectionLen, err := d.DecodeUvarint32()
		if err != nil {
			return nil, fmt.Errorf("decode section byte count: %w", err)
		}
		ell := d.Len()

		if err := d.decodeNonCustomSectionIntoModule(sectionID, out); err != nil {
			return nil, fmt.Errorf("bad section(%d): %w", sectionID, err)
		}

		if d.Len()+int(sectionLen) != ell {
			return nil, fmt.Errorf("section(%d) size mismatch: %w", sectionID, err)
		}
	}

	return out, nil
}

func (d *Decoder) DecodeName() (string, error) {
	buf, err := d.DecodeBytes()
	if err != nil {
		return "", err
	}

	return string(buf), nil
}

func (d *Decoder) DecodeUint32() (uint32, error) {
	var out uint32
	if err := binary.Read(d.Reader, binary.LittleEndian, &out); err != nil {
		return 0, err
	}

	return out, nil
}

func (d *Decoder) DecodeUvarint32() (uint32, error) {
	out, err := localBinaryPkg.ReadUvarint(d.Reader, localBinaryPkg.BitsLen32)
	return uint32(out), err
}

func (d *Decoder) DecodeVarint32() (int32, error) {
	out, err := localBinaryPkg.ReadVarint(d.Reader, localBinaryPkg.BitsLen32)
	return int32(out), err
}

func (d *Decoder) DecodeVarint64() (int64, error) {
	out, err := localBinaryPkg.ReadVarint(d.Reader, localBinaryPkg.BitsLen64)
	return int64(out), err
}

func NewDecoder(buf []byte) *Decoder {
	return &Decoder{Reader: bytes.NewReader(buf)}
}
