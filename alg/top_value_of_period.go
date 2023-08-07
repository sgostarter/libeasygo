package alg

import (
	"time"

	"github.com/sgostarter/libeasygo/i"
)

type TopCheck[T any] interface {
	NewIsTop(oldV, newV T) bool
}

func NewTopValueOfPeriod[T any, L i.RWLock](lock L, period time.Duration,
	check TopCheck[T]) TopValueOfPeriod[T, L] {
	return TopValueOfPeriod[T, L]{
		lock:    lock,
		period:  period,
		check:   check,
		nowTops: make([]TWithTime[T], 0, 10),
	}
}

type TWithTime[T any] struct {
	At time.Time
	D  T
}

type TopValueOfPeriod[T any, L i.RWLock] struct {
	lock    L
	period  time.Duration
	check   TopCheck[T]
	nowTops []TWithTime[T]
}

func (a *TopValueOfPeriod[T, L]) expire() {
	a.lock.Lock()
	defer a.lock.Unlock()

	if len(a.nowTops) == 0 {
		return
	}

	var idx int

	for ; idx < len(a.nowTops); idx++ {
		if time.Since(a.nowTops[idx].At) <= a.period {
			break
		}
	}

	if idx > 0 {
		a.nowTops = append([]TWithTime[T]{}, a.nowTops[idx:]...)
	}
}

func (a *TopValueOfPeriod[T, L]) set(v T) {
	a.lock.Lock()
	defer a.lock.Unlock()

	var idx int

	for idx = len(a.nowTops) - 1; idx >= 0; idx-- {
		if !a.check.NewIsTop(a.nowTops[idx].D, v) {
			break
		}
	}

	a.nowTops = append(a.nowTops[:idx+1], TWithTime[T]{
		At: time.Now(),
		D:  v,
	})
}

func (a *TopValueOfPeriod[T, L]) get() (d T, exists bool) {
	a.lock.RLock()
	defer a.lock.RUnlock()

	if len(a.nowTops) == 0 {
		return
	}

	d = a.nowTops[0].D
	exists = true

	return
}

func (a *TopValueOfPeriod[T, L]) Set(v T) {
	a.expire()
	a.set(v)
}

func (a *TopValueOfPeriod[T, L]) Get() (T, bool) {
	a.expire()

	return a.get()
}
