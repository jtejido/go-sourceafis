package matcher

import "github.com/jtejido/sourceafis/templates"

type PairingGraph struct {
	pool                                *MinutiaPairPool
	Count                               int
	Tree, byProbe, byCandidate, Support []*MinutiaPair
	supportEnabled                      bool
}

func NewPairingGraph(pool *MinutiaPairPool) *PairingGraph {
	return &PairingGraph{
		pool:        pool,
		Tree:        make([]*MinutiaPair, 1),
		byProbe:     make([]*MinutiaPair, 1),
		byCandidate: make([]*MinutiaPair, 1),
		Support:     make([]*MinutiaPair, 0),
	}
}

func (g *PairingGraph) ReserveProbe(probe *Probe) {
	capacity := len(probe.template.Minutiae)
	if capacity > len(g.Tree) {
		g.Tree = make([]*MinutiaPair, capacity)
		g.byProbe = make([]*MinutiaPair, capacity)
	}
}
func (g *PairingGraph) ReserveCandidate(candidate *templates.SearchTemplate) {
	capacity := len(candidate.Minutiae)
	if len(g.byCandidate) < capacity {
		g.byCandidate = make([]*MinutiaPair, capacity)
	}
}

func (g *PairingGraph) AddPair(pair *MinutiaPair) {
	g.Tree[g.Count] = pair
	g.byProbe[pair.Probe] = pair
	g.byCandidate[pair.Candidate] = pair
	g.Count++
}

func (g *PairingGraph) SupportPair(pair *MinutiaPair) {
	if g.byProbe[pair.Probe] != nil && g.byProbe[pair.Probe].Candidate == pair.Candidate {
		g.byProbe[pair.Probe].supportingEdges++
		g.byProbe[pair.ProbeRef].supportingEdges++
		if g.supportEnabled {
			g.Support = append(g.Support, pair)
		} else {
			g.pool.Release(pair)
		}
	} else {
		g.pool.Release(pair)
	}
}

func (g *PairingGraph) Clear() {
	for i := 0; i < g.Count; i++ {
		g.byProbe[g.Tree[i].Probe] = nil
		g.byCandidate[g.Tree[i].Candidate] = nil
		if i > 0 {
			g.pool.Release(g.Tree[i])
		} else {
			g.Tree[0].supportingEdges = 0
		}
		g.Tree[i] = nil
	}
	g.Count = 0
	if g.supportEnabled {
		for _, pair := range g.Support {
			g.pool.Release(pair)
		}
		g.Support = make([]*MinutiaPair, 0)
	}
}
