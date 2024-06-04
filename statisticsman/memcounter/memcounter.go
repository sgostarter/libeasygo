package memcounter

import (
	"time"

	"github.com/sgostarter/libeasygo/timespan"
)

type tsData struct {
	v  int
	tm int64
}

type MemCounter struct {
	maxN     int
	duration time.Duration
	ds       []*tsData
	ts       *timespan.TimeSpan
}

func NewMemCounter(duration time.Duration, maxN int) *MemCounter {
	if duration <= 0 || maxN <= 0 {
		return nil
	}

	ts := timespan.NewTimeSpan(duration)

	ds := make([]*tsData, maxN)
	for idx := 0; idx < maxN; idx++ {
		ds[idx] = &tsData{}
	}

	return &MemCounter{
		maxN:     maxN,
		duration: duration,
		ds:       ds,
		ts:       ts,
	}
}

func (mc *MemCounter) reorder() {
	timeNow := time.Now()

	fnGetTimeLabel := func(t time.Time) int64 {
		n, _ := mc.ts.LabelString2Int(mc.ts.GetLabel(t))

		return n
	}

	label := fnGetTimeLabel(timeNow)

	if mc.ds[0].tm == label {
		return
	}

	idx := 1
	for ; idx < mc.maxN; idx++ {
		lt := timeNow
		if mc.ds[0].tm == fnGetTimeLabel(lt.Add(-time.Duration(idx)*mc.duration)) {
			break
		}
	}

	for n, m, f := mc.maxN-1, mc.maxN-idx-1, 0; f < mc.maxN-idx; n, m, f = n-1, m-1, f+1 {
		mc.ds[n].tm = mc.ds[m].tm
		mc.ds[n].v = mc.ds[m].v
	}

	for n := 0; n < idx; n++ {
		lt := timeNow
		mc.ds[n].tm = fnGetTimeLabel(lt.Add(-time.Duration(n) * mc.duration))
		mc.ds[n].v = 0
	}
}

func (mc *MemCounter) Inc() {
	mc.reorder()

	mc.ds[0].v++
}

func (mc *MemCounter) Count() int {
	v := 0
	for idx := 0; idx < len(mc.ds); idx++ {
		v += mc.ds[idx].v
	}

	return v
}
