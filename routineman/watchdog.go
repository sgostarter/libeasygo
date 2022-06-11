package routineman

import (
	"sync"
	"time"

	"github.com/godruoyi/go-snowflake"
)

type Ob func(runnerTag interface{}, id uint64)

type WatchDog interface {
	StopAsync(rm RoutineMan)
	StopAndWait(rm RoutineMan)

	Run(runner func(), runnerTag interface{}, timeoutCheckedOb, timeoutRunnerFinishedOb Ob)
	RunWithTimeoutCheck(runner func(), runnerTag interface{}, to time.Duration, timeoutCheckedOb, timeoutRunnerFinishedOb Ob)
}

func NewWatchDog() WatchDog {
	return &watchDogImpl{}
}

var (
	_initOnce sync.Once
	_defPool  = &watchDogImpl{}
)

func GetDefaultWatchDog() WatchDog {
	return _defPool
}

func DefaultWatchDogInitOnce(disableCheck bool, timeoutCheckedOb, timeoutRunnerFinishedOb Ob) {
	_initOnce.Do(func() {
		_defPool.disableCheck = disableCheck
		_defPool.timeoutCheckedOb = timeoutCheckedOb
		_defPool.timeoutRunnerFinishedOb = timeoutRunnerFinishedOb
	})
}

type watchDogImpl struct {
	disableCheck            bool
	timeoutCheckedOb        Ob
	timeoutRunnerFinishedOb Ob
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

func (impl *watchDogImpl) Run(runner func(), runnerTag interface{}, timeoutCheckedOb, timeoutRunnerFinishedOb Ob) {
	impl.RunWithTimeoutCheck(runner, runnerTag, 0, timeoutCheckedOb, timeoutRunnerFinishedOb)
}

func (impl *watchDogImpl) RunWithTimeoutCheck(runner func(), runnerTag interface{}, to time.Duration, timeoutCheckedOb, timeoutRunnerFinishedOb Ob) {
	if runner == nil {
		return
	}

	if impl.disableCheck {
		runner()

		return
	}

	if to <= 0 {
		to = time.Second
	}

	if timeoutCheckedOb == nil {
		timeoutCheckedOb = impl.timeoutCheckedOb
	}

	if timeoutRunnerFinishedOb == nil {
		timeoutRunnerFinishedOb = impl.timeoutRunnerFinishedOb
	}

	ch := make(chan interface{}, 1)

	go func() {
		id := snowflake.ID()

		select {
		case <-ch:
			return
		case <-time.After(to):
			if timeoutCheckedOb != nil {
				timeoutCheckedOb(runnerTag, id)
			}
		}

		<-ch

		if timeoutRunnerFinishedOb != nil {
			timeoutRunnerFinishedOb(runnerTag, id)
		}
	}()

	runner()

	ch <- 1
}
