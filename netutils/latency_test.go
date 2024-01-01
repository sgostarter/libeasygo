package netutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLatency(t *testing.T) {
	latency, err := Latency("www.ymicj.com")
	assert.Nil(t, err)
	t.Log(latency)

	latency, err = Latency("127.0.0.1")
	assert.Nil(t, err)
	t.Log(latency)

	latency, err = Latency("www.ymicj.com")
	assert.Nil(t, err)

	t.Log(latency)
}
