package servicewrapper

import (
	"context"

	"github.com/sgostarter/i/commerr"
	"github.com/sgostarter/i/l"
	"github.com/sgostarter/libeasygo/routineman"
	"github.com/sgostarter/libeasygo/servicewrapper/ssinterface"
	"go.uber.org/atomic"
)

type ServiceWrapper struct {
	routineMan routineman.RoutineMan
	logger     l.Wrapper

	startFlag atomic.Bool
}

func NewServiceWrapper(ctx context.Context, logger l.Wrapper) *ServiceWrapper {
	if logger == nil {
		logger = l.NewNopLoggerWrapper()
	}

	return &ServiceWrapper{
		routineMan: routineman.NewRoutineMan(ctx, logger),
		logger:     logger.WithFields(l.StringField(l.ClsKey, "ServiceWrapper")),
	}
}

func (sw *ServiceWrapper) Start(serviceImpl ssinterface.ServiceStub) error {
	if serviceImpl == nil {
		return commerr.ErrInvalidArgument
	}

	if !sw.startFlag.CompareAndSwap(false, true) {
		sw.logger.Fatal("initAgain")

		return commerr.ErrAlreadyExists
	}

	sw.routineMan.StartRoutine(func(ctx context.Context, _ func() bool) {
		serviceImpl.Run(ctx, sw.logger)
	}, "")

	return nil
}

func (sw *ServiceWrapper) Stop() {
	sw.routineMan.TriggerStop()
}

func (sw *ServiceWrapper) Wait() {
	sw.routineMan.Wait()
}
