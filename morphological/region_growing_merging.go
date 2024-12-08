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

func randomColor() color.Color {
	return color.RGBA{
		R: uint8(rand.Intn(256)),
		G: uint8(rand.Intn(256)),
		B: uint8(rand.Intn(256)),
		A: 255,
	}
}

func getNeighbors(p Point, rows, cols int) []Point {
	neighbors := []Point{}
	directions := []Point{
		{-1, 0}, {1, 0}, {0, -1}, {0, 1},
		{-1, -1}, {-1, 1}, {1, -1}, {1, 1},
	}
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

	normalizedR := float64(r) / 257.0
	normalizedG := float64(g) / 257.0
	normalizedB := float64(b) / 257.0

	if math.Abs(normalizedR-normalizedG) < 1e-6 && math.Abs(normalizedG-normalizedB) < 1e-6 {
		return []float64{normalizedR} // Single channel for grayscale
	}

	return []float64{normalizedR, normalizedG, normalizedB} // RGB
}

func calculateDistance(criterion int, p1, p2 []float64) float64 {
	if len(p1) != len(p2) {
		panic("Pixel values must have the same dimension")
	}

	switch criterion {
	case 0: // Euclidean Distance
		sum := 0.0
		for i := 0; i < len(p1); i++ {
			diff := p1[i] - p2[i]
			sum += diff * diff
		}
		return math.Sqrt(sum)
	case 1: // Manhattan Distance
		sum := 0.0
		for i := 0; i < len(p1); i++ {
			sum += math.Abs(p1[i] - p2[i])
		}
		return sum
	case 2: // Chebyshev Distance
		maxDiff := 0.0
		for i := 0; i < len(p1); i++ {
			diff := math.Abs(p1[i] - p2[i])
			if diff > maxDiff {
				maxDiff = diff
			}
		}
		return maxDiff

	default:
		return 0
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
		if segmented[seed.Y][seed.X] != -1 {
			continue
		}

		queue := []Point{seed}
		seedValue := getPixelValue(img.At(seed.X, seed.Y))
		regionColor := randomColor()

		for len(queue) > 0 {
			current := queue[0]
			queue = queue[1:]

			if segmented[current.Y][current.X] != -1 {
				continue
			}

			segmented[current.Y][current.X] = label
			outputImage.Set(current.X, current.Y, regionColor)

			for _, neighbor := range getNeighbors(current, width, height) {
				if segmented[neighbor.Y][neighbor.X] == -1 {
					neighborValue := getPixelValue(img.At(neighbor.X, neighbor.Y))
					distance := calculateDistance(criterion, neighborValue, seedValue)

					if distance <= threshold {
						queue = append(queue, neighbor)
					}
				}
			}
		}

		label++
	}

	return segmented, outputImage
}
