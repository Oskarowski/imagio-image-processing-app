package orthogonal_transforms

import (
	"image-processing/morphological"
	"math"
)

//https://www.clear.rice.edu/elec301/Projects01/image_filt/concept.html

func BandPassFilter2D(input [][]complex128, lowCutoff, highCutoff float64) [][]complex128 {
	n := len(input)
	m := len(input[0])

	// Create filter mask
	filter := make([][]float64, n)
	for i := range filter {
		filter[i] = make([]float64, m)
	}

	centerX, centerY := n/2, m/2

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			// Compute the distance from the center
			dist := math.Sqrt(float64((i-centerX)*(i-centerX) + (j-centerY)*(j-centerY)))

			// Pass frequencies within the band
			if dist >= lowCutoff && dist <= highCutoff {
				filter[i][j] = 1.0
			} else {
				filter[i][j] = 0.0
			}
		}
	}

	filtered := make([][]complex128, n)
	for i := range filtered {
		filtered[i] = make([]complex128, m)
		for j := 0; j < m; j++ {
			filtered[i][j] = input[i][j] * complex(filter[i][j], 0)
		}
	}

	return filtered
}

func LowPassFilter2D(input [][]complex128, cutoff float64) [][]complex128 {
	n := len(input)
	m := len(input[0])

	// Create filter mask
	filter := make([][]float64, n)
	for i := range filter {
		filter[i] = make([]float64, m)
	}

	centerX, centerY := n/2, m/2

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			// Compute the distance from the center
			dist := math.Sqrt(float64((i-centerX)*(i-centerX) + (j-centerY)*(j-centerY)))

			// Allow frequencies below the cutoff
			if dist <= cutoff {
				filter[i][j] = 1.0
			} else {
				filter[i][j] = 0.0
			}
		}
	}

	filtered := make([][]complex128, n)
	for i := range filtered {
		filtered[i] = make([]complex128, m)
		for j := 0; j < m; j++ {
			filtered[i][j] = input[i][j] * complex(filter[i][j], 0)
		}
	}

	return filtered
}

func HighPassFilter2D(input [][]complex128, cutoff float64) [][]complex128 {
	n := len(input)
	m := len(input[0])

	// Create filter mask
	filter := make([][]float64, n)
	for i := range filter {
		filter[i] = make([]float64, m)
	}

	centerX, centerY := n/2, m/2

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			// Compute the distance from the center
			dist := math.Sqrt(float64((i-centerX)*(i-centerX) + (j-centerY)*(j-centerY)))

			// Allow frequencies above the cutoff
			if dist >= cutoff {
				filter[i][j] = 1.0
			} else {
				filter[i][j] = 0.0
			}
		}
	}

	filtered := make([][]complex128, n)
	for i := range filtered {
		filtered[i] = make([]complex128, m)
		for j := 0; j < m; j++ {
			filtered[i][j] = input[i][j] * complex(filter[i][j], 0)
		}
	}

	return filtered
}

func BandCutFilter2D(input [][]complex128, lowCutoff, highCutoff float64) [][]complex128 {
	n := len(input)
	m := len(input[0])

	// Create filter mask
	filter := make([][]float64, n)
	for i := range filter {
		filter[i] = make([]float64, m)
	}

	centerX, centerY := n/2, m/2

	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			// Compute the distance from the center
			dist := math.Sqrt(float64((i-centerX)*(i-centerX) + (j-centerY)*(j-centerY)))

			// Block frequencies within the band (lowCutoff <= dist <= highCutoff)
			if dist >= lowCutoff && dist <= highCutoff {
				filter[i][j] = 0.0
			} else {
				filter[i][j] = 1.0
			}
		}
	}

	filtered := make([][]complex128, n)
	for i := range filtered {
		filtered[i] = make([]complex128, m)
		for j := 0; j < m; j++ {
			filtered[i][j] = input[i][j] * complex(filter[i][j], 0)
		}
	}

	return filtered
}

func ScaleMask(mask [][]int, targetHeight, targetWidth int) [][]int {
	sourceHeight := len(mask)
	sourceWidth := len(mask[0])

	scaledMask := make([][]int, targetHeight)
	for y := range scaledMask {
		scaledMask[y] = make([]int, targetWidth)
		for x := 0; x < targetWidth; x++ {
			// Calculate source coordinates
			sourceX := float64(x) * float64(sourceWidth) / float64(targetWidth)
			sourceY := float64(y) * float64(sourceHeight) / float64(targetHeight)

			// Bilinear interpolation
			x0, y0 := int(math.Floor(sourceX)), int(math.Floor(sourceY))
			x1, y1 := x0+1, y0+1
			if x1 >= sourceWidth {
				x1 = sourceWidth - 1
			}
			if y1 >= sourceHeight {
				y1 = sourceHeight - 1
			}

			topLeft := mask[y0][x0]
			topRight := mask[y0][x1]
			bottomLeft := mask[y1][x0]
			bottomRight := mask[y1][x1]

			fx, fy := sourceX-float64(x0), sourceY-float64(y0)
			top := float64(topLeft)*(1-fx) + float64(topRight)*fx
			bottom := float64(bottomLeft)*(1-fx) + float64(bottomRight)*fx

			scaledMask[y][x] = int(math.Round(top*(1-fy) + bottom*fy))
		}
	}

	return scaledMask
}

func HighPassFilterWithEdgeDetection2D(imageSpectrum [][]complex128, mask morphological.BinaryImage, cutoff int, edgeDirection int) [][]complex128 {
	height := len(imageSpectrum)
	width := len(imageSpectrum[0])

	scaledMask := ScaleMask(mask, height, width)

	spectrum := make([][]complex128, height)
	for y := range spectrum {
		spectrum[y] = make([]complex128, width)
		copy(spectrum[y], imageSpectrum[y])
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if scaledMask[y][x] == 0 {
				spectrum[y][x] = complex(0, 0)
			}
		}
	}

	return spectrum
}
