package ssinterface

import (
	"context"
	"time"

	"github.com/sgostarter/i/l"
)

type CycleJobService interface {
	OnStart(logger l.Wrapper)
	DoJob(ctx context.Context, logger l.Wrapper) (time.Duration, error)
	OnFinish(logger l.Wrapper)
}
