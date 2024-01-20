package patch

import (
	"sync"

	"golang.org/x/exp/constraints"
)

type AtomicNumber[T constraints.Integer | constraints.Float] struct {
	lock sync.Mutex
	v    T
}

func (n *AtomicNumber[T]) Load() T {
	n.lock.Lock()
	defer n.lock.Unlock()

	return n.v
}

func (n *AtomicNumber[T]) Store(v T) {
	n.lock.Lock()
	defer n.lock.Unlock()

	n.v = v
}

func (n *AtomicNumber[T]) Inc() T {
	n.lock.Lock()
	defer n.lock.Unlock()

	n.v++

	return n.v
}

func (n *AtomicNumber[T]) Dec() T {
	n.lock.Lock()
	defer n.lock.Unlock()

	n.v--

	return n.v
}
