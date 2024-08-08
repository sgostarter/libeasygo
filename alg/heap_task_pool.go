package alg

import (
	"container/heap"
	"sync"
	"time"

	"github.com/sgostarter/i/commerr"
)

// TaskFunc .
type TaskFunc func(key string, args ...any)

type taskItem struct {
	at       time.Time
	exec     TaskFunc
	params   []any
	canceled bool
	key      string
}

type taskHeap []*taskItem

func (th *taskHeap) Len() int {
	return len(*th)
}

func (th *taskHeap) Less(i, j int) bool {
	return (*th)[i].at.Unix() < (*th)[j].at.Unix()
}
func (th *taskHeap) Swap(i, j int) {
	(*th)[i], (*th)[j] = (*th)[j], (*th)[i]
}

func (th *taskHeap) Push(x interface{}) {
	*th = append(*th, x.(*taskItem)) // nolint: forcetypeassert
}

func (th *taskHeap) Pop() interface{} {
	old := *th
	n := len(old)
	x := old[n-1]
	*th = old[0 : n-1]

	return x
}

type taskOpInfo struct {
	opType int
	key    string
	at     time.Time
	exec   TaskFunc
	params []any
}

const (
	taskOpAdd = iota
	taskOpDel
	taskOpUpdate
)

// HeapTaskPool .
type HeapTaskPool struct {
	wg     sync.WaitGroup
	closed chan bool
	taskOp chan *taskOpInfo
	keys   map[string]*taskItem

	tasks *taskHeap
}

// NewHeapTaskPool .
func NewHeapTaskPool() *HeapTaskPool {
	pool := &HeapTaskPool{
		closed: make(chan bool),
		taskOp: make(chan *taskOpInfo, 100),
		keys:   make(map[string]*taskItem),
		tasks:  &taskHeap{},
	}

	pool.init()

	return pool
}

func (tp *HeapTaskPool) init() {
	heap.Init(tp.tasks)

	tp.wg.Add(1)
	go tp.loop()
}

func (tp *HeapTaskPool) Stop() {
	close(tp.closed)
	tp.wg.Wait()
}

func (tp *HeapTaskPool) Wait() {
	tp.wg.Wait()
}

func (tp *HeapTaskPool) process() time.Duration {
	for {
		if tp.tasks.Len() <= 0 {
			return 24 * time.Hour
		}

		timeNow := time.Now()

		task := (*tp.tasks)[0]

		if task.canceled {
			heap.Pop(tp.tasks)

			continue
		}

		if task.at.After(timeNow) {
			return task.at.Sub(timeNow)
		}

		go task.exec(task.key, task.params...)

		heap.Pop(tp.tasks)

		delete(tp.keys, task.key)
	}
}

func (tp *HeapTaskPool) addOrUpdateTask(t time.Time, key string, exec TaskFunc, params []any) {
	tp.cancelTask(key)

	ti := &taskItem{
		at:     t,
		exec:   exec,
		params: params,
		key:    key,
	}
	tp.keys[key] = ti
	heap.Push(tp.tasks, ti)
}

func (tp *HeapTaskPool) cancelTask(key string) {
	if task, ok := tp.keys[key]; ok {
		task.canceled = true
	}
}

func (tp *HeapTaskPool) loop() {
	defer tp.wg.Done()

	nextTaskInterval := time.Second

	for {
		select {
		case <-tp.closed:
			return
		case taskI := <-tp.taskOp:
			if taskI.opType == taskOpAdd || taskI.opType == taskOpUpdate {
				tp.addOrUpdateTask(taskI.at, taskI.key, taskI.exec, taskI.params)
			} else if taskI.opType == taskOpDel {
				tp.cancelTask(taskI.key)
			}

			nextTaskInterval = tp.process()
		case <-time.After(nextTaskInterval):
			nextTaskInterval = tp.process()
		}
	}
}

func (tp *HeapTaskPool) AddTask(key string, t time.Time, exec TaskFunc, params ...any) (err error) {
	if key == "" || exec == nil {
		err = commerr.ErrInvalidArgument

		return
	}

	select {
	case tp.taskOp <- &taskOpInfo{
		opType: taskOpAdd,
		key:    key,
		at:     t,
		exec:   exec,
		params: params,
	}:
	default:
		err = commerr.ErrTimeout
	}

	return
}

func (tp *HeapTaskPool) RemoveTask(key string) (err error) {
	if key == "" {
		err = commerr.ErrInvalidArgument

		return
	}

	select {
	case tp.taskOp <- &taskOpInfo{
		opType: taskOpDel,
		key:    key,
	}:
	default:
		err = commerr.ErrTimeout
	}

	return
}
