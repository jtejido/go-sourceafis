package matcher

import (
	"container/heap"
	"math/rand"
	"sync"
)

var threads = sync.Map{}

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
	Queue   PriorityQueue
	Score   *ScoringData
}

func CurrentThread() *MatcherThread {
	thread, ok := threads.Load(goroutineID())
	if !ok {
		thread = createAndStoreMatcherThread()
	}
	return thread.(*MatcherThread)
}

func kill() {
	threads.Delete(goroutineID())
}

func goroutineID() int64 {
	return rand.Int63()
}

func createAndStoreMatcherThread() *MatcherThread {
	pool := NewMinutiaPairPool()
	thread := &MatcherThread{
		Roots:   NewRootList(pool),
		Pairing: NewPairingGraph(pool),
		Queue:   make(PriorityQueue, 0),
		Score:   new(ScoringData),
	}
	heap.Init(&thread.Queue)
	threads.Store(goroutineID(), thread)
	return thread
}
