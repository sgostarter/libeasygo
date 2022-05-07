package routineman

import (
	"context"
	"sync"

	"github.com/sgostarter/i/l"
	"go.uber.org/atomic"
)

func NewRoutineMan(ctx context.Context, logger l.Wrapper) RoutineMan {
	if logger == nil {
		logger = l.NewNopLoggerWrapper()
	}

	ctx, cancel := context.WithCancel(ctx)

	return &routineManImpl{
		ctx:       ctx,
		ctxCancel: cancel,
		logger:    logger,
	}
}

type routineManImpl struct {
	wg        sync.WaitGroup
	ctx       context.Context
	ctxCancel context.CancelFunc
	logger    l.Wrapper

	exiting atomic.Bool
}

func (impl *routineManImpl) Context() context.Context {
	return impl.ctx
}

func (impl *routineManImpl) Exiting() bool {
	return impl.exiting.Load()
}

func (impl *routineManImpl) StartRoutine(routine func(ctx context.Context), _ string) {
	impl.wg.Add(1)

	go func() {
		defer impl.wg.Done()

		routine(impl.ctx)
	}()
}

func (impl *routineManImpl) Wait() {
	impl.wg.Wait()
}

func (impl *routineManImpl) StopAndWait() {
	impl.exiting.Store(true)
	impl.ctxCancel()

	impl.wg.Wait()
}