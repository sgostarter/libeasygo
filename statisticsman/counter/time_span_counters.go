package counter

import (
	"sync"
	"time"
)

type TimeSpan interface {
	GetNowTimeString() string
	GetTimeStringFromTime(t time.Time) string
	GetTimeFromTimeString(s string) (time.Time, error)
	GetInterval() time.Duration
}

type TimeSpanCounters struct {
	lock sync.RWMutex

	ts       TimeSpan
	counters map[string]*KeyCounters
}

func NewTimeSpanCounters(ts TimeSpan) *TimeSpanCounters {
	return &TimeSpanCounters{
		ts:       ts,
		counters: make(map[string]*KeyCounters),
	}
}

func (tsc *TimeSpanCounters) GetTimeSpan() TimeSpan {
	return tsc.ts
}

func (tsc *TimeSpanCounters) get(ts string) *KeyCounters {
	tsc.lock.RLock()
	defer tsc.lock.RUnlock()

	if c, ok := tsc.counters[ts]; ok {
		return c
	}

	return nil
}

func (tsc *TimeSpanCounters) GetByTimeSpanS(ts string) *KeyCounters {
	c := tsc.get(ts)
	if c != nil {
		return c
	}

	tsc.lock.Lock()
	defer tsc.lock.Unlock()

	var ok bool
	if c, ok = tsc.counters[ts]; ok {
		return c
	}

	c = NewKeyCounters()
	tsc.counters[ts] = c

	return c
}

func (tsc *TimeSpanCounters) Get() *KeyCounters {
	return tsc.GetByTimeSpanS(tsc.ts.GetNowTimeString())
}

func (tsc *TimeSpanCounters) GetCounters() map[string]*KeyCounters {
	ret := make(map[string]*KeyCounters)

	tsc.lock.RLock()
	defer tsc.lock.RUnlock()

	for s, counters := range tsc.counters {
		ret[s] = counters
	}

	return ret
}

func (tsc *TimeSpanCounters) Remove(timeSs []string) {
	tsc.lock.Lock()
	defer tsc.lock.Unlock()

	for _, s := range timeSs {
		delete(tsc.counters, s)
	}
}

func (tsc *TimeSpanCounters) CanSafeRemove(ts string) bool {
	t, _ := tsc.ts.GetTimeFromTimeString(ts)

	return time.Since(t) > tsc.ts.GetInterval()
}
