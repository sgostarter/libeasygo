package simencrypt

import (
	"math"
	"strconv"

	"github.com/sgostarter/libeasygo/hash"
)

func GetStringIndex(s string) int {
	n, _ := strconv.ParseInt(hash.MD5(s), 16, 64)
	if n < 0 {
		n = -n
	}

	n = n % math.MaxInt

	return int(n)
}
