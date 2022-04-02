package aes

import (
	"bytes"
	"errors"
)

var errOverflow = errors.New("overflow")

func PKCS7Padding(plaintext []byte, blockSize int) []byte {
	padding := blockSize - len(plaintext)%blockSize

	return append(plaintext, bytes.Repeat([]byte{byte(padding)}, padding)...)
}

func PKCS5Padding(plaintext []byte) []byte {
	return PKCS7Padding(plaintext, 8)
}

func PKCSUnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	unPadding := int(origData[length-1])

	if length < unPadding {
		return nil, errOverflow
	}

	return origData[:(length - unPadding)], nil
}

func PKCS5UnPadding(origData []byte) ([]byte, error) {
	return PKCSUnPadding(origData)
}

func PKCS7UnPadding(origData []byte) ([]byte, error) {
	return PKCSUnPadding(origData)
}
