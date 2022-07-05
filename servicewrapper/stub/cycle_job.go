package stub

import (
	"context"
	"time"

	"github.com/sgostarter/i/l"
	"github.com/sgostarter/libeasygo/servicewrapper/ssinterface"
)

type cycleJobService struct {
	serviceImpl ssinterface.CycleJobService
}

func NewCycleJobServiceStub(serviceImpl ssinterface.CycleJobService) ssinterface.ServiceStub {
	if serviceImpl == nil {
		return nil
	}

	return &cycleJobService{
		serviceImpl: serviceImpl,
	}
}

func (ss *cycleJobService) Run(ctx context.Context, logger l.Wrapper) {
	if logger == nil {
		logger = l.NewNopLoggerWrapper()
	}

	loop := true

	logger.Debug("enter cycle job loop")

	duration, err := ss.serviceImpl.DoJob(ctx, logger)
	if err != nil {
		logger.Errorf("do job failed: %v", err)

		return
	}

	for loop {
		select {
		case <-ctx.Done():
			loop = false

			logger.Debug("check ctx done, try exit loop")
		case <-time.After(duration):
			duration, err = ss.serviceImpl.DoJob(ctx, logger)
			if err != nil {
				logger.Errorf("do job failed: %v", err)

				loop = false

				break
			}
		}
	}

	logger.Debug("leave cycle job loop")
}
