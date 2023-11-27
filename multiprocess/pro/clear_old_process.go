package pro

import (
	"os"
	"path/filepath"

	"github.com/shirou/gopsutil/v3/process"
)

func ClearOldProcesses() (err error) {
	ex, err := os.Executable()
	if err != nil {
		return
	}

	fileName := filepath.Base(ex)

	processes, _ := process.Processes()
	for _, p := range processes {
		n, _ := p.Name()
		if n == fileName && int(p.Pid) != os.Getpid() {
			_ = p.Kill()
		}
	}

	return
}
