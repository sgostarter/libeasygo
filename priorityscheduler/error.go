package priorityscheduler

import "github.com/sgostarter/libeasygo/cuserror"

var (
	ErrCancelled = cuserror.NewWithErrorMsg("cancelled")
)
