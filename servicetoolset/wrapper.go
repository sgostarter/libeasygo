package servicetoolset

import (
	"sync/atomic"

	"github.com/kardianos/service"
)

type Runner interface {
	Run(logger service.Logger)
	StopAndWait()
}

type OnceWrapper struct {
	chRunExited chan interface{}
	startCalled atomic.Bool
	runner      Runner
}

func NewOnceWrapper(runner Runner) *OnceWrapper {
	if runner == nil {
		return nil
	}

	return &OnceWrapper{
		chRunExited: make(chan interface{}, 10),
		runner:      runner,
	}
}

func (w *OnceWrapper) Start(logger service.Logger) {
	if !w.startCalled.CompareAndSwap(false, true) {
		return
	}

	go func() {
		w.runner.Run(logger)

		select {
		case w.chRunExited <- true:
		default:
		}
	}()
}

func (w *OnceWrapper) ExitRunnerAndWait() {
	if !w.startCalled.Load() {
		return
	}

	w.runner.StopAndWait()
}

func (w *OnceWrapper) Wait() {
	<-w.chRunExited
}
