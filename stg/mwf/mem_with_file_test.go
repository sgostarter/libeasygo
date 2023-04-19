package mwf_test

import (
	"testing"

	"github.com/sgostarter/libeasygo/stg/mwf"
	"github.com/stretchr/testify/assert"
)

func TestMemAndFile1(t *testing.T) {
	var mm *mwf.MemWithFile[map[int]string, mwf.Serial, mwf.Lock]
	mm = mwf.NewMemWithFile[map[int]string, mwf.Serial, mwf.Lock](make(map[int]string), &mwf.JSONSerial{}, &mwf.NoLock{}, "utStorage.txt")
	t.Log(mm)

	memWithFile := mwf.NewMemWithFile(make(map[int]string), &mwf.JSONSerial{}, &mwf.NoLock{}, "utStorage.txt")
	assert.NotNil(t, memWithFile)

	_ = memWithFile.Change(func(m map[int]string) (map[int]string, error) {
		if m == nil {
			m = make(map[int]string)
		}

		m[1] = "1xx"

		return m, nil
	})

	memWithFile.Read(func(m map[int]string) {
		t.Log(m[1])
	})
}
