package simencrypt

import (
	"encoding/hex"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// nolint
func TestCryptor(t *testing.T) {
	t.SkipNow()

	rand.Seed(time.Now().UnixNano())

	crypto := NewCryptor()

	for idx := 0; idx < 100000; idx++ {
		// nolint: gosec
		id := rand.Int63()
		eID := crypto.EncryptInt64(id)
		dID, err := crypto.DecryptInt64(eID)
		assert.Nil(t, err)
		assert.EqualValues(t, id, dID)
	}

	for idx := 0; idx < 100000; idx++ {
		// nolint: gosec
		l := rand.Int63() % 0xff
		buf := make([]byte, l)
		rand.Read(buf)
		s := hex.EncodeToString(buf)
		es := crypto.EncryptString(s)
		ds, err := crypto.DecryptString(es)
		assert.Nil(t, err)
		assert.EqualValues(t, s, ds)
	}
}
