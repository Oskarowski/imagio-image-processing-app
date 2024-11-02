package manipulations

import (
	"image"
	"image/color"
	"image/draw"
	"math"
)

func toGrayscale(img image.Image) *image.Gray {
	bounds := img.Bounds()
	grayImg := image.NewGray(bounds)
	draw.Draw(grayImg, bounds, img, bounds.Min, draw.Src)
	return grayImg
}

func calculateHistogram(img *image.Gray) []int {
	bounds := img.Bounds()
	histogram := make([]int, 256)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			gray := img.GrayAt(x, y).Y
			histogram[gray]++
		}
	}
	return histogram
}

func rayleighTransform(hist []int, fmin, fmax, gmin, gmax, totalPixels int) []float64 {
	transformed := make([]float64, 256)

	cumulativeHist := make([]int, 256)
	cumulativeSum := 0
	for f := 0; f < len(hist); f++ {
		cumulativeSum += hist[f]
		cumulativeHist[f] = cumulativeSum
	}

	// TODO: how to calculate alpha? + completely different formula: https://www.sciencedirect.com/science/article/pii/S2092678216302588?ref=pdf_download&fr=RR-2&rr=8dc38c4dada6bf3f
	alpha := float64(gmax-gmin) / math.Sqrt(2)

	// Apply the Rayleigh formula
	for f := fmin; f <= fmax; f++ {
		cumulativeSum := cumulativeHist[f]
		prob := float64(cumulativeSum) / float64(totalPixels)
		logTerm := 1 - prob
		if logTerm <= 0 {
			logTerm = 1e-10 // Avoid log(0) error by using a small positive value
		}
		term := math.Log(logTerm)
		transformed[f] = float64(gmin) + math.Sqrt(-2*alpha*alpha*term)
	}

	// Find min and max of the transformed values for normalization
	minVal, maxVal := transformed[fmin], transformed[fmin]
	for f := fmin; f <= fmax; f++ {
		if transformed[f] < minVal {
			minVal = transformed[f]
		}
		if transformed[f] > maxVal {
			maxVal = transformed[f]
		}
	}

	// Normalize transformed values to ensure they fit within [gmin, gmax]
	for f := fmin; f <= fmax; f++ {
		transformed[f] = float64(gmin) + (transformed[f]-minVal)*(float64(gmax-gmin)/(maxVal-minVal))
	}

	return transformed
}

func applyHistogramTransform(img *image.Gray, transform []float64) *image.RGBA {
	bounds := img.Bounds()
	outImg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			oldGray := img.GrayAt(x, y).Y
			newGrayValue := uint8(transform[oldGray])
			rgbaColor := color.RGBA{R: newGrayValue, G: newGrayValue, B: newGrayValue, A: 255}
			outImg.SetRGBA(x, y, rgbaColor)
		}
	}
	return outImg
}

// ApplyRayleighTransform applies a Rayleigh transformation to the given image.
// The transformation adjusts the pixel values based on the Rayleigh distribution
// to enhance the contrast of the image.
//
// Parameters:
//   - img: The input image to be transformed.
//   - gMin: The minimum value for the output image.
//   - gMax: The maximum value for the output image.
//
// Returns:
//   - *image.RGBA: The transformed image with enhanced contrast.
func ApplyRayleighTransform(img image.Image, gMin int, gMax int) *image.RGBA {

	// TODO: make it work for color images.
	grayImg := toGrayscale(img)
	histogram := calculateHistogram(grayImg)
	totalPixels := grayImg.Bounds().Dx() * grayImg.Bounds().Dy()

	fMin, fMax := FindMinMax(histogram)

	transform := rayleighTransform(histogram, fMin, fMax, gMin, gMax, totalPixels)

	outputImg := applyHistogramTransform(grayImg, transform)

	return outputImg

}
