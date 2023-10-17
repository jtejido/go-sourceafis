package primitives

type IntPointIterator interface {
	HasNext() bool
	Next() IntPoint
}

type ListIterator[V any] interface {
	HasNext() bool
	Next() (V, error)
}

type iterator struct {
	IntPoint
	atX, atY int
}

func newIterator(s IntPoint) *iterator {
	return &iterator{s, 0, 0}
}

// Next returns the next item in the collection.
func (i *iterator) Next() IntPoint {
	result := IntPoint{i.atX, i.atY}
	i.atX++
	if i.atX >= i.X {
		i.atX = 0
		i.atY++
	}

	return result
}

// HasNext return true if there are values to be read.
func (i *iterator) HasNext() bool {
	return i.atY < i.Y && i.atX < i.X
}

type blockIterator struct {
	IntRect
	atX, atY int
}

func newBlockIterator(s IntRect) *blockIterator {
	return &blockIterator{s, 0, 0}
}

// Next returns the next item in the collection.
func (i *blockIterator) Next() IntPoint {
	result := IntPoint{i.X + i.atX, i.Y + i.atY}
	i.atX++
	if i.atX >= i.Width {
		i.atX = 0
		i.atY++
	}

	return result
}

// HasNext return true if there are values to be read.
func (i *blockIterator) HasNext() bool {
	return i.atY < i.Height && i.atX < i.Width
}
