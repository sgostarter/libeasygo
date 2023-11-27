//go:build !windows && !linux

package pro

import (
	"os/exec"
)

func GroupCmdStart(cmd *exec.Cmd, fn func(cmd *exec.Cmd)) (err error) {
	err = cmd.Start()
	if err != nil {
		return
	}

	fn(cmd)

	return
}
