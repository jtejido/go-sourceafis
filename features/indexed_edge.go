package features

type IndexedEdge struct {
	*EdgeShape
	reference byte
	neighbor  byte
}

func NewIndexedEdge(minutiae []*SearchMinutia, reference, neighbor int) *IndexedEdge {
	es := NewEdgeShape(minutiae[reference], minutiae[neighbor])
	return &IndexedEdge{
		EdgeShape: es,
		reference: byte(reference),
		neighbor:  byte(neighbor),
	}
}

func (i *IndexedEdge) Reference() int {
	return int(i.reference)
}

func (i *IndexedEdge) Neighbor() int {
	return int(i.neighbor)
}
