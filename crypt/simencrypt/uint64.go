package simencrypt

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/sgostarter/i/commerr"
)

var (
	_xorN  = uint64(0x000000000073C3B3)
	_xor   = NewUint64Xor()
	_crypt = NewUint64Crypt()
)

func NewUint64Xor() *Uint64Xor {
	return NewUint64XorEx(0)
}

func NewUint64XorEx(n uint64) *Uint64Xor {
	if n == 0 {
		n = _xorN
	}

	return &Uint64Xor{xorN: new(big.Int).SetUint64(n)}
}

type Uint64Xor struct {
	xorN *big.Int
}

func (xor *Uint64Xor) Xor(n uint64) uint64 {
	return new(big.Int).Xor(new(big.Int).SetUint64(n), xor.xorN).Uint64()
}

//
//
//

func XorUInt64(i uint64) uint64 {
	return _xor.Xor(i)
}

func Uint64N2S(n uint64) string {
	return strconv.FormatUint(n, 31)
}

//
//
//

func UnConfuseUint64(i uint64, trim00Bytes int8) uint64 {
	s := fmt.Sprintf("%04x", i)
	if len(s)%2 != 0 {
		s = "0" + s
	}

	s = strings.Repeat("00", int(trim00Bytes)) + s

	firstS := s[0:2]
	lastS := s[len(s)-2:]

	firstSN, err := strconv.ParseUint(firstS, 16, 8)
	if err != nil {
		panic(err)
	}

	lastSN, err := strconv.ParseUint(lastS, 16, 8)
	if err != nil {
		panic(err)
	}

	firstSN8 := uint8(firstSN)
	lastSN8 := uint8(lastSN)

	firstSN16 := uint16(firstSN8)
	firstSN16 += 0x100
	firstSN16 -= uint16(lastSN8 + 0x27)

	firstSN8 = uint8(firstSN16)

	s = fmt.Sprintf("%02x", firstSN8) + s[2:]

	n, err := strconv.ParseUint(s, 16, 64)
	if err != nil {
		panic(err)
	}

	return n
}

func ConfuseUint64(i uint64) (n uint64, trim00Bytes int8) {
	// 12 0012 1212
	s := fmt.Sprintf("%04x", i)
	if len(s)%2 != 0 {
		s = "0" + s
	}

	firstS := s[0:2]
	lastS := s[len(s)-2:]

	firstSN, err := strconv.ParseUint(firstS, 16, 8)
	if err != nil {
		panic(err)
	}

	lastSN, err := strconv.ParseUint(lastS, 16, 8)
	if err != nil {
		panic(err)
	}

	fistSN16 := uint16(firstSN)
	lastSN16 := uint16(lastSN)

	fistSN16 += lastSN16 + 0x27
	firstN8 := uint8(fistSN16)

	s2 := fmt.Sprintf("%02x", firstN8)
	s = s2 + s[2:]

	for idx := 0; idx < len(s); idx += 2 {
		if s[idx] != '0' || s[idx+1] != '0' {
			break
		}

		trim00Bytes++
	}

	n, err = strconv.ParseUint(s, 16, 64)
	if err != nil {
		panic(err)
	}

	return
}

//
//
//

func Uint64S2N(s string) (uint64, error) {
	return strconv.ParseUint(s, 31, 64)
}

func NewUint64Crypt() *Uint64Crypt {
	return NewUint64CryptEx(nil)
}

func NewUint64CryptEx(xor *Uint64Xor) *Uint64Crypt {
	if xor == nil {
		xor = _xor
	}

	return &Uint64Crypt{
		_xor: xor,
	}
}

type Uint64Crypt struct {
	_xor *Uint64Xor
}

func (crypt *Uint64Crypt) Encrypt(n uint64) string {
	n, trim00Bytes := ConfuseUint64(n)

	return Uint64N2S(crypt._xor.Xor(n)) + fmt.Sprintf("%01x", trim00Bytes)
}

func (crypt *Uint64Crypt) Decrypt(s string) (uint64, error) {
	if len(s) <= 1 {
		return 0, commerr.ErrInvalidArgument
	}

	trimS := s[len(s)-1:]

	trim00Bytes, err := strconv.ParseUint(trimS, 16, 8)
	if err != nil {
		return 0, err
	}

	s = s[:len(s)-1]

	n, err := Uint64S2N(s)
	if err != nil {
		return 0, err
	}

	n = crypt._xor.Xor(n)

	return UnConfuseUint64(n, int8(trim00Bytes)), nil
}

func EncryptUInt64(i uint64) string {
	return _crypt.Encrypt(i)
}

func DecryptUint64(s string) (uint64, error) {
	return _crypt.Decrypt(s)
}
