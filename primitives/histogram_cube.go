package primitives

import "math"

type HistogramCube struct {
	Width, Height, Bins int
	Counts              []int
}

func NewHistogramCube(width, height, bins int) *HistogramCube {
	return &HistogramCube{
		Width:  width,
		Height: height,
		Bins:   bins,
		Counts: make([]int, width*height*bins),
	}
}

func NewHistogramCubeFromPoint(size IntPoint, bins int) *HistogramCube {
	return NewHistogramCube(size.X, size.Y, bins)
}

func (h *HistogramCube) Constrain(z int) int {
	return int(math.Max(0, math.Min(float64(h.Bins-1), float64(z))))
}

func (h *HistogramCube) Get(x, y, z int) int {
	return h.Counts[h.offset(x, y, z)]
}
func (h *HistogramCube) GetPoint(at IntPoint, z int) int {
	return h.Get(at.X, at.Y, z)
}

func (h *HistogramCube) Set(x, y, z, value int) {
	h.Counts[h.offset(x, y, z)] = value
}

func (h *HistogramCube) SetPoint(at IntPoint, z, value int) {
	h.Set(at.X, at.Y, z, value)
}

func (h *HistogramCube) Add(x, y, z, value int) {
	h.Counts[h.offset(x, y, z)] += value
}

func (h *HistogramCube) AddPoint(at IntPoint, z, value int) {
	h.Add(at.X, at.Y, z, value)
}

func (h *HistogramCube) Increment(x, y, z int) {
	h.Add(x, y, z, 1)
}
func (h *HistogramCube) IncrementPoint(at IntPoint, z int) {
	h.Increment(at.X, at.Y, z)
}
func (h *HistogramCube) Sum(x, y int) int {
	var sum int
	for i := 0; i < h.Bins; i++ {
		sum += h.Get(x, y, i)
	}
	return sum
}
func (h *HistogramCube) SumPoint(at IntPoint) int {
	return h.Sum(at.X, at.Y)
}

func (h *HistogramCube) offset(x, y, z int) int {
	return (y*h.Width+x)*h.Bins + z
}

func (histogram *HistogramCube) Merge(other *HistogramCube) {
	if histogram.Width != other.Width || histogram.Height != other.Height || histogram.Bins != other.Bins {
		return
	}

	for i := 0; i < len(histogram.Counts); i++ {
		histogram.Counts[i] += other.Counts[i]
	}
}
