package transparency

import "sourceafis/features"

type ConsistentHashEntry struct {
	Key   int
	Edges []*features.IndexedEdge
}
