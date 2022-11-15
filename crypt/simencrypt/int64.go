package simencrypt

import (
	"encoding/hex"
	"strconv"

	"github.com/sgostarter/i/commerr"
)

const (
	_DefMinLen     = 6
	_DefStringSeed = "this is stw"
)

func EncodeID(id int64) string {
	return EncodeIDEx(id, _DefMinLen, _DefStringSeed)
}

func DecodeID(s string) int64 {
	id, _ := DecodeIDE(s)

	return id
}

func DecodeIDE(s string) (int64, error) {
	return DecodeIDEx(s)
}

func EncodeIDEx(id int64, minLen int, seed string) string {
	if len(seed) == 0 {
		seed = _DefStringSeed
	}

	s := strconv.FormatInt(id, 16)
	if len(s) > 0xFF {
		panic("longID")
	}

	l := hex.EncodeToString([]byte{byte(len(s))})

	e := ""

	if len(l)+len(s) < minLen {
		for idx := 0; idx < minLen-len(l)-len(s); idx++ {
			e += seed[(int(id)+idx)%len(seed) : (int(id)+idx)%len(seed)+1]
		}
	}

	return l + s + e
}

func DecodeIDEx(s string) (id int64, err error) {
	if len(s) <= 2 {
		err = commerr.ErrOverflow

		return
	}

	lb, err := hex.DecodeString(s[:2])
	if err != nil {
		return
	}

	l := int(lb[0])

	if 2+l > len(s) {
		err = commerr.ErrOverflow

		return
	}

	id, err = strconv.ParseInt(s[2:2+l], 16, 64)

	return
}
