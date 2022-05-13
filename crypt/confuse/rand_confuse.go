package confuse

import (
	"math/rand"
)

func NewRandConfuse() Confuse {
	return NewRandConfuseEx(NewDefN(), NewDefN())
}

func NewRandConfuseEx(rN, wN N) Confuse {
	return &randConfuseImpl{
		rN: rN,
		wN: wN,
	}
}

type randConfuseImpl struct {
	rN N
	wN N
}

func (impl randConfuseImpl) Seal(d []byte) ([]byte, error) {
	cd := make([]byte, 0, len(d)*3)

	randBuf := make([]byte, 300)

	// nolint: gosec
	_, _ = rand.Read(randBuf)

	randIdx := 0

	fnRandN := func(n int) []byte {
		if randIdx+n >= len(randBuf) {
			randIdx = 0
		}

		return randBuf[randIdx : randIdx+n]
	}

	for i := range d {
		cd = append(cd, fnRandN(impl.wN.NextN())...)
		cd = append(cd, d[i])
	}

	return cd, nil
}

func (impl randConfuseImpl) Open(d []byte) ([]byte, error) {
	cd := make([]byte, 0, len(d))

	idx := 0

	for {
		if idx >= len(d) {
			break
		}

		n := impl.rN.NextN()
		idx = idx + n

		cd = append(cd, d[idx])
		idx++
	}

	return cd, nil
}
