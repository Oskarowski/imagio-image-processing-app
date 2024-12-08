package morphological

import (
	"image"
	"image/color"
	"math"
	"math/rand"
)

type Point struct {
	X, Y int
}

type Region struct {
	Pixels []Point
	Label  int // ID of the region
}

func getNeighbors(p Point, rows, cols int) []Point {
	neighbors := []Point{}
	directions := []Point{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}

	for _, d := range directions {
		neighbor := Point{X: p.X + d.X, Y: p.Y + d.Y}

		if neighbor.X >= 0 && neighbor.X < rows && neighbor.Y >= 0 && neighbor.Y < cols {
			neighbors = append(neighbors, neighbor)
		}
	}

	return neighbors
}

func getPixelValue(c color.Color) []float64 {
	r, g, b, _ := c.RGBA()
	return []float64{float64(r >> 8), float64(g >> 8), float64(b >> 8)}
}

func calculateDistance(criterion int, p1, p2 []float64) float64 {
	rDiff := p1[0] - p2[0]
	gDiff := p1[1] - p2[1]
	bDiff := p1[2] - p2[2]

	switch criterion {
	case 0: // Euclidean Distance - Reference: https://medium.com/@khushihp7903/exploring-image-segmentation-with-the-region-growing-algorithm-4972dae63680
		return math.Sqrt(rDiff*rDiff + gDiff*gDiff + bDiff*bDiff)
	case 1: // Manhattan Distance
		return math.Abs(rDiff) + math.Abs(gDiff) + math.Abs(bDiff)
	case 2: // Chebyshev Distance
		return math.Max(math.Max(math.Abs(rDiff), math.Abs(gDiff)), math.Abs(bDiff))
	default:
		return 0
	}
}
func randomColor() color.Color {
	return color.RGBA{
		R: uint8(rand.Intn(256)),
		G: uint8(rand.Intn(256)),
		B: uint8(rand.Intn(256)),
		A: 255,
	}
}

func RegionGrowing(img image.Image, seeds []Point, criterion int, threshold float64) ([][]int, *image.RGBA) {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	segmented := make([][]int, height)
	for y := range segmented {
		segmented[y] = make([]int, width)

		for x := range segmented[y] {
			segmented[y][x] = -1
		}
	}

	outputImage := image.NewRGBA(bounds)
	label := 0

	for _, seed := range seeds {
		if segmented[seed.X][seed.Y] != -1 {
			continue
		}

		queue := []Point{seed}
		seedValue := getPixelValue(img.At(seed.X, seed.Y))
		regionColor := randomColor()

		for len(queue) > 0 {
			current := queue[0]
			queue = queue[1:]
			if segmented[current.X][current.Y] != -1 {
				continue
			}

			segmented[current.X][current.Y] = label
			outputImage.Set(current.X, current.Y, regionColor)

			for _, neighbor := range getNeighbors(current, width, height) {
				neighborValue := getPixelValue(img.At(neighbor.X, neighbor.Y))

				distance := calculateDistance(criterion, neighborValue, seedValue)

				if segmented[neighbor.X][neighbor.Y] == -1 && distance <= threshold {
					queue = append(queue, neighbor)
				}
			}
		}

		label++
	}

	return segmented, outputImage

}
