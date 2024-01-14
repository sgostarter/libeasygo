package ptl

import (
	"fmt"
	"sync/atomic"
)

type Code int

const (
	CodeSuccess Code = iota
	CodeErrCommunication
	CodeErrInvalidArgs
	CodeErrInternal
	CodeErrBadToken
	CodeErrNeedAuth
	CodeErrDisabled

	CodeErrCustomStart = 1000
	CodeErrCustomEnd   = 3000
)

type FNCode2Message func(code Code) (msg string, ok bool)

type fnCode2MessageWrapper struct {
	fn FNCode2Message
}

var (
	_exCode2Message atomic.Pointer[fnCode2MessageWrapper]
)

// InstallCode2Message warning, not thread safe
func InstallCode2Message(f FNCode2Message) {
	_exCode2Message.Store(&fnCode2MessageWrapper{
		fn: f,
	})
}

func getCode2MessageFn() FNCode2Message {
	wrapper := _exCode2Message.Load()
	if wrapper == nil {
		return nil
	}

	return wrapper.fn
}

func (c Code) String() string {
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
		if fn := getCode2MessageFn(); fn != nil {
			msg, ok := fn(c)
			if ok {
				return msg
			}
		}
	}

	return fmt.Sprintf("未知错误%d", c)
}

func CodeToMessage(code Code, msg string) string {
	codeMsg := code.String()

	if msg != "" {
		codeMsg += ":" + msg
	}

	return codeMsg
}
