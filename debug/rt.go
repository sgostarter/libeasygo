package debug

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type RTProbe interface {
	Leave()
}

type RT interface {
	Enter() RTProbe
	GetAVG() time.Duration
	DestroyAndWait()
}

func NewRT(ctx context.Context, internal time.Duration) RT {
	ctx, cancel := context.WithCancel(ctx)

	var avg atomic.Value

	avg.Store(float64(0))

	rt := &rtImpl{
		ctx:       ctx,
		ctxCancel: cancel,
		chV:       make(chan *rtProbeData, 1000),
		avg:       avg,
		internal:  internal,
	}

	rt.wg.Add(1)

	go rt.chRoutine()

	return rt
}

type rtProbeData struct {
	start  time.Time
	finish time.Time
}

type rtImpl struct {
	wg        sync.WaitGroup
	ctx       context.Context
	ctxCancel context.CancelFunc

	chV      chan *rtProbeData
	avg      atomic.Value
	internal time.Duration

	quit bool
}

func (impl *rtImpl) chRoutine() {
	defer impl.wg.Done()

	loop := true

	bufferSize := 1000
	buffer := make([]*rtProbeData, bufferSize)
	bufferIdx := bufferSize * 10

	var count, rtSum int64

	fnTrim := func() {
		if impl.internal <= 0 {
			return
		}

		for start := (bufferIdx - 1) % bufferSize; start != bufferIdx%bufferSize; start-- {
			if buffer[start%bufferSize] == nil {
				break
			}

			if buffer[start%bufferSize].finish.After(time.Now().Add(-impl.internal)) {
				break
			}

			count--

			rtSum -= buffer[start%bufferSize].finish.UnixNano() - buffer[start%bufferSize].start.UnixNano()
		}
	}

	for loop {
		select {
		case <-impl.ctx.Done():
			loop = false

			continue
		case <-time.After(time.Second):
			fnTrim()
		case v := <-impl.chV:
			idx := bufferIdx % bufferSize
			if buffer[idx] != nil {
				count--

				rtSum -= buffer[idx].finish.UnixNano() - buffer[idx].start.UnixNano()
			}

			buffer[idx] = v
			count++

			rtSum += buffer[idx].finish.UnixNano() - buffer[idx].start.UnixNano()

			impl.avg.Store(float64(rtSum/1e6) / float64(count))

			bufferIdx++
			if bufferIdx > 0xFFFFFFF {
				bufferIdx = bufferSize * 10
			}

			fnTrim()
		}
	}
}

func (impl *rtImpl) Enter() RTProbe {
	if impl.quit {
		return &rtProbeNull{}
	}

	return newRtProbeImpl(impl.chV)
}

func (impl *rtImpl) GetAVG() time.Duration {
	if v, ok := impl.avg.Load().(float64); ok {
		return time.Duration(v) * time.Millisecond
	}

	return 0
}

func (impl *rtImpl) DestroyAndWait() {
	impl.quit = true
	impl.ctxCancel()
	impl.wg.Wait()
}

func newRtProbeImpl(ch chan *rtProbeData) RTProbe {
	return &rtProbeImpl{
		start: time.Now(),
		ch:    ch,
	}
}

type rtProbeImpl struct {
	start time.Time
	ch    chan *rtProbeData
}

func (impl *rtProbeImpl) Leave() {
	impl.ch <- &rtProbeData{
		start:  impl.start,
		finish: time.Now(),
	}
}

type rtProbeNull struct {
}

func (impl *rtProbeNull) Leave() {

}
