package counter

import "sync"

type KeyCounters struct {
	lock sync.RWMutex

	counters map[string]*Counter
}

func NewKeyCounters() *KeyCounters {
	return &KeyCounters{
		counters: make(map[string]*Counter),
	}
}

func (kc *KeyCounters) getCounter(key string) *Counter {
	kc.lock.RLock()
	defer kc.lock.RUnlock()

	if c, ok := kc.counters[key]; ok {
		return c
	}

	return nil
}

func (kc *KeyCounters) GetCounter(key string) *Counter {
	c := kc.getCounter(key)
	if c != nil {
		return c
	}

	kc.lock.Lock()
	defer kc.lock.Unlock()

	var ok bool
	if c, ok = kc.counters[key]; ok {
		return c
	}

	c = &Counter{}
	kc.counters[key] = c

	return c
}

func (kc *KeyCounters) GetCounters() map[string]*Counter {
	ret := make(map[string]*Counter)

	kc.lock.RLock()
	defer kc.lock.RUnlock()

	for s, counters := range kc.counters {
		ret[s] = counters
	}

	return ret
}
