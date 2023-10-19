package primitives

import "sync"

type IntMatrix struct {
	sync.RWMutex
	Width, Height int
	cells         []int
}

func NewIntMatrix(width, height int) *IntMatrix {
	return &IntMatrix{
		Width:  width,
		Height: height,
		cells:  make([]int, width*height),
	}
}

func NewIntMatrixFromPoint(size IntPoint) *IntMatrix {
	return NewIntMatrix(size.X, size.Y)
}

func (m *IntMatrix) Get(x, y int) int {
	m.RLock()
	defer m.RUnlock()
	return m.cells[m.offset(x, y)]
}

func (m *IntMatrix) GetPoint(at IntPoint) int {
	return m.Get(at.X, at.Y)
}

func (m *IntMatrix) Set(x, y, value int) {
	m.Lock()
	defer m.Unlock()
	m.cells[m.offset(x, y)] = value
}

func (m *IntMatrix) Size() IntPoint {
	return IntPoint{m.Width, m.Height}
}

func (m *IntMatrix) offset(x, y int) int {
	return y*m.Width + x
}
