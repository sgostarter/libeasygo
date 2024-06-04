package counter

import (
	"context"
	"sync"
	"time"

	"github.com/sgostarter/i/l"
	"github.com/sgostarter/libeasygo/statisticsman/inters"
)

type AsyncStore struct {
	ctx       context.Context
	ctxCancel context.CancelFunc
	wg        sync.WaitGroup

	stg          inters.Storage
	counters     *TimeSpanCounters
	logger       l.Wrapper
	syncDuration time.Duration
	tsPre        string
}

func NewAsyncStore(ctx context.Context, stg inters.Storage, counters *TimeSpanCounters, logger l.Wrapper) *AsyncStore {
	return NewAsyncStoreEx(ctx, stg, counters, logger, 0, "")
}

func NewAsyncStoreEx(ctx context.Context, stg inters.Storage, counters *TimeSpanCounters, logger l.Wrapper,
	syncDuration time.Duration, tsPre string) *AsyncStore {
	ctx, cancel := context.WithCancel(ctx)

	if logger == nil {
		logger = l.NewNopLoggerWrapper()
	}

	if syncDuration <= 0 {
		syncDuration = time.Second
	}

	logger = logger.WithFields(l.StringField("cls", "counter_store"))

	s := &AsyncStore{
		ctx:          ctx,
		ctxCancel:    cancel,
		stg:          stg,
		counters:     counters,
		logger:       logger,
		syncDuration: syncDuration,
		tsPre:        tsPre,
	}

	s.wg.Add(1)

	go s.startRoutine()

	return s
}

func (s *AsyncStore) Wait() {
	s.wg.Wait()
}

func (s *AsyncStore) startRoutine() {
	defer s.wg.Done()

	s.logger.Info("enter AsyncStore::startRoutine")
	defer s.logger.Info("leave AsyncStore::startRoutine")

	storeCounter := 0
	loop := true

	for loop {
		select {
		case <-s.ctx.Done():
			t := time.Now()

			s.store()

			d := time.Since(t)
			s.logger.WithFields(l.AnyField("usedTimeMS", d/time.Millisecond)).
				Info("event_store_dump_final")

			loop = false

			continue
		case <-time.After(s.syncDuration):
			t := time.Now()

			s.store()

			d := time.Since(t)

			if storeCounter%600 == 0 || d > 5*s.syncDuration {
				s.logger.WithFields(l.AnyField("usedTimeMS", d/time.Millisecond)).Info("event_store_dump")
			}

			storeCounter++

			if storeCounter >= 0x00FFFFFF {
				storeCounter = 0
			}
		}
	}
}

func (s *AsyncStore) store() {
	removedTimeS := make([]string, 0)

	for timeS, counters := range s.counters.GetCounters() {
		for kv, c := range counters.GetCounters() {
			cnt := c.HC()
			if cnt > 0 {
				s.stg.Inc(s.tsPre+timeS, kv, cnt)
			}
		}

		if s.counters.CanSafeRemove(timeS) {
			removedTimeS = append(removedTimeS, timeS)
		}
	}

	if len(removedTimeS) > 0 {
		s.counters.Remove(removedTimeS)
	}
}

func (s *AsyncStore) Add(timeS string, k inters.DataKey, v int64) {
	s.stg.Inc(s.tsPre+timeS, k.Key(), v)
}
