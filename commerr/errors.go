package commerr

import "errors"

var (
	ErrCanceled        = errors.New("cancelled")
	ErrUnknown         = errors.New("unknown")
	ErrInvalidArgument = errors.New("invalidArgument")
)
