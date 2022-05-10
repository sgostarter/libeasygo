package routineman

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/godruoyi/go-snowflake"
	"github.com/sgostarter/i/l"
	"github.com/spf13/cast"
	"go.uber.org/atomic"
)

func NewRoutineManWithTimeoutCheck(ctx context.Context, name string, timeout time.Duration, logger l.Wrapper) DebugRoutineMan {
	if timeout <= 0 {
		timeout = time.Second
	}

	if logger == nil {
		logger = l.NewNopLoggerWrapper()
	}

	ctx, cancel := context.WithCancel(ctx)

	impl := &routineManWithTimeoutCheckImpl{
		ctx:       ctx,
		ctxCancel: cancel,
		name:      name,
		timeout:   timeout,
		logger:    logger,
	}

	impl.ob.Store(obWrapper{})

	return impl
}

type obWrapper struct {
	ob DebugRoutineManTimeoutObserver
}

type routineManWithTimeoutCheckImpl struct {
	wg        sync.WaitGroup
	ctx       context.Context
	ctxCancel context.CancelFunc
	name      string
	timeout   time.Duration
	logger    l.Wrapper

	exiting  atomic.Bool
	routines sync.Map

	ob atomic.Value
}

func (impl *routineManWithTimeoutCheckImpl) Context() context.Context {
	return impl.ctx
}

func (impl *routineManWithTimeoutCheckImpl) Exiting() bool {
	return impl.exiting.Load()
}

func (impl *routineManWithTimeoutCheckImpl) routineName(name string) string {
	return fmt.Sprintf("%s-%d", name, snowflake.ID())
}

func (impl *routineManWithTimeoutCheckImpl) StartRoutine(routine func(ctx context.Context, exiting func() bool), name string) {
	impl.wg.Add(1)

	name = impl.routineName(name)
	impl.routines.Store(name, time.Now())

	go func() {
		defer func() {
			impl.wg.Done()
			impl.routines.Delete(name)
		}()

		routine(impl.ctx, func() bool {
			return impl.exiting.Load()
		})
	}()
}

func (impl *routineManWithTimeoutCheckImpl) Wait() {
	impl.wg.Wait()
}

func (impl *routineManWithTimeoutCheckImpl) StopAndWait() {
	impl.exiting.Store(true)
	impl.ctxCancel()

	ch := make(chan interface{}, 2)

	go func() {
		select {
		case <-time.After(impl.timeout):
			msg := impl.dumpRoutineMsg()
			impl.logger.Warn(msg)

			if obW, ok := impl.ob.Load().(obWrapper); ok && obW.ob != nil {
				obW.ob(msg)
			}
		case <-ch:
		}
	}()

	impl.wg.Wait()

	ch <- true
}

func (impl *routineManWithTimeoutCheckImpl) dumpRoutineMsg() string {
	ss := strings.Builder{}
	ss.WriteString("!!ROUTINE TERMINATE TIMEOUT CHECKED\n")
	ss.WriteString(fmt.Sprintf("function %s\n", impl.name))
	impl.routines.Range(func(key, value interface{}) bool {
		ss.WriteString(fmt.Sprintf(" %s: %s\n", cast.ToString(key), cast.ToTime(value).String()))

		return true
	})
	ss.WriteString("\n")

	return ss.String()
}

func (impl *routineManWithTimeoutCheckImpl) TriggerStop() {
	impl.exiting.Store(true)
	impl.ctxCancel()
}

func (impl *routineManWithTimeoutCheckImpl) SetExitTimeoutObserver(ob DebugRoutineManTimeoutObserver) {
	impl.ob.Store(obWrapper{ob: ob})
}

func (impl *routineManWithTimeoutCheckImpl) DoWithTimeout(label string, do func(), d time.Duration) {
	if do == nil {
		return
	}

	if d <= 0 {
		do()

		return
	}

	ch := make(chan interface{}, 2)

	go func() {
		do()

		ch <- 1
	}()

	select {
	case <-ch:
		return
	case <-time.After(d):
		if obW, ok := impl.ob.Load().(obWrapper); ok && obW.ob != nil {
			var ss strings.Builder

			ss.WriteString("!!DO WITH TIMEOUT CHECKED\n")
			ss.WriteString(fmt.Sprintf("operation label %s\n", label))
			obW.ob(ss.String())
		}
	}
}
