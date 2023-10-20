package transparency

import "github.com/jtejido/sourceafis/primitives"

type ConsistentSkeletonRidge struct {
	Start, End int
	Points     primitives.List[primitives.IntPoint]
}
