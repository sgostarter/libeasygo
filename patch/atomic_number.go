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

func (n *AtomicNumber[T]) Add(delta T) T {
	n.lock.Lock()
	defer n.lock.Unlock()

	n.v += delta

	return n.v
}

func (n *AtomicNumber[T]) Sub(delta T) T {
	n.lock.Lock()
	defer n.lock.Unlock()

	n.v -= delta

	return n.v
}

func (n *AtomicNumber[T]) Swap(val T) (old T) {
	n.lock.Lock()
	defer n.lock.Unlock()

	old = n.v
	n.v = val

	return
}
