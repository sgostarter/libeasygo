package debug

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectTCP(t *testing.T) {
	assert.Nil(t, TestTCPConnect("www.ymipro.com:80", false))
	assert.Nil(t, TestTCPConnect("www.ymipro.com:443", true))
	assert.NotNil(t, TestTCPConnect("abc.com:80", true))
}
