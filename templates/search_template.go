package templates

import (
	"sort"

	"github.com/jtejido/sourceafis/features"
)

const PRIME = 1610612741

type SearchTemplate struct {
	Width, Height int
	Minutiae      []*features.SearchMinutia
	Edges         [][]*features.NeighborEdge
}

func NewSearchTemplate(logger features.EdgeTableLogger, feat *FeatureTemplate) *SearchTemplate {
	t := new(SearchTemplate)
	t.Width = feat.Size.X
	t.Height = feat.Size.Y
	minutiae := make([]*features.SearchMinutia, feat.Minutiae.Len())

	var i int
	for temp := feat.Minutiae.Front(); temp != nil; temp = temp.Next() {
		minutia := temp.Value.(*features.FeatureMinutia)
		minutiae[i] = features.NewSearchMinutia(minutia)
		i++
	}

	sort.Slice(minutiae, func(i, j int) bool {
		return (minutiae[i].X*PRIME+minutiae[i].Y)*PRIME < (minutiae[j].X*PRIME+minutiae[j].Y)*PRIME ||
			minutiae[i].X < minutiae[j].X ||
			minutiae[i].Y < minutiae[j].Y ||
			minutiae[i].Direction < minutiae[j].Direction ||
			minutiae[i].T < minutiae[j].T
	})
	t.Minutiae = minutiae
	b := features.NewNeighborhoodBuilder(logger)
	t.Edges = b.Build(minutiae)
	return t
}
