package commerr

import "errors"

var (
	ErrCanceled          = errors.New("cancelled")
	ErrUnknown           = errors.New("unknown")
	ErrInvalidArgument   = errors.New("invalidArgument")
	ErrInternal          = errors.New("internal")
	ErrNotFound          = errors.New("notFound")
	ErrAlreadyExists     = errors.New("alreadyExists")
	ErrPermissionDenied  = errors.New("permissionDenied")
	ErrAborted           = errors.New("aborted")
	ErrOutOfRange        = errors.New("outOfRange")
	ErrUnimplemented     = errors.New("unimplemented")
	ErrUnavailable       = errors.New("unavailable")
	ErrUnauthenticated   = errors.New("unauthenticated")
	ErrResourceExhausted = errors.New("resourceExhausted")
	ErrReject            = errors.New("reject")
	ErrCrash             = errors.New("crash")
	ErrOverflow          = errors.New("overflow")
	ErrUnknownBadFormat  = errors.New("unknownFlrmat")
	ErrTimeout           = errors.New("timeout")
	ErrExiting           = errors.New("exiting")
)
