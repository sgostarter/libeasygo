package i

type Lock interface {
	Lock()
	Unlock()
}

type RWLock interface {
	Lock
	RLock()
	RUnlock()
}

type NopLock struct{}

func (NopLock) Lock() {}

func (NopLock) Unlock() {}

type NopRWLock struct{}

func (NopRWLock) Lock() {}

func (NopRWLock) Unlock() {}

func (NopRWLock) RLock() {}

func (NopRWLock) RUnlock() {}
