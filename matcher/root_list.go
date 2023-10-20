package matcher

import (
	"github.com/jtejido/sourceafis/config"

	"github.com/emirpasic/gods/sets/hashset"
)

type RootList struct {
	pool       *MinutiaPairPool
	count      int
	pairs      []*MinutiaPair
	duplicates *hashset.Set
}

func NewRootList(pool *MinutiaPairPool) *RootList {
	return &RootList{
		pool:       pool,
		pairs:      make([]*MinutiaPair, config.Config.MaxTriedRoots),
		duplicates: hashset.New(),
	}
}

func (r *RootList) Add(pair *MinutiaPair) {
	r.pairs[r.count] = pair
	r.count++
}
func (r *RootList) Discard() {
	for i := 0; i < r.count; i++ {
		r.pool.Release(r.pairs[i])
		r.pairs[i] = nil
	}
	r.count = 0
	r.duplicates.Clear()
}
