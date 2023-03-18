package simencrypt

import (
	"encoding/hex"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestXor(t *testing.T) {
	t.SkipNow()

	s := "E12"
	s1 := XorEncrypt(s)
	t.Log(s1)
	t.Log(XorDecrypt(s1))

	s = "E12344589"
	s1 = XorEncrypt(s)
	t.Log(s1)
	t.Log(XorDecrypt(s1))

	var lastS string

	for idx := 0; idx < 100000; idx++ {
		buf := make([]byte, 30)
		// nolint: gosec
		rand.Read(buf)
		s = hex.EncodeToString(buf)
		es := XorEncrypt(s)
		s2 := XorDecrypt(es)
		assert.Equal(t, s, s2)
		assert.NotEqual(t, lastS, s)
		lastS = s
	}

	xor := NewXorEx([]byte{0x11, 0x11})

	for idx := 0; idx < 1000; idx++ {
		buf := make([]byte, 30)
		// nolint: gosec
		rand.Read(buf)
		s = hex.EncodeToString(buf)
		es := xor.Encrypt(s)
		s2 := xor.Decrypt(es)
		assert.Equal(t, s, s2)
		assert.NotEqual(t, lastS, s)
		lastS = s
	}
}
