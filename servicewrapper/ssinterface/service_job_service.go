package ssinterface

import (
	"context"
	"time"

	"github.com/sgostarter/i/l"
)

type CycleJobService interface {
	DoJob(ctx context.Context, logger l.Wrapper) (time.Duration, error)
}
