package noise

import (
	"image"
	"image-processing/manipulations"
	"image/color"
	"slices"
)

func getWindowPixelsRGB(img image.Image, x, y, windowSize int) ([]int, []int, []int) {
	bounds := img.Bounds()

	halfSize := windowSize / 2

	size := (2*halfSize + 1) * (2*halfSize + 1)
	reds, greens, blues := make([]int, 0, size), make([]int, 0, size), make([]int, 0, size)

	for j := y - halfSize; j <= y+halfSize; j++ {
		for i := x - halfSize; i <= x+halfSize; i++ {
			if i >= bounds.Min.X && i < bounds.Max.X && j >= bounds.Min.Y && j < bounds.Max.Y {
				r, g, b, _ := img.At(i, j).RGBA()
				reds = append(reds, int(r>>8))
				greens = append(greens, int(g>>8))
				blues = append(blues, int(b>>8))
			}
		}
	}

	return reds, greens, blues
}

func minMaxMedian(pixels []int) (int, int, int) {
	if len(pixels) == 0 {
		return 0, 0, 0
	}

	slices.Sort(pixels)

	min := pixels[0]
	max := pixels[len(pixels)-1]
	median := pixels[len(pixels)/2]

	return min, max, median
}

func applyAdaptiveMedian(B1, B2, zMed, zxy int) int {
	if B1 > 0 && B2 < 0 {
		return zxy
	}
	return zMed

}
func adaptiveChannel(newVal, xyVal, min, max, med int) int {
	A1, A2 := med-min, med-max
	if A1 > 0 && A2 < 0 {
		B1, B2 := xyVal-min, xyVal-max
		return applyAdaptiveMedian(B1, B2, med, xyVal)
	}
	return newVal
}

func AdaptiveMedianFilter(img image.Image, sMin, sMax int) *image.RGBA {
	// https://www.irjet.net/archives/V6/i10/IRJET-V6I10148.pdf
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {

			windowSize := sMin
			var newR, newG, newB = 0, 0, 0

			urxy, ugxy, ubxy, _ := img.At(x, y).RGBA()
			rxy, gxy, bxy := int(urxy>>8), int(ugxy>>8), int(ubxy>>8)
			newR, newG, newB = rxy, gxy, bxy

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

			newImg.Set(x, y, color.RGBA{manipulations.ClampUint8(newR), manipulations.ClampUint8(newG), manipulations.ClampUint8(newB), 255})
		}
	}

	return newImg
}

func MinFilter(img image.Image, windowSize int) *image.RGBA {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			reds, greens, blues := getWindowPixelsRGB(img, x, y, windowSize)
			rMin, _, _ := minMaxMedian(reds)
			gMin, _, _ := minMaxMedian(greens)
			bMin, _, _ := minMaxMedian(blues)

			newImg.Set(x, y, color.RGBA{uint8(rMin), uint8(gMin), uint8(bMin), 255})
		}
	}

	return newImg
}

func MaxFilter(img image.Image, windowSize int) *image.RGBA {
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			reds, greens, blues := getWindowPixelsRGB(img, x, y, windowSize)

			_, rMax, _ := minMaxMedian(reds)
			_, gMax, _ := minMaxMedian(greens)
			_, bMax, _ := minMaxMedian(blues)

			newImg.Set(x, y, color.RGBA{uint8(rMax), uint8(gMax), uint8(bMax), 255})
		}
	}

	return newImg
}
