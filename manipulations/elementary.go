package manipulations

import (
	"image"
	"image/color"
)

func ClampUint8(value int) uint8 {
	if value > 255 {
		return 255
	}
	if value < 0 {
		return 0
	}
	return uint8(value)
}

// The function returns a new image with the adjusted brightness.
//
// Parameters:
//
//	img       - The input image to adjust.
//	brightness - Percentage level of brightness.
//
// Returns:
//
//	*image.RGBA - A new image with the adjusted brightness.
func AdjustBrightness(img image.Image, brightness int) *image.RGBA {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	factor := (brightness * 255) / 100

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			r8 := int(r >> 8)
			g8 := int(g >> 8)
			b8 := int(b >> 8)

			newR := ClampUint8(r8 + factor)
			newG := ClampUint8(g8 + factor)
			newB := ClampUint8(b8 + factor)

			newImg.Set(x, y, color.RGBA{newR, newG, newB, uint8(a >> 8)})
		}
	}

	return newImg
}

// The function returns a new image with the adjusted contrast.
//
// Parameters:
//
//	img      - The input image to adjust.
//	contrast - Desired level of contrast.
//
// Returns:
//
//	*image.RGBA - A new image with the adjusted contrast.
func AdjustContrast(img image.Image, contrast int) *image.RGBA {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	// Contrast correction factor formula:
	// https://www.dfstudios.co.uk/articles/programming/image-programming-algorithms/image-processing-algorithms-part-5-contrast-adjustment/
	// https://ie.nitk.ac.in/blog/2020/01/19/algorithms-for-adjusting-brightness-and-contrast-of-an-image/
	var contrastCorrectionFactor float64 = (259.0 * float64(contrast+255)) / (255.0 * float64(259-contrast))

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			r8 := float64(r >> 8)
			g8 := float64(g >> 8)
			b8 := float64(b >> 8)

			newR := ClampUint8(int(contrastCorrectionFactor*(r8-128) + 128))
			newG := ClampUint8(int(contrastCorrectionFactor*(g8-128) + 128))
			newB := ClampUint8(int(contrastCorrectionFactor*(b8-128) + 128))

			newImg.Set(x, y, color.RGBA{newR, newG, newB, uint8(a >> 8)})
		}
	}

	return newImg
}

// The function returns a new image with the negative of the input image.
// Parameters:
//
//	img - The input image to negate.
//
// Returns:
//
//	*image.RGBA - A new image with the negative of the input image.
func NegativeImage(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			newR := 255 - uint8(r>>8)
			newG := 255 - uint8(g>>8)
			newB := 255 - uint8(b>>8)

			newImg.Set(x, y, color.RGBA{newR, newG, newB, uint8(a >> 8)})
		}
	}

	return newImg
}
