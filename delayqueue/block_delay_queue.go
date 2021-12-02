package delayqueue

import (
	"context"
	"time"
)

func NewBlockDelayQueue(dq DelayQueue) BlockDelayQueue {
	if dq == nil || dq.GetReadyPool() == nil {
		return nil
	}

	rp := dq.GetReadyPool()
	rpFetcher, ok := rp.(ReadyPoolFetcher)

	if !ok || rpFetcher == nil {
		return nil
	}

	return &blockDelayQueueImpl{
		dq:        dq,
		rpFetcher: rpFetcher,
	}
}

type blockDelayQueueImpl struct {
	dq        DelayQueue
	rpFetcher ReadyPoolFetcher
}

func (impl *blockDelayQueueImpl) GetDelayQueue() DelayQueue {
	return impl.dq
}

func (impl *blockDelayQueueImpl) PushJob(job *Job) error {
	if job.TTR != 0 {
		return ErrSafeJob
	}

	return impl.dq.JobPush(job)
}

func (impl *blockDelayQueueImpl) PushSafeJob(job *Job) error {
	if job.TTR <= 0 {
		return ErrNoSafeJob
	}

	return impl.dq.JobPush(job)
}

func (impl *blockDelayQueueImpl) BlockProcessJobOnce(f FNProcessJob, timeout time.Duration, jobIn *Job, topics ...string) (ok bool, err error) {
	job, err := impl.jobBPopEx(timeout, jobIn, topics...)
	if err != nil {
		return
	}

	if job == nil {
		return
	}

	ok = true
	jobID := job.ID
	newJob, err := f(job)

	if err != nil {
		return
	}

	impl.dq.JobDone(jobID)

	if newJob != nil {
		err = impl.dq.JobPush(newJob)
	}

	return
}

func (impl *blockDelayQueueImpl) jobBPopEx(timeout time.Duration, jobIn *Job, topics ...string) (job *Job, err error) {
	b := time.Now()
	tm := timeout

	for {
		job, err = impl.bPopEx(tm, jobIn, topics...)
		if err != nil {
			return
		}

		if job != nil {
			return
		}

		if timeout <= 0 {
			return
		}

		tm = timeout - time.Since(b)
		if tm <= 0 {
			return
		}
	}
}

func (impl *blockDelayQueueImpl) bPopEx(timeout time.Duration, jobIn *Job, topics ...string) (job *Job, err error) {
	start := time.Now()

	jid, err := impl.rpFetcher.GetReadyJob(timeout, topics...)
	if err != nil {
		return
	}

	if jid == nil {
		return
	}

	to := timeout - time.Since(start)
	if to <= 0 {
		err = ErrTimeout

		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), to)
	defer cancel()

	job, err = impl.dq.GetJobPool().GetJob(ctx, jid.ID, jobIn)
	if err != nil {
		return
	}

	if job == nil {
		impl.dq.JobDone(jid.ID)
	}

	return
}
