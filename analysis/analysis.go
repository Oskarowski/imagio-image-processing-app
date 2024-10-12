package analysis

import (
	"image"
	"math"
)

func pixelDifference(img1, img2 image.Image) (diffR, diffG, diffB float64) {
	bounds := img1.Bounds()

	var totalR, totalG, totalB float64

	numPixels := bounds.Dx() * bounds.Dy()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r1, g1, b1, _ := img1.At(x, y).RGBA()
			r2, g2, b2, _ := img2.At(x, y).RGBA()

			totalR += float64(r1>>8) - float64(r2>>8)
			totalG += float64(g1>>8) - float64(g2>>8)
			totalB += float64(b1>>8) - float64(b2>>8)
		}
	}

	diffR = totalR / float64(numPixels)
	diffG = totalG / float64(numPixels)
	diffB = totalB / float64(numPixels)

	return diffR, diffG, diffB
}

func MeanSquareError(img1, img2 image.Image) float64 {
	diffR, diffG, diffB := pixelDifference(img1, img2)

	return (diffR + diffG + diffB) / 3
}

func PeakMeanSquareError(img1, img2 image.Image) float64 {
	maxVal := 255.0

	msaValue := MeanSquareError(img1, img2)

	return msaValue / (maxVal * maxVal)
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

func MaxDifference(img1, img2 image.Image) float64 {
	bounds := img1.Bounds()

	var maxDiff float64

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

			diffR := math.Abs(fr1 - fr2)
			diffG := math.Abs(fg1 - fg2)
			diffB := math.Abs(fb1 - fb2)

			maxDiffForPixel := math.Max(diffR, math.Max(diffG, diffB))
			maxDiff = math.Max(maxDiff, maxDiffForPixel)
		}
	}

	return maxDiff
}
