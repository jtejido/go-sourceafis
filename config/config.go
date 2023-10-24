package config

import (
	"github.com/BurntSushi/toml"
	"github.com/mcuadros/go-defaults"
)

// const (
// 	BLOCK_SIZE                         15
// 	HISTOGRAM_DEPTH                    256
// 	CLIPPED_CONTRAST                   0.08
// 	MIN_ABSOLUTE_CONTRAST              17 / 255.0
// 	MIN_RELATIVE_CONTRAST              0.34
// 	RELATIVE_CONTRAST_SAMPLE           168568
// 	RELATIVE_CONTRAST_PERCENTILE       0.49
// 	MASK_VOTE_RADIUS                   7
// 	MASK_VOTE_MAJORITY                 0.51
// 	MASK_VOTE_BORDER_DISTANCE          4
// 	BLOCK_ERRORS_VOTE_RADIUS           1
// 	BLOCK_ERRORS_VOTE_MAJORITY         0.7
// 	BLOCK_ERRORS_VOTE_BORDER_DISTANCE  4
// 	MAX_EQUALIZATION_SCALING           3.99
// 	MIN_EQUALIZATION_SCALING           0.25
// 	MIN_ORIENTATION_RADIUS             2
// 	MAX_ORIENTATION_RADIUS             6
// 	ORIENTATION_SPLIT                  50
// 	ORIENTATIONS_CHECKED               20
// 	ORIENTATION_SMOOTHING_RADIUS       1
// 	PARALLEL_SMOOTHING_RESOLUTION      32
// 	PARALLEL_SMOOTHING_RADIUS          7
// 	PARALLEL_SMOOTHING_STEP            1.59
// 	ORTHOGONAL_SMOOTHING_RESOLUTION    11
// 	ORTHOGONAL_SMOOTHING_RADIUS        4
// 	ORTHOGONAL_SMOOTHING_STEP          1.11
// 	BINARIZED_VOTE_RADIUS              2
// 	BINARIZED_VOTE_MAJORITY            0.61
// 	BINARIZED_VOTE_BORDER_DISTANCE     17
// 	INNER_MASK_BORDER_DISTANCE         14
// 	MASK_DISPLACEMENT                  10.06
// 	MINUTIA_CLOUD_RADIUS               20
// 	MAX_CLOUD_SIZE                     4
// 	MAX_MINUTIAE                       100
// 	SORT_BY_NEIGHBOR                   5
// 	EDGE_TABLE_NEIGHBORS               9
// 	THINNING_ITERATIONS                26
// 	MAX_PORE_ARM                       41
// 	SHORTEST_JOINED_ENDING             7
// 	MAX_RUPTURE_SIZE                   5
// 	MAX_GAP_SIZE                       20
// 	GAP_ANGLE_OFFSET                   22
// 	TOLERATED_GAP_OVERLAP              2
// 	MIN_TAIL_LENGTH                    21
// 	MIN_FRAGMENT_LENGTH                22
// 	MAX_DISTANCE_ERROR                 13
// 	MAX_ANGLE_ERROR                    math.Pi / 180 * 10
// 	MAX_GAP_ANGLE                      45 * math.Pi / 180
// 	RIDGE_DIRECTION_SAMPLE             21
// 	RIDGE_DIRECTION_SKIP               1
// 	MAX_TRIED_ROOTS                    70
// 	MIN_ROOT_EDGE_LENGTH               58
// 	MAX_ROOT_EDGE_LOOKUPS              1633
// 	MIN_SUPPORTING_EDGES               1
// 	DISTANCE_ERROR_FLATNESS            0.69
// 	ANGLE_ERROR_FLATNESS               0.27
// 	MINUTIA_SCORE                      0.032
// 	MINUTIA_FRACTION_SCORE             8.98
// 	MINUTIA_TYPE_SCORE                 0.629
// 	SUPPORTED_MINUTIA_SCORE            0.193
// 	EDGE_SCORE                         0.265
// 	DISTANCE_ACCURACY_SCORE            9.9
// 	ANGLE_ACCURACY_SCORE               2.79
// 	THRESHOLD_FMR_MAX                  8.48
// 	THRESHOLD_FMR_2                    11.12
// 	THRESHOLD_FMR_10                   14.15
// 	THRESHOLD_FMR_100                  18.22
// 	THRESHOLD_FMR_1000                 22.39
// 	THRESHOLD_FMR_10_000               27.24
// 	THRESHOLD_FMR_100_000              32.01
// )

