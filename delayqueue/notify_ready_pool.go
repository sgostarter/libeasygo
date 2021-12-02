package delayqueue

func NewNotifyReadyPool(maxCap int) ReadyPool {
	if maxCap <= 0 {
		maxCap = 100
	}

	return &notifyReadyPool{
		jobChan: make(chan *JobIdentify, maxCap),
	}
}

type notifyReadyPool struct {
	jobChan chan *JobIdentify
}

func (impl *notifyReadyPool) NewReadyJob(topic, jobID string) (err error) {
	impl.jobChan <- &JobIdentify{
		Topic: topic,
		ID:    jobID,
	}

	return
}

func (impl *notifyReadyPool) JobChan() <-chan *JobIdentify {
	return impl.jobChan
}
