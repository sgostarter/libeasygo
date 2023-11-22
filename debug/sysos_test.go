package debug

import (
	"testing"
	"time"

	"github.com/sgostarter/libeasygo/helper"
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
	t.Logf("Total: %s, Available:%s, Used: %s, UsedPercent:%f",
		helper.FormatVToString(float64(vii.Total), 1024, "B"),
		helper.FormatVToString(float64(vii.Available), 1024, "B"),
		helper.FormatVToString(float64(vii.Used), 1024, "B"),
		vii.UsedPercent)
}

func TestGetRootDisk(t *testing.T) {
	i, err := GetRootDisk()
	assert.Nil(t, err)
	t.Logf("Total: %s, Free:%s, UsedPercent:%f",
		helper.FormatVToString(float64(i.Total), 1024, "B"),
		helper.FormatVToString(float64(i.Free), 1024, "B"),
		i.UsedPercent)
}

// nolint
func TestGetNetIOInfo(t *testing.T) {
	i1, err := GetNetIOInfo()
	assert.Nil(t, err)

	time.Sleep(time.Second * 3)

	i2, err := GetNetIOInfo()
	assert.Nil(t, err)

	t.Logf("send bytes: %d, recv bytes: %d",
		i2.BytesSent-i1.BytesSent, i2.BytesRecv-i1.BytesRecv)
}
