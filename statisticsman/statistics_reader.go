package statisticsman

import (
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sgostarter/libeasygo/statisticsman/counter"
	"github.com/sgostarter/libeasygo/statisticsman/impl"
	"github.com/sgostarter/libeasygo/statisticsman/inters"
)

type StatisticsScanResultCB func(timeSpanS string, k inters.DataKey, v int64, err error) error

type StatisticsReader interface {
	Scan4Current(dataKey inters.DataKey, cb StatisticsScanResultCB) error
	Scan4TimeSpanString(timeSpanS string, dataKey inters.DataKey, cb StatisticsScanResultCB) error
	FlushAndRemoveLastHourData(dataKey inters.DataKey, deepLevel int, cb StatisticsScanResultCB) error
}

func NewStatisticsReader(redisCli *redis.Client) StatisticsReader {
	return NewStatisticsReaderEx(impl.NewRedisDataProvider(redisCli), impl.NewHourTimeSpan())
}

func NewStatisticsReaderEx(dataProvider inters.DataProvider, timeSpan counter.TimeSpan) StatisticsReader {
	if dataProvider == nil || timeSpan == nil {
		return nil
	}

	return &statisticsReaderImpl{
		dataProvider: dataProvider,
		timeSpan:     timeSpan,
	}
}

type statisticsReaderImpl struct {
	dataProvider inters.DataProvider
	timeSpan     counter.TimeSpan
}

func (impl *statisticsReaderImpl) Scan4Current(dataKey inters.DataKey, cb StatisticsScanResultCB) error {
	return impl.Scan4TimeSpanString(impl.timeSpan.GetNowTimeString(), dataKey, cb)
}

func (impl *statisticsReaderImpl) Scan4TimeSpanString(timeSpanS string, dataKey inters.DataKey, cb StatisticsScanResultCB) error {
	return impl.dataProvider.Scan(timeSpanS, func(rKey, k string, v int64, err error) error {
		if err != nil {
			err = cb(rKey, nil, v, err)

			return err
		}

		err = dataKey.From(k)
		if err != nil {
			err = cb(rKey, nil, 0, err)
		} else {
			err = cb(rKey, dataKey, v, nil)
		}

		return err
	})
}

func (impl *statisticsReaderImpl) FlushAndRemoveLastHourData(dataKey inters.DataKey, deepLevel int, cb StatisticsScanResultCB) (err error) {
	if deepLevel <= 0 {
		deepLevel = 1
	}

	for loop := 1; loop <= deepLevel; loop++ {
		// nolint: durationcheck
		t := time.Now().Add(-time.Duration(loop) * impl.timeSpan.GetInterval())
		timeSpanS := impl.timeSpan.GetTimeStringFromTime(t)
		exists, _ := impl.dataProvider.Exists(timeSpanS)

		if !exists {
			continue
		}

		err = impl.Scan4TimeSpanString(timeSpanS, dataKey, cb)
		if err != nil {
			break
		}

		err = impl.dataProvider.Delete(timeSpanS)
		if err != nil {
			break
		}
	}

	return
}
