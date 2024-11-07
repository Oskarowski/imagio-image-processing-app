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

func calculateCumulativeHistogram(histogram []int) [256]float64 {
	var cumulativeHistogram [256]float64
	cumulativeHistogram[0] = float64(histogram[0])

	for i := 1; i < 256; i++ {
		cumulativeHistogram[i] = cumulativeHistogram[i-1] + float64(histogram[i])
	}

	return cumulativeHistogram
}

func rayleighTransform(img image.Gray, baseHistogram []int, fmin, fmax, gmin, gmax, totalPixels int, alpha float64) *image.Gray {
	cumulativeHistogram := calculateCumulativeHistogram(baseHistogram)

	transform := make([]float64, 256)

	for f := 0; f < 256; f++ {
		cumulativeProbability := cumulativeHistogram[f] / float64(totalPixels)
		if cumulativeProbability > 0 {
			alphaPow := math.Pow(alpha, 2)
			transform[f] = float64(gmin) + math.Sqrt(2*alphaPow*math.Log(1/cumulativeProbability)-1)
			if transform[f] < float64(gmin) {
				transform[f] = float64(gmin)
			} else if transform[f] > float64(gmax) {
				transform[f] = float64(gmax)
			}
		} else {
			transform[f] = float64(gmin)
		}
	}

	bounds := img.Bounds()
	newImg := image.NewGray(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			oldPixel := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
			newBrightness := uint8(transform[oldPixel.Y])
			newImg.SetGray(x, y, color.Gray{Y: newBrightness})
		}
	}

	return newImg
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
//   - alpha: Parameter alpha for Rayleigh formula.
//
// Returns:
//   - *image.RGBA: The transformed image with enhanced contrast.
func ApplyRayleighTransform(img image.Image, gMin int, gMax int, alpha float64) *image.RGBA {

	// TODO: make it work for color images.
	grayImg := toGrayscale(img)
	histogram := calculateHistogram(grayImg)
	totalPixels := grayImg.Bounds().Dx() * grayImg.Bounds().Dy()

	fMin, fMax := FindMinMax(histogram)

	transformedImg := rayleighTransform(*grayImg, histogram, fMin, fMax, gMin, gMax, totalPixels, alpha)

	transformedRGBAImg := image.NewRGBA(transformedImg.Bounds())
	draw.Draw(transformedRGBAImg, transformedRGBAImg.Bounds(), transformedImg, transformedImg.Bounds().Min, draw.Src)

	// outputImg := applyHistogramTransform(grayImg, transform)

	return transformedRGBAImg

}
