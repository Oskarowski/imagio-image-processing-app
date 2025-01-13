package orthogonal_transforms

import "math"

func LowPassFilter(data [][]complex128, cutoff float64) [][]complex128 {
	N, M := len(data), len(data[0])
	output := make([][]complex128, N)
	for u := 0; u < N; u++ {
		output[u] = make([]complex128, M)
		for v := 0; v < M; v++ {
			dist := math.Sqrt(float64((u-N/2)*(u-N/2) + (v-M/2)*(v-M/2)))
			if dist <= cutoff {
				output[u][v] = data[u][v]
			}
		}
	}

	return output
}

// BandPassFilter2D applies a band-pass filter in the frequency domain
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

	// Apply the filter in the frequency domain
	filtered := make([][]complex128, n)
	for i := range filtered {
		filtered[i] = make([]complex128, m)
		for j := 0; j < m; j++ {
			filtered[i][j] = input[i][j] * complex(filter[i][j], 0)
		}
	}

	return filtered
}

// BandCutFilter - Apply a band-cut (band-stop) filter to the 2D frequency-domain data.
func BandCutFilter(fftData [][]complex128, lowCutoff, highCutoff int) [][]complex128 {
	rows := len(fftData)
	cols := len(fftData[0])

	// Create a new matrix to store the filtered data
	filtered := make([][]complex128, rows)
	for i := range filtered {
		filtered[i] = make([]complex128, cols)
	}

	// Apply the band-cut filter
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			// Calculate distance from the center (0,0) in frequency space
			u := i - rows/2
			v := j - cols/2
			distance := math.Sqrt(float64(u*u + v*v))

			// Handle the DC component (at (0,0))
			if i == rows/2 && j == cols/2 {
				// Optional: Zero-out the DC component (remove average intensity)
				filtered[i][j] = complex(0, 0) // Zero-out DC component
			} else if int(distance) >= lowCutoff && int(distance) <= highCutoff {
				// Apply the band-cut filter: zero out frequencies within the cutoff range
				filtered[i][j] = complex(0, 0)
			} else {
				// Keep other frequencies
				filtered[i][j] = fftData[i][j]
			}
		}
	}

	return filtered
}
