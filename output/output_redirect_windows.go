//go:build windows

package output

import (
	"os"

	"golang.org/x/sys/windows"
)

func Redirect(file *os.File) {
	_ = windows.SetStdHandle(windows.STD_OUTPUT_HANDLE, windows.Handle(file.Fd()))
	os.Stdout = file

	_ = windows.SetStdHandle(windows.STD_ERROR_HANDLE, windows.Handle(file.Fd()))

	os.Stderr = file

	return
}
