package cloud

import (
	"github.com/jtejido/sourceafis/config"
	"github.com/jtejido/sourceafis/features"
	"github.com/jtejido/sourceafis/primitives"
	"github.com/jtejido/sourceafis/utils"
)

func Apply(minutiae *primitives.GenericList[*features.FeatureMinutia]) {
	radiusSq := utils.SquareInt(config.Config.MinutiaCloudRadius)

	// Create a map to count the number of neighbors within the radius for each minutia
	neighborCount := make(map[*features.FeatureMinutia]int)

	for e := minutiae.Front(); e != nil; e = e.Next() {
		minutia := e.Value.(*features.FeatureMinutia)
		neighborCount[minutia] = 0

		for e2 := minutiae.Front(); e2 != nil; e2 = e2.Next() {
			neighbor := e2.Value.(*features.FeatureMinutia)
			if neighbor != minutia {
				if neighbor.Position.Minus(minutia.Position).LengthSq() <= radiusSq {
					neighborCount[minutia]++
				}
			}
		}
	}

	// Remove minutiae that don't meet the criteria
	for e := minutiae.Front(); e != nil; {
		minutia := e.Value.(*features.FeatureMinutia)
		if neighborCount[minutia] >= config.Config.MaxCloudSize {
			next := e.Next()
			minutiae.Remove(e)
			e = next
		} else {
			e = e.Next()
		}
	}
}
