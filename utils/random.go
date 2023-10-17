package utils

type OrientationRandom struct {
	PRIME   int
	BITS    int
	MASK    int
	SCALING float64
	state   int64
}

func NewOrientationRandom() *OrientationRandom {
	return &OrientationRandom{
		PRIME:   1610612741,
		BITS:    30,
		MASK:    (1 << 30) - 1,
		SCALING: 1.0 / float64(1<<30),
		state:   int64(2017612753857937533), // Specify the value directly as int64
	}
}

func (or *OrientationRandom) Float64() float64 {
	or.state *= int64(or.PRIME)
	return (float64((or.state & int64(or.MASK))) + 0.5) * or.SCALING
}
