package certutils

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/sgostarter/libeasygo/cuserror"
)

func GetCertPublicKeyHash(certFile string) (string, error) {
	certPEMBlock, err := os.ReadFile(certFile)
	if err != nil {
		return "", cuserror.NewWithErrorMsg(fmt.Sprintf("%v cert load failed: %v", certFile, err))
	}

	certDERBlock, _ := pem.Decode(certPEMBlock)

	x509Cert, err := x509.ParseCertificate(certDERBlock.Bytes)
	if err != nil {
		return "", cuserror.NewWithErrorMsg(fmt.Sprintf("%v cert ParseCertificate failed: %v", certFile, err))
	}

	pubKey, err := GetPublicKeyHash(x509Cert.PublicKey)
	if err != nil {
		return "", cuserror.NewWithErrorMsg(fmt.Sprintf("%v cert GetPublicKeyHash failed: %v", certFile, err))
	}

	return pubKey, nil
}

// GetPublicKeyHash .
func GetPublicKeyHash(publicKey interface{}) (string, error) {
	pubKey, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", err
	}

	h := hmac.New(sha256.New, []byte("1"))
	_, _ = h.Write(pubKey)

	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}
