package utils

const (
	PRIME   = 1610612741
	BITS    = 30
	MASK    = (1 << BITS) - 1
	SCALING = 1.0 / float64(1<<BITS)
)

func discardHighOrderBits(value int64) int64 {
	bitmask := (1 << 32) - 1
	result := value & int64(bitmask)

	return result
}

type OrientationRandom struct {
	state int64
}

func NewOrientationRandom() *OrientationRandom {
	return &OrientationRandom{
		state: discardHighOrderBits(PRIME) * discardHighOrderBits(PRIME) * discardHighOrderBits(PRIME),
	}
}

func (or *OrientationRandom) Float64() float64 {
	or.state *= PRIME
	return (float64(or.state&MASK) + 0.5) * SCALING
}
