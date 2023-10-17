package features

import (
	"math"
	"math/bits"
	"sourceafis/primitives"
	"sourceafis/utils"
)

const (
	POLAR_CACHE_BITS   = 8
	POLAR_CACHE_RADIUS = 1 << POLAR_CACHE_BITS
)

var (
	POLAR_DISTANCE_CACHE = make([]int, POLAR_CACHE_RADIUS*POLAR_CACHE_RADIUS)
	POLAR_ANGLE_CACHE    = make([]float64, POLAR_CACHE_RADIUS*POLAR_CACHE_RADIUS)
)

func init() {
	for y := 0; y < POLAR_CACHE_RADIUS; y++ {
		for x := 0; x < POLAR_CACHE_RADIUS; x++ {
			v := math.Sqrt(float64(utils.SquareInt(x) + utils.SquareInt(y)))
			POLAR_DISTANCE_CACHE[y*POLAR_CACHE_RADIUS+x] = int(math.Floor(v + .5))
			if y > 0 || x > 0 {
				POLAR_ANGLE_CACHE[y*POLAR_CACHE_RADIUS+x] = float64(primitives.AtanFromFloatPointVector(primitives.FloatPoint{X: float64(x), Y: float64(y)}))
			} else {
				POLAR_ANGLE_CACHE[y*POLAR_CACHE_RADIUS+x] = 0
			}
		}
	}

}

type EdgeShape struct {
	Length                        int
	ReferenceAngle, NeighborAngle primitives.FloatAngle
}

func NewEdgeShape(reference, neighbor *SearchMinutia) *EdgeShape {
	e := new(EdgeShape)
	var quadrant float64
	x := neighbor.X - reference.X
	y := neighbor.Y - reference.Y
	if y < 0 {
		x = -x
		y = -y
		quadrant = primitives.Pi
	}
	if x < 0 {
		tmp := -x
		x = y
		y = tmp
		quadrant += primitives.HalfPi
	}

	shift := 32 - bits.LeadingZeros32(uint32(x|y)>>POLAR_CACHE_BITS)

	offset := (y>>shift)*POLAR_CACHE_RADIUS + (x >> shift)
	e.Length = (POLAR_DISTANCE_CACHE[offset] << shift)
	angle := POLAR_ANGLE_CACHE[offset] + quadrant
	e.ReferenceAngle = primitives.FloatAngle(reference.Direction).Difference(primitives.FloatAngle(angle))
	e.NeighborAngle = primitives.FloatAngle(neighbor.Direction).Difference(primitives.FloatAngle(angle).Opposite())

	return e
}
