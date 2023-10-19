package knot

import (
	"sourceafis/extractor/skeletons/filters/dot"
	"sourceafis/features"
)

func Apply(skeleton *features.Skeleton) error {
	for _, minutia := range skeleton.Minutiae {
		if len(minutia.Ridges) == 2 && minutia.Ridges[0].Reversed != minutia.Ridges[1] {
			extended := minutia.Ridges[0].Reversed
			removed := minutia.Ridges[1]
			if extended.Points.Size() < removed.Points.Size() {
				tmp := extended
				extended = removed
				removed = tmp
				extended = extended.Reversed
				removed = removed.Reversed
			}
			if err := extended.Points.Remove(extended.Points.Size() - 1); err != nil {
				return err
			}
			it := removed.Points.Iterator()
			for it.HasNext() {
				point, err := it.Next()
				if err != nil {
					return err
				}
				if err := extended.Points.Add(point); err != nil {
					return err
				}
			}
			extended.SetEnd(removed.End())
			removed.Detach()
		}
	}

	return dot.Apply(skeleton)
}
