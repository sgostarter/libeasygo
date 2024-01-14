package ptl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAtomic(t *testing.T) {
	fn := getCode2MessageFn()
	assert.Nil(t, fn)
}
