package routineman

type WatchDog interface {
	StopAsync(rm RoutineMan)
	StopAndWait(rm RoutineMan)
}

func NewWatchDog() WatchDog {
	return &watchDogImpl{}
}

var _defPool = NewWatchDog()

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
