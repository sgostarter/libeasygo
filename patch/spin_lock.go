package patch

import (
	"runtime"
	"sync"
	"sync/atomic"
)

type SpinLock struct {
	_    sync.Mutex // for copy protection compiler warning
	lock uint32
}

func (sl *SpinLock) Lock() {
	counter := 1

	for !sl.TryLock() {
		for i := 0; i < counter; i++ {
			runtime.Gosched()
		}

		if counter < 256 { // Limit the maximum counter time
			counter *= 2
		}
	}
}

func (sl *SpinLock) Unlock() {
	atomic.StoreUint32(&sl.lock, 0)
}

func (sl *SpinLock) TryLock() bool {
	return atomic.CompareAndSwapUint32(&sl.lock, 0, 1)
}
