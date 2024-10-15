package noise

import (
	"image"
	"image/color"
	"sort"
)

func getWindowPixelsRGB(img image.Image, x, y, windowSize int) ([]int, []int, []int) {
	bounds := img.Bounds()
	var reds, greens, blues []int

	halfSize := windowSize / 2

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

	sort.Ints(pixels)

	min := pixels[0]
	max := pixels[len(pixels)-1]
	median := pixels[len(pixels)/2]

	return min, max, median
}

func applyAdaptiveMedian(A1, A2, B1, B2, zMed, zxy int) int {
	if A1 > 0 && A2 < 0 {
		if B1 > 0 && B2 < 0 {
			return zxy
		}
		return zMed
	}
	return zxy
}

func AdaptiveMedianFilter(img image.Image, sMin, sMax int) *image.RGBA {
	// https://www.irjet.net/archives/V6/i10/IRJET-V6I10148.pdf
	bounds := img.Bounds()
	newImg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {

			windowSize := sMin
			rxy, gxy, bxy, _ := img.At(x, y).RGBA()
			rxy, gxy, bxy = rxy>>8, gxy>>8, bxy>>8

			for windowSize <= sMax {
				reds, greens, blues := getWindowPixelsRGB(img, x, y, windowSize)

				rMin, rMax, rMed := minMaxMedian(reds)
				gMin, gMax, gMed := minMaxMedian(greens)
				bMin, bMax, bMed := minMaxMedian(blues)

				// Stage A - Red Channel
				A1r, A2r := rMed-rMin, rMed-rMax
				if A1r > 0 && A2r < 0 {
					// Stage B - Red Channel
					B1r, B2r := int(rxy)-rMin, int(rxy)-rMax
					newR := applyAdaptiveMedian(A1r, A2r, B1r, B2r, rMed, int(rxy))

					// Stage A - Green Channel
					A1g, A2g := gMed-gMin, gMed-gMax

					if A1g > 0 && A2g < 0 {
						// Stage B - Green Channel
						B1g, B2g := int(gxy)-gMin, int(gxy)-gMax
						newG := applyAdaptiveMedian(A1g, A2g, B1g, B2g, gMed, int(gxy))

						// Stage A - Blue Channel
						A1b, A2b := bMed-bMin, bMed-bMax
						if A1b > 0 && A2b < 0 {
							// Stage B - Blue Channel
							B1b, B2b := int(bxy)-bMin, int(bxy)-bMax
							newB := applyAdaptiveMedian(A1b, A2b, B1b, B2b, bMed, int(bxy))

							newImg.Set(x, y, color.RGBA{uint8(newR), uint8(newG), uint8(newB), 255})
							break
						}
					}
				}

				windowSize += 2
			}

			if windowSize > sMax {
				newImg.Set(x, y, img.At(x, y))
			}
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
