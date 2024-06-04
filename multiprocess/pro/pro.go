package pro

import (
	"os/exec"
	"sync"

	"go.uber.org/atomic"
)

type Logger interface {
	OnStopped(tag interface{})
	OnRestarted(tag interface{})
	OnError(tag interface{}, op string, err error)

	WriteString(s string)
}

func NewWrapper(path string, args []string, dir string, tag interface{}, logger Logger) *Wrapper {
	if len(args) > 0 {
		if args[0] != path {
			args = append([]string{path}, args...)
		}
	}

	w := &Wrapper{
		path:   path,
		args:   args,
		dir:    dir,
		tag:    tag,
		logger: logger,
	}

	w.init()

	return w
}

type Wrapper struct {
	path   string
	args   []string
	dir    string
	tag    interface{}
	logger Logger

	wg         sync.WaitGroup
	stop       atomic.Bool
	cmdWrapper atomic.Value
}

type cmdHandle struct {
	cmd *exec.Cmd
}

func (w *Wrapper) init() {
	w.wg.Add(1)
	w.cmdWrapper.Store(cmdHandle{})

	go w.mainRoutine()
}

func (w *Wrapper) mainRoutine() {
	defer w.wg.Done()

	counter := 0

	for !w.stop.Load() {
		cmd := exec.Cmd{
			Path: w.path,
			Args: w.args,
			Dir:  w.dir,
		}

		counter++

		err := GroupCmdStart(&cmd, func(_ *exec.Cmd) {
			if counter > 1 {
				if w.logger != nil {
					w.logger.OnRestarted(w.tag)
				}
			}

			w.cmdWrapper.Store(cmdHandle{
				cmd: &cmd,
			})

			if w.stop.Load() {
				_ = cmd.Process.Kill()
			}

			_ = cmd.Wait()
		})

		if err != nil {
			if w.logger != nil {
				w.logger.OnError(w.tag, "start_process", err)
			}

			break
		}

		if w.logger != nil {
			w.logger.OnStopped(w.tag)
		}
	}
}

func (w *Wrapper) Running() bool {
	return !w.stop.Load()
}

func (w *Wrapper) StopAndWait() {
	w.stop.Store(true)

	i := w.cmdWrapper.Load()

	h, _ := i.(cmdHandle)
	if h.cmd != nil && h.cmd.Process != nil {
		_ = h.cmd.Process.Kill()
	}

	w.wg.Wait()
}
