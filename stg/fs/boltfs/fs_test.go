package boltfs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoltFS(t *testing.T) {
	fs, err := NewFileStorage("./tmp/aa/a.db")
	assert.Nil(t, err)

	err = fs.WriteFile("abcd/aa", []byte("1234"))
	assert.Nil(t, err)

	err = fs.WriteFile("abcd/aa", []byte("abc"))
	assert.Nil(t, err)

	d, err := fs.ReadFile("abcd/aa")
	assert.Nil(t, err)

	assert.EqualValues(t, "abc", d)
}
