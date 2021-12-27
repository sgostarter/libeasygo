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
}

func NewStatisticsWriter(redisCli *redis.Client, logger l.Wrapper) StatisticsWriter {
	return NewStatisticsWriterEx(context.Background(), counter.NewTimeSpanCounters(impl.NewHourTimeSpan()),
		impl.NewRedisCounterStorage(redisCli, logger), logger)
}

func NewStatisticsWriterEx(ctx context.Context, tsCounters *counter.TimeSpanCounters, storage inters.Storage, logger l.Wrapper) StatisticsWriter {
	return &statisticsWriterImpl{
		asyncStore: counter.NewAsyncStore(ctx, storage, tsCounters, logger),
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
