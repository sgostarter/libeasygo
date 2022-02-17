package dcfg

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/sgostarter/libeasygo/cuserror"
)

type ConfigFetcher interface {
	GetCfg() interface{}
}

type DCfg interface {
	GetCfg() interface{}

	Destroy()
}

func NewDCfg(ctx context.Context, fetcher ConfigFetcher, flushInterval time.Duration) (DCfg, error) {
	if fetcher == nil {
		return nil, cuserror.NewWithErrorMsg("no fetcher")
	}

	if flushInterval <= 0 {
		return nil, cuserror.NewWithErrorMsg("invalid flush interval")
	}

	ctx, cancel := context.WithCancel(ctx)

	dCfg := &dCfgImp{
		ctx:           ctx,
		ctxCancel:     cancel,
		fetcher:       fetcher,
		flushInterval: flushInterval,
	}

	if err := dCfg.init(); err != nil {
		return nil, err
	}

	return dCfg, nil
}

type dCfgImp struct {
	wg        sync.WaitGroup
	ctx       context.Context
	ctxCancel context.CancelFunc

	fetcher       ConfigFetcher
	flushInterval time.Duration

	cfg atomic.Value
}

func (impl *dCfgImp) init() error {
	cfg := impl.fetcher.GetCfg()
	if cfg == nil {
		return cuserror.NewWithErrorMsg("no config")
	}

	impl.cfg.Store(cfg)

	impl.wg.Add(1)

	go impl.flushRoutine()

	return nil
}

func (impl *dCfgImp) doFlush() error {
	cfg := impl.fetcher.GetCfg()
	if cfg == nil {
		return cuserror.NewWithErrorMsg("no config")
	}

	impl.cfg.Store(cfg)

	return nil
}

func (impl *dCfgImp) flushRoutine() {
	defer impl.wg.Done()

	loop := true

	for loop {
		select {
		case <-impl.ctx.Done():
			loop = false

			break
		case <-time.After(impl.flushInterval):
			_ = impl.doFlush()
		}
	}
}

func (impl *dCfgImp) GetCfg() interface{} {
	return impl.cfg.Load()
}

func (impl *dCfgImp) Destroy() {
	impl.ctxCancel()
	impl.wg.Wait()
}
