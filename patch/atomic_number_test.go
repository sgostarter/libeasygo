package patch

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// nolint
func TestAtomicNumber(t *testing.T) {
	var n64 AtomicNumber[int64]
	assert.EqualValues(t, 1, n64.Inc())
	assert.EqualValues(t, 2, n64.Inc())
	assert.EqualValues(t, 1, n64.Dec())
	assert.EqualValues(t, 0, n64.Dec())
	assert.EqualValues(t, 0, n64.Load())
	n64.Store(10008)
	assert.EqualValues(t, 10008, n64.Load())

	n64.Store(10)
	n := n64
	t.Log(n.Load())
}
