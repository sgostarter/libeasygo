package confuse

type Confuse interface {
	Seal(d []byte) ([]byte, error)
	Open(d []byte) ([]byte, error)
}