type config struct {
	Workers                       int     `default:"1"`
	BlockSize                     int     `default:"15"`
	HistogramDepth                int     `default:"256"`
	ClippedContrast               float64 `default:"0.08"`
	MinAbsoluteContrast           float64 `default:"0.06666666666"` // 17 / 255.0
	MinRelativeContrast           float64 `default:"0.34"`
	RelativeContrastSample        int     `default:"168568"`
	RelativeContrastPercentile    float64 `default:"0.49"`
	MaskVoteRadius                int     `default:"7"`
	MaskVoteMajority              float64 `default:"0.51"`
	MaskVoteBorderDistance        int     `default:"4"`
	BlockErrorsVoteRadius         int     `default:"1"`
	BlockErrorsVoteMajority       float64 `default:"0.7"`
	BlockErrorsVoteBorderDistance int     `default:"4"`
	MaxEqualizationScaling        float64 `default:"3.99"`
	MinEqualizationScaling        float64 `default:"0.25"`
	MinOrientationRadius          int     `default:"2"`
	MaxOrientationRadius          int     `default:"6"`
	OrientationSplit              int     `default:"50"`
	OrientationsChecked           int     `default:"20"`
	OrientationSmoothingRadius    int     `default:"1"`
	ParallelSmoothingResolution   int     `default:"32"`
	ParallelSmoothingRadius       int     `default:"7"`
	ParallelSmoothingStep         float64 `default:"1.59"`
	OrthogonalSmoothingResolution int     `default:"11"`
	OrthogonalSmoothingRadius     int     `default:"4"`
	OrthogonalSmoothingStep       float64 `default:"1.11"`
	BinarizedVoteRadius           int     `default:"2"`
	BinarizedVoteMajority         float64 `default:"0.61"`
	BinarizedVoteBorderDistance   int     `default:"17"`
	InnerMaskBorderDistance       int     `default:"14"`
	MaskDisplacement              float64 `default:"10.06"`
	MinutiaCloudRadius            int     `default:"20"`
	MaxCloudSize                  int     `default:"4"`
	MaxMinutiae                   int     `default:"100"`
	SortByNeighbor                int     `default:"5"`
	EdgeTableNeighbors            int     `default:"9"`
	ThinningIterations            int     `default:"26"`
	MaxPoreArm                    int     `default:"41"`
	ShortestJoinedEnding          int     `default:"7"`
	MaxRuptureSize                int     `default:"5"`
	MaxGapSize                    int     `default:"20"`
	GapAngleOffset                int     `default:"22"`
	ToleratedGapOverlap           int     `default:"2"`
	MinTailLength                 int     `default:"21"`
	MinFragmentLength             int     `default:"22"`
	MaxDistanceError              int     `default:"13"`
	MaxAngleError                 float64 `default:"0.17453292519"` // "Math.Pi / 180 * 10"
	MaxGapAngle                   float64 `default:"0.78539816339"` // 45 * Math.Pi / 180
	RidgeDirectionSample          int     `default:"21"`
	RidgeDirectionSkip            int     `default:"1"`
	MaxTriedRoots                 int     `default:"70"`
	MinRootEdgeLength             int     `default:"58"`
	MaxRootEdgeLookups            int     `default:"1633"`
	MinSupportingEdges            int     `default:"1"`
	DistanceErrorFlatness         float64 `default:"0.69"`
	AngleErrorFlatness            float64 `default:"0.27"`
	MinutiaScore                  float64 `default:"0.032"`
	MinutiaFractionScore          float64 `default:"8.98"`
	MinutiaTypeScore              float64 `default:"0.629"`
	SupportedMinutiaScore         float64 `default:"0.193"`
	EdgeScore                     float64 `default:"0.265"`
	DistanceAccuracyScore         float64 `default:"9.9"`
	AngleAccuracyScore            float64 `default:"2.79"`
	ThresholdFMRMax               float64 `default:"8.48"`
	ThresholdFMR2                 float64 `default:"11.12"`
	ThresholdFMR10                float64 `default:"14.15"`
	ThresholdFMR100               float64 `default:"18.22"`
	ThresholdFMR1000              float64 `default:"22.39"`
	ThresholdFMR10000             float64 `default:"27.24"`
	ThresholdFMR100000            float64 `default:"32.01"`
}

var Config *config

func LoadConfig(configPath string) error {
	LoadDefaultConfig()
	if configPath == "" {
		return nil
	}
	if _, err := toml.DecodeFile(configPath, Config); err != nil {
		return err
	}

	return nil
}

func LoadDefaultConfig() {
	c := new(config)
	defaults.SetDefaults(c)
	Config = c
}
