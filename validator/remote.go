package validator

import (
	"strconv"
	"strings"
)

func ValidateRemoteAddress(s string) (host string, port int, ok bool) {
	ps := strings.Split(s, ":")
	if len(ps) != 2 {
		return
	}

	if ps[0] == "" || ps[1] == "" {
		return
	}

	port, err := strconv.Atoi(ps[1])
	if err != nil {
		return
	}

	if port <= 0 || port >= 65535 {
		return
	}

	host = strings.TrimLeft(ps[0], "\r\n\t ")
	ok = true

	return
}
