package simencrypt

import (
	"strconv"

	"github.com/sgostarter/i/commerr"
)

var _DefXorKey = []byte{0xB2, 0x09, 0xBB, 0x55, 0x93, 0x6D, 0x44, 0x47}

type StringEnDecrypter interface {
	Encrypt(src string) string
	Decrypt(src string) string
}

func NewXor() *Xor {
	return NewXorEx(_DefXorKey)
}

func NewXorEx(key []byte) *Xor {
	if len(key) == 0 {
		key = _DefXorKey
	}

	return &Xor{
		key: key,
	}
}

type Xor struct {
	key []byte
}

func (xor *Xor) Encrypt(src string) string {
	s, _ := XorEncryptEx(src, xor.key)

	return s
}

func (xor *Xor) Decrypt(src string) string {
	s, _ := XorDecryptEx(src, xor.key)

	return s
}

func XorEncrypt(src string) string {
	s, _ := XorEncryptEx(src, _DefXorKey)

	return s
}

func XorEncryptE(src string) (string, error) {
	return XorEncryptEx(src, _DefXorKey)
}

func XorEncryptEx(src string, xorKey []byte) (result string, err error) {
	if len(xorKey) == 0 {
		err = commerr.ErrInvalidArgument

		return
	}

	j := 0

	bt := []rune(src)
	for i := 0; i < len(bt); i++ {
		s := strconv.FormatInt(int64(byte(bt[i])^xorKey[j]), 16)
		if len(s) == 1 {
			s = "0" + s
		}

		result = result + (s)
		j = (j + 1) % len(xorKey)
	}

	return
}

func XorDecrypt(src string) string {
	s, _ := XorDecryptEx(src, _DefXorKey)

	return s
}

func XorDecryptE(src string) (string, error) {
	return XorDecryptEx(src, _DefXorKey)
}

func XorDecryptEx(src string, xorKey []byte) (result string, err error) {
	if len(xorKey) == 0 {
		err = commerr.ErrInvalidArgument

		return
	}

	var s int64

	j := 0

	bt := []rune(src)

	for i := 0; i < len(src)/2; i++ {
		s, _ = strconv.ParseInt(string(bt[i*2:i*2+2]), 16, 0)
		result = result + string(byte(s)^xorKey[j])
		j = (j + 1) % len(xorKey)
	}

	return
}
