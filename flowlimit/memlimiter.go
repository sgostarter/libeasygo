package flowlimit

import (
	"context"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

func NewMemLimiter(r rate.Limit, b int) Limiter {
	ctx, cancel := context.WithCancel(context.Background())

	impl := &limitImpl{
		r:         r,
		b:         b,
		ctx:       ctx,
		ctxCancel: cancel,
		limiters:  make(map[string]*limitInfo),
	}

	impl.init()

	return impl
}

type limitImpl struct {
	r rate.Limit
	b int

	wg        sync.WaitGroup
	ctx       context.Context
	ctxCancel context.CancelFunc

	limiters map[string]*limitInfo
	lock     sync.Mutex
}

func (impl *limitImpl) Wait() {
	impl.wg.Wait()
}

func (impl *limitImpl) Close() {
	impl.ctxCancel()
}

type limitInfo struct {
	limiter    *rate.Limiter
	lastAccess time.Time
}

func (impl *limitImpl) init() {
	impl.wg.Add(1)

	go func() {
		defer impl.wg.Done()

		loop := true

		for loop {
			select {
			case <-impl.ctx.Done():
				loop = false

				continue
			case <-time.After(time.Minute * 10):
				impl.clean()
			}
		}
	}()
}

func (impl *limitImpl) clean() {
	impl.lock.Lock()
	defer impl.lock.Unlock()

	for k, li := range impl.limiters {
		if time.Since(li.lastAccess) > time.Hour {
			delete(impl.limiters, k)
		}
	}
}

func (impl *limitImpl) getLimiterForKey(key string) *rate.Limiter {
	impl.lock.Lock()
	defer impl.lock.Unlock()

	li := impl.limiters[key]
	if li == nil {
		li = &limitInfo{
			limiter:    rate.NewLimiter(impl.r, impl.b),
			lastAccess: time.Now(),
		}

		impl.limiters[key] = li
	}

	li.lastAccess = time.Now()

	return li.limiter
}

func (impl *limitImpl) Allow(key string) bool {
	return impl.getLimiterForKey(key).Allow()
}
