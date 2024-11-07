package manipulations

import (
	"image"
	"image/color"
	"image/draw"
	"math"
)

// EnhanceImageWithRayleigh applies a Rayleigh transformation to the given image
// to enhance the contrast of the image.
//
// Parameters:
//
//	img       - The input image to be transformed.
//	gmin      - The minimum value for the output image.
//	gmax      - The maximum value for the output image.
//	alpha     - Parameter alpha for Rayleigh formula.
//
// Returns:
//
//	*image.RGBA - The transformed image with enhanced contrast.
func EnhanceImageWithRayleigh(img image.Image, gmin, gmax, alpha float64) *image.RGBA {
	baseHistogram := CalculateHistogram(img)
	bounds := img.Bounds()
	N := float64(bounds.Dx() * bounds.Dy())

	// Calculate cumulative histogram
	var cumulativeHistogram [256]float64
	cumulativeHistogram[0] = float64(baseHistogram[0]) / N
	for i := 1; i < 256; i++ {
		cumulativeHistogram[i] = cumulativeHistogram[i-1] + float64(baseHistogram[i])/N
	}

	// Transform brightness levels using the Rayleigh PDF formula
	output := image.NewGray(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixel := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
			f := int(pixel.Y)

			// g := gmin + math.Sqrt(2*alpha*alpha*math.Log(1/(1-cumulativeHistogram[f])))

			cumulativeValueAdjustment := 1 - cumulativeHistogram[f]
			// cumulativeValueAdjustment := float64(cumulativeHistogram[f])

			if cumulativeValueAdjustment <= 0 {
				cumulativeValueAdjustment = 1e-10
			}

			logPart := math.Log(1 / cumulativeValueAdjustment)

			squaredAlpha := alpha * alpha
			transformedValue := 2 * squaredAlpha * logPart
			sqrtValue := math.Sqrt(transformedValue)

			g := gmin + sqrtValue

			if g > gmax {
				g = gmax
			} else if g < gmin {
				g = gmin
			}

			output.SetGray(x, y, color.Gray{Y: uint8(g)})
		}
	}

	RGBA := image.NewRGBA(output.Bounds())
	draw.Draw(RGBA, RGBA.Bounds(), output, output.Bounds().Min, draw.Src)

	return RGBA
}
