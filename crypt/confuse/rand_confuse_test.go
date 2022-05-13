package confuse

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRandConfuse(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	c := NewRandConfuseEx(NewDefN(), NewDefN())

	for idx := 0; idx < 10000; idx++ {
		// nolint: gosec
		src := make([]byte, rand.Int31n(100))
		// nolint: gosec
		_, _ = rand.Read(src)
		dst, err := c.Seal(src)
		assert.Nil(t, err)
		src2, err := c.Open(dst)
		assert.Nil(t, err)
		assert.EqualValues(t, src, src2)
	}
}
