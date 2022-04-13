package simencrypt

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInt64(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	for idx := 0; idx < 1000000; idx++ {
		// nolint: gosec
		id := rand.Int63()
		es := EncodeID(id)
		id2 := DecodeID(es)
		assert.EqualValues(t, id, id2)
	}
}
