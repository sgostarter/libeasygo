package ptl

import (
	"context"
	"fmt"
	"sync/atomic"
)

type Code int

const (
	CodeSuccess Code = iota
	CodeErrUnknown
	CodeErrCommunication
	CodeErrInvalidArgs
	CodeErrInternal
	CodeErrBadToken
	CodeErrNeedAuth
	CodeErrDisabled
	CodeErrExists
	CodeErrNotExists
	CodeErrLogic
	CodeErrConflict
	CodeErrUnimplemented
	CodeErrInvalidToken
	CodeErrVerify
	CodeErrResourceExhausted
	CodeErrPartSuccess
	CodeErrFail

	CodeErrCustomStart = 1000
	CodeErrCustomEnd   = 3000
)

type FNCode2Message func(code Code) (msg string, ok bool)
type FNCode2MessageWithContext func(ctx context.Context, code Code) (msg string, ok bool)

type fnCode2MessageWrapper struct {
	fnPre            FNCode2Message
	fnEx             FNCode2Message
	fnPreWithContext FNCode2MessageWithContext
	fnExWithContext  FNCode2MessageWithContext
}

var (
	_exCode2Message atomic.Pointer[fnCode2MessageWrapper]
)

// InstallCode2Message warning, not thread safe
func InstallCode2Message(fnPre, fnEx FNCode2Message) {
	InstallCode2MessageEx(fnPre, fnEx, nil, nil)
}

// InstallCode2MessageEx warning, not thread safe
func InstallCode2MessageEx(fnPre, fnEx FNCode2Message, fnPreWithContext, fnExWithContext FNCode2MessageWithContext) {
	_exCode2Message.Store(&fnCode2MessageWrapper{
		fnPre:            fnPre,
		fnEx:             fnEx,
		fnPreWithContext: fnPreWithContext,
		fnExWithContext:  fnExWithContext,
	})
}

func getCode2MessageFn() (fnPre, fnEx FNCode2Message, fnPreWithContext, fnExWithContext FNCode2MessageWithContext) {
	wrapper := _exCode2Message.Load()
	if wrapper == nil {
		return
	}

	fnPre = wrapper.fnPre
	fnEx = wrapper.fnEx
	fnPreWithContext = wrapper.fnPreWithContext
	fnExWithContext = wrapper.fnExWithContext

	return
}

func (c Code) Key() string {
	switch c {
	case CodeSuccess:
		return "CodeSuccess"
	case CodeErrCommunication:
		return "CodeErrCommunication"
	case CodeErrInvalidArgs:
		return "CodeErrInvalidArgs"
	case CodeErrInternal:
		return "CodeErrInternal"
	case CodeErrBadToken:
		return "CodeErrBadToken"
	case CodeErrNeedAuth:
		return "CodeErrNeedAuth"
	case CodeErrDisabled:
		return "CodeErrDisabled"
	case CodeErrExists:
		return "CodeErrExists"
	case CodeErrNotExists:
		return "CodeErrNotExists"
	case CodeErrLogic:
		return "CodeErrLogic"
	case CodeErrConflict:
		return "CodeErrConflict"
	case CodeErrUnimplemented:
		return "CodeErrUnimplemented"
	case CodeErrInvalidToken:
		return "CodeErrInvalidToken"
	case CodeErrVerify:
		return "CodeErrVerify"
	case CodeErrResourceExhausted:
		return "CodeErrResourceExhausted"
	case CodeErrPartSuccess:
		return "CodeErrPartSuccess"
	case CodeErrFail:
		return "CodeErrFail"
	}

	return ""
}

func (c Code) String() string {
	return c.StringWithContext(context.Background())
}

func (c Code) StringWithContext(ctx context.Context) string {
	fnPre, fnEx, fnPreWithContext, fnExWithContext := getCode2MessageFn()

	if fnPreWithContext != nil {
		t, ok := fnPreWithContext(ctx, c)
		if ok {
			return t
		}
	}

	if fnPre != nil {
		t, ok := fnPre(c)
		if ok {
			return t
		}
	}

	t := c.Key()
	if t != "" {
		return t
	}

	if fnExWithContext != nil {
		msg, ok := fnExWithContext(ctx, c)
		if ok {
			return msg
		}
	}

	if fnEx != nil {
		msg, ok := fnEx(c)
		if ok {
			return msg
		}
	}

	return fmt.Sprintf("Unknow error: %d", c)
}

func CodeToMessage(code Code, msg string) string {
	return CodeToMessageWithContext(context.Background(), code, msg)
}

func CodeToMessageWithContext(ctx context.Context, code Code, msg string) string {
	if msg != "" {
		return msg
	}

	return code.StringWithContext(ctx)
}
