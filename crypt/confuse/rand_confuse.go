package confuse

import "math/rand"

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

func (impl randConfuseImpl) randN(n int) []byte {
	if n == 0 {
		return []byte{}
	}

	r := make([]byte, n)
	for i := range r {
		// nolint: gosec
		r[i] = byte(rand.Intn(0xff))
	}

	return r
}

func (impl randConfuseImpl) Seal(d []byte) ([]byte, error) {
	cd := make([]byte, 0, len(d)*3)

	for i := range d {
		cd = append(cd, impl.randN(impl.wN.NextN())...)
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
