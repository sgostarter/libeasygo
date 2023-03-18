package simencrypt

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// nolint
func TestUint64(t *testing.T) {
	t.SkipNow()

	rand.Seed(time.Now().UnixNano())

	var n1 = uint64(1)
	t.Log(n1, XorUInt64(n1), XorUInt64(XorUInt64(n1)))
	n1 = 1000
	t.Log(n1, XorUInt64(n1), XorUInt64(XorUInt64(n1)))
	n1 = 1001
	t.Log(n1, XorUInt64(n1), XorUInt64(XorUInt64(n1)))
	n1 = 1002
	t.Log(n1, XorUInt64(n1), XorUInt64(XorUInt64(n1)))
	n1 = 9999
	t.Log(n1, XorUInt64(n1), XorUInt64(XorUInt64(n1)))
	n1 = 0xF923897619098765
	t.Log(n1, XorUInt64(n1), XorUInt64(XorUInt64(n1)))

	for idx := 0; idx < 10000008; idx++ {
		// nolint: gosec
		id := rand.Uint64()
		es := XorUInt64(id)
		id2 := XorUInt64(es)
		assert.EqualValues(t, id, id2)
	}
}

// nolint
func TestUint64Crypt(t *testing.T) {
	t.SkipNow()

	rand.Seed(time.Now().UnixNano())

	fnDecrypt := func(s string) uint64 {
		n, err := DecryptUint64(s)
		assert.Nil(t, err)
		return n
	}

	var n1 = uint64(1)
	t.Log(n1, EncryptUInt64(n1), fnDecrypt(EncryptUInt64(n1)))
	n1 = 1000
	t.Log(n1, EncryptUInt64(n1), fnDecrypt(EncryptUInt64(n1)))
	n1 = 1001
	t.Log(n1, EncryptUInt64(n1), fnDecrypt(EncryptUInt64(n1)))
	n1 = 1002
	t.Log(n1, EncryptUInt64(n1), fnDecrypt(EncryptUInt64(n1)))
	n1 = 9999
	t.Log(n1, EncryptUInt64(n1), fnDecrypt(EncryptUInt64(n1)))
	n1 = 0xF923897619098765
	t.Log(n1, EncryptUInt64(n1), fnDecrypt(EncryptUInt64(n1)))

	for idx := 0; idx < 10000008; idx++ {
		// nolint: gosec
		id := rand.Uint64()
		es := EncryptUInt64(id)
		id2 := fnDecrypt(es)
		assert.EqualValues(t, id, id2)
	}
}

func TestChangeUint64(t *testing.T) {
	t.SkipNow()

	var n uint64
	n = 0x1203
	m, f := ConfuseUint64(n)
	n2 := UnConfuseUint64(m, f)
	t.Log(n, m, n2)
	assert.EqualValues(t, n, n2)

	n = 0xAE99
	m, f = ConfuseUint64(n)
	n2 = UnConfuseUint64(m, f)
	t.Log(n, m, n2)
	assert.EqualValues(t, n, n2)

	n = 1000
	m, f = ConfuseUint64(n)
	n2 = UnConfuseUint64(m, f)
	t.Log(n, m, n2)
	assert.EqualValues(t, n, n2)

	n = 1001
	m, f = ConfuseUint64(n)
	n2 = UnConfuseUint64(m, f)
	t.Log(n, m, n2)
	assert.EqualValues(t, n, n2)

	n = 1002
	m, f = ConfuseUint64(n)
	n2 = UnConfuseUint64(m, f)
	t.Log(n, m, n2)
	assert.EqualValues(t, n, n2)

	n = 0x25d5d21a0c8814b4
	m, f = ConfuseUint64(n)
	n2 = UnConfuseUint64(m, f)
	t.Log(n, m, n2)
	assert.EqualValues(t, n, n2)

	n = 0xfc001183bea872dd
	m, f = ConfuseUint64(n)
	n2 = UnConfuseUint64(m, f)
	t.Log(n, m, n2)
	assert.EqualValues(t, n, n2)

	for idx := 0; idx < 10000008; idx++ {
		// nolint: gosec
		id := rand.Uint64()
		es, f := ConfuseUint64(id)
		id2 := UnConfuseUint64(es, f)
		assert.EqualValues(t, id, id2)
	}
}
