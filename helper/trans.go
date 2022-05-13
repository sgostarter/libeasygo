package helper

import (
	"bytes"
	"encoding/binary"
)

func IntToBytes(n int) []byte {
	buf := bytes.NewBuffer([]byte{})
	_ = binary.Write(buf, binary.BigEndian, uint32(n))

	return buf.Bytes()
}

func BytesToInt(d []byte) int {
	buf := bytes.NewBuffer(d)

	var data uint32

	_ = binary.Read(buf, binary.BigEndian, &data)

	return int(data)
}
