package validator

import (
	"strconv"
	"strings"
)

func ValidateHostAndPort(s string) bool {
	ps := strings.Split(s, ":")
	if len(ps) != 2 {
		return false
	}

	if ps[0] == "" || ps[1] == "" {
		return false
	}

	port, err := strconv.Atoi(ps[1])
	if err != nil {
		return false
	}

	return port > 0 && port < 65535
}
