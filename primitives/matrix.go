package primitives

type Matrix struct {
	Width, Height int
	Cells         []float64
}

func NewMatrix(width, height int) *Matrix {
	return &Matrix{
		Width:  width,
		Height: height,
		Cells:  make([]float64, width*height),
	}
}

func NewMatrixFromPoint(size IntPoint) *Matrix {
	return NewMatrix(size.X, size.Y)
}

func (m *Matrix) Get(x, y int) float64 {
	return m.Cells[m.offset(x, y)]
}

func (m *Matrix) Add(x, y int, value float64) {
	m.Cells[m.offset(x, y)] += value
}

func (m *Matrix) AddPoint(at IntPoint, value float64) {
	m.Add(at.X, at.Y, value)
}

func (m *Matrix) Multiply(x, y int, value float64) {
	m.Cells[m.offset(x, y)] *= value
}

func (m *Matrix) MultiplyPoint(at IntPoint, value float64) {
	m.Multiply(at.X, at.Y, value)
}

func (m *Matrix) GetPoint(at IntPoint) float64 {
	return m.Get(at.X, at.Y)
}

func (m *Matrix) Set(x, y int, value float64) {
	m.Cells[m.offset(x, y)] = value
}
func (m *Matrix) SetPoint(at IntPoint, value float64) {
	m.Set(at.X, at.Y, value)
}

func (m *Matrix) Size() IntPoint {
	return IntPoint{m.Width, m.Height}
}

func (m *Matrix) offset(x, y int) int {
	return y*m.Width + x
}
