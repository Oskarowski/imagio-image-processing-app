package analysis

import (
	"image"
	"math"
)

func pixelSquaredDifference(img1, img2 image.Image) (totalR, totalG, totalB float64) {
	bounds := img1.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {

			r1, g1, b1, _ := img1.At(x, y).RGBA()
			r2, g2, b2, _ := img2.At(x, y).RGBA()

			diffR := float64(r1>>8) - float64(r2>>8)
			diffG := float64(g1>>8) - float64(g2>>8)
			diffB := float64(b1>>8) - float64(b2>>8)

			totalR += diffR * diffR
			totalG += diffG * diffG
			totalB += diffB * diffB
		}
	}

	return totalR, totalG, totalB
}

func MeanSquareError(img1, img2 image.Image) float64 {
	totalR, totalG, totalB := pixelSquaredDifference(img1, img2)

	pixelsInTotal := float64(img1.Bounds().Dx() * img1.Bounds().Dy())

	return (totalR + totalG + totalB) / (3 * pixelsInTotal)
}

// Helper function to get the maximum pixel value in an image.
func maxPixelValue(img image.Image) int {
	maxVal := 0
	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()

			maxVal = max(maxVal, int(r>>8))
			maxVal = max(maxVal, int(g>>8))
			maxVal = max(maxVal, int(b>>8))
		}
	}

	return maxVal
}

func PeakMeanSquareError(img1, img2 image.Image) float64 {
	bounds := img1.Bounds()
	totalError := 0.0

	// calculate the max pixel value in the original image rather than hardcoding it to 255
	maxVal := float64(maxPixelValue(img1))
	maxValSquared := maxVal * maxVal

	// Loop over all pixels to calculate the squared error normalized by the max pixel value squared.
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r1, g1, b1, _ := img1.At(x, y).RGBA()
			r2, g2, b2, _ := img2.At(x, y).RGBA()

			diffR := float64(int(r1>>8) - int(r2>>8))
			diffG := float64(int(g1>>8) - int(g2>>8))
			diffB := float64(int(b1>>8) - int(b2>>8))

			// Sum of squared differences, normalized by the max value squared.
			totalError += (diffR * diffR) / maxValSquared
			totalError += (diffG * diffG) / maxValSquared
			totalError += (diffB * diffB) / maxValSquared
		}
	}

	pixelsInTotal := float64(img1.Bounds().Dx() * img1.Bounds().Dy())

	return totalError / (3 * pixelsInTotal)
}

func SignalToNoiseRatio(img1, img2 image.Image) float64 {
	bounds := img1.Bounds()

	var signalSum, noiseSum float64

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r1, g1, b1, _ := img1.At(x, y).RGBA()
			r2, g2, b2, _ := img2.At(x, y).RGBA()

			fr1 := float64(r1 >> 8)
			fg1 := float64(g1 >> 8)
			fb1 := float64(b1 >> 8)

			fr2 := float64(r2 >> 8)
			fg2 := float64(g2 >> 8)
			fb2 := float64(b2 >> 8)

			signalSum += math.Pow(fr1, 2) + math.Pow(fg1, 2) + math.Pow(fb1, 2)
			noiseSum += math.Pow(fr1-fr2, 2) + math.Pow(fg1-fg2, 2) + math.Pow(fb1-fb2, 2)
		}
	}

	if signalSum == 0 {
		return math.Inf(1)
	}

	return 10 * math.Log10(signalSum/noiseSum)
}

func PeakSignalToNoiseRatio(img1, img2 image.Image) float64 {
	maxValue := 255.0
	mseValue := MeanSquareError(img1, img2)

	if mseValue == 0 {
		return math.Inf(1)
	}

	return 10 * math.Log10(maxValue*maxValue/mseValue)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func MaxDifference(img1, img2 image.Image) int {
	bounds := img1.Bounds()

	var maxDiff int

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r1, g1, b1, _ := img1.At(x, y).RGBA()
			r2, g2, b2, _ := img2.At(x, y).RGBA()

			fr1 := int(r1 >> 8)
			fg1 := int(g1 >> 8)
			fb1 := int(b1 >> 8)

			fr2 := int(r2 >> 8)
			fg2 := int(g2 >> 8)
			fb2 := int(b2 >> 8)

			diffR := abs(fr1 - fr2)
			diffG := abs(fg1 - fg2)
			diffB := abs(fb1 - fb2)

			maxDiffForPixel := max(diffR, max(diffG, diffB))

			if maxDiffForPixel > maxDiff {
				maxDiff = maxDiffForPixel
			}
		}
	}

	return maxDiff
}
