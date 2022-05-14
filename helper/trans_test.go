package helper

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBytesToInt(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	for idx := 0; idx < 100000; idx++ {
		// nolint: gosec
		n := int(rand.Int31())
		assert.EqualValues(t, n, BytesToInt(IntToBytes(n)))
	}
}

func TestTrans(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	for idx := 0; idx < 100000; idx++ {
		// nolint: gosec
		n := int(rand.Int31())
		s := strconv.Itoa(n)
		n2, err := strconv.Atoi(s)
		assert.Nil(t, err)
		assert.EqualValues(t, n, n2)
	}
}

func TestBytesToInt2(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	for idx := 0; idx < 100000; idx++ {
		// nolint: gosec
		n := int(rand.Int31())
		assert.EqualValues(t, n, BytesToInt2(IntToBytes2(n)))
	}
}
