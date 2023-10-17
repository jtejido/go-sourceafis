package matcher

type ScoringData struct {
	MinutiaCount               int
	MinutiaScore               float64
	MinutiaFractionInProbe     float64
	MinutiaFractionInCandidate float64
	MinutiaFraction            float64
	MinutiaFractionScore       float64
	SupportingEdgeSum          int
	EdgeCount                  int
	EdgeScore                  float64
	SupportedMinutiaCount      int
	SupportedMinutiaScore      float64
	MinutiaTypeHits            int
	MinutiaTypeScore           float64
	DistanceErrorSum           int
	DistanceAccuracySum        int
	DistanceAccuracyScore      float64
	AngleErrorSum              float64
	AngleAccuracySum           float64
	AngleAccuracyScore         float64
	TotalScore                 float64
	ShapedScore                float64
}
