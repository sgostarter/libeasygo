package query

import "fmt"

type parseError struct {
	error
	action string
	tag    string
	tp     string
	val    string
	key    string
}

func newParseError(err error, action, tag, tp, val, key string) parseError {
	return parseError{
		error:  err,
		action: action,
		tag:    tag,
		tp:     tp,
		val:    val,
		key:    key,
	}
}

func (e parseError) Error() string {
	return fmt.Sprintf(" %s failed: type %s, tag: %s, val %s, key %s %v", e.action, e.tag, e.tp, e.val, e.key, e.error)
}
