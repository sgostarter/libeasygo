package validator

import (
	"strings"

	"github.com/sgostarter/libeasygo/commerr"
)

const (
	recognizeS = "://"

	ShouldUseUnknown = 0
	ShouldUseTLS     = 1
	ShouldUsePlain   = 2
)

func PoolAddressShouldTLS(s string) int {
	s = strings.ToLower(s)

	if strings.HasPrefix(s, "ssl"+recognizeS) || strings.HasSuffix(s, "stratum+ssl"+recognizeS) {
		return ShouldUseTLS
	}

	if strings.HasPrefix(s, "tcp"+recognizeS) || strings.HasSuffix(s, "stratum+tcp"+recognizeS) {
		return ShouldUsePlain
	}

	return ShouldUseUnknown
}

func PoolAddressValidate(s string) (hostAndPort string, shouldUse int, err error) {
	s = strings.TrimLeft(s, "\r\n\t ")

	hostAndPort, shouldUse = poolAddressValidate(s)
	if !ValidateHostAndPort(hostAndPort) {
		err = commerr.ErrUnknownBadFormat

		return
	}

	return
}

func poolAddressValidate(s string) (string, int) {
	ps := strings.SplitN(s, recognizeS, 2)
	if len(ps) != 2 {
		return s, PoolAddressShouldTLS(s)
	}

	return ps[1], PoolAddressShouldTLS(ps[0] + recognizeS)
}
