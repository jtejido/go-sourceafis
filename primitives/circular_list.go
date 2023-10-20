package primitives

import "fmt"

type CircularList[V any] struct {
	inner *CircularArray
}

func NewCircularList[V any]() *CircularList[V] {
	return &CircularList[V]{
		inner: NewCircularArray(16),
	}
}

func (l *CircularList[V]) Get(index int) (V, error) {
	v, err := l.inner.get(index)
	return v.(V), err
}

func (l *CircularList[V]) Add(item V) error {
	if err := l.inner.insert(l.inner.size, 1); err != nil {
		return err
	}
	if err := l.inner.set(l.inner.size-1, item); err != nil {
		return err
	}

	return nil
}

func (l *CircularList[V]) AddAt(index int, item V) error {
	if err := l.inner.insert(index, 1); err != nil {
		return err
	}
	if err := l.inner.set(index, item); err != nil {
		return err
	}

	return nil
}

func (l *CircularList[V]) Size() int {
	return l.inner.size
}

func (l *CircularList[V]) Remove(index int) error {
	return l.inner.remove(index, 1)
}

func (l *CircularList[V]) Iterator() ListIterator[V] {
	return newArrayIterator(l)
}

type arrayIterator[V any] struct {
	*CircularList[V]
	index int
}

func newArrayIterator[V any](s *CircularList[V]) *arrayIterator[V] {
	return &arrayIterator[V]{s, 0}
}

// Next returns the next item in the collection.
func (i *arrayIterator[V]) Next() (val V, err error) {
	if i.index >= i.Size() {
		err = fmt.Errorf("no element exist.")
	}
	if err != nil {
		return
	}
	i.index++
	return i.Get(i.index - 1)
}

// HasNext return true if there are values to be read.
func (i *arrayIterator[V]) HasNext() bool {
	return i.index < i.Size()
}
