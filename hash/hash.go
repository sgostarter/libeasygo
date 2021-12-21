package hash

import (
	// nolint: gosec
	"crypto/md5"
	// nolint: gosec
	"crypto/sha1"
	"encoding/hex"
)

func MD5(s string) string {
	// nolint: gosec
	sum := md5.Sum([]byte(s))

	return hex.EncodeToString(sum[:])
}

func SHA1(s string) string {
	// nolint: gosec
	sum := sha1.Sum([]byte(s))

	return hex.EncodeToString(sum[:])
}
