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

	// nolint: lll
	encodeByteSeeds = []int{43, 208, 35, 155, 132, 214, 104, 210, 254, 27, 112, 84, 136, 72, 125, 175, 180, 145, 164, 24, 228, 69, 234, 157, 117, 172, 107, 130, 80, 76, 149, 148, 73, 81, 85, 90, 56, 249, 115, 255, 241, 200, 120, 247, 191, 93, 88, 131, 64, 217, 183, 239, 42, 21, 95, 86, 171, 54, 20, 109, 126, 31, 22, 223, 252, 92, 140, 33, 177, 219, 0, 6, 75, 70, 68, 3, 161, 197, 170, 11, 55, 110, 102, 179, 243, 224, 119, 181, 248, 133, 221, 212, 162, 17, 244, 186, 53, 152, 12, 250, 82, 146, 44, 40, 29, 118, 114, 209, 36, 153, 207, 105, 66, 213, 182, 143, 169, 206, 134, 123, 26, 233, 101, 204, 245, 62, 30, 34, 227, 189, 122, 78, 91, 160, 187, 18, 237, 61, 251, 159, 97, 194, 9, 47, 236, 23, 253, 229, 16, 174, 96, 13, 67, 7, 188, 154, 65, 128, 240, 163, 10, 168, 211, 203, 226, 137, 39, 220, 74, 116, 165, 141, 238, 147, 235, 176, 100, 52, 198, 192, 108, 37, 142, 113, 4, 218, 48, 190, 225, 89, 106, 150, 83, 151, 178, 1, 242, 94, 205, 77, 63, 144, 193, 184, 138, 98, 124, 201, 199, 216, 99, 195, 58, 49, 185, 19, 215, 28, 196, 103, 246, 167, 231, 156, 60, 38, 14, 87, 46, 230, 59, 50, 127, 121, 202, 139, 111, 25, 2, 166, 232, 41, 45, 158, 8, 71, 57, 79, 173, 135, 15, 129, 51, 5, 32, 222}
	// nolint: lll
	decodeByteSeeds = []int{70, 195, 238, 75, 184, 253, 71, 153, 244, 142, 160, 79, 98, 151, 226, 250, 148, 93, 135, 215, 58, 53, 62, 145, 19, 237, 120, 9, 217, 104, 126, 61, 254, 67, 127, 2, 108, 181, 225, 166, 103, 241, 52, 0, 102, 242, 228, 143, 186, 213, 231, 252, 177, 96, 57, 80, 36, 246, 212, 230, 224, 137, 125, 200, 48, 156, 112, 152, 74, 21, 73, 245, 13, 32, 168, 72, 29, 199, 131, 247, 28, 33, 100, 192, 11, 34, 55, 227, 46, 189, 35, 132, 65, 45, 197, 54, 150, 140, 205, 210, 176, 122, 82, 219, 6, 111, 190, 26, 180, 59, 81, 236, 10, 183, 106, 38, 169, 24, 105, 86, 42, 233, 130, 119, 206, 14, 60, 232, 157, 251, 27, 47, 4, 89, 118, 249, 12, 165, 204, 235, 66, 171, 182, 115, 201, 17, 101, 173, 31, 30, 191, 193, 97, 109, 155, 3, 223, 23, 243, 139, 133, 76, 92, 159, 18, 170, 239, 221, 161, 116, 78, 56, 25, 248, 149, 15, 175, 68, 194, 83, 16, 87, 114, 50, 203, 214, 95, 134, 154, 129, 187, 44, 179, 202, 141, 211, 218, 77, 178, 208, 41, 207, 234, 163, 123, 198, 117, 110, 1, 107, 7, 162, 91, 113, 5, 216, 209, 49, 185, 69, 167, 90, 255, 63, 85, 188, 164, 128, 20, 147, 229, 222, 240, 121, 22, 174, 144, 136, 172, 51, 158, 40, 196, 84, 94, 124, 220, 43, 88, 37, 99, 138, 64, 146, 8, 39}
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
	lastSN8 := uint8(decodeByteSeeds[lastSN])

	firstSN16 := uint16(firstSN8)
	firstSN16 += 0x100
	firstSN16 -= uint16(lastSN8 + 0x27)

	firstSN8 = uint8(firstSN16)

	s = fmt.Sprintf("%02x", firstSN8) + s[2:len(s)-2] + fmt.Sprintf("%02x", lastSN8)

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

	lastSN8 := encodeByteSeeds[lastSN16]

	firstS = fmt.Sprintf("%02x", firstN8)
	lastS = fmt.Sprintf("%02x", lastSN8)
	s = firstS + s[2:len(s)-2] + lastS

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
