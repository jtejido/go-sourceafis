package smoothing

import (
	"sourceafis/config"
	"sourceafis/extractor/logger"
	"sourceafis/primitives"
	"sync"
)

type OrientedSmoothing struct {
	logger logger.TransparencyLogger
}

func New(logger logger.TransparencyLogger) *OrientedSmoothing {
	return &OrientedSmoothing{
		logger: logger,
	}
}

func (s *OrientedSmoothing) Parallel(input, orientation *primitives.Matrix, mask *primitives.BooleanMatrix, blocks *primitives.BlockMap) (*primitives.Matrix, error) {
	lines := lines(config.Config.ParallelSmoothingResolution, config.Config.ParallelSmoothingRadius, config.Config.ParallelSmoothingStep)
	smoothed := smooth(input, orientation, mask, blocks, 0, lines)
	return smoothed, s.logger.Log("parallel-smoothing", smoothed)
}

func (s *OrientedSmoothing) Orthogonal(input, orientation *primitives.Matrix, mask *primitives.BooleanMatrix, blocks *primitives.BlockMap) (*primitives.Matrix, error) {
	lines := lines(config.Config.OrthogonalSmoothingResolution, config.Config.OrthogonalSmoothingRadius, config.Config.OrthogonalSmoothingStep)
	smoothed := smooth(input, orientation, mask, blocks, primitives.Pi, lines)
	return smoothed, s.logger.Log("orthogonal-smoothing", smoothed)
}

func lines(resolution, radius int, step float64) [][]primitives.IntPoint {
	result := make([][]primitives.IntPoint, resolution)
	var wg sync.WaitGroup
	var mu sync.Mutex

	numCPU := config.Config.Workers

	chunkSize := resolution / numCPU
	remainder := resolution % numCPU

	for i := 0; i < numCPU; i++ {
		wg.Add(1)
		startIndex := i * chunkSize
		endIndex := (i + 1) * chunkSize
		if i == numCPU-1 {
			endIndex += remainder
		}
		go func(start, end int) {
			defer wg.Done()
			for orientationIndex := start; orientationIndex < end; orientationIndex++ {
				line := computeLine(resolution, orientationIndex, radius, step)
				mu.Lock()
				result[orientationIndex] = line
				mu.Unlock()
			}
		}(startIndex, endIndex)
	}
	wg.Wait()
	return result
}

func computeLine(resolution, orientationIndex, radius int, step float64) []primitives.IntPoint {
	line := []primitives.IntPoint{primitives.ZeroIntPoint()}
	direction := primitives.BucketCenter(orientationIndex, resolution).FromOrientation().ToVector()
	for r := float64(radius); r >= 0.5; r /= step {
		sample := direction.Multiply(r).Round()
		var isFound bool
		for _, samp := range line {
			if samp.Equals(sample) {
				isFound = true
			}
		}
		if !isFound {
			line = append(line, sample)
			line = append(line, sample.Negate())
		}
	}
	return line
}

func smooth(input, orientation *primitives.Matrix, mask *primitives.BooleanMatrix, blocks *primitives.BlockMap, angle float64, lines [][]primitives.IntPoint) *primitives.Matrix {
	output := primitives.NewMatrixFromPoint(input.Size())
	it := blocks.Primary.Blocks.Iterator()
	var wg sync.WaitGroup
	ch := make(chan primitives.IntPoint)

	for i := 0; i < config.Config.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for block := range ch {
				processBlock(input, orientation, lines, block, output, blocks, angle)
			}
		}()
	}

	for it.HasNext() {
		block := it.Next()
		if mask.GetPoint(block) {
			ch <- block
		}
	}

	close(ch)
	wg.Wait()
	return output
}

func processBlock(input, orientation *primitives.Matrix, lines [][]primitives.IntPoint, block primitives.IntPoint, output *primitives.Matrix, blocks *primitives.BlockMap, angle float64) {
	line := lines[primitives.AngleAdd(orientation.GetPoint(block), angle).Quantize(len(lines))]

	target := blocks.Primary.BlockPoint(block)
	blockArea := blocks.Primary.BlockPoint(block)

	for _, linePoint := range line {
		source := target.Move(linePoint).Intersect(primitives.IntRect{
			X:      0,
			Y:      0,
			Width:  blocks.Pixels.X,
			Height: blocks.Pixels.Y,
		})
		target = source.Move(linePoint.Negate())
		for y := target.Top(); y < target.Bottom(); y++ {
			for x := target.Left(); x < target.Right(); x++ {
				output.Add(x, y, input.Get(x+linePoint.X, y+linePoint.Y))
			}
		}
	}

	for y := blockArea.Top(); y < blockArea.Bottom(); y++ {
		for x := blockArea.Left(); x < blockArea.Right(); x++ {
			output.Multiply(x, y, 1.0/float64(len(line)))
		}
	}
}
