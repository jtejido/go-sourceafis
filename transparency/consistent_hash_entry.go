package transparency

import "github.com/jtejido/sourceafis/features"

type ConsistentHashEntry struct {
	Key   int
	Edges []*features.IndexedEdge
}
