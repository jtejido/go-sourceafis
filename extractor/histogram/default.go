package histogram

import (
	"sourceafis/config"
	"sourceafis/extractor/logger"
	"sourceafis/primitives"
	"sync"
)

type LocalHistograms struct {
	logger logger.TransparencyLogger
}

func New(logger logger.TransparencyLogger) *LocalHistograms {
	return &LocalHistograms{
		logger: logger,
	}
}

func (h *LocalHistograms) Create(blocks *primitives.BlockMap, image *primitives.Matrix) (*primitives.HistogramCube, error) {
	var wg sync.WaitGroup
	numWorkers := config.Config.Workers
	resultChan := make(chan *primitives.HistogramCube, numWorkers)
	blockIterator := blocks.Primary.Blocks.Iterator()
	tasks := make(chan primitives.IntPoint, numWorkers)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			localHistogram := primitives.NewHistogramCubeFromPoint(blocks.Primary.Blocks, config.Config.HistogramDepth)

			for block := range tasks {
				area := blocks.Primary.BlockPoint(block)
				for y := area.Top(); y < area.Bottom(); y++ {
					for x := area.Left(); x < area.Right(); x++ {
						depth := int(image.Get(x, y) * float64(localHistogram.Bins))
						localHistogram.IncrementPoint(block, localHistogram.Constrain(depth))
					}
				}
			}

			resultChan <- localHistogram
		}()
	}

	go func() {
		for blockIterator.HasNext() {
			block := blockIterator.Next()
			tasks <- block
		}
		close(tasks)
	}()

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	mergedHistogram := primitives.NewHistogramCubeFromPoint(blocks.Primary.Blocks, config.Config.HistogramDepth)
	for partialHistogram := range resultChan {
		mergedHistogram.Merge(partialHistogram)
	}

	return mergedHistogram, h.logger.Log("histogram", mergedHistogram)
}

func (h *LocalHistograms) Smooth(blocks *primitives.BlockMap, input *primitives.HistogramCube) (*primitives.HistogramCube, error) {
	var wg sync.WaitGroup
	numWorkers := config.Config.Workers
	resultChan := make(chan *primitives.HistogramCube, numWorkers)
	blockIterator := blocks.Secondary.Blocks.Iterator()
	tasks := make(chan primitives.IntPoint, numWorkers)
	blocksAround := []primitives.IntPoint{
		{X: 0, Y: 0},
		{X: -1, Y: 0},
		{X: 0, Y: -1},
		{X: -1, Y: -1},
	}

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			localOutput := primitives.NewHistogramCubeFromPoint(blocks.Secondary.Blocks, input.Bins)
			for corner := range tasks {
				for _, relative := range blocksAround {
					block := corner.Plus(relative)
					if blocks.Primary.Blocks.Contains(block) {
						for k := 0; k < input.Bins; k++ {
							localOutput.AddPoint(corner, k, input.GetPoint(block, k))
						}
					}
				}
			}

			resultChan <- localOutput
		}()
	}

	go func() {
		for blockIterator.HasNext() {
			block := blockIterator.Next()
			tasks <- block
		}
		close(tasks)
	}()

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	finalOutput := primitives.NewHistogramCubeFromPoint(blocks.Secondary.Blocks, input.Bins)
	for i := 0; i < numWorkers; i++ {
		partialOutput := <-resultChan
		finalOutput.Merge(partialOutput)
	}

	return finalOutput, h.logger.Log("smoothed-histogram", finalOutput)
}
