package orthogonal_transforms

import "math"

func scaleMask(mask [][]int, targetHeight, targetWidth int) [][]int {
	sourceHeight := len(mask)
	sourceWidth := len(mask[0])

	scaledMask := make([][]int, targetHeight)
	rowScale := float64(sourceHeight) / float64(targetHeight)
	colScale := float64(sourceWidth) / float64(targetWidth)

	for y := 0; y < targetHeight; y++ {
		scaledMask[y] = make([]int, targetWidth)
		sourceY := rowScale * float64(y)

		y0 := int(sourceY)
		y1 := min(y0+1, sourceHeight-1)
		yFrac := sourceY - float64(y0)

		for x := 0; x < targetWidth; x++ {
			sourceX := colScale * float64(x)

			x0 := int(sourceX)
			x1 := min(x0+1, sourceWidth-1)
			xFrac := sourceX - float64(x0)

			topLeft := float64(mask[y0][x0])
			topRight := float64(mask[y0][x1])
			bottomLeft := float64(mask[y1][x0])
			bottomRight := float64(mask[y1][x1])

			// Bilinear interpolation
			top := (1-xFrac)*topLeft + xFrac*topRight
			bottom := (1-xFrac)*bottomLeft + xFrac*bottomRight
			scaledMask[y][x] = int(math.Round((1-yFrac)*top + yFrac*bottom))
		}
	}

	return scaledMask
}
