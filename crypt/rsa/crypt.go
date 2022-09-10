package rsa

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
)

func Encrypt(publicKey *rsa.PublicKey, d []byte) ([]byte, error) {
	return rsa.EncryptOAEP(sha256.New(), rand.Reader, publicKey, d, nil)
}

func Decrypt(privateKey *rsa.PrivateKey, d []byte) ([]byte, error) {
	return privateKey.Decrypt(nil, d, &rsa.OAEPOptions{Hash: crypto.SHA256})
}

func Signature(privateKey *rsa.PrivateKey, d []byte) (signature []byte, err error) {
	msgHash := sha256.New()

	_, err = msgHash.Write(d)
	if err != nil {
		return
	}

	msgHashSum := msgHash.Sum(nil)

	signature, err = rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, msgHashSum, nil)

	return
}

func VerifySignature(publicKey *rsa.PublicKey, data, signature []byte) (err error) {
	msgHash := sha256.New()

	_, err = msgHash.Write(data)
	if err != nil {
		return
	}

	msgHashSum := msgHash.Sum(nil)

	err = rsa.VerifyPSS(publicKey, crypto.SHA256, msgHashSum, signature, nil)

	return
}
