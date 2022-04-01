package decodebuffer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeBuffer(t *testing.T) {
	buf := NewBuffer(nil)

	rule := NewRuleTerminator([]byte("1"))
	buf.SetDecodeRule(rule)

	buf.Append([]byte("a1bc2def3def3"))
	d, ok := buf.FindTerminator()
	assert.True(t, ok)
	assert.Equal(t, string(d), "a")

	rule.SetTerminator([]byte("2"))

	d, ok = buf.FindTerminator()
	assert.True(t, ok)
	assert.Equal(t, string(d), "bc")

	rule.SetTerminator([]byte("3"))

	d, ok = buf.FindTerminator()
	assert.True(t, ok)
	assert.Equal(t, string(d), "def")

	buf.Append([]byte("xxx3yy"))

	d, ok = buf.FindTerminator()
	assert.True(t, ok)
	assert.Equal(t, string(d), "def")

	d, ok = buf.FindTerminator()
	assert.True(t, ok)
	assert.Equal(t, string(d), "xxx")

	buf.Append([]byte("zz3"))

	d, ok = buf.FindTerminator()
	assert.True(t, ok)
	assert.Equal(t, string(d), "yyzz")
}
