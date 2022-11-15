package priorityscheduler

import (
	"context"
	"testing"
	"time"

	"github.com/sgostarter/i/commerr"
	"github.com/stretchr/testify/assert"
)

type TestPriorityRunner struct {
	delay time.Duration
	idx   string
}

func (r *TestPriorityRunner) Run(ctx context.Context) (interface{}, error) {
	select {
	case <-time.After(r.delay):
		return r.idx, nil
	case <-ctx.Done():
		return nil, commerr.ErrCanceled
	}
}

func TestPS(t *testing.T) {
	//
	scheduler := New()
	scheduler.AddRunner(1, &TestPriorityRunner{
		delay: 2 * time.Second,
		idx:   "1-1",
	})
	scheduler.AddRunner(2, &TestPriorityRunner{
		delay: 1 * time.Second,
		idx:   "2-1",
	})

	resp := scheduler.Do(context.Background())
	assert.EqualValues(t, resp, "1-1")

	//
	scheduler = New()
	scheduler.AddRunner(1, &TestPriorityRunner{
		delay: 2 * time.Second,
		idx:   "1-1",
	})
	scheduler.AddRunner(1, &TestPriorityRunner{
		delay: 500 * time.Millisecond,
		idx:   "1-2",
	})
	scheduler.AddRunner(2, &TestPriorityRunner{
		delay: 1 * time.Second,
		idx:   "2-1",
	})

	resp = scheduler.Do(context.Background())
	assert.EqualValues(t, resp, "1-2")
	//
	scheduler = New()
	scheduler.AddRunner(1, &TestPriorityRunner{
		delay: 2 * time.Second,
		idx:   "1-1",
	})
	scheduler.AddRunner(1, &TestPriorityRunner{
		delay: 500 * time.Millisecond,
		idx:   "1-2",
	})
	scheduler.AddRunner(2, &TestPriorityRunner{
		delay: 1 * time.Second,
		idx:   "2-1",
	})
	scheduler.AddRunner(20, &TestPriorityRunner{
		delay: 10 * time.Millisecond,
		idx:   "20-1",
	})

	resp = scheduler.Do(context.Background())
	assert.EqualValues(t, resp, "1-2")

	//
	scheduler.AddRunner(0, &TestPriorityRunner{
		delay: 3 * time.Second,
		idx:   "0-1",
	})
	scheduler.AddRunner(0, &TestPriorityRunner{
		delay: 1 * time.Second,
		idx:   "0-2",
	})

	resp = scheduler.Do(context.Background())
	assert.EqualValues(t, resp, "0-2")
}
