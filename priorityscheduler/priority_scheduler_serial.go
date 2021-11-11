package priorityscheduler

import "context"

type Serial struct {
	schedulers []*PriorityScheduler
}

func NewPrioritySchedulerSerial() *Serial {
	return &Serial{}
}

func (ps *Serial) NewSerialPriorityScheduler() *PriorityScheduler {
	scheduler := New()
	ps.schedulers = append(ps.schedulers, scheduler)

	return scheduler
}

func (ps *Serial) Do(ctx context.Context) interface{} {
	for _, scheduler := range ps.schedulers {
		ret := scheduler.Do(ctx)
		if ret != nil {
			return ret
		}
	}

	return nil
}
