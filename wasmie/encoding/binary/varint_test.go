package binary_test

import (
	"bytes"
	"errors"
	"strconv"
	"testing"
	"unsafe"

	"github.com/sammyne/mastering-wasm/wasmie/encoding/binary"
)

func TestReadUvarint32(t *testing.T) {
	type expect struct {
		val uint64
		err error
	}

	testVector := []struct {
		data   []byte
		expect expect
	}{
		{
			[]byte{0b11100101, 0b10001110, 0b00100110},
			expect{val: 0x098765},
		},
		{
			[]byte{0b1111_1111, 0b1111_1111, 0b1110_0101, 0b1000_1110, 0b0010_0110},
			expect{err: binary.ErrOverflow},
		},
	}

	for i, c := range testVector {
		got, err := binary.ReadUvarint(bytes.NewReader(c.data), binary.BitsLen32)
		if !errors.Is(err, c.expect.err) {
			t.Fatalf("#%d unexpected error: expect %v, got %v", i, c.expect.err, err)
		}

		if c.expect.val != got {
			t.Fatalf("#%d failed: expect %08x, got %08x", i, c.expect.val, got)
		}
	}
}

func TestReadUvarint64(t *testing.T) {
	type expect struct {
		val uint64
		err error
	}

	testVector := []struct {
		data   []byte
		expect expect
	}{
		{
			[]byte{0b11100101, 0b10001110, 0b00100110},
			expect{val: 0x098765},
		},
		{
			[]byte{0b1111_1111, 0b1111_1111, 0b1110_0101, 0b1000_1110, 0b0010_0110},
			expect{val: 0x261d97fff},
		},
	}

	for i, c := range testVector {
		got, err := binary.ReadUvarint(bytes.NewReader(c.data), binary.BitsLen64)
		if !errors.Is(err, c.expect.err) {
			t.Fatalf("#%d unexpected error: expect %v, got %v", i, c.expect.err, err)
		}

		if c.expect.val != got {
			t.Fatalf("#%d failed: expect %08x, got %08x", i, c.expect.val, got)
		}
	}
}

func TestReadVarint32(t *testing.T) {
	type expect struct {
		val int64
		err error
	}

	testVector := []struct {
		data   []byte
		expect expect
	}{
		{
			[]byte{0b11000000, 0b10111011, 0b01111000},
			expect{val: -0x001e240},
		},
	}

	for i, c := range testVector {
		got, err := binary.ReadVarint(bytes.NewReader(c.data), binary.BitsLen32)
		if !errors.Is(err, c.expect.err) {
			t.Fatalf("#%d unexpected error: expect %v, got %v", i, c.expect.err, err)
		}

		if c.expect.val != got {
			expect2 := strconv.FormatUint(uint64(c.expect.val), 2)
			got2 := strconv.FormatUint(uint64(got), 2)
			t.Fatalf("#%d failed:\n expect %s,\n    got %s", i, expect2, got2)
		}
	}
}

func TestReadVarint64(t *testing.T) {
	type expect struct {
		val int64
		err error
	}

	testVector := []struct {
		data   []byte
		expect expect
	}{
		{
			[]byte{0b11000000, 0b10111011, 0b01111000},
			expect{
				val: Int64frombits(
					0b11111111_11111111_11111111_11111111_11111111_11111110_00011101_11000000),
			},
		},
	}

	for i, c := range testVector {
		got, err := binary.ReadVarint(bytes.NewReader(c.data), binary.BitsLen64)
		if !errors.Is(err, c.expect.err) {
			t.Fatalf("#%d unexpected error: expect %v, got %v", i, c.expect.err, err)
		}

		if c.expect.val != got {
			//expect2 := strconv.FormatUint(uint64(c.expect.val), 2)
			//got2 := strconv.FormatUint(uint64(got), 2)
			t.Fatalf("#%d failed:\n expect %064b,\n    got %064b", i, c.expect.val, got)
		}
	}
}

func Int64frombits(b uint64) int64 {
	return *(*int64)(unsafe.Pointer(&b))
}
