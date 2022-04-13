package simencrypt

import (
	"strconv"
)

var __defXorKey = []byte{0xB2, 0x09, 0xBB, 0x55, 0x93, 0x6D, 0x44, 0x47}

func NewXor() *Xor {
	return NewXorEx(__defXorKey)
}

func NewXorEx(key []byte) *Xor {
	return &Xor{
		key: key,
	}
}

type Xor struct {
	key []byte
}

func (xor *Xor) Encrypt(src string) string {
	return XorEncryptEx(src, xor.key)
}

func (xor *Xor) Decrypt(src string) string {
	return XorDecryptEx(src, xor.key)
}

func XorEncrypt(src string) string {
	return XorEncryptEx(src, __defXorKey)
}

func XorEncryptEx(src string, xorKey []byte) string {
	var result string
	j := 0
	s := ""
	bt := []rune(src)
	for i := 0; i < len(bt); i++ {
		s = strconv.FormatInt(int64(byte(bt[i])^xorKey[j]), 16)
		if len(s) == 1 {
			s = "0" + s
		}
		result = result + (s)
		j = (j + 1) % len(xorKey)
	}
	return result
}

func XorDecrypt(src string) string {
	return XorDecryptEx(src, __defXorKey)
}

func XorDecryptEx(src string, xorKey []byte) string {
	var result string
	var s int64
	j := 0
	bt := []rune(src)
	for i := 0; i < len(src)/2; i++ {
		s, _ = strconv.ParseInt(string(bt[i*2:i*2+2]), 16, 0)
		result = result + string(byte(s)^xorKey[j])
		j = (j + 1) % len(xorKey)
	}
	return result
}
