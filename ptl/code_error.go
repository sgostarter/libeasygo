package ptl

import "errors"

func NewCodeError(code Code) CodeError {
	return CodeError{
		code: code,
	}
}

func NewCodeErrorEx(code Code, msg string) CodeError {
	return CodeError{
		code: code,
		msg:  msg,
	}
}

func NewCodeErrorFromError(err error) CodeError {
	ce := CodeError{}

	ce.code, ce.msg = CodeFromError(err)

	return ce
}

func NewCodeErrorFromCodeAndError(code Code, err error) CodeError {
	ce := CodeError{}

	ce.code = code
	if err != nil {
		ce.msg = err.Error()
	}

	return ce
}

type CodeError struct {
	code Code
	msg  string
}

func (c CodeError) Error() string {
	return c.code.String()
}

func (c CodeError) GetCode() Code {
	return c.code
}

func (c CodeError) GetMsg() string {
	return c.msg
}

func (c CodeError) Success() bool {
	return c.code == CodeSuccess
}

func CodeFromError(err error) (Code, string) {
	if err == nil {
		return CodeSuccess, ""
	}

	var ce CodeError

	if errors.As(err, &ce) {
		return ce.code, ce.Error()
	}

	return CodeErrInternal, ""
}
