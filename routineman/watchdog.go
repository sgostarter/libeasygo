package routineman

import (
	"context"
	"time"
)

type WatchDog interface {
	StopAsync(rm RoutineMan)
	StopAndWait(rm RoutineMan)

	Run(label string, runner func())
	RunWthCustomTimeout(label string, runner func(), to time.Duration)
}

func NewWatchDog() WatchDog {
	return &watchDogImpl{}
}

var _defPool = NewWatchDog()
var _defRoutineMan = NewRoutineManWithTimeoutCheck(context.Background(), "__def", time.Second*30, nil)

func GetDefaultWatchDog() WatchDog {
	return _defPool
}

type watchDogImpl struct {
}

func (impl *watchDogImpl) StopAsync(rm RoutineMan) {
	if rm == nil {
		return
	}

	go func() {
		rm.StopAndWait()
	}()
}

func (impl *watchDogImpl) StopAndWait(rm RoutineMan) {
	if rm == nil {
		return
	}

	rm.StopAndWait()
}

func (impl *watchDogImpl) Run(label string, runner func()) {
	_defRoutineMan.Run(label, runner)
}

func (impl *watchDogImpl) RunWthCustomTimeout(label string, runner func(), to time.Duration) {
	_defRoutineMan.RunWthCustomTimeout(label, runner, to)
}
