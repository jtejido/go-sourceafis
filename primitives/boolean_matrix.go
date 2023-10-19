package primitives

import (
	"fmt"
	"sync"
)

type BooleanMatrix struct {
	sync.RWMutex
	Width, Height int
	Cells         []bool
}

func NewBooleanMatrix(width, height int) *BooleanMatrix {
	return &BooleanMatrix{
		Width:  width,
		Height: height,
		Cells:  make([]bool, width*height),
	}
}

func NewBooleanMatrixFromPoint(size IntPoint) *BooleanMatrix {
	return NewBooleanMatrix(size.X, size.Y)
}

func NewBooleanMatrixFromBooleanMatrix(other *BooleanMatrix) *BooleanMatrix {
	m := NewBooleanMatrixFromPoint(other.Size())
	for i := 0; i < len(m.Cells); i++ {
		m.Cells[i] = other.Cells[i]
	}

	return m
}

func (m *BooleanMatrix) BlockPoint(blockX, blockY, blockSize int) (int, int) {
	if blockX < 0 || blockY < 0 || blockX*blockSize >= m.Width || blockY*blockSize >= m.Height {
		// Handle out-of-bounds blocks or invalid input
		return -1, -1
	}
	x := blockX * blockSize
	y := blockY * blockSize
	return x, y
}

func (m *BooleanMatrix) Get(x, y int) bool {
	m.RLock()
	defer m.RUnlock()
	return m.Cells[m.offset(x, y)]
}
func (m *BooleanMatrix) GetPoint(at IntPoint) bool {
	return m.Get(at.X, at.Y)
}
func (m *BooleanMatrix) GetWithFallback(x, y int, fallback bool) bool {
	m.RLock()
	defer m.RUnlock()
	if x < 0 || y < 0 || x >= m.Width || y >= m.Height {
		return fallback
	}

	return m.Cells[m.offset(x, y)]
}
func (m *BooleanMatrix) GetPointWithFallback(at IntPoint, fallback bool) bool {
	return m.GetWithFallback(at.X, at.Y, fallback)
}
func (m *BooleanMatrix) Set(x, y int, value bool) {
	m.Lock()
	defer m.Unlock()
	m.Cells[m.offset(x, y)] = value
}
func (m *BooleanMatrix) SetPoint(at IntPoint, value bool) {
	m.Set(at.X, at.Y, value)
}

func (m *BooleanMatrix) Invert() {
	m.Lock()
	defer m.Unlock()
	for i := 0; i < len(m.Cells); i++ {
		m.Cells[i] = !m.Cells[i]
	}
}
func (m *BooleanMatrix) Merge(other *BooleanMatrix) error {
	m.Lock()
	defer m.Unlock()
	if other.Width != m.Width || other.Height != m.Height {
		return fmt.Errorf("unable to merge.")
	}
	for i := 0; i < len(m.Cells); i++ {
		m.Cells[i] = m.Cells[i] || other.Cells[i]
	}

	return nil
}

func (m *BooleanMatrix) offset(x, y int) int {
	return y*m.Width + x
}

func (m *BooleanMatrix) Size() IntPoint {
	return IntPoint{m.Width, m.Height}
}
