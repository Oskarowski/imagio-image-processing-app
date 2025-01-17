package orthogonal_transforms

import (
	"image"
	"image/color"
	"math"
)

func magnitude(c complex128) float64 {
	return math.Sqrt(real(c)*real(c) + imag(c)*imag(c))
}

// Create a 2D FFT magnitude spectrum with log scaling
func FFTMagnitudeSpectrum(fftResult [][]complex128) [][]float64 {
	rows := len(fftResult)
	cols := len(fftResult[0])
	magnitudeSpectrum := make([][]float64, rows)

	// Calculate the magnitude of each complex number
	for i := 0; i < rows; i++ {
		magnitudeSpectrum[i] = make([]float64, cols)
		for j := 0; j < cols; j++ {
			magnitudeSpectrum[i][j] = magnitude(fftResult[i][j])
		}
	}

	// Apply log scaling
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			magnitudeSpectrum[i][j] = math.Log(1 + magnitudeSpectrum[i][j])
		}
	}

	return magnitudeSpectrum
}

// Normalize the magnitude spectrum to [0, 255]
func NormalizeMagnitude(magnitude [][]float64) [][]uint8 {
	rows := len(magnitude)
	cols := len(magnitude[0])

	// Find the min and max values in the magnitude spectrum
	min, max := math.Inf(1), -math.Inf(1)
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			if magnitude[i][j] < min {
				min = magnitude[i][j]
			}
			if magnitude[i][j] > max {
				max = magnitude[i][j]
			}
		}
	}

	// Normalize values to the range [0, 255]
	normalized := make([][]uint8, rows)
	scale := 255 / (max - min)
	for i := 0; i < rows; i++ {
		normalized[i] = make([]uint8, cols)
		for j := 0; j < cols; j++ {
			normalized[i][j] = uint8((magnitude[i][j] - min) * scale)
		}
	}

	return normalized
}

// Convert the normalized magnitude to an image
func MagnitudeToImage(magnitude [][]uint8) *image.RGBA {
	rows := len(magnitude)
	cols := len(magnitude[0])

	img := image.NewRGBA(image.Rect(0, 0, cols, rows))

	for y, row := range magnitude {
		for x, pixelValue := range row {
			img.Set(x, y, color.RGBA{pixelValue, pixelValue, pixelValue, 255})
		}
	}

	return img
}
