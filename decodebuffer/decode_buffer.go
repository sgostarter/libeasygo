package decodebuffer

type Buffer struct {
	dataBuffer []byte
	rule       Rule
}

func NewBuffer(rule Rule) *Buffer {
	return &Buffer{
		rule: rule,
	}
}

func (buf *Buffer) SetDecodeRule(rule Rule) {
	buf.rule = rule
}

func (buf *Buffer) Append(data []byte) {
	buf.dataBuffer = append(buf.dataBuffer, data...)
}

func (buf *Buffer) FindTerminator() (valid []byte, found bool) {
	rule := buf.rule
	if rule == nil {
		return
	}

	valid, remain, found := rule.FindTerminator(buf.dataBuffer)

	if !found {
		return
	}

	buf.dataBuffer = remain

	return
}

func (buf *Buffer) Len() int {
	return len(buf.dataBuffer)
}
