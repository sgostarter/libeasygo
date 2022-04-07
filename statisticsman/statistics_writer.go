package statisticsman

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/sgostarter/i/l"
	"github.com/sgostarter/libeasygo/statisticsman/counter"
	"github.com/sgostarter/libeasygo/statisticsman/impl"
	"github.com/sgostarter/libeasygo/statisticsman/inters"
)

type StatisticsWriter interface {
	Inc(k inters.DataKey)
	Add(k inters.DataKey, v int64)
	AddSync(k inters.DataKey, v int64)
}

func NewStatisticsWriter(redisCli *redis.Client, logger l.Wrapper) StatisticsWriter {
	tsCounter := counter.NewTimeSpanCounters(impl.NewHourTimeSpan())

	return NewStatisticsWriterEx(tsCounter, counter.NewAsyncStore(context.Background(), impl.NewRedisCounterStorage(redisCli, logger),
		tsCounter, logger))
}

func NewStatisticsWriterEx(tsCounters *counter.TimeSpanCounters, asyncStore *counter.AsyncStore) StatisticsWriter {
	return &statisticsWriterImpl{
		asyncStore: asyncStore,
		tsCounters: tsCounters,
	}
}

type statisticsWriterImpl struct {
	asyncStore *counter.AsyncStore
	tsCounters *counter.TimeSpanCounters
}

func (impl *statisticsWriterImpl) Inc(k inters.DataKey) {
	impl.tsCounters.Get().GetCounter(k.Key()).Inc()
}

func (impl *statisticsWriterImpl) Add(k inters.DataKey, v int64) {
	impl.tsCounters.Get().GetCounter(k.Key()).Add(v)
}

func (impl *statisticsWriterImpl) AddSync(k inters.DataKey, v int64) {
	impl.asyncStore.Add(impl.tsCounters.GetTimeSpan().GetNowTimeString(), k, v)
}
