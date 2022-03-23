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

	Scan4CurrentEx(dataKey inters.DataKey, cb StatisticsScanResultCB, reset bool) error
	Scan4TimeSpanStringEx(timeSpanS string, dataKey inters.DataKey, cb StatisticsScanResultCB, reset bool) error
	FlushAndRemoveLastHourDataEx(dataKey inters.DataKey, deepLevel int, cb StatisticsScanResultCB, reset bool) error
}

func NewStatisticsReader(redisCli *redis.Client) StatisticsReader {
	return NewStatisticsReaderEx(impl.NewRedisDataProvider(redisCli), impl.NewHourTimeSpan(), "")
}

func NewStatisticsReaderEx(dataProvider inters.DataProvider, timeSpan counter.TimeSpan, tsPre string) StatisticsReader {
	if dataProvider == nil || timeSpan == nil {
		return nil
	}

	return &statisticsReaderImpl{
		dataProvider: dataProvider,
		timeSpan:     timeSpan,
		tsPre:        tsPre,
	}
}

type statisticsReaderImpl struct {
	dataProvider inters.DataProvider
	timeSpan     counter.TimeSpan
	tsPre        string
}

func (impl *statisticsReaderImpl) Scan4Current(dataKey inters.DataKey, cb StatisticsScanResultCB) error {
	return impl.Scan4CurrentEx(dataKey, cb, false)
}

func (impl *statisticsReaderImpl) Scan4CurrentEx(dataKey inters.DataKey, cb StatisticsScanResultCB, reset bool) error {
	return impl.Scan4TimeSpanStringEx(impl.timeSpan.GetNowTimeString(), dataKey, cb, reset)
}

func (impl *statisticsReaderImpl) Scan4TimeSpanString(timeSpanS string, dataKey inters.DataKey, cb StatisticsScanResultCB) error {
	return impl.Scan4TimeSpanStringEx(timeSpanS, dataKey, cb, false)
}

func (impl *statisticsReaderImpl) Scan4TimeSpanStringEx(timeSpanS string, dataKey inters.DataKey, cb StatisticsScanResultCB, reset bool) error {
	return impl.dataProvider.ScanEx(impl.tsPre+timeSpanS, func(rKey, k string, v int64, err error) error {
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
	}, reset)
}

func (impl *statisticsReaderImpl) FlushAndRemoveLastHourDataEx(dataKey inters.DataKey, deepLevel int, cb StatisticsScanResultCB, reset bool) error {
	return impl.flushAndRemoveLastHourData(dataKey, deepLevel, cb, reset)
}

func (impl *statisticsReaderImpl) FlushAndRemoveLastHourData(dataKey inters.DataKey, deepLevel int, cb StatisticsScanResultCB) (err error) {
	return impl.FlushAndRemoveLastHourDataEx(dataKey, deepLevel, cb, false)
}

func (impl *statisticsReaderImpl) flushAndRemoveLastHourData(dataKey inters.DataKey, deepLevel int, cb StatisticsScanResultCB, reset bool) (err error) {
	if deepLevel <= 0 {
		deepLevel = 1
	}

	for loop := 1; loop <= deepLevel; loop++ {
		t := time.Now().Add(-time.Duration(loop) * impl.timeSpan.GetInterval())
		timeSpanS := impl.timeSpan.GetTimeStringFromTime(t)
		exists, _ := impl.dataProvider.Exists(timeSpanS)

		if !exists {
			continue
		}

		err = impl.Scan4TimeSpanStringEx(timeSpanS, dataKey, cb, reset)
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
