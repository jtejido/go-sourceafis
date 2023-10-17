package matcher

import (
	"sourceafis/config"
	"sourceafis/features"
	"sourceafis/templates"
)

func Enumerate(probe *Probe, candidate *templates.SearchTemplate, roots *RootList) {
	cminutiae := candidate.Minutiae
	var lookups, tried int
	b := []bool{false, true}
	for _, shortEdges := range b {
		for period := 1; period < len(cminutiae); period++ {
			for phase := 0; phase <= period; phase++ {
				for creference := phase; creference < len(cminutiae); creference += period + 1 {
					cneighbor := (creference + period) % len(cminutiae)
					cedge := features.NewEdgeShape(cminutiae[creference], cminutiae[cneighbor])
					if (cedge.Length >= config.Config.MinRootEdgeLength) != shortEdges {
						matches := probe.hash[Hash(cedge)]
						if matches != nil {
							for _, match := range matches {
								if Matching(match, cedge) {
									duplicateKey := (match.Reference() << 16) | creference
									roots.duplicates.Add(duplicateKey)
									pair := roots.pool.Allocate()
									pair.Probe = match.Reference()
									pair.Candidate = creference
									roots.Add(pair)

									tried++
									if tried >= config.Config.MaxTriedRoots {
										return
									}

								}
							}
						}
						lookups++
						if lookups >= config.Config.MaxRootEdgeLookups {
							return
						}
					}
				}
			}
		}
	}
}
