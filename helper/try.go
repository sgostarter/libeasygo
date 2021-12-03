package helper

import (
	"context"
	"time"

	"github.com/sgostarter/libeasygo/cuserror"
)

const (
	TimeoutInfinity = time.Duration(-1)
)

// TryWithTimeout timeout=0 代表不重试，timeout<0代表无限重试
func TryWithTimeout(timeout time.Duration, try func(timeout time.Duration) bool) (ok bool, err error) {
	if try == nil {
		err = cuserror.NewWithErrorMsg("noTryFunction")

		return
	}

	start := time.Now()
	to := timeout

	for {
		ok = try(to)
		if ok {
			break
		}

		if timeout <= 0 {
			break
		}

		to = timeout - time.Since(start)
		if to <= 0 {
			break
		}
	}

	return
}

func TryWithTimeoutContext(ctx context.Context, try func(ctx context.Context) bool) (ok bool, err error) {
	if try == nil {
		err = cuserror.NewWithErrorMsg("noTryFunction")

		return
	}

	loop := true
	for loop {
		ok = try(ctx)
		if ok {
			break
		}

		select {
		case <-ctx.Done():
			loop = false

			continue
		default:
		}
	}

	return
}
