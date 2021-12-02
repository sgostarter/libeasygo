package helper

import (
	"context"
	"time"
)

func RunWithTimeoutEx(ctx context.Context, timeOut time.Duration, fn func(ctx context.Context) error) error {
	ctx, cancel := context.WithTimeout(ctx, timeOut)
	defer cancel()

	return fn(ctx)
}

func RunWithTimeout4Redis(ctx context.Context, fn func(ctx context.Context) error) error {
	return RunWithTimeoutEx(ctx, 2*time.Second, fn)
}
