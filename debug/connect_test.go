package debug

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnectTCP(t *testing.T) {
	assert.Nil(t, ConnectTCP("www.ymipro.com:80"))
	assert.Nil(t, ConnectTLS("www.ymipro.com:443", false))
}
