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

func HighPassFilterWithEdgeDetection2D(imageSpectrum [][]complex128, mask morphological.BinaryImage) [][]complex128 {
	height := len(imageSpectrum)
	width := len(imageSpectrum[0])

	scaledMask := scaleMask(mask, height, width)

	spectrum := make([][]complex128, height)
	for y := range spectrum {
		spectrum[y] = make([]complex128, width)
		copy(spectrum[y], imageSpectrum[y])
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if scaledMask[y][x] == 0 {
				spectrum[y][x] = 0
			}
		}
	}

	return spectrum
}
