package priorityscheduler

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type TestPriorityRunner2 struct {
	delay time.Duration
}

func (r *TestPriorityRunner2) Run(ctx context.Context) (interface{}, error) {
	select {
	case <-time.After(r.delay):
		return nil, nil
	case <-ctx.Done():
		return nil, ErrCancelled
	}
}

func TestNewPrioritySchedulerSerial(t *testing.T) {
	schedulerSerial := NewPrioritySchedulerSerial()
	resp := schedulerSerial.Do(context.Background())
	assert.Equal(t, resp, nil)

	scheduler := schedulerSerial.NewSerialPriorityScheduler()
	scheduler.AddRunner(1, &TestPriorityRunner{
		delay: 2 * time.Second,
		idx:   "1-1",
	})
	scheduler.AddRunner(2, &TestPriorityRunner{
		delay: 1 * time.Second,
		idx:   "2-1",
	})

	resp = schedulerSerial.Do(context.Background())
	assert.Equal(t, resp.(string), "1-1")

	resp = schedulerSerial.Do(context.Background())
	assert.Equal(t, resp.(string), "1-1")

	scheduler = schedulerSerial.NewSerialPriorityScheduler()
	scheduler.AddRunner(1, &TestPriorityRunner{
		delay: 2 * time.Millisecond,
		idx:   "1-2",
	})

	resp = schedulerSerial.Do(context.Background())
	assert.Equal(t, resp.(string), "1-1")
}

func TestNewPrioritySchedulerSerial2(t *testing.T) {
	schedulerSerial := NewPrioritySchedulerSerial()
	resp := schedulerSerial.Do(context.Background())
	assert.Equal(t, resp, nil)

	scheduler := schedulerSerial.NewSerialPriorityScheduler()
	scheduler.AddRunner(1, &TestPriorityRunner2{
		delay: 2 * time.Second,
	})

	scheduler = schedulerSerial.NewSerialPriorityScheduler()
	scheduler.AddRunner(1, &TestPriorityRunner{
		delay: 2 * time.Millisecond,
		idx:   "1-2",
	})

	timeStart := time.Now()
	resp = schedulerSerial.Do(context.Background())

	assert.True(t, time.Since(timeStart) > 2*time.Second)
	assert.Equal(t, resp.(string), "1-2")
}

func TestNewPrioritySchedulerSerial3(t *testing.T) {
	schedulerSerial := NewPrioritySchedulerSerial()
	resp := schedulerSerial.Do(context.Background())
	assert.Equal(t, resp, nil)

	scheduler := schedulerSerial.NewSerialPriorityScheduler()
	scheduler.AddRunner(1, &TestPriorityRunner2{
		delay: 10 * time.Second,
	})

	scheduler = schedulerSerial.NewSerialPriorityScheduler()
	scheduler.AddRunner(1, &TestPriorityRunner{
		delay: 1 * time.Millisecond,
		idx:   "1-2",
	})

	timeStart := time.Now()
	ctx, ctxCancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer ctxCancel()

	resp = schedulerSerial.Do(ctx)

	assert.True(t, time.Since(timeStart) > 3*time.Second)

	assert.Equal(t, resp, nil)
}
