package impl

import (
	"time"

	"github.com/sgostarter/libeasygo/statisticsman/counter"
	"github.com/sgostarter/libeasygo/timespan"
)

func NewTimeSpan(duration time.Duration) counter.TimeSpan {
	ts := timespan.NewTimeSpan(duration)

	return &timeSpanImpl{
		duration: duration,
		ts:       ts,
	}
}

type timeSpanImpl struct {
	duration time.Duration
	ts       *timespan.TimeSpan
}

func (impl *timeSpanImpl) GetNowTimeString() string {
	return impl.ts.GetCurrentLabel()
}

func (impl *timeSpanImpl) GetTimeStringFromTime(t time.Time) string {
	return impl.ts.GetLabel(t)
}

func (impl *timeSpanImpl) GetTimeFromTimeString(s string) (time.Time, error) {
	return impl.ts.Label2Time(s)
}

func (impl *timeSpanImpl) GetInterval() time.Duration {
	return impl.duration
}
