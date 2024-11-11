package manipulations

import (
	"image"
	"image/color"
	"image/draw"
	"math"
)

func ApplyKirshEdgeDetection(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	edgeImage := image.NewGray(bounds)
	draw.Draw(edgeImage, bounds, img, bounds.Min, draw.Src)

	neighborhood := [][]int{
		{-1, 0}, {-1, 1}, {0, 1}, {1, 1}, {1, 0}, {1, -1}, {0, -1}, {-1, -1},
	}

	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {

			maxGradient := 0.0
			neighborPixels := make([]int, 8)

			for i := 0; i < 8; i++ {
				dx, dy := neighborhood[i][0], neighborhood[i][1]
				pixel := color.GrayModel.Convert(img.At(x+dx, y+dy)).(color.Gray)
				neighborPixels[i] = int(pixel.Y)
			}

			for i := 0; i < 8; i++ {
				S := neighborPixels[i] +
					neighborPixels[(i+1)%8] +
					neighborPixels[(i+2)%8]

				T := neighborPixels[(i+3)%8] +
					neighborPixels[(i+4)%8] +
					neighborPixels[(i+5)%8] +
					neighborPixels[(i+6)%8] +
					neighborPixels[(i+7)%8]

				gradient := math.Abs(5*float64(S) - 3*float64(T))

				if gradient > maxGradient {
					maxGradient = gradient
				}
			}

			edgeValue := uint8(math.Min(maxGradient, 255))
			edgeImage.SetGray(x, y, color.Gray{Y: uint8(edgeValue)})
		}
	}

	rgbaImg := image.NewRGBA(bounds)
	draw.Draw(rgbaImg, bounds, edgeImage, bounds.Min, draw.Src)

	return rgbaImg
}
