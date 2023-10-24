package matcher

import (
	"container/heap"
	"context"
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
	old[n-1] = nil
	*pq = old[0 : n-1]
	return item
}

type matcherThreadContextKey struct{}

type MatcherThread struct {
	Roots   *RootList
	Pairing *PairingGraph
	Queue   PriorityQueue
	Score   *ScoringData
}

func CurrentThread(ctx context.Context) *MatcherThread {
	if thread, ok := ctx.Value(matcherThreadContextKey{}).(*MatcherThread); ok {
		return thread
	}
	return createAndStoreMatcherThread(ctx)
}

func createAndStoreMatcherThread(ctx context.Context) *MatcherThread {
	pool := NewMinutiaPairPool()
	thread := &MatcherThread{
		Roots:   NewRootList(pool),
		Pairing: NewPairingGraph(pool),
		Queue:   make(PriorityQueue, 0),
		Score:   new(ScoringData),
	}
	heap.Init(&thread.Queue)
	ctx = context.WithValue(ctx, matcherThreadContextKey{}, thread)
	return thread
}
