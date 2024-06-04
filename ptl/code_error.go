package ptl

func NewCodeError(code Code) CodeError {
	return CodeError{
		code: code,
	}
}

type CodeError struct {
	code Code
}

func (c CodeError) Error() string {
	return c.code.String()
}
