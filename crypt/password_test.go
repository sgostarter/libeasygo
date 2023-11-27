package crypt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	hashedPassword, err := HashPassword("pwd1", 4096)
	assert.Nil(t, err)

	assert.True(t, CheckHashedPassword("pwd1", hashedPassword, 4096))
	assert.False(t, CheckHashedPassword("pwd2", hashedPassword, 4096))
}
