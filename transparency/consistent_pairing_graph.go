package transparency

import "sourceafis/matcher"

type ConsistentPairingGraph struct {
	Root    *ConsistentMinutiaPair
	Tree    []*ConsistentEdgePair
	Support []*ConsistentEdgePair
}

func NewConsistentPairingGraph(pairing *matcher.PairingGraph) *ConsistentPairingGraph {
	tree := make([]*ConsistentEdgePair, 0)
	for i := 1; i < len(pairing.Tree) && i < pairing.Count; i++ {
		tree = append(tree, NewConsistentEdgePair(pairing.Tree[i]))
	}
	support := make([]*ConsistentEdgePair, len(pairing.Support))
	for i, edge := range pairing.Support {
		support[i] = NewConsistentEdgePair(edge)
	}
	return &ConsistentPairingGraph{
		Root:    &ConsistentMinutiaPair{Probe: pairing.Tree[0].Probe, Candidate: pairing.Tree[0].Candidate},
		Tree:    tree,
		Support: support,
	}
}
