package counter

import "sync/atomic"

type Counter struct {
	cnt int64
}

func (c *Counter) Inc() int64 {
	return c.Add(1)
}

func (c *Counter) Add(v int64) int64 {
	return atomic.AddInt64(&c.cnt, v)
}

func (c *Counter) HC() int64 {
	return atomic.SwapInt64(&c.cnt, 0)
}
