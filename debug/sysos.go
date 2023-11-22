package debug

import (
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
)

type HostInfo struct {
	Hostname        string
	HostID          string
	UptimeAt        time.Time
	BootTime        time.Duration
	OS              string
	Platform        string
	PlatformVersion string
	KernelVersion   string
}

func GetHostInfo() (hostInfo *HostInfo, err error) {
	hi, err := host.Info()
	if err != nil {
		return
	}

	hostInfo = &HostInfo{
		Hostname:        hi.Hostname,
		HostID:          hi.HostID,
		UptimeAt:        time.Unix(int64(hi.BootTime), 0),
		BootTime:        time.Duration(hi.Uptime) * time.Second,
		OS:              hi.OS,
		Platform:        hi.Platform,
		PlatformVersion: hi.PlatformVersion,
		KernelVersion:   hi.KernelVersion,
	}

	return
}

type LoadAvg struct {
	Load1  float64
	Load5  float64
	Load15 float64
}

func GetLoadAvg() (loadAvg *LoadAvg, err error) {
	la, err := load.Avg()
	if err != nil {
		return
	}

	loadAvg = &LoadAvg{
		Load1:  la.Load1,
		Load5:  la.Load5,
		Load15: la.Load15,
	}

	return
}

func GetCPUCount() (physicalCount, logicalCount int, err error) {
	physicalCount, err = cpu.Counts(false)
	if err != nil {
		return
	}

	logicalCount, err = cpu.Counts(true)

	return
}

func GetCPUPercent(interval time.Duration) (v float64, err error) {
	vs, err := cpu.Percent(interval, false)
	if err != nil {
		return
	}

	if len(vs) > 0 {
		v = vs[0]
	}

	return
}

type VirtualMemoryInfo struct {
	Total       uint64  `json:"total"`
	Available   uint64  `json:"available"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	UsedPercent float64 `json:"usedPercent"`
}

func GetVirtualMemInfo() (virtualMemoryInfo *VirtualMemoryInfo, err error) {
	vii, err := mem.VirtualMemory()
	if err != nil {
		return
	}

	virtualMemoryInfo = &VirtualMemoryInfo{
		Total:       vii.Total,
		Available:   vii.Available,
		Used:        vii.Used,
		Free:        vii.Free,
		UsedPercent: vii.UsedPercent,
	}

	return
}

type DiskInfo struct {
	Total       uint64  `json:"total"`
	Free        uint64  `json:"free"`
	UsedPercent float64 `json:"usedPercent"`
}

func GetRootDisk() (info *DiskInfo, err error) {
	rootPath := "/"

	if runtime.GOOS == "windows" {
		rootPath = "c:/"
	}

	di, err := disk.Usage(rootPath)
	if err != nil {
		return
	}

	info = &DiskInfo{
		Total:       di.Total,
		Free:        di.Free,
		UsedPercent: di.UsedPercent,
	}

	return
}

type NetIOInfo struct {
	BytesSent uint64
	BytesRecv uint64
}

func GetNetIOInfo() (info *NetIOInfo, err error) {
	countersStats, err := net.IOCounters(false)
	if err != nil {
		return
	}

	if len(countersStats) == 0 {
		return
	}

	info = &NetIOInfo{
		BytesSent: countersStats[0].BytesSent,
		BytesRecv: countersStats[0].BytesRecv,
	}

	return
}
