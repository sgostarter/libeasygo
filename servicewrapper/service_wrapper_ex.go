package servicewrapper

import (
	"context"

	"github.com/sgostarter/i/l"
	"github.com/sgostarter/libeasygo/servicewrapper/ssinterface"
	"github.com/sgostarter/libeasygo/servicewrapper/stub"
)

type CycleServiceWrapper struct {
	*ServiceWrapper
}

func NewCycleServiceWrapper(ctx context.Context, logger l.Wrapper) *CycleServiceWrapper {
	return &CycleServiceWrapper{
		ServiceWrapper: NewServiceWrapper(ctx, logger),
	}
}

func (sw *CycleServiceWrapper) Start(serviceImpl ssinterface.CycleJobService) error {
	return sw.ServiceWrapper.Start(stub.NewCycleJobServiceStub(serviceImpl))
}
