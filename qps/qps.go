package qps

import (
	"sync/atomic"
	"time"
)

type QPS struct {
	ds []*qpsData
}

func New(cnt int) *QPS {
	if cnt <= 0 {
		cnt = 2
	}

	ds := make([]*qpsData, cnt)
	for idx := 0; idx < len(ds); idx++ {
		ds[idx] = &qpsData{}
	}

	return &QPS{
		ds: ds,
	}
}

type qpsData struct {
	unix int64
	v    int64
}

func (qps *QPS) Inc() {
	timeNow := time.Now()
	d := qps.ds[timeNow.Unix()%int64(len(qps.ds))]

	if d.unix != timeNow.Unix() {
		d.unix = timeNow.Unix()
		atomic.StoreInt64(&d.v, 0)
	}

	atomic.AddInt64(&d.v, 1)
}

func (qps *QPS) Gets() []int64 {
	vs := make([]int64, 0, len(qps.ds)-1)
	curIndex := time.Now().Unix()

	for idx := curIndex - 1; idx > curIndex-int64(len(qps.ds)); idx-- {
		vs = append(vs, atomic.LoadInt64(&qps.ds[idx%int64(len(qps.ds))].v))
	}

	return vs
}
