package simencrypt

import (
	"encoding/hex"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestString(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	for idx := 0; idx < 1000000; idx++ {
		// nolint: gosec
		l := rand.Int63() % 0xff
		buf := make([]byte, l)
		// nolint: gosec
		rand.Read(buf)
		s := hex.EncodeToString(buf)
		es := EncodeString(s)
		s2, err := DecodeString(es)
		assert.Nil(t, err)
		assert.EqualValues(t, s, s2)
	}
}
