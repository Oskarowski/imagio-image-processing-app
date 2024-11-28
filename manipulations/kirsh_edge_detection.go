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

	edgeImage := image.NewRGBA(bounds)
	draw.Draw(edgeImage, bounds, img, bounds.Min, draw.Src)

	neighborhood := [][]int{
		{-1, 0}, {-1, 1}, {0, 1}, {1, 1}, {1, 0}, {1, -1}, {0, -1}, {-1, -1},
	}

	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {

			maxGradient := 0.0
			neighborValues := make([]float64, 8)

			for i := 0; i < 8; i++ {
				dx, dy := neighborhood[i][0], neighborhood[i][1]
				pixel := color.RGBAModel.Convert(img.At(x+dx, y+dy)).(color.RGBA)
				_, _, v := RGBToHSV(pixel.R, pixel.G, pixel.B)

				neighborValues[i] = v
			}

			for i := 0; i < 8; i++ {
				S := neighborValues[i] +
					neighborValues[(i+1)%8] +
					neighborValues[(i+2)%8]

				T := neighborValues[(i+3)%8] +
					neighborValues[(i+4)%8] +
					neighborValues[(i+5)%8] +
					neighborValues[(i+6)%8] +
					neighborValues[(i+7)%8]

				gradient := math.Abs(5*S - 3*T)

				if gradient > maxGradient {
					maxGradient = gradient
				}
			}

			edgeValue := uint8(math.Min(maxGradient*255, 255))
			edgeImage.SetRGBA(x, y, color.RGBA{R: edgeValue, G: edgeValue, B: edgeValue, A: 255})
		}
	}

	return edgeImage
}
