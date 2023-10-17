package primitives

import (
	"sourceafis/utils"
)

type BlockMap struct {
	Pixels             IntPoint
	Primary, Secondary *BlockGrid
}

func NewBlockMap(width, height, maxBlockSize int) *BlockMap {
	pixels := IntPoint{X: width, Y: height}
	primary := NewBlockGrid(IntPoint{
		X: utils.RoundUpDiv(pixels.X, maxBlockSize),
		Y: utils.RoundUpDiv(pixels.Y, maxBlockSize),
	})
	for y := 0; y <= primary.Blocks.Y; y++ {
		primary.Y[y] = y * pixels.Y / primary.Blocks.Y
	}
	for x := 0; x <= primary.Blocks.X; x++ {
		primary.X[x] = x * pixels.X / primary.Blocks.X
	}
	secondary := NewBlockGrid(primary.Corners)
	secondary.Y[0] = 0
	for y := 0; y < primary.Blocks.Y; y++ {
		secondary.Y[y+1] = primary.Block(0, y).Center().Y
	}
	secondary.Y[secondary.Blocks.Y] = pixels.Y
	secondary.X[0] = 0
	for x := 0; x < primary.Blocks.X; x++ {
		secondary.X[x+1] = primary.Block(x, 0).Center().X
	}
	secondary.X[secondary.Blocks.X] = pixels.X

	return &BlockMap{
		Pixels:    pixels,
		Primary:   primary,
		Secondary: secondary,
	}
}
