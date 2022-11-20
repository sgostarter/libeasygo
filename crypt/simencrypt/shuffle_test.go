package simencrypt

import (
	"math/rand"
	"strconv"
	"testing"
	"time"
)

// 3 4 1 2 0
// 4 2 3 0 1

// nolint
func TestShuffle(t *testing.T) {
	a := make([]uint8, 0xFF+1)
	for idx := 0; idx < len(a); idx++ {
		a[idx] = uint8(idx)
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(a), func(i, j int) {
		a[i], a[j] = a[j], a[i]
	})

	m := make(map[uint8]uint8)

	for i, u := range a {
		m[u] = uint8(i)
	}

	b := make([]uint8, len(a))
	for idx := 0; idx < len(b); idx++ {
		b[idx] = m[uint8(idx)]
	}

	s := `encodeByteSeeds = []int{`
	for idx := 0; idx < len(a); idx++ {
		s += strconv.Itoa(int(a[idx])) + ","
	}
	s = s[:len(s)-1]
	s += `}`

	t.Log(s)

	s = `decodeByteSeeds = []int{`
	for idx := 0; idx < len(b); idx++ {
		s += strconv.Itoa(int(b[idx])) + ","
	}
	s = s[:len(s)-1]
	s += `}`

	t.Log(s)
}
