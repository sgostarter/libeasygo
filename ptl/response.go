package ptl

type ResponseWrapper struct {
	Code    Code        `json:"code"`
	Message string      `json:"message"`
	Resp    interface{} `json:"resp,omitempty"`
}

func (wr *ResponseWrapper) Apply(code Code, msg string) bool {
	wr.Code = code
	wr.Message = CodeToMessage(code, msg)

	return wr.Code == CodeSuccess
}
