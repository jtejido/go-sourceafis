package primitives

type BlockGrid struct {
	Blocks, Corners IntPoint
	X, Y            []int
}

func NewBlockGrid(size IntPoint) *BlockGrid {
	return &BlockGrid{
		Blocks:  size,
		Corners: IntPoint{X: size.X + 1, Y: size.Y + 1},
		X:       make([]int, size.X+1),
		Y:       make([]int, size.Y+1),
	}
}

func (b *BlockGrid) Corner(atX, atY int) IntPoint {
	return IntPoint{X: b.X[atX], Y: b.Y[atY]}
}

func (b *BlockGrid) Block(atX, atY int) IntRect {
	return IntRectBetweenIntPoint(b.Corner(atX, atY), b.Corner(atX+1, atY+1))
}

func (b *BlockGrid) BlockPoint(at IntPoint) IntRect {
	return b.Block(at.X, at.Y)
}
