package helper

import (
	"context"
	"time"
)

func DoWithTimeout(ctx context.Context, timeout time.Duration, cb func(ctx context.Context)) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cb(ctx)
}
