package simencrypt

import (
	"encoding/hex"
	"math"

	"github.com/sgostarter/i/commerr"
)

var _DefSeed = []byte{0xc1, 0x98, 0x96, 0x76, 0x91, 0x72, 0x2, 0x17, 0x65, 0xa8, 0xf7, 0xee, 0xdd, 0x62, 0xf1, 0xef}

type Cryptor struct {
	crypt   StringEnDecrypter
	defSeed []byte
}

func NewCryptor() *Cryptor {
	return NewCryptorEx(NewXor(), _DefSeed)
}

func NewCryptorEx(crypt StringEnDecrypter, defSeed []byte) *Cryptor {
	if crypt == nil {
		crypt = NewXor()
	}

	if len(defSeed) == 0 {
		defSeed = _DefSeed
	}

	return &Cryptor{
		crypt:   crypt,
		defSeed: defSeed,
	}
}

func (c *Cryptor) getPrefix(id int, seed []byte) string {
	return hex.EncodeToString([]byte{seed[id%len(seed)]})
}

func (c *Cryptor) EncryptString(s string) string {
	return c.EncryptStringEx(s, c.defSeed)
}

func (c *Cryptor) EncryptStringEx(s string, seed []byte) string {
	return c.getPrefix(GetStringIndex(s), seed) + c.crypt.Encrypt(EncodeString(s))
}

func (c *Cryptor) DecryptString(s string) (ret string, err error) {
	return c.DecryptStringEx(s, c.defSeed)
}

func (c *Cryptor) DecryptStringEx(s string, seed []byte) (ret string, err error) {
	ret, err = DecodeString(c.crypt.Decrypt(s[2:]))

	if err != nil {
		return
	}

	if s[:2] != c.getPrefix(GetStringIndex(s), seed) {
		err = commerr.ErrOutOfRange
	}

	return
}

func (c *Cryptor) EncryptInt64(n int64) string {
	return c.EncryptInt64Ex(n, c.defSeed)
}

func (c *Cryptor) EncryptInt64Ex(n int64, seed []byte) string {
	return c.getPrefix(int(n%math.MaxInt), seed) + c.crypt.Encrypt(EncodeID(n))
}

func (c *Cryptor) DecryptInt64(s string) (ret int64, err error) {
	return c.DecryptInt64Ex(s, c.defSeed)
}

func (c *Cryptor) DecryptInt64Ex(s string, seed []byte) (ret int64, err error) {
	ret = DecodeID(c.crypt.Decrypt(s[2:]))

	if s[:2] != c.getPrefix(int(ret%math.MaxInt), seed) {
		err = commerr.ErrOutOfRange
	}

	return
}
