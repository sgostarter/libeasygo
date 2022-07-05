package ssinterface

import (
	"context"

	"github.com/sgostarter/i/l"
)

type ServiceStub interface {
	Run(ctx context.Context, logger l.Wrapper)
}
