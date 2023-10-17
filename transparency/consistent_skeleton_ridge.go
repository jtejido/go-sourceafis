package transparency

import "sourceafis/primitives"

type ConsistentSkeletonRidge struct {
	Start, End int
	Points     primitives.List[primitives.IntPoint]
}
