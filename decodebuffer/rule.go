package decodebuffer

import "bytes"

type Rule interface {
	FindTerminator(d []byte) (valid, remain []byte, found bool)
}

type RuleTerminator struct {
	terminator []byte
}

func NewRuleTerminator(terminator []byte) *RuleTerminator {
	return &RuleTerminator{terminator: terminator}
}

func (rule *RuleTerminator) SetTerminator(terminator []byte) {
	rule.terminator = terminator
}

func (rule *RuleTerminator) FindTerminator(d []byte) (valid, remain []byte, found bool) {
	if len(rule.terminator) <= 0 {
		return
	}

	idx := bytes.Index(d, rule.terminator)

	if idx == -1 {
		return
	}

	valid = d[:idx]
	remain = d[idx+len(rule.terminator):]
	found = true

	return
}
