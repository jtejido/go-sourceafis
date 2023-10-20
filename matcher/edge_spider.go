package matcher

import (
	"container/heap"

	"github.com/jtejido/sourceafis/config"
	"github.com/jtejido/sourceafis/features"
)

func Crawl(pedges, cedges [][]*features.NeighborEdge, pairing *PairingGraph, root *MinutiaPair, queue *PriorityQueue) {
	heap.Push(queue, root)
	for queue.Len() > 0 {
		item := heap.Pop(queue).(*MinutiaPair)
		pairing.AddPair(item)
		collectEdges(pedges, cedges, pairing, queue)
		skipPaired(pairing, queue)
	}

}

func collectEdges(pedges, cedges [][]*features.NeighborEdge, pairing *PairingGraph, queue *PriorityQueue) {
	reference := pairing.Tree[pairing.Count-1]
	pstar := pedges[reference.Probe]
	cstar := cedges[reference.Candidate]

	for _, pair := range matchPairs(pstar, cstar, pairing.pool) {
		pair.ProbeRef = reference.Probe
		pair.CandidateRef = reference.Candidate
		if pairing.byCandidate[pair.Candidate] == nil && pairing.byProbe[pair.Probe] == nil {
			heap.Push(queue, pair)
		} else {
			pairing.SupportPair(pair)
		}
	}
}
func skipPaired(pairing *PairingGraph, queue *PriorityQueue) {
	for queue.Len() > 0 && (pairing.byProbe[(*queue)[0].Probe] != nil || pairing.byCandidate[(*queue)[0].Candidate] != nil) {
		item := heap.Pop(queue).(*MinutiaPair)
		pairing.SupportPair(item)
	}
}

func matchPairs(pstar, cstar []*features.NeighborEdge, pool *MinutiaPairPool) []*MinutiaPair {
	results := make([]*MinutiaPair, 0)
	var start, end int
	for cindex := 0; cindex < len(cstar); cindex++ {
		cedge := cstar[cindex]
		for start < len(pstar) && pstar[start].Length < cedge.Length-config.Config.MaxDistanceError {
			start++
		}
		if end < start {
			end = start
		}
		for end < len(pstar) && pstar[end].Length <= cedge.Length+config.Config.MaxDistanceError {
			end++
		}
		for pindex := start; pindex < end; pindex++ {
			pedge := pstar[pindex]
			rdiff := float64(pedge.ReferenceAngle.Difference(cedge.ReferenceAngle))
			if rdiff <= config.Config.MaxAngleError || rdiff >= ComplementaryMaxAngleError() {
				ndiff := float64(pedge.NeighborAngle.Difference(cedge.NeighborAngle))
				if ndiff <= config.Config.MaxAngleError || ndiff >= ComplementaryMaxAngleError() {
					pair := pool.Allocate()
					pair.Probe = pedge.Neighbor
					pair.Candidate = cedge.Neighbor
					pair.distance = cedge.Length
					results = append(results, pair)
				}
			}
		}
	}
	return results
}
