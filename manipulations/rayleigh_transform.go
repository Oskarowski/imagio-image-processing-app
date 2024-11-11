package manipulations

import (
	"image"
	"image/color"
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

	var cumulativeHistogram [256]float64
	cumulativeHistogram[0] = float64(baseHistogram[0]) / N
	for i := 1; i < 256; i++ {
		cumulativeHistogram[i] = cumulativeHistogram[i-1] + float64(baseHistogram[i])/N
	}

	outputImg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			h, s, v := RGBToHSV(uint8(r>>8), uint8(g>>8), uint8(b>>8))
			f := int(v * 255)

			// g := gmin + math.Sqrt(2*alpha*alpha*math.Log(1/(1-cumulativeHistogram[f])))

			cumulativeValueAdjustment := 1 - cumulativeHistogram[f]

			if cumulativeValueAdjustment <= 0 {
				cumulativeValueAdjustment = 1e-10
			}

			logPart := math.Log(1 / cumulativeValueAdjustment)

			squaredAlpha := alpha * alpha
			transformedValue := 2 * squaredAlpha * logPart
			newV := gmin + math.Sqrt(transformedValue)

			if newV > gmax {
				newV = gmax
			} else if newV < gmin {
				newV = gmin
			}
			newV /= 255.0

			rOut, gOut, bOut := HSVToRGB(h, s, newV)
			outputImg.Set(x, y, color.RGBA{
				R: uint8(rOut),
				G: uint8(gOut),
				B: uint8(bOut),
				A: 255,
			})
		}
	}

	return outputImg
}
