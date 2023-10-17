package primitives

type FloatPointMatrix struct {
	Width, Height int
	Vectors       []float64
}

func NewFloatPointMatrix(width, height int) *FloatPointMatrix {
	return &FloatPointMatrix{
		Width:   width,
		Height:  height,
		Vectors: make([]float64, 2*width*height),
	}
}

func NewFloatPointMatrixFromPoint(size IntPoint) *FloatPointMatrix {
	return NewFloatPointMatrix(size.X, size.Y)
}

func (m *FloatPointMatrix) Size() IntPoint {
	return IntPoint{X: m.Width, Y: m.Height}
}

func (m *FloatPointMatrix) Get(x, y int) FloatPoint {
	i := m.offset(x, y)
	return FloatPoint{X: m.Vectors[i], Y: m.Vectors[i+1]}
}

func (m *FloatPointMatrix) GetPoint(at IntPoint) FloatPoint {
	return m.Get(at.X, at.Y)
}

func (m *FloatPointMatrix) offset(x, y int) int {
	return 2 * (y*m.Width + x)
}

func (m *FloatPointMatrix) Set(x, y int, px, py float64) {
	i := m.offset(x, y)
	m.Vectors[i] = px
	m.Vectors[i+1] = py
}
func (m *FloatPointMatrix) SetFloatPoint(x, y int, point FloatPoint) {
	m.Set(x, y, point.X, point.Y)
}
func (m *FloatPointMatrix) SetPoint(at IntPoint, point FloatPoint) {
	m.SetFloatPoint(at.X, at.Y, point)
}
func (m *FloatPointMatrix) Add(x, y int, px, py float64) {
	i := m.offset(x, y)
	m.Vectors[i] += px
	m.Vectors[i+1] += py
}
func (m *FloatPointMatrix) AddFloatPoint(x, y int, point FloatPoint) {
	m.Add(x, y, point.X, point.Y)
}
func (m *FloatPointMatrix) AddPoint(at IntPoint, point FloatPoint) {
	m.AddFloatPoint(at.X, at.Y, point)
}
