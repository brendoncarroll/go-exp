// package sbe implements simple binary encoding formats for serializing and deserializing data.
package sbe

import (
	"encoding/binary"
	"fmt"
)

func AppendUint64(out []byte, x uint64) []byte {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], x)
	return append(out, buf[:]...)
}

func ReadUint64(data []byte) (uint64, []byte, error) {
	if len(data) < 8 {
		return 0, nil, fmt.Errorf("too short to contain uint64")
	}
	return binary.LittleEndian.Uint64(data[:8]), data[8:], nil
}

// AppendUint32 appends a little enddian Uint32 to out and returns the new slice.
func AppendUint32(out []byte, x uint32) []byte {
	var buf [4]byte
	binary.LittleEndian.PutUint32(buf[:], x)
	return append(out, buf[:]...)
}

// ReadUint32 reads a uint32 from x
func ReadUint32(x []byte) (uint32, []byte, error) {
	if len(x) < 4 {
		return 0, nil, fmt.Errorf("too short to contain uint32")
	}
	return binary.LittleEndian.Uint32(x[:4]), x[4:], nil
}

func AppendUint16(out []byte, x uint16) []byte {
	var buf [2]byte
	binary.LittleEndian.PutUint16(buf[:], x)
	return append(out, buf[:]...)
}

func ReadUint16(data []byte) (uint16, []byte, error) {
	if len(data) < 2 {
		return 0, nil, fmt.Errorf("too short to contain uint16")
	}
	return binary.LittleEndian.Uint16(data[:2]), data[2:], nil
}

// AppendUvarint appends a varint to out and returns the new slice.
func AppendUVarint(out []byte, x uint64) []byte {
	return binary.AppendUvarint(out, x)
}

// ReadUVarint attempts to read a varint encoded integer from x.
// If there is an error, then it is returned, otherwise the integer
// and the rest of the buffer are returned.
func ReadUVarint(x []byte) (uint64, []byte, error) {
	i, n := binary.Uvarint(x)
	if n <= 0 {
		return 0, nil, fmt.Errorf("too short to contain uvarint")
	}
	return i, x[n:], nil
}

// AppendLP appends a length prefixed buffer to out and returns the new slice.
// The length is varint encoded
func AppendLP(out []byte, x []byte) []byte {
	out = AppendUVarint(out, uint64(len(x)))
	return append(out, x...)
}

// ReadLP reads a byte slice from x that was encoded using AppendLP.
func ReadLP(x []byte) ([]byte, []byte, error) {
	l, x, err := ReadUVarint(x)
	if err != nil {
		return nil, nil, err
	}
	if l > uint64(len(x)) {
		return nil, nil, fmt.Errorf("buffer too short to contain %d bytes. only %d", l, len(x))
	}
	return x[:int(l)], x[int(l):], nil
}

// Uint64Bytes returns x as bytes in little endian order.
func Uint64Bytes(x uint64) []byte {
	var buf [8]byte
	binary.LittleEndian.PutUint64(buf[:], x)
	return buf[:]
}

// AppendLP16 appends a length-prefixed byte slice to out.
// the length is encoded as a 16-bit little-endian integer.
func AppendLP16(out []byte, x []byte) []byte {
	out = AppendUint16(out, uint16(len(x)))
	return append(out, x...)
}

// ReadLP16 reads a length-prefixed byte slice from data.
// ReadLP16 reads the format output by AppendLP16.
func ReadLP16(x []byte) (ret []byte, rest []byte, _ error) {
	n, rest, err := ReadUint16(x)
	if err != nil {
		return nil, nil, err
	}
	if len(rest) < int(n) {
		return nil, nil, fmt.Errorf("too short to contain lp16")
	}
	return rest[:n], rest[n:], nil
}

// ReadN reads n bytes from data.
func ReadN(data []byte, n int) ([]byte, []byte, error) {
	if len(data) < n {
		return nil, nil, fmt.Errorf("cannot read %d bytes from %d bytes", n, len(data))
	}
	return data[:n], data[n:], nil
}
