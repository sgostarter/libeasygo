package debug

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetHostInfo(t *testing.T) {
	hi, err := GetHostInfo()
	assert.Nil(t, err)
	t.Log(hi)
}

func TestGetLoadAvg(t *testing.T) {
	la, err := GetLoadAvg()
	assert.Nil(t, err)
	t.Log(la.Load1, la.Load5, la.Load15)
}

func TestGetCPUCount(t *testing.T) {
	pCores, lCores, err := GetCPUCount()
	assert.Nil(t, err)
	t.Log(pCores, lCores)
}

func TestGetCPUPercent(t *testing.T) {
	for idx := 0; idx < 5; idx++ {
		t.Log(GetCPUPercent(time.Second))
	}
}

func TestGetVirtualMemInfo(t *testing.T) {
	vii, err := GetVirtualMemInfo()
	assert.Nil(t, err)
	t.Log(vii)
}
