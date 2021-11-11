package cuserror

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var errI = NewWithCode(1)

func errIReturn() error {
	return errI
}

func TestError1(t *testing.T) {
	err1 := NewWithCode(1)
	assert.NotNil(t, As(err1))

	// nolint: goerr113
	err2 := errors.New("xx")
	assert.Nil(t, As(err2))

	assert.False(t, Is(err1, errI))
	assert.True(t, Is(errIReturn(), errI))

	err3 := NewWithError(2, errI)
	assert.True(t, Is(err3, errI))
}
