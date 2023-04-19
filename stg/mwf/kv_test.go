package mwf

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type utKVItem struct {
	N int
	S string
}

func TestSimpleKV(t *testing.T) {
	_ = os.Remove("ut.txt")

	kv := NewKV("ut.txt")

	var item utKVItem

	ok, err := kv.Get("key", &item)
	assert.Nil(t, err)
	assert.False(t, ok)

	err = kv.Set("key", &utKVItem{
		N: 10,
		S: "S20S",
	})
	assert.Nil(t, err)

	ok, err = kv.Get("key", &item)
	assert.Nil(t, err)
	assert.True(t, ok)
	assert.EqualValues(t, 10, item.N)
	assert.EqualValues(t, "S20S", item.S)

	kv2 := NewKV("ut.txt")

	ok, err = kv2.Get("key", &item)
	assert.Nil(t, err)
	assert.True(t, ok)
	assert.EqualValues(t, 10, item.N)
	assert.EqualValues(t, "S20S", item.S)
}
