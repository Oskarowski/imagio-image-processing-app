package orthogonal_transforms

import (
	"image"
	"image/color"
	"math"
	"math/cmplx"
)

func VisualizeSpectrumInImage(data [][]complex128, width, height int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			magnitude := cmplx.Abs(data[y][x])
			intensity := uint8(255 * math.Log(1+magnitude) / math.Log(256))
			img.SetRGBA(x, y, color.RGBA{intensity, intensity, intensity, 255})
		}
	}

	return img
}

func VisualizeSpectrum(input [][]complex128) [][]float64 {
	N, M := len(input), len(input[0])
	output := make([][]float64, N)
	for i := range output {
		output[i] = make([]float64, M)
		for j := range output[i] {
			output[i][j] = cmplx.Abs(input[i][j])
		}
	}
	return output
}
