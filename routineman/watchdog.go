package routineman

import (
	"time"

	"github.com/godruoyi/go-snowflake"
)

type TimeoutObserver interface {
	OnTimeoutRunnerChecked(runnerTag interface{}, id uint64)
	OnTimeoutRunnerFinished(runnerTag interface{}, id uint64)
}

type WatchDog interface {
	StopAsync(rm RoutineMan)
	StopAndWait(rm RoutineMan)

	RunWithDefaultCheck(runner func(), runnerTag interface{}, ob TimeoutObserver)
	RunWithTimeoutCheck(runner func(), runnerTag interface{}, to time.Duration, ob TimeoutObserver)
}

func NewWatchDog() WatchDog {
	return &watchDogImpl{}
}

var _defPool = &watchDogImpl{}

func GetDefaultWatchDog() WatchDog {
	return _defPool
}

func WatchDogReInit(disableCheck bool) {
	_defPool.disableCheck = disableCheck
}

type watchDogImpl struct {
	disableCheck bool
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

func (impl *watchDogImpl) RunWithDefaultCheck(runner func(), runnerTag interface{}, ob TimeoutObserver) {
	impl.RunWithTimeoutCheck(runner, runnerTag, 0, ob)
}

func (impl *watchDogImpl) RunWithTimeoutCheck(runner func(), runnerTag interface{}, to time.Duration, ob TimeoutObserver) {
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

	ch := make(chan interface{}, 1)

	go func() {
		id := snowflake.ID()

		select {
		case <-ch:
			return
		case <-time.After(to):
			if ob != nil {
				ob.OnTimeoutRunnerChecked(runnerTag, id)
			}
		}

		<-ch

		if ob != nil {
			ob.OnTimeoutRunnerFinished(runnerTag, id)
		}
	}()

	runner()

	ch <- 1
}
