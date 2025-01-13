package orthogonal_transforms

import (
	"image"
	"image/color"
	"math"
	"math/cmplx"
)

func ComputeMagnitude(frequencySignal []complex128) []float64 {
	N := len(frequencySignal)
	magnitude := make([]float64, N)
	for i, value := range frequencySignal {
		magnitude[i] = math.Log(1 + cmplx.Abs(value))
	}
	return magnitude
}

// Calculate magnitude of a complex number
func magnitude(c complex128) float64 {
	return math.Sqrt(real(c)*real(c) + imag(c)*imag(c))
}

// Create a 2D FFT magnitude spectrum from the FFT result
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

	// Apply log scaling to enhance visualization
	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			magnitudeSpectrum[i][j] = math.Log(1 + magnitudeSpectrum[i][j])
		}
	}

	return magnitudeSpectrum
}

// Normalize the magnitude spectrum to fit in the range [0, 255]
func NormalizeMagnitude(magnitude [][]float64) [][]uint8 {
	rows := len(magnitude)
	cols := len(magnitude[0])

	// Find the min and max values in the magnitude spectrum
	var min, max float64
	min = math.Inf(1)
	max = -math.Inf(1)

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
	for i := 0; i < rows; i++ {
		normalized[i] = make([]uint8, cols)
		for j := 0; j < cols; j++ {
			normalized[i][j] = uint8((magnitude[i][j] - min) / (max - min) * 255)
		}
	}

	return normalized
}

// Convert the normalized magnitude to an image
func MagnitudeToImage(magnitude [][]uint8) *image.RGBA {
	rows := len(magnitude)
	cols := len(magnitude[0])

	img := image.NewRGBA(image.Rect(0, 0, cols, rows))

	for y := 0; y < rows; y++ {
		for x := 0; x < cols; x++ {
			pixelValue := uint8(magnitude[y][x])

			img.Set(x, y, color.RGBA{pixelValue, pixelValue, pixelValue, 255})
		}
	}

	return img
}
