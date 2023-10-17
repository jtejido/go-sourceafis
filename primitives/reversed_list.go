package primitives

import "fmt"

type ReversedList[V any] struct {
	inner List[V]
}

func NewReversedList[V any](items List[V]) *ReversedList[V] {
	return &ReversedList[V]{
		inner: items,
	}
}

func (l *ReversedList[V]) Add(item V) error {
	return l.inner.AddAt(0, item)
}

func (l *ReversedList[V]) AddAt(index int, item V) error {
	return l.inner.AddAt(l.Size()-index, item)
}

func (l *ReversedList[V]) Get(index int) (V, error) {
	return l.inner.Get(l.Size() - index - 1)
}

func (l *ReversedList[V]) Size() int {
	return l.inner.Size()
}

func (l *ReversedList[V]) Remove(index int) error {
	return l.inner.Remove(l.Size() - index - 1)
}

func (l *ReversedList[V]) Iterator() ListIterator[V] {
	return newReversedIterator(l)
}

type reversedIterator[V any] struct {
	*ReversedList[V]
	index int
}

func newReversedIterator[V any](s *ReversedList[V]) *reversedIterator[V] {
	return &reversedIterator[V]{s, 0}
}

// Next returns the next item in the collection.
func (i *reversedIterator[V]) Next() (val V, err error) {
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
func (i *reversedIterator[V]) HasNext() bool {
	return i.index < i.Size()
}
