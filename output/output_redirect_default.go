//go:build !windows

package output

import (
	"os"
	"syscall"
)

func Redirect(file *os.File) {
	_ = syscall.Dup2(int(file.Fd()), int(os.Stdout.Fd()))
	_ = syscall.Dup2(int(file.Fd()), int(os.Stderr.Fd()))
}
