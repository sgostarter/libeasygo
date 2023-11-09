package netutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAvailableTCPPort(t *testing.T) {
	port, err := GetAvailableTCPPort()
	assert.Nil(t, err)
	assert.True(t, port > 0)
	t.Log("GetAvailableTCPPort:", port)
}
