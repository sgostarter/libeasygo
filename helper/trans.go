package helper

import (
	"bytes"
	"encoding/binary"
)

func IntToBytes2(n int) []byte {
	buf := bytes.NewBuffer([]byte{})
	_ = binary.Write(buf, binary.BigEndian, uint32(n))

	return buf.Bytes()
}

func BytesToInt2(d []byte) int {
	buf := bytes.NewBuffer(d)

	var data uint32

	_ = binary.Read(buf, binary.BigEndian, &data)

	return int(data)
}

func IntToBytes(n int) []byte {
	buf := make([]byte, 4)

	buf[0] = uint8(n)
	buf[1] = uint8(n >> 8)
	buf[2] = uint8(n >> 16)
	buf[3] = uint8(n >> 24)

	return buf
}

func BytesToInt(d []byte) int {
	return int(d[0]) + int(d[1])<<8 + int(d[2])<<16 + int(d[3])<<24
}
