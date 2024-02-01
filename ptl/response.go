package ptl

type ResponseWrapper struct {
	Code       Code        `json:"code"`
	Message    string      `json:"message"`
	Resp       interface{} `json:"resp,omitempty"`
	RawMessage string      `json:"-" yaml:"-"`
}

func (wr *ResponseWrapper) Apply(code Code, msg string) bool {
	wr.Code = code
	wr.RawMessage = msg
	wr.Message = CodeToMessage(code, msg)

	return wr.Code == CodeSuccess
}

func (wr *ResponseWrapper) Clone(wro ResponseWrapper) bool {
	wr.Code = wro.Code
	wr.RawMessage = wro.RawMessage
	wr.Message = wro.Message

	return wr.Code == CodeSuccess
}
