package crypt

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"math/rand"

	"golang.org/x/crypto/pbkdf2"
)

const (
	saltMinLen = 8
	saltMaxLen = 32
	keyLen     = 32
)

func HashPassword(password string, iterCount int) (hashedPassword string, err error) {
	salt, err := _passwordGenRandSalt()
	if err != nil {
		return
	}

	d := _passwordHashPasswordWithSalt([]byte(password), salt, iterCount)
	d = append(d, salt...)

	hashedPassword = base64.StdEncoding.EncodeToString(d)

	return
}

func CheckHashedPassword(password, hashedPassword string, iterCount int) (ok bool) {
	if len(hashedPassword) == 0 {
		return
	}

	d, err := base64.StdEncoding.DecodeString(hashedPassword)
	if err != nil {
		return
	}

	salt := d[keyLen:]

	ok = bytes.Equal(_passwordHashPasswordWithSalt([]byte(password), salt, iterCount), d[:keyLen])

	return
}

//
//
//

func _passwordGenRandSalt() (salt []byte, err error) {
	salt = make([]byte, rand.Intn(saltMaxLen-saltMinLen)+saltMinLen) // nolint:gosec

	_, err = rand.Read(salt) // nolint:gosec

	return
}

func _passwordHashPasswordWithSalt(password, salt []byte, iterCount int) []byte {
	return pbkdf2.Key(password, salt, iterCount, keyLen, sha256.New)
}
