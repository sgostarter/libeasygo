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

// nolint
func TestDecodeBuffer2(t *testing.T) {
	buf := NewBuffer(nil)

	rule := NewRuleTerminator([]byte("1"))
	buf.SetDecodeRule(rule)

	buf.Append([]byte("a1b1c2def3def3"))
	d, ok := buf.FindTerminator()
	assert.True(t, ok)
	assert.EqualValues(t, "a", string(d))

	d = append(d, []byte("xxxxxxxxxxxxxxxxxxxxxxxxxxxxx")...)
	e, ok := buf.FindTerminator()
	assert.True(t, ok)
	assert.EqualValues(t, "b", string(e))

}

func TestDecodeBuffer3(t *testing.T) {
	buf := NewBufferEx(nil, []byte("abc"))

	rule := NewRuleTerminator([]byte("1"))
	buf.SetDecodeRule(rule)

	_, ok := buf.PeekTerminator()
	assert.False(t, ok)

	buf.Append([]byte("1d"))

	d, ok := buf.PeekTerminator()
	assert.True(t, ok)
	assert.EqualValues(t, "abc", string(d))

	d = buf.GetBuffer()
	assert.EqualValues(t, "abc1d", string(d))
}
