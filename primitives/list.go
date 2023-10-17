package primitives

type List[V any] interface {
	Add(item V) error
	AddAt(index int, item V) error
	Get(index int) (V, error)
	Remove(index int) error
	Size() int
	Iterator() ListIterator[V]
}
