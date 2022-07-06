package certutils

import (
	"fmt"

	"github.com/sgostarter/libeasygo/commerr"
	"github.com/sgostarter/libeasygo/cuserror"
)

// CertNameToSignatureKey .
func CertNameToSignatureKey(certName string) string {
	return "cert_" + certName + "_sig"
}

// QueryPublicKeySignature .
func QueryPublicKeySignature(certName string, secureOption *SecureOption) (string, error) {
	if certName == "" || secureOption == nil || len(secureOption.CertSignatures) <= 0 {
		return "", commerr.ErrInvalidArgument
	}

	if signature, ok := secureOption.CertSignatures[CertNameToSignatureKey(certName)]; ok {
		return signature, nil
	}

	return "", cuserror.NewWithErrorMsg(fmt.Sprintf("cert %s no signature record", certName))
}

// VerifyCertPublicKey .
func VerifyCertPublicKey(publicKey interface{}, certName string, secureOption *SecureOption) (ok bool, err error) {
	publicKeyHash, err := GetPublicKeyHash(publicKey)
	if err != nil {
		return
	}

	signature, err := QueryPublicKeySignature(certName, secureOption)
	if err != nil {
		return
	}

	if publicKeyHash != signature {
		return
	}

	ok = true

	return
}
