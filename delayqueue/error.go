package delayqueue

import (
	"github.com/sgostarter/libeasygo/cuserror"
)

var (
	ErrSafeJob   = cuserror.NewWithErrorMsg("isSafeJob")
	ErrNoSafeJob = cuserror.NewWithErrorMsg("isNoSafeJob")
	ErrTimeout   = cuserror.NewWithErrorMsg("timeout")
)
