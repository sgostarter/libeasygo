package debug

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLatency(t *testing.T) {
	latency, err := Latency("www.ymicj.com")
	assert.Nil(t, err)

	t.Log(latency)

	latency, err = Latency("www.ymicj.com")
	assert.Nil(t, err)

	t.Log(latency)
}
