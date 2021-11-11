package priorityscheduler

import (
	"context"
	"sort"
	"sync"
)

type PriorityRunner interface {
	Run(ctx context.Context) (interface{}, error)
}

type PriorityScheduler struct {
	runnerGroupsMap map[uint]*runnerGroupContext
}

func New() *PriorityScheduler {
	return &PriorityScheduler{
		runnerGroupsMap: make(map[uint]*runnerGroupContext),
	}
}

// AddRunner priority 数值越低，优先级越高.
func (ps *PriorityScheduler) AddRunner(priority uint, runner PriorityRunner) {
	if _, ok := ps.runnerGroupsMap[priority]; !ok {
		ps.runnerGroupsMap[priority] = &runnerGroupContext{
			priority: priority,
		}
	}

	ps.runnerGroupsMap[priority].runner = append(ps.runnerGroupsMap[priority].runner, runner)
}

func (ps *PriorityScheduler) Do(ctx context.Context) interface{} {
	rgcs := ps.doEx(ctx)
	for _, rgc := range rgcs {
		if len(rgc.resp) > 0 {
			return rgc.resp[0]
		}
	}

	return nil
}

func (ps *PriorityScheduler) doEx(ctx context.Context) runnerGroupContexts {
	if len(ps.runnerGroupsMap) == 0 {
		return nil
	}

	rgcs := runnerGroupContexts{}

	for _, rgc := range ps.runnerGroupsMap {
		rgc.ctx = nil
		rgc.ctxCancels = nil
		rgc.resp = nil
		rgcs = append(rgcs, rgc)
	}

	sort.Sort(rgcs)

	wg := sync.WaitGroup{}

	var lastCtxCancels []context.CancelFunc

	for idx := rgcs.Len() - 1; idx >= 0; idx-- {
		ctxP, cancel := context.WithCancel(ctx)
		rgcs[idx].ctx = ctxP
		rgcs[idx].ctxCancels = append(rgcs[idx].ctxCancels, cancel)

		if len(lastCtxCancels) > 0 {
			rgcs[idx].ctxCancels = append(rgcs[idx].ctxCancels, lastCtxCancels...)
		}

		lastCtxCancels = rgcs[idx].ctxCancels
	}

	for _, rgc := range rgcs {
		for _, runner := range rgc.runner {
			wg.Add(1)

			go func(rgc *runnerGroupContext, runner PriorityRunner) {
				defer wg.Done()

				resp, err := runner.Run(rgc.ctx)
				if err != nil {
					// TODO LOG
					return
				}

				if resp == nil {
					// TODO LOG
					return
				}

				rgc.Lock()
				rgc.resp = append(rgc.resp, resp)
				rgc.Unlock()

				rgc.Cancel()
			}(rgc, runner)
		}
	}

	wg.Wait()

	return rgcs
}

//
//
//
type runnerGroupContext struct {
	sync.Mutex

	priority uint
	runner   []PriorityRunner

	ctx        context.Context
	ctxCancels []context.CancelFunc

	resp []interface{}
}

func (rgc *runnerGroupContext) Cancel() {
	for _, cancel := range rgc.ctxCancels {
		cancel()
	}
}

type runnerGroupContexts []*runnerGroupContext

func (rgc runnerGroupContexts) Len() int {
	return len(rgc)
}

func (rgc runnerGroupContexts) Less(i, j int) bool {
	return rgc[i].priority < rgc[j].priority
}

func (rgc runnerGroupContexts) Swap(i, j int) {
	rgc[i], rgc[j] = rgc[j], rgc[i]
}
