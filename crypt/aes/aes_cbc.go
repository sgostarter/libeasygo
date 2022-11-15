package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"fmt"

	"github.com/sgostarter/i/commerr"
)

func CBCEncrypt(origData, key []byte) (crypted []byte, err error) {
	defer func() {
		if errR := recover(); errR != nil {
			err = commerr.ErrCrash
		}
	}()

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("new cipher failed: %w", err)
	}

	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted = make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)

	return
}

func CBCDecrypt(encryptedData, key []byte) (decryptedData []byte, err error) {
	defer func() {
		if errR := recover(); errR != nil {
			err = commerr.ErrCrash
		}
	}()

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("new cipher failed: %w", err)
	}

	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(encryptedData))
	blockMode.CryptBlocks(origData, encryptedData)

	decryptedData, err = PKCS5UnPadding(origData)

	return
}
