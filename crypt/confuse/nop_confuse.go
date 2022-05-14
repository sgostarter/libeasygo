package confuse

func NewNopConfuse() Confuse {
	return &nopConfuseImpl{}
}

type nopConfuseImpl struct {
}

func (impl *nopConfuseImpl) Seal(d []byte) ([]byte, error) {
	return d, nil
}

func (impl *nopConfuseImpl) Open(d []byte) ([]byte, error) {
	return d, nil
}
