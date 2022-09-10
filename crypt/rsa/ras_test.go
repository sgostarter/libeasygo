package rsa

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrypt(t *testing.T) {
	rsaKey, err := GenRSAKey()
	assert.Nil(t, err)

	rawD := []byte("hoho")
	ed, err := Encrypt(GetPublicKeyFromPrivateKey(rsaKey), rawD)
	assert.Nil(t, err)

	sig, err := Signature(rsaKey, rawD)
	assert.Nil(t, err)

	pd, err := Decrypt(rsaKey, ed)
	assert.Nil(t, err)
	assert.EqualValues(t, pd, rawD)

	err = VerifySignature(GetPublicKeyFromPrivateKey(rsaKey), rawD, sig)
	assert.Nil(t, err)

	publicKey2, err := DecodeRSAPublicKey(ConvRSAPublicKeyToBytes(GetPublicKeyFromPrivateKey(rsaKey)))
	assert.Nil(t, err)

	privateKey2, err := DecodeRSAPrivateKey(ConvRSAPrivateKeyToBytes(rsaKey))
	assert.Nil(t, err)

	pd, err = Decrypt(privateKey2, ed)
	assert.Nil(t, err)
	assert.EqualValues(t, pd, rawD)

	err = VerifySignature(publicKey2, rawD, sig)
	assert.Nil(t, err)
}
