package noise

import (
	"image"
	"image/color"
	"runtime"
	"sync"
)

func processPixel(img image.Image, x, y, sMin, sMax int) (int, int, int) {
	windowSize := sMin
	urxy, ugxy, ubxy, _ := img.At(x, y).RGBA()
	rxy, gxy, bxy := int(urxy>>8), int(ugxy>>8), int(ubxy>>8)
	newR, newG, newB := rxy, gxy, bxy

	for windowSize <= sMax {
		reds, greens, blues := getWindowPixelsRGB(img, x, y, windowSize)
		rMin, rMax, rMed := minMaxMedian(reds)
		gMin, gMax, gMed := minMaxMedian(greens)
		bMin, bMax, bMed := minMaxMedian(blues)

		newR = adaptiveChannel(newR, rxy, rMin, rMax, rMed)
		newG = adaptiveChannel(newG, gxy, gMin, gMax, gMed)
		newB = adaptiveChannel(newB, bxy, bMin, bMax, bMed)

		if newR != rxy || newG != gxy || newB != bxy {
			break
		}
		windowSize += 2
	}

	return newR, newG, newB
}

func AdaptiveMedianFilterParallel(img image.Image, sMin, sMax int) *image.RGBA {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	workersNumber := runtime.NumCPU()

	var wg sync.WaitGroup

	pixelChan := make(chan struct{ x, y int }, workersNumber)

	for i := 0; i < workersNumber; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for pixel := range pixelChan {
				x, y := pixel.x, pixel.y
				newR, newG, newB := processPixel(img, x, y, sMin, sMax)
				newImg.Set(x, y, color.RGBA{
					uint8(newR),
					uint8(newG),
					uint8(newB),
					255,
				})
			}
		}()
	}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixelChan <- struct{ x, y int }{x, y}
		}
	}
	close(pixelChan)

	wg.Wait()

	return newImg

}
