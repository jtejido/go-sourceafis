package matcher

import "container/heap"

var (
	current *MatcherThread
)

type PriorityQueue []*MinutiaPair

func (pq PriorityQueue) Len() int { return len(pq) }

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].distance < pq[j].distance
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x any) {
	item := x.(*MinutiaPair)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	*pq = old[0 : n-1]
	return item
}

type MatcherThread struct {
	Roots   *RootList
	Pairing *PairingGraph
	Queue   *PriorityQueue
	Score   *ScoringData
}

func NewMatcherThread(pool *MinutiaPairPool) *MatcherThread {
	pq := make(PriorityQueue, 0)
	heap.Init(&pq)
	return &MatcherThread{
		Roots:   NewRootList(pool),
		Pairing: NewPairingGraph(pool),
		Queue:   &pq,
		Score:   new(ScoringData),
	}
}

func CurrentThread() *MatcherThread {
	if current == nil {
		current = NewMatcherThread(NewMinutiaPairPool())
	}

	return current
}
