package mwf

type NoLock struct {
}

func (lock NoLock) RLock() {

}

func (lock NoLock) RUnlock() {

}

func (lock NoLock) Lock() {

}

func (lock NoLock) Unlock() {

}
