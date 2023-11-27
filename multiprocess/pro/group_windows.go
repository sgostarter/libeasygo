//go:build windows

package pro

import (
	"os"
	"os/exec"
	"unsafe"

	"golang.org/x/sys/windows"
)

// We use this struct to retreive process handle(which is unexported)
// from os.Process using unsafe operation.
type winProcess struct {
	Pid    int
	Handle uintptr
}

type ProcessExitGroup windows.Handle

func NewProcessExitGroup() (ProcessExitGroup, error) {
	handle, err := windows.CreateJobObject(nil, nil)
	if err != nil {
		return 0, err
	}

	info := windows.JOBOBJECT_EXTENDED_LIMIT_INFORMATION{
		BasicLimitInformation: windows.JOBOBJECT_BASIC_LIMIT_INFORMATION{
			LimitFlags: windows.JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE,
		},
	}
	if _, err := windows.SetInformationJobObject(
		handle,
		windows.JobObjectExtendedLimitInformation,
		uintptr(unsafe.Pointer(&info)),
		uint32(unsafe.Sizeof(info))); err != nil {
		return 0, err
	}

	return ProcessExitGroup(handle), nil
}

func (g ProcessExitGroup) Dispose() error {
	return windows.CloseHandle(windows.Handle(g))
}

func (g ProcessExitGroup) AddProcess(p *os.Process) error {
	return windows.AssignProcessToJobObject(
		windows.Handle(g),
		windows.Handle((*winProcess)(unsafe.Pointer(p)).Handle))
}

func GroupCmdStart(cmd *exec.Cmd, fn func(cmd *exec.Cmd)) (err error) {
	g, err := NewProcessExitGroup()
	if err != nil {
		return
	}

	defer g.Dispose()

	err = cmd.Start()
	if err != nil {
		return
	}

	if err = g.AddProcess(cmd.Process); err != nil {
		return
	}

	fn(cmd)

	return
}
