package manipulations

import (
	"image"
	"image/color"
	"image/draw"
	"math"
)

// source of the masks: https://www.tutorialspoint.com/dip/krisch_compass_mask.htm
var kirschDirections = [][]int{
	{-3, -3, 5, -3, 0, 5, -3, -3, 5}, // North
	{-3, 5, 5, -3, 0, 5, -3, -3, -3}, // Northwest
	{5, 5, 5, -3, 0, -3, -3, -3, -3}, // West
	{5, 5, -3, 5, 0, -3, -3, -3, -3}, // Southwest
	{5, -3, -3, 5, 0, -3, 5, -3, -3}, // South
	{-3, -3, -3, 5, 0, -3, 5, 5, -3}, // Southeast
	{-3, -3, -3, -3, 0, -3, 5, 5, 5}, // East
	{-3, -3, -3, -3, 0, 5, -3, 5, 5}, // Northeast
}

func ApplyKirshEdgeDetection(img image.Image) *image.RGBA {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	edgeImage := image.NewGray(bounds)
	draw.Draw(edgeImage, bounds, img, bounds.Min, draw.Src)

	for y := 1; y < height-1; y++ {
		for x := 1; x < width-1; x++ {
			maxGradient := float64(0)

			for _, direction := range kirschDirections {
				g := 0.0
				for dy := -1; dy <= 1; dy++ {
					for dx := -1; dx <= 1; dx++ {
						pixel := color.GrayModel.Convert(img.At(x+dx, y+dy)).(color.Gray)
						g += float64(pixel.Y) * float64(direction[(dy+1)*3+(dx+1)])
					}
				}
				maxGradient = math.Max(maxGradient, math.Abs(g))
			}

			edgeColor := uint8(math.Min(maxGradient, 255))

			edgeImage.SetGray(x, y, color.Gray{Y: uint8(edgeColor)})
		}
	}

	rgbaImg := image.NewRGBA(bounds)
	draw.Draw(rgbaImg, bounds, edgeImage, bounds.Min, draw.Src)

	return rgbaImg
}
