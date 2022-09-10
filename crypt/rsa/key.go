package rsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"

	"github.com/sgostarter/libeasygo/cuserror"
)

func GenRSAKeyEx(bits int) (priKey *rsa.PrivateKey, err error) {
	return rsa.GenerateKey(rand.Reader, bits)
}

func GenRSAKey() (*rsa.PrivateKey, error) {
	return GenRSAKeyEx(2048)
}

func GetPublicKeyFromPrivateKey(key *rsa.PrivateKey) *rsa.PublicKey {
	publicKey, _ := key.Public().(*rsa.PublicKey)

	return publicKey
}

func ConvRSAPrivateKeyToBytes(key *rsa.PrivateKey) []byte {
	return x509.MarshalPKCS1PrivateKey(key)
}

func ConvRSAPublicKeyToBytes(key *rsa.PublicKey) []byte {
	return x509.MarshalPKCS1PublicKey(key)
}

func GenRSAKeyPairEx(bits int) (priKey, pubKey []byte, err error) {
	key, err := GenRSAKeyEx(bits)
	if err != nil {
		return
	}

	pKey := GetPublicKeyFromPrivateKey(key)
	if pKey == nil {
		err = cuserror.NewWithErrorMsg("invalidKey")

		return
	}

	priKey = ConvRSAPrivateKeyToBytes(key)
	pubKey = ConvRSAPublicKeyToBytes(pKey)

	return
}

func GenRSAKeyPair() (priKey, pubKey []byte, err error) {
	return GenRSAKeyPairEx(2048)
}

func GenRSAKeyPairStringEx(bits int) (privateKey, publicKey string, err error) {
	priKeyD, pubKeyD, err := GenRSAKeyPairEx(bits)
	if err != nil {
		return
	}

	privateKey = hex.EncodeToString(priKeyD)
	publicKey = hex.EncodeToString(pubKeyD)

	return
}

func GenRSAKeyPairString() (privateKey, publicKey string, err error) {
	return GenRSAKeyPairStringEx(2048)
}

func DecodeRSAPrivateKey(key []byte) (*rsa.PrivateKey, error) {
	return x509.ParsePKCS1PrivateKey(key)
}

func DecodeRSAPublicKey(key []byte) (*rsa.PublicKey, error) {
	return x509.ParsePKCS1PublicKey(key)
}

func DecodeRSAPrivateKeyString(key string) (*rsa.PrivateKey, error) {
	d, err := hex.DecodeString(key)
	if err != nil {
		return nil, err
	}

	return DecodeRSAPrivateKey(d)
}

func DecodeRSAPublicKeyString(key string) (*rsa.PublicKey, error) {
	d, err := hex.DecodeString(key)
	if err != nil {
		return nil, err
	}

	return DecodeRSAPublicKey(d)
}
