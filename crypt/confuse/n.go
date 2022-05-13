package confuse

type N interface {
	NextN() int
}

func NewDefN() N {
	return &defNImpl{}
}

type defNImpl struct {
	n int
}

func (impl *defNImpl) NextN() int {
	switch impl.n {
	case 0:
		impl.n = 2
	case 2:
		impl.n = 1
	case 1:
		impl.n = 0
	}

	return impl.n
}
