//go:build linux

package pro

import (
	"os/exec"
	"syscall"
)

func GroupCmdStart(cmd *exec.Cmd, fn func(cmd *exec.Cmd)) (err error) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Pdeathsig: syscall.SIGTERM,
	}

	err = cmd.Start()
	if err != nil {
		return
	}

	fn(cmd)

	return
}
