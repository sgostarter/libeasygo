package simencrypt

import (
	"encoding/hex"

	"github.com/sgostarter/libeasygo/commerr"
)

func EncodeString(s string) string {
	return EncodeStringEx(s, _DefMinLen, _DefStringSeed)
}

func EncodeStringEx(s string, minLen int, seed string) string {
	if seed == "" {
		seed = _DefStringSeed
	}

	if len(s) >= 0xFF {
		return hex.EncodeToString([]byte{0xff}) + s
	}

	n := GetStringIndex(s)

	l := hex.EncodeToString([]byte{byte(len(s))})

	e := ""

	if len(l)+len(s) < minLen {
		for idx := 0; idx < minLen-len(l)-len(s); idx++ {
			e += seed[(n+idx)%len(seed) : (n+idx)%len(seed)+1]
		}
	}

	return l + s + e
}

func DecodeString(s string) (string, error) {
	return DecodeStringEx(s)
}

func DecodeStringEx(s string) (ret string, err error) {
	if len(s) <= 2 {
		err = commerr.ErrOverflow

		return
	}

	lb, err := hex.DecodeString(s[:2])
	if err != nil {
		return
	}

	if lb[0] == 0xFF {
		ret = s[2:]

		return
	}

	l := int(lb[0])

	if 2+l > len(s) {
		err = commerr.ErrOverflow

		return
	}

	ret = s[2 : 2+l]

	return
}
