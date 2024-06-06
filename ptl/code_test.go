package ptl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAtomic(t *testing.T) {
	fnPre, fnEx, _, _ := getCode2MessageFn()
	assert.Nil(t, fnPre)
	assert.Nil(t, fnEx)
}
