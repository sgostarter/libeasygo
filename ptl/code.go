package ptl

import (
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
	CodeLogic
	CodeConflict

	CodeErrCustomStart = 1000
	CodeErrCustomEnd   = 3000
)

type FNCode2Message func(code Code) (msg string, ok bool)

type fnCode2MessageWrapper struct {
	fnPre FNCode2Message
	fnEx  FNCode2Message
}

var (
	_exCode2Message atomic.Pointer[fnCode2MessageWrapper]
)

// InstallCode2Message warning, not thread safe
func InstallCode2Message(fnPre, fnEx FNCode2Message) {
	_exCode2Message.Store(&fnCode2MessageWrapper{
		fnPre: fnPre,
		fnEx:  fnEx,
	})
}

func getCode2MessageFn() (fnPre, fnEx FNCode2Message) {
	wrapper := _exCode2Message.Load()
	if wrapper == nil {
		return
	}

	fnPre = wrapper.fnPre
	fnEx = wrapper.fnEx

	return
}

func (c Code) String() string {
	fnPre, fnEx := getCode2MessageFn()

	if fnPre != nil {
		t, ok := fnPre(c)
		if ok {
			return t
		}
	}

	switch c {
	case CodeSuccess:
		return "成功"
	case CodeErrCommunication:
		return "通信出错"
	case CodeErrInvalidArgs:
		return "参数非法"
	case CodeErrInternal:
		return "内部错误"
	case CodeErrBadToken:
		return "凭证非法"
	case CodeErrNeedAuth:
		return "需要授权"
	case CodeErrDisabled:
		return "被禁止"
	default:
		if fnEx != nil {
			msg, ok := fnEx(c)
			if ok {
				return msg
			}
		}
	}

	return fmt.Sprintf("Unknow error: %d", c)
}

func CodeToMessage(code Code, msg string) string {
	codeMsg := code.String()

	if msg != "" {
		codeMsg += ":" + msg
	}

	return codeMsg
}
